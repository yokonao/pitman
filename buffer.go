package main

import (
	"fmt"
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
		// str として許容されない文字列
		if c == ']' || c == '/' || c == '>' {
			b.unreadChar()
			break
		}
		builder.WriteByte(c)
	}
	return builder.String()
}

func (b *Buffer) readName() string {
	builder := strings.Builder{}
	builder.WriteByte('/')
	for {
		c := b.readChar()
		if !isNameChar(c) {
			b.unreadChar()
			break
		}
		builder.WriteByte(c)
	}
	return builder.String()
}

// Name として許容される文字列かどうか a-z, A-Z, 0-9, - を確認している
func isNameChar(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z' || ('0' <= c && c <= '9')) || c == '-'
}

func (b *Buffer) toTokenBuffer() []*Token {
	var res []*Token

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
			res = append(res, newToken(b.readLiteralStr(), true))
		case '[', ']':
			res = append(res, newToken(string(c), false))
		case '<':
			c2 := b.readChar()
			if c2 == '<' {
				res = append(res, newToken("<<", false))
			} else {
				panic(fmt.Sprintf("unexpected character %c", c2))
			}
		case '>':
			c2 := b.readChar()
			if c2 == '>' {
				res = append(res, newToken(">>", false))
			} else {
				panic(fmt.Sprintf("unexpected character %c", c2))
			}
		case '/':
			name := b.readName()
			res = append(res, newToken(name, false))
		default:
			if isSpace(c) {
				continue
			}
			b.unreadChar()
			res = append(res, newToken(b.readStr(), false))
		}
	}

	return res
}
