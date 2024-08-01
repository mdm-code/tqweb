package lexer

import (
	"errors"
	"strings"

	"github.com/mdm-code/scanner"
)

var (
	// ErrNilScanner indicates that the provided Scanner is nil.
	ErrNilScanner = errors.New("provided Scanner is nil")

	// ErrKeyCharUnsupported indicates that the key character is unsupported.
	ErrKeyCharUnsupported = errors.New("unsupported key character")

	// ErrUnterminatedString indicates that the string literal is not terminated.
	ErrUnterminatedString = errors.New("unterminated string literal")

	// ErrDisallowedChar indcates the the character is disallowed.
	ErrDisallowedChar = errors.New("disallowed character")
)

// Error wraps a concrete lexer error to represent its query context in the
// error message. It stores references to the Lexer buffer context and the
// Lexer token start offset.
type Error struct {
	buffer *[]scanner.Token // Lexer buffer context pointer
	offset int              // Lexer token start offset
	err    error            // wrapped Lexer error
}

// Is allows to check if Error.err matches the target error.
func (e *Error) Is(target error) bool {
	return e.err == target
}

// Error reports the Lexer error wrapped inside the Lexer buffer context with
// a marker indicating the start of the Lexer token at which the occurred.
func (e *Error) Error() string {
	line := e.getErrorLine()
	if e.buffer == nil || len(*e.buffer) < 1 {
		return line
	}
	pointer := "^"
	indentChar := " "
	result := e.wrapErrorLine(line, pointer, indentChar)
	return result
}

// getErrorLine provides the error line with the nil error as default.
func (e *Error) getErrorLine() string {
	var b strings.Builder
	b.WriteString("Lexer error: ")
	if e.err != nil {
		b.WriteString(e.err.Error())
		return b.String()
	}
	b.WriteString("nil")
	return b.String()
}

// wrapErrorLine wraps the error message inside the Lexer buffer context.
func (e *Error) wrapErrorLine(line, pointer, indentChar string) string {
	var b strings.Builder
	b.Grow(len(*e.buffer)*2 + 1)
	for _, t := range *e.buffer {
		b.WriteRune(t.Rune)
	}
	b.WriteString("\n")
	indent := e.getIndent(indentChar)
	b.WriteString(indent)
	b.WriteString(pointer)
	b.WriteString("\n")
	b.WriteString(line)
	return b.String()
}

// getIndent constructs the pointer indentation. Negative offsets result in an
// empty string.
func (e *Error) getIndent(indentChar string) string {
	if e.offset > 0 {
		return strings.Repeat(indentChar, e.offset)
	}
	return ""
}
