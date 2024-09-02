package toml

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/pelletier/go-toml/v2"
)

// Decoder defines the interface for decoding TOML data.
type Decoder interface {
	Decode(io.Reader, any) error
}

// Encoder defines the interface for encoding TOML data.
type Encoder interface {
	Encode(any) ([]byte, error)
}

// DecodeEncoder defines the combined interface for both decoding and encoding
// TOML data.
type DecodeEncoder interface {
	Decoder
	Encoder
}

// Adapter unifies the external TOML library interface to confine any changes
// to external libraries confined to a particular place in code.
type Adapter struct {
	adapted DecodeEncoder
}

// NewAdapter returns the adapted external library TOML functionalities.
func NewAdapter(adapted DecodeEncoder) Adapter {
	return Adapter{adapted: adapted}
}

// Unmarshal unmarshals the input r into the reference pointer argument passed
// to the parameter v.
func (a *Adapter) Unmarshal(r io.Reader, v any) error {
	err := a.adapted.Decode(r, v)
	if err != nil {
		err = errors.Join(ErrTOMLUnmarshal, err)
		err = fmt.Errorf("TOML error: %w", err)
	}
	return err
}

// Marshal marshals the argument passed to the parameter v to a slice of bytes.
func (a *Adapter) Marshal(v any) ([]byte, error) {
	bytes, err := a.adapted.Encode(v)
	if err != nil {
		err = errors.Join(ErrTOMLMarshal, err)
		err = fmt.Errorf("TOML error: %w", err)
	}
	return bytes, err
}

// GoTOML exposes the go-toml/v2 package functionality to that satisfies the
// decodeEncoder interface.
type GoTOML struct {
	conf GoTOMLConf
}

// NewGoTOML returns a struct exposing the go-toml/v2 package functionality to
// that satisfies the decodeEncoder interface.
func NewGoTOML(c GoTOMLConf) GoTOML {
	return GoTOML{conf: c}
}

// Decode decodes the input r into the reference pointer argument passed to the
// parameter v.
func (t GoTOML) Decode(r io.Reader, v any) error {
	d := toml.NewDecoder(r)
	return d.Decode(v)
}

// Encode encodes the argument passed to the parameter v as a slice of bytes.
func (t GoTOML) Encode(v any) ([]byte, error) {
	var buf bytes.Buffer
	e := toml.NewEncoder(&buf)
	e.SetTablesInline(t.conf.Encoder.TablesInline)
	e.SetArraysMultiline(t.conf.Encoder.ArraysMultiline)
	e.SetIndentSymbol(t.conf.Encoder.IndentSymbol)
	e.SetIndentTables(t.conf.Encoder.IndentTables)
	err := e.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GoTOMLConf specifies sensible configuration for the go-toml/v2 package.
type GoTOMLConf struct {
	Encoder struct {
		TablesInline    bool
		ArraysMultiline bool
		IndentSymbol    string
		IndentTables    bool
	}
}
