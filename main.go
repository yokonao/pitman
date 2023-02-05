package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	f := flag.String("path", "sample1", "")
	flag.Parse()
	path := fmt.Sprintf("samples/%s.pdf", *f)
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	b := newBuffer(string(content))

	tokens := &TokenBuffer{tokens: b.toTokenBuffer()}
	doc := parse(tokens)
	fmt.Print(doc)
}
