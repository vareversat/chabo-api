package utils

import "unicode"

func StartWithUpperCase(s string) bool {
	return len(s) > 0 && isUpperRune(rune(s[0]))
}

func EndWithLowerCase(s string) bool {
	return len(s) > 0 && isLowerRune(rune(s[len(s)-1]))
}

func isUpperRune(r rune) bool {
	return unicode.IsUpper(r) && unicode.IsLetter(r)
}

func isLowerRune(r rune) bool {
	return unicode.IsLower(r) && unicode.IsLetter(r)
}
