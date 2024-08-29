package ast

import (
	"fmt"
	"strconv"
	"strings"
)

// Expr defines the expression interface for the visitor to operate on the
// contents of the expression node.
type Expr interface {
	Accept(v Visitor)
}

// Root stands for the top-level root node of the tq query. This version of
// the tq parser allows a single query, but extending the root to span
// multiple queries in always an option.
type Root struct {
	Query Expr
}

// Query represents a single tq query that can be run against a de-serialized
// TOML data object. It potentially comprises of zero or more filters used to
// filter the TOML data. Although filters are stored in a slice implying a
// sequence, the order in not enforced neither by the expression nor the
// parser. It is the responsibility of the visiting interpreter run against the
// AST to provide the filtering mechanism.
type Query struct {
	Filters []Expr
}

// Filter stands for a single tq filter. It the fundamental building block of
// the tq query.
type Filter struct {
	Kind Expr
}

// Identity specifies the identity data transformation that returns the
// filtered data argument unchanged.
type Identity struct{}

// Selector represents a select-driven data filter.
type Selector struct {
	Value Expr
}

// Span represents a filter that takes a slice of a list-like sequence.
type Span struct {
	Left, Right *Integer
}

// Iterator represents a sequeced iterator. The implementation of the iterator
// for TOML data types is to be provided by the visiting interpreter.
type Iterator struct{}

// String represents the key selector that can be used, for instance, in a form
// of dictionary lookup.
type String struct {
	Value string
}

// Integer represents your everyday integer. It can be used, for example, as an
// index of a data point in a sequence or a start/stop index of a span.
type Integer struct {
	Value string
}

// Accept implements the Expr interface for the visitor design pattern.
func (r *Root) Accept(v Visitor) {
	v.VisitRoot(r)
}

// String provides the string representation of the AST expression.
func (*Root) String() string {
	return "root"
}

// Accept implements the Expr interface for the visitor design pattern.
func (q *Query) Accept(v Visitor) {
	v.VisitQuery(q)
}

// String provides the string representation of the AST expression.
func (*Query) String() string {
	return "query"
}

// Accept implements the Expr interface for the visitor design pattern.
func (f *Filter) Accept(v Visitor) {
	v.VisitFilter(f)
}

// String provides the string representation of the AST expression.
func (*Filter) String() string {
	return "filter"
}

// Accept implements the Expr interface for the visitor design pattern.
func (i *Identity) Accept(v Visitor) {
	v.VisitIdentity(i)
}

// String provides the string representation of the AST expression.
func (*Identity) String() string {
	return "identity"
}

// Accept implements the Expr interface for the visitor design pattern.
func (s *Selector) Accept(v Visitor) {
	v.VisitSelector(s)
}

// String provides the string representation of the AST expression.
func (*Selector) String() string {
	return "selector"
}

// Accept implements the Expr interface for the visitor design pattern.
func (s *Span) Accept(v Visitor) {
	v.VisitSpan(s)
}

// String provides the string representation of the AST expression.
func (s *Span) String() string {
	var l, r string
	if s.Left != nil {
		l = s.Left.Value
	}
	if s.Right != nil {
		r = s.Right.Value
	}
	return fmt.Sprintf("span [%s:%s]", l, r)
}

// GetLeft returns the value of the left-hand side expression node of the Span.
func (s *Span) GetLeft(def int) int {
	return s.asInt(s.Left, def)
}

// GetRight returns the value of the right-hand side expression node of the Span.
func (s *Span) GetRight(def int) int {
	return s.asInt(s.Right, def)
}

func (s *Span) asInt(i *Integer, def int) int {
	var result = def
	if i != nil {
		integer, err := i.Vtoi()
		if err != nil {
			return result
		}
		result = integer
	}
	return result
}

// Accept implements the Expr interface for the visitor design pattern.
func (i *Iterator) Accept(v Visitor) {
	v.VisitIterator(i)
}

// String provides the string representation of the AST expression.
func (i *Iterator) String() string {
	return "iterator"
}

// Accept implements the Expr interface for the visitor design pattern.
func (s *String) Accept(v Visitor) {
	v.VisitString(s)
}

// String provides the string representation of the AST expression.
func (s *String) String() string {
	return fmt.Sprintf("string %q", s.Trim())
}

// Trim returns the value of the String expression node with the surrounding
// quotation marks stripped off.
func (s *String) Trim() string {
	result := s.Value
	for _, c := range `'"` {
		result = strings.Trim(result, string(c))
	}
	return result
}

// Accept implements the Expr interface for the visitor design pattern.
func (i *Integer) Accept(v Visitor) {
	v.VisitInteger(i)
}

// String provides the string representation of the AST expression.
func (i *Integer) String() string {
	return fmt.Sprintf("integer %s", i.Value)
}

// Vtoi returns the value of the Integer expression node converted to a value
// of the type int.
func (i *Integer) Vtoi() (int, error) {
	return strconv.Atoi(i.Value)
}
