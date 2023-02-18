package main

// This file is based on PDF Reference 1.7 7.2.2 "Character Set".
// That chapter describes PDF character classification.

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
	case 0, '\t', '\n', 12, '\r', ' ': // 12 = FormFeed
		return true
	default:
		return false
	}
}