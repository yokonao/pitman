package main

import (
	"fmt"
	"strconv"
	"strings"
)

type TokenType int

const (
	UnknownToken TokenType = iota
	RegularToken
	LiteralToken
	StreamToken
)

type Token struct {
	str       string
	tokenType TokenType
}

func newToken(str string, tokenType TokenType) *Token {
	return &Token{str: str, tokenType: tokenType}
}

func (t *Token) isLiteral() bool {
	return t.tokenType == LiteralToken
}

type TokenBuffer struct {
	tokens  []*Token
	current int
}

func (tb *TokenBuffer) String() string {
	var arr []string
	for _, t := range tb.tokens {
		arr = append(arr, t.str)
	}
	return fmt.Sprint(arr)
}

func (tb *TokenBuffer) isEmpty() bool {
	return tb.current >= len(tb.tokens)
}

func (tb *TokenBuffer) readToken() *Token {
	if tb.isEmpty() {
		panic("out of range")
	}

	token := tb.tokens[tb.current]
	tb.current++
	return token
}

func (tb *TokenBuffer) unreadToken() {
	if tb.current > 0 {
		tb.current--
	}
}

func (tb *TokenBuffer) mustNum() int {
	i, err := tb.expectNum()
	if err != nil {
		panic(err)
	}
	return i
}

func (tb *TokenBuffer) expectNum() (int, error) {
	s := tb.readToken()
	if s.isLiteral() {
		tb.unreadToken()
		return 0, fmt.Errorf("unexpecte token type")
	}

	i, err := strconv.Atoi(s.str)
	if err != nil {
		tb.unreadToken()
		return 0, err
	}
	return i, nil
}

func (tb *TokenBuffer) expectBool() (bool, error) {
	s := tb.readToken()
	if s.isLiteral() {
		tb.unreadToken()
		return false, fmt.Errorf("unexpecte token type")
	}

	switch s.str {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		tb.unreadToken()
		return false, fmt.Errorf("unexpected string: %s", s.str)
	}
}

func (tb *TokenBuffer) mustStr(cmp string) string {
	s, err := tb.expectStr(cmp)
	if err != nil {
		panic(err)
	}
	return s
}

func (tb *TokenBuffer) expectStr(cmp string) (string, error) {
	s := tb.readToken()
	if s.isLiteral() {
		tb.unreadToken()
		return "", fmt.Errorf("unexpecte token type")
	}

	if s.str != cmp {
		tb.unreadToken()
		return "", fmt.Errorf("unexpected string: %s, expected: %s", s.str, cmp)
	}
	return s.str, nil
}

func (tb *TokenBuffer) expectLiteral() (string, error) {
	s := tb.readToken()
	if !s.isLiteral() {
		tb.unreadToken()
		return "", fmt.Errorf("unexpecte token type")
	}
	return s.str, nil
}

func (tb *TokenBuffer) mustName() string {
	s, err := tb.expectName()
	if err != nil {
		panic(err)
	}
	return s
}

func (tb *TokenBuffer) expectName() (string, error) {
	s := tb.readToken()
	if s.isLiteral() {
		tb.unreadToken()
		return "", fmt.Errorf("unexpecte token type")
	}
	if !strings.HasPrefix(s.str, "/") {
		tb.unreadToken()
		return "", fmt.Errorf("unexpected prefix")
	}
	return s.str, nil
}

func (tb *TokenBuffer) expectRef() (*PDFReference, error) {
	ref, err := tb.expectNum()
	if err != nil {
		return nil, err
	}

	version, err := tb.expectNum()
	if err != nil {
		tb.unreadToken()
		return nil, err
	}

	_, err = tb.expectStr("R")
	if err != nil {
		tb.unreadToken()
		tb.unreadToken()
		return nil, err
	}

	return &PDFReference{ref: ref, version: version}, nil
}

func (tb *TokenBuffer) expectArrayElement() (interface{}, error) {
	ref, err := tb.expectRef()
	if err == nil {
		return ref, nil
	}
	name, err := tb.expectName()
	if err == nil {
		return name, nil
	}
	num, err := tb.expectNum()
	if err == nil {
		return num, nil
	}
	return tb.expectDict()
}

