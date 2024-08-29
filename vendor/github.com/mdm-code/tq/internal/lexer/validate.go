package lexer

import "unicode"

// isNewline verifies if the rune r is a newline character.
func isNewline(r rune) bool {
	return r == '\n' || r == '\r'
}

// isWhitespace verifies if the rune r is a whitespace character.
func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isDigit verifies if the rune r is a digit character.
func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

// isKeyChar verifies if the rune r is a key character.
func isKeyChar(r rune) bool {
	_, ok := keyCharMap[r]
	return ok
}

// isQuote verifies if the rune r is a quote character.
func isQuote(r rune) bool {
	return r == '"' || r == '\''
}
