/*
Package tq encapsulates the logic of the Tq program behind a public interface.
It reads input data from the input reader, processes it with the interpreted
query string, and writes output data to the output writer.
*/
package tq

import (
	"fmt"
	"io"
	"strings"

	"github.com/mdm-code/scanner"
	"github.com/mdm-code/tq/internal/interpreter"
	"github.com/mdm-code/tq/internal/lexer"
	"github.com/mdm-code/tq/internal/parser"
	"github.com/mdm-code/tq/toml"
)

// Tq accepts TOML data from input and produces the result TOML data to output.
// The process of data decoding and encoding is handled by the adapter. The
// query passed to the Run method string is interpreted and executed against
// the input data to produce the output data.
type Tq struct {
	adapter toml.Adapter
}

// New returns a new Tq struct with the provided TOML adapter.
func New(adapter toml.Adapter) *Tq {
	return &Tq{
		adapter: adapter,
	}
}

// Run executes the query string against the input data and writes the output
// data to the output writer.
func (t *Tq) Run(input io.Reader, output io.Writer, query string) error {
	reader := strings.NewReader(query)
	scanner, err := scanner.New(reader)
	if err != nil {
		return err
	}
	lexer, err := lexer.New(scanner)
	if err != nil {
		return err
	}
	parser, err := parser.New(lexer)
	if err != nil {
		return err
	}
	ast, err := parser.Parse()
	if err != nil {
		return err
	}
	interpreter := interpreter.New()
	exec := interpreter.Interpret(ast)
	var data any
	err = t.adapter.Unmarshal(input, &data)
	if err != nil {
		return err
	}
	filteredData, err := exec(data)
	if err != nil {
		return err
	}
	for _, d := range filteredData {
		bytes, err := t.adapter.Marshal(d)
		if err != nil {
			return err
		}
		if len(bytes) == 0 {
			continue
		}
		fmt.Fprintln(output, string(bytes))
	}
	return nil
}
