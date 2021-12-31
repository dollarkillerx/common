package logger

import "strings"

func sanitize(s string) string {
	return strings.Map(sanitizeRune, s)
}

func sanitizeRune(r rune) rune {
	switch {
	case 'a' <= r && r <= 'z':
		return r
	case '0' <= r && r <= '9':
		return r
	case 'A' <= r && r <= 'Z':
		return r
	default:
		return '_'
	}
}
