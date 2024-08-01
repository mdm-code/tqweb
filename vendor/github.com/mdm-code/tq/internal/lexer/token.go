package lexer

import (
	"strings"

	"github.com/mdm-code/scanner"
)

const (
	// Undefined represents an undefined token type.
	Undefined TokenType = iota

	// String represents a string token type.
	String

	// Integer represents an integer token type.
	Integer

	// Dot represents a full stop token type.
	Dot

	// Colon represents a colon token type.
	Colon

	// ArrayOpen represents an opening bracket token type.
	ArrayOpen

	// ArrayClose represents a closing bracket token type.
	ArrayClose

	// Whitespace represents a white space token type.
	Whitespace
)

// keyCharMap maps runes onto TokenTypes.
var keyCharMap = map[rune]TokenType{
	'.': Dot,
	':': Colon,
	'[': ArrayOpen,
	']': ArrayClose,
}

// TokenType indicates the type of the lexer Token.
type TokenType uint8

// Token represents a single lexeme read from the Scanner token buffer.
type Token struct {
	Type       TokenType
	Buffer     *[]scanner.Token
	Start, End int
}

// Lexeme returns the string representation of the Token.
func (t Token) Lexeme() string {
	if t.Buffer == nil || len(*t.Buffer) < 1 || t.Start > t.End {
		return ""
	}
	end := t.End
	if end > len(*t.Buffer) {
		end = len(*t.Buffer)
	}
	chars := make([]string, end-t.Start)
	for _, t := range (*t.Buffer)[t.Start:end] {
		chars = append(chars, string(t.Rune))
	}
	return strings.Join(chars, "")
}