func (tb *TokenBuffer) expectArray() (*PDFArray, error) {
	cur := tb.current
	var err error
	defer func() {
		if err != nil {
			tb.current = cur
		}
	}()

	_, err = tb.expectStr("[")
	if err != nil {
		return nil, err
	}

	var arr []interface{}
	for {
		el, err := tb.expectArrayElement()
		if err != nil {
			break
		}
		arr = append(arr, el)
	}

	_, err = tb.expectStr("]")
	if err != nil {
		return nil, err
	}
	return &PDFArray{array: arr}, nil
}

func parseDictValue(tb *TokenBuffer) interface{} {
	ref, err := tb.expectRef()
	if err == nil {
		return ref
	}
	i, err := tb.expectNum()
	if err == nil {
		return i
	}

	b, err := tb.expectBool()
	if err == nil {
		return b
	}

	arr, err := tb.expectArray()
	if err == nil {
		return arr
	}

	name, err := tb.expectName()
	if err == nil {
		return name
	}

	literalStr, err := tb.expectLiteral()

	if err == nil {
		return literalStr
	}

	return parseDict(tb)
}

func parseDict(tb *TokenBuffer) *PDFDict {
	// ignore error
	dict, _ := tb.expectDict()
	return dict
}

func (tb *TokenBuffer) expectDict() (*PDFDict, error) {
	cur := tb.current
	var err error
	defer func() {
		if err != nil {
			tb.current = cur
		}
	}()

	d := map[string]interface{}{}
	_, err = tb.expectStr("<<")
	if err != nil {
		return nil, err
	}
	for {
		name, err := tb.expectName()
		if err != nil {
			break
		}
		d[name] = parseDictValue(tb)
	}
	tb.mustStr(">>")
	return &PDFDict{dict: d}, nil
}

func (tb *TokenBuffer) expectStream() (*PDFStream, error) {
	var err error

	_, err = tb.expectStr("stream")
	if err != nil {
		return nil, err
	}

	var tokens []string
	for {
		t := tb.readToken()
		if !t.isLiteral() && t.str == "endstream" {
			break
		}
		tokens = append(tokens, t.str)
	}

	return &PDFStream{tokens: tokens}, nil
}

func parseStream(tb *TokenBuffer) *PDFStream {
	// ignore error
	stream, _ := tb.expectStream()
	return stream
}

func parseObj(tb *TokenBuffer) *PDFObject {
	var dict *PDFDict
	var stream *PDFStream
	var array *PDFArray

	ref := tb.mustNum()
	version := tb.mustNum()
	tb.mustStr("obj")
	dict, err := tb.expectDict()

	if err != nil {
		array, err = tb.expectArray()
		if err != nil {
			panic(err)
		}
	} else {
		stream = parseStream(tb)
	}

	obj := &PDFObject{
		PDFReference: &PDFReference{ref: ref, version: version},
		dict:         dict,
		stream:       stream,
		array:        array,
	}
	tb.mustStr("endobj")
	return obj
}

func parseXRef(tb *TokenBuffer) *PDFXRef {
	i := tb.mustNum()
	j := tb.mustNum()

	if i != 0 {
		panic("expect 0")
	}

	xref := make([]string, j)
	for idx := 0; idx < j; idx++ {
		offset := tb.readToken()
		_ = tb.readToken()
		if idx == 0 {
			tb.mustStr("f")
		} else {
			tb.mustStr("n")
		}
		xref[idx] = offset.str
	}
	return &PDFXRef{xref: xref}
}

func parse(tb *TokenBuffer) *PDFDocument {
	doc := &PDFDocument{}
	var objects []*PDFObject
	for {
		if tb.isEmpty() {
			break
		}

		_, err := tb.expectStr("trailer")
		if err == nil {
			dict := parseDict(tb)
			doc.trailer = &PDFTrailer{dict: dict}
			continue
		}

		_, err = tb.expectStr("xref")
		if err == nil {
			doc.xref = parseXRef(tb)
			continue
		}

		_, err = tb.expectStr("startxref")
		if err == nil {
			if doc.xref == nil {
				panic("expect xref")
			}
			doc.startxref = tb.mustNum()
			continue
		}

		objects = append(objects, parseObj(tb))
	}
	doc.objects = objects
	return doc
}
