package interpreter

import (
	"fmt"

	"github.com/mdm-code/tq/internal/ast"
)

// FilterFunc specifies the data transformation function type.
type FilterFunc func(data ...any) ([]any, error)

type filter struct {
	name  string
	inner FilterFunc
}

func (f *filter) call(data ...any) ([]any, error) {
	return f.inner(data...)
}

// Interpreter interprets the tq query AST into a pipe-like sequence of
// filtering functions processing TOML input data as specified in the query.
type Interpreter struct {
	filters []filter
}

// New returns a new instance of Interpreter.
func New() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) eval(es ...ast.Expr) {
	for _, e := range es {
		e.Accept(i)
	}
}

// Interpret extracts a sequence of filtering functions by traversing the AST.
// It returns an entry function that takes in deserialized TOML data and
// applies filtering functions in the sequence provided by the Interpreter.
func (i *Interpreter) Interpret(root ast.Expr) FilterFunc {
	i.filters = nil // clear out previously accumulated filtering functions
	i.eval(root)
	return func(data ...any) ([]any, error) {
		var err error
		for _, f := range i.filters {
			data, err = f.call(data...)
			if err != nil {
				return data, err
			}
		}
		return data, nil
	}
}

// VisitRoot interprets the Root AST node.
func (i *Interpreter) VisitRoot(e ast.Expr) {
	r := e.(*ast.Root)
	i.eval(r.Query)
}

// VisitQuery interprets the Query AST node.
func (i *Interpreter) VisitQuery(e ast.Expr) {
	q := e.(*ast.Query)
	i.eval(q.Filters...)
}

// VisitFilter interprets the Filter AST node.
func (i *Interpreter) VisitFilter(e ast.Expr) {
	f := e.(*ast.Filter)
	i.eval(f.Kind)
}

// VisitIdentity interprets the Identity AST node.
func (i *Interpreter) VisitIdentity(e ast.Expr) {
	f := filter{
		name: "identity",
		inner: func(data ...any) ([]any, error) {
			return data, nil
		},
	}
	i.filters = append(i.filters, f)
}

// VisitSelector interprets the Selector AST node.
func (i *Interpreter) VisitSelector(e ast.Expr) {
	s := e.(*ast.Selector)
	i.eval(s.Value)
}

// VisitSpan interprets the Span AST node.
func (i *Interpreter) VisitSpan(e ast.Expr) {
	span := e.(*ast.Span)
	f := filter{
		name: "span",
		inner: func(data ...any) ([]any, error) {
			result := make([]any, 0, len(data))
			var err error
			for _, d := range data {
				switch v := d.(type) {
				case []any:
					l, r := span.GetLeft(0), span.GetRight(len(v))
					if r > len(v) {
						r = len(v)
					}
					if l > r || l >= len(v) {
						continue
					}
					result = append(result, v[l:r])
				default:
					err = &Error{
						data:   d,
						filter: fmt.Sprintf("%s", span),
						err:    ErrTOMLDataType,
					}
				}
			}
			return result, err
		},
	}
	i.filters = append(i.filters, f)
}

// VisitIterator interprets the Iterator AST node.
func (i *Interpreter) VisitIterator(e ast.Expr) {
	iter := e.(*ast.Iterator)
	f := filter{
		name: "iterator",
		inner: func(data ...any) ([]any, error) {
			result := make([]any, 0, len(data))
			var err error
			for _, d := range data {
				switch v := d.(type) {
				case map[string]any:
					for _, val := range v {
						result = append(result, val)
					}
				case []any:
					for _, val := range v {
						result = append(result, val)
					}
				default:
					err = &Error{
						data:   d,
						filter: fmt.Sprintf("%s", iter),
						err:    ErrTOMLDataType,
					}
				}
			}
			return result, err
		},
	}
	i.filters = append(i.filters, f)
}

// VisitString interprets the String AST node.
func (i *Interpreter) VisitString(e ast.Expr) {
	str := e.(*ast.String)
	f := filter{
		name: "string",
		inner: func(data ...any) ([]any, error) {
			result := make([]any, 0, len(data))
			var err error
			for _, d := range data {
				switch v := d.(type) {
				case map[string]any:
					key := str.Trim()
					res, ok := v[key]
					if ok {
						result = append(result, res)
					}
				default:
					err = &Error{
						data:   d,
						filter: fmt.Sprintf("%s", str),
						err:    ErrTOMLDataType,
					}
				}
			}
			return result, err
		},
	}
	i.filters = append(i.filters, f)
}

// VisitInteger interprets the Integer AST node.
func (i *Interpreter) VisitInteger(e ast.Expr) {
	integer := e.(*ast.Integer)
	f := filter{
		name: "integer",
		inner: func(data ...any) ([]any, error) {
			result := make([]any, 0, len(data))
			var err error
			for _, d := range data {
				switch v := d.(type) {
				case []any:
					idx, _ := integer.Vtoi()
					if idx >= 0 && idx < len(v) {
						result = append(result, v[idx])
					}
				default:
					err = &Error{
						data:   d,
						filter: fmt.Sprintf("%s", integer),
						err:    ErrTOMLDataType,
					}
				}
			}
			return result, err
		},
	}
	i.filters = append(i.filters, f)
}
