package main

// This file is based on PDF Reference 1.7 7.2.2 "Character Set".
// That chapter describes PDF character classification.
// @see https://opensource.adobe.com/dc-acrobat-sdk-docs/pdfstandards/PDF32000_2008.pdf

func isRegularChar(c byte) bool {
	return !isDelemiterChar(c) && !isWhiteSpaceChar(c)
}

func isDelemiterChar(c byte) bool {
	switch c {
	case '(', ')', '<', '>', '[', ']', '{', '}', '/', '%':
		return true
	default:
		return false
	}
}
func isWhiteSpaceChar(c byte) bool {
	switch c {
	case 0, '\t', '\n', 12, '\r', ' ': // 12 = 改ページ
		return true
	default:
		return false
	}
}
