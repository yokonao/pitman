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
	token string
}

func (s *PDFStream) String() string {
	return fmt.Sprint(s.token)
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
	return fmt.Sprintf("trailer\n%s\n\n", t.dict.String())
}

type PDFXRef struct {
	xref []string
}

func (xref *PDFXRef) String() string {
	builder := strings.Builder{}
	builder.WriteString("xref\n")
	for i := 0; i < len(xref.xref); i++ {
		builder.WriteString(fmt.Sprintln(i, xref.xref[i]))

	}
	return builder.String()
}

type PDFDocument struct {
	objects   []*PDFObject
	xref      *PDFXRef
	trailer   *PDFTrailer
	startxref int
}

func (doc *PDFDocument) String() string {
	builder := strings.Builder{}
	for _, obj := range doc.objects {
		builder.WriteString(fmt.Sprintln(obj.String()))
		builder.WriteByte('\n')
	}
	if doc.xref != nil {
		builder.WriteString(doc.xref.String())
	}
	if doc.trailer != nil {
		builder.WriteString(doc.trailer.String())
	}
	builder.WriteString(fmt.Sprintf("startxref %d", doc.startxref))
	return builder.String()
}
