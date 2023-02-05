package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
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

func (b *Buffer) isEOF()bool{
	return b.current >= len(b.content)
}

func (b *Buffer) readChar() byte {
	if b.isEOF(){
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
		if b.isEOF(){
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

type PDFDict struct{
	// key: literal 
	// value: PDFDict | PDFNode | Reference | literal | number | array 
	dict map[string] interface{}
}

type PDFObject struct{
	ref int
	version int
	dict PDFDict
}

type Tokens struct {
	tokens []string
	current int
}

func (t *Tokens) readToken()string{
	// TODO out of range
	token := t.tokens[t.current]
	t.current++
	return token
}

func (t *Tokens) mustNum()int{
	i, err := strconv.Atoi(t.readToken())
	if err!= nil{
		panic(err)
	}
	return i
}

func (t *Tokens) mustStr(cmp string)string{
	s := t.readToken()
	if s != cmp{
		panic("unexpected string'")
	}
	return s
}

func (t *Tokens) mustName()string{
	s := t.readToken()
	strings.Has
	if s != {
		panic("unexpected string'")
	}
	return s
}

func parseDict(t *Tokens) *PDFDict{
	mustStr("<<")
	
}

func parse(t *Tokens) []PDFObject{
	ref := t.mustNum()
	version := t.mustNum()
	t.mustStr("obj")
	//obj := &PDFObject{ref: ref, version: version}


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
