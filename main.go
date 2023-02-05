package main

import (
	"fmt"
	"os"
)

func main() {
	content, err := os.ReadFile("samples/sample1.pdf")
	if err != nil {
		panic(err)
	}
	b := newBuffer(string(content))

	tokens := &Tokens{tokens: b.toTokens()}
	fmt.Println(tokens.tokens)
	doc := parse(tokens)
	fmt.Println(doc)
}
