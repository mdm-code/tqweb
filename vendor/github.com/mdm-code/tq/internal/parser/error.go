package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mdm-code/scanner"
)

var (
	// ErrQueryElement indicates an unprocessable query filter element.
	ErrQueryElement = errors.New("expected '.' or '[' to parse query element")

	// ErrSelectorUnterminated indicates an unterminated selector element.
	ErrSelectorUnterminated = errors.New("expected ']' to terminate selector")

	// ErrParserBufferOutOfRange indicates the end of the parser buffer has
	// been reached.
	ErrParserBufferOutOfRange = errors.New("reached the end of the buffer")
)

// Error wraps a concrete parser error to represent its context. It reports the
// token where the error has occurred.
type Error struct {
	lexeme string
	buffer *[]scanner.Token
	offset int
	err    error
}

// Is allows to check if Error.err matches the target error.
func (e *Error) Is(target error) bool {
	return e.err == target
}

// Error reports the parser error wrapped inside of the custom context.
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

// wrapErrorLine wraps the error message inside the Parser buffer context.
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

// getErrorLine provides the error line with the nil error as default.
func (e *Error) getErrorLine() string {
	var b strings.Builder
	b.WriteString("Parser error: ")
	if e.err != nil {
		b.WriteString(e.err.Error())
		b.WriteString(fmt.Sprintf(" but got '%s'", e.lexeme))
		return b.String()
	}
	b.WriteString("nil")
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
