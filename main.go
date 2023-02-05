package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func isSpace(c byte) bool {
	switch c {
	case '\n', '\r', ' ':
		return true
	}
	return false
}

type Buffer struct {
	content string
	current int
}

func newBuffer(content string) Buffer {
	return Buffer{content: content}
}

func (b *Buffer) isEOF() bool {
	return b.current >= len(b.content)
}

func (b *Buffer) readChar() byte {
	if b.isEOF() {
		panic("out of range")
	}

	c := b.content[b.current]
	b.current++
	return c
}

func (b *Buffer) unreadChar() {
	if b.current > 0 {
		b.current--
	}
}

func (b *Buffer) readLiteralStr() string {
	builder := strings.Builder{}
	for {
		c := b.readChar()
		// TODO "(())"のパターンを考慮する必要がある
		if c == ')' {
			break
		} else {
			builder.WriteByte(c)
		}
	}
	return builder.String()
}

func (b *Buffer) readStr() string {
	builder := strings.Builder{}
	for {
		c := b.readChar()
		if isSpace(c) {
			break
		}
		builder.WriteByte(c)
	}
	return builder.String()
}

func (b *Buffer) toTokens() []string {
	var res []string

	for {
		if b.isEOF() {
			break
		}

		c := b.readChar()
		switch c {
		case '%':
			for {
				t := b.readChar()
				if t == '\n' {
					break
				}
			}
		case '(':
			res = append(res, b.readLiteralStr())
		default:
			if isSpace(c) {
				continue
			}
			b.unreadChar()
			res = append(res, b.readStr())
		}
	}

	return res
}

type PDFArray struct {
	array []interface{}
}

func (arr *PDFArray) String() string {
	return fmt.Sprint(arr.array...)
}

type PDFDict struct {
	// key: literal
	// value: PDFDict | PDFNode | Reference | literal | number | array
	dict map[string]interface{}
}

func (d *PDFDict) String() string {
	builder := strings.Builder{}
	for k, v := range d.dict {
		builder.WriteString(
			fmt.Sprintln(k, v))
	}
	return builder.String()
}

type PDFReference struct {
	ref     int
	version int
}

func (r *PDFReference) String() string {
	return fmt.Sprintf("ref: %d, version: %d", r.ref, r.version)
}

type PDFObject struct {
	*PDFReference
	dict *PDFDict
}

func (obj *PDFObject) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintln(obj.PDFReference.String()))
	builder.WriteString(obj.dict.String())
	return builder.String()
}

type Tokens struct {
	tokens  []string
	current int
}

func (t *Tokens) readToken() string {
	// TODO out of range
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
		return "", fmt.Errorf("unexpected string")
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
	// [0 0 R 1 0 R]
	ref, err := t.expectRef()
	if err == nil {
		return ref, nil
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

	return t.mustName()
}

func parseDict(t *Tokens) *PDFDict {
	d := map[string]interface{}{}
	t.mustStr("<<")
	for {
		name, err := t.expectName()
		if err != nil {
			break
		}
		d[name] = parseDictValue(t)
	}
	t.mustStr(">>")
	return &PDFDict{dict: d}
}

func parseObj(t *Tokens) *PDFObject {
	ref := t.mustNum()
	version := t.mustNum()
	t.mustStr("obj")
	dict := parseDict(t)
	obj := &PDFObject{PDFReference: &PDFReference{ref: ref, version: version}, dict: dict}
	t.mustStr("endobj")
	fmt.Println(obj)
	return obj
}

func parse(t *Tokens) []*PDFObject {
	parseObj(t)
	parseObj(t)
	parseObj(t)

	return nil
}

func main() {
	content, err := os.ReadFile("samples/sample1.pdf")
	if err != nil {
		panic(err)
	}
	b := newBuffer(string(content))

	tokens := &Tokens{tokens: b.toTokens()}
	fmt.Println(tokens.tokens)
	parse(tokens)
}
