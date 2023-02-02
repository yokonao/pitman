package main

import (
	"fmt"
	"os"
)

type Buffer struct {
	content string
	current int
}

func newBuffer(content string) Buffer {
	return Buffer{content: content}
}

func (b *Buffer) toChar(n int) byte {
	return b.content[b.current+n]
}

func (b *Buffer) toStr(n int) string {
	return b.content[b.current : b.current+n]
}

func (b *Buffer) next(n int) {
	b.current += n
}

func (b *Buffer) toTokens() []string {
	var res []string
	l := len(b.content)

	for {
		if b.current == l {
			break
		}

		if b.toChar(0) == '%' {
			for {
				b.next(1)
				if b.toChar(0) == '\n' {
					b.next(1)
					break
				}
			}
		} else if b.toChar(0) == ' ' || b.toChar(0) == '\n' {
			b.next(1)
		} else if b.toChar(0) == '(' {
			b.next(1)
			i := 0
			for {
				if b.toChar(i) == ')' {
					if i == 0 {
						res = append(res, "")
					} else {
						res = append(res, b.toStr(i))
						b.next(i + 1)
					}
					break
				}
				i++
			}
		} else {
			i := 1
			for {
				if b.content[b.current+i] == ' ' || b.content[b.current+i] == '\n' {
					res = append(res, b.toStr(i))
					b.next(i)
					break
				}
				i++
			}
		}
	}

	return res
}

var b Buffer

func main() {
	content, err := os.ReadFile("samples/sample.pdf")
	if err != nil {
		panic(err)
	}
	b = newBuffer(string(content))

	fmt.Println(b.toTokens())
}
