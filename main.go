package main

import (
	"fmt"
	"os"
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


func main() {
	content, err := os.ReadFile("samples/sample1.pdf")
	if err != nil {
		panic(err)
	}
	b := newBuffer(string(content))

	fmt.Println(b.toTokens())


}
