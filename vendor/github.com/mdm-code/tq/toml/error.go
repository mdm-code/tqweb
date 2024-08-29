package toml

import (
	"errors"
)

// ErrTOMLUnmarshal means an error encountered when unmarshalling the input.
var ErrTOMLUnmarshal = errors.New("failed to unmarshal TOML input")

// ErrTOMLMarshal mean an error was raised when marshalling the output.
var ErrTOMLMarshal = errors.New("failed to marshal TOML output")
