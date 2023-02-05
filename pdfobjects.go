package main

import (
	"fmt"
	"sort"
	"strings"
)

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
	var keys []string
	for k := range d.dict {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	builder.WriteString("<<\n")
	for _, k := range keys {
		builder.WriteString(
			fmt.Sprintln(k, d.dict[k]))
	}
	builder.WriteString(">>")
	return builder.String()
}

type PDFReference struct {
	ref     int
	version int
}

func (r *PDFReference) String() string {
	return fmt.Sprintf("ref: %d, version: %d", r.ref, r.version)
}

type PDFStream struct {
	tokens []string
}

func (s *PDFStream) String() string {
	return fmt.Sprint(s.tokens)
}

type PDFObject struct {
	*PDFReference
	dict   *PDFDict
	stream *PDFStream
	array  *PDFArray
}

func (obj *PDFObject) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintln(obj.PDFReference.String()))
	if obj.dict != nil {
		builder.WriteString(obj.dict.String())
		if obj.stream != nil {
			builder.WriteByte('\n')
			builder.WriteString(obj.stream.String())
		}
	}

	if obj.array != nil {
		builder.WriteString(obj.array.String())
	}

	return builder.String()
}

type PDFTrailer struct {
	dict *PDFDict
}

func (t *PDFTrailer) String() string {
	return fmt.Sprintf("trailer\n%s", t.dict.String())
}

type PDFDocument struct {
	objects []*PDFObject
	trailer *PDFTrailer
}

func (doc *PDFDocument) String() string {
	builder := strings.Builder{}
	for _, obj := range doc.objects {
		builder.WriteString(fmt.Sprintln(obj.String()))
		builder.WriteByte('\n')
	}
	builder.WriteString(doc.trailer.String())
	return builder.String()
}
