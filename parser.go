package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Tokens struct {
	tokens  []string
	current int
}

func (t *Tokens) isEmpty() bool {
	return t.current >= len(t.tokens)
}

func (t *Tokens) readToken() string {
	if t.isEmpty() {
		panic("out of range")
	}

	token := t.tokens[t.current]
	t.current++
	return token
}

func (t *Tokens) unreadToken() {
	if t.current > 0 {
		t.current--
	}
}

func (t *Tokens) mustNum() int {
	i, err := t.expectNum()
	if err != nil {
		panic(err)
	}
	return i
}

func (t *Tokens) expectNum() (int, error) {
	i, err := strconv.Atoi(t.readToken())
	if err != nil {
		t.unreadToken()
		return 0, err
	}
	return i, nil
}

func (t *Tokens) mustStr(cmp string) string {
	s, err := t.expectStr(cmp)
	if err != nil {
		panic(err)
	}
	return s
}

func (t *Tokens) expectStr(cmp string) (string, error) {
	s := t.readToken()
	if s != cmp {
		t.unreadToken()
		return "", fmt.Errorf("unexpected string: %s", s)
	}
	return s, nil
}

func (t *Tokens) mustName() string {
	s, err := t.expectName()
	if err != nil {
		panic(err)
	}
	return s
}

func (t *Tokens) expectName() (string, error) {
	s := t.readToken()
	if !strings.HasPrefix(s, "/") {
		t.unreadToken()
		return "", fmt.Errorf("unexpected prefix")
	}
	return s, nil
}

func (t *Tokens) expectRef() (*PDFReference, error) {
	ref, err := t.expectNum()
	if err != nil {
		return nil, err
	}

	version, err := t.expectNum()
	if err != nil {
		t.unreadToken()
		return nil, err
	}

	_, err = t.expectStr("R")
	if err != nil {
		t.unreadToken()
		t.unreadToken()
		return nil, err
	}

	return &PDFReference{ref: ref, version: version}, nil
}

func (t *Tokens) expectArrayElement() (interface{}, error) {
	ref, err := t.expectRef()
	if err == nil {
		return ref, nil
	}
	name, err := t.expectName()
	if err == nil {
		return name, nil
	}
	return t.expectNum()
}

func (t *Tokens) expectArray() (*PDFArray, error) {
	cur := t.current
	var err error
	defer func() {
		if err != nil {
			t.current = cur
		}
	}()

	_, err = t.expectStr("[")
	if err != nil {
		return nil, err
	}

	var arr []interface{}
	for {
		el, err := t.expectArrayElement()
		if err != nil {
			break
		}
		arr = append(arr, el)
	}

	_, err = t.expectStr("]")
	if err != nil {
		return nil, err
	}
	return &PDFArray{array: arr}, nil
}

func parseDictValue(t *Tokens) interface{} {
	ref, err := t.expectRef()
	if err == nil {
		return ref
	}
	i, err := t.expectNum()
	if err == nil {
		return i
	}

	arr, err := t.expectArray()
	if err == nil {
		return arr
	}

	name, err := t.expectName()
	if err == nil {
		return name
	}

	return parseDict(t)
}

func parseDict(t *Tokens) *PDFDict {
	// ignore error
	dict, _ := t.expectDict()
	return dict
}

func (t *Tokens) expectDict() (*PDFDict, error) {
	cur := t.current
	var err error
	defer func() {
		if err != nil {
			t.current = cur
		}
	}()

	d := map[string]interface{}{}
	_, err = t.expectStr("<<")
	if err != nil {
		return nil, err
	}
	for {
		name, err := t.expectName()
		if err != nil {
			break
		}
		d[name] = parseDictValue(t)
	}
	t.mustStr(">>")
	return &PDFDict{dict: d}, nil
}

func (t *Tokens) expectStream() (*PDFStream, error) {
	var err error

	_, err = t.expectStr("stream")
	if err != nil {
		return nil, err
	}

	var tokens []string
	for {
		t := t.readToken()
		if t == "endstream" {
			break
		}
		tokens = append(tokens, t)
	}

	return &PDFStream{tokens: tokens}, nil
}

func parseStream(t *Tokens) *PDFStream {
	// ignore error
	stream, _ := t.expectStream()
	return stream
}

func parseObj(t *Tokens) *PDFObject {
	var dict *PDFDict
	var stream *PDFStream
	var array *PDFArray

	ref := t.mustNum()
	version := t.mustNum()
	t.mustStr("obj")
	dict, err := t.expectDict()

	if err != nil {
		array, err = t.expectArray()
		if err != nil {
			panic(err)
		}
	} else {
		stream = parseStream(t)
	}

	obj := &PDFObject{
		PDFReference: &PDFReference{ref: ref, version: version},
		dict:         dict,
		stream:       stream,
		array:        array,
	}
	t.mustStr("endobj")
	return obj
}

func parse(t *Tokens) *PDFDocument {
	var objects []*PDFObject
	var trailer *PDFTrailer
	for {
		_, err := t.expectStr("trailer")
		if err == nil {
			dict := parseDict(t)
			trailer = &PDFTrailer{dict: dict}
			break
		}

		objects = append(objects, parseObj(t))
	}

	return &PDFDocument{objects: objects, trailer: trailer}
}
