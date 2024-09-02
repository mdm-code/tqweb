package parser

import (
	"errors"

	"github.com/mdm-code/tq/internal/ast"
	"github.com/mdm-code/tq/internal/lexer"
)

// Parser encapsulates the logic of parsing tq queries into valid expressions.
type Parser struct {
	buffer  []lexer.Token
	current int
}

// New returns a new Parser with the buffer populated with lexer tokens read
// from the Lexer l.
func New(l *lexer.Lexer) (*Parser, error) {
	buf, ok := l.ScanAll(true)
	if !ok {
		err := errors.Join(l.Errors...)
		return nil, err
	}
	p := Parser{
		buffer:  buf,
		current: 0,
	}
	return &p, nil
}

// Parse the abstract syntax tree given the buffer of tq lexer tokens.
func (p *Parser) Parse() (*ast.Root, error) {
	root, err := p.root()
	return &root, err
}

func (p *Parser) root() (ast.Root, error) {
	q, err := p.query()
	expr := ast.Root{Query: &q}
	return expr, err
}

func (p *Parser) query() (ast.Query, error) {
	var expr ast.Query
	var err error
	for !p.isAtEnd() {
		var f ast.Filter
		f, err = p.filter()
		expr.Filters = append(expr.Filters, &f)
		if err != nil {
			break
		}
	}
	return expr, err
}

func (p *Parser) filter() (ast.Filter, error) {
	var expr ast.Filter
	var err error
	switch {
	case p.match(lexer.Dot):
		var i ast.Identity
		i, err = p.identity()
		expr.Kind = &i
	case p.match(lexer.ArrayOpen):
		var s ast.Selector
		s, err = p.selector()
		expr.Kind = &s
	default:
		v, _ := p.peek()
		err = &Error{v.Lexeme(), v.Buffer, v.Start, ErrQueryElement}
	}
	return expr, err
}

func (p *Parser) identity() (ast.Identity, error) {
	return ast.Identity{}, nil
}

func (p *Parser) selector() (ast.Selector, error) {
	var expr ast.Selector
	var err error
	switch {
	case p.check(lexer.ArrayClose):
		i, _ := p.iterator()
		expr.Value = &i
	case p.match(lexer.String):
		s, _ := p.string()
		expr.Value = &s
	case p.match(lexer.Colon):
		s, _ := p.span(nil)
		expr.Value = &s
	case p.match(lexer.Integer):
		i, _ := p.integer()
		if p.match(lexer.Colon) {
			s, _ := p.span(&i)
			expr.Value = &s
		} else {
			expr.Value = &i
		}
	}
	_, err = p.consume(lexer.ArrayClose, ErrSelectorUnterminated)
	return expr, err
}

func (p *Parser) iterator() (ast.Iterator, error) {
	return ast.Iterator{}, nil
}

func (p *Parser) string() (ast.String, error) {
	return ast.String{Value: p.previous().Lexeme()}, nil
}

func (p *Parser) integer() (ast.Integer, error) {
	return ast.Integer{Value: p.previous().Lexeme()}, nil
}

func (p *Parser) span(left *ast.Integer) (ast.Span, error) {
	s := ast.Span{Left: left}
	if p.match(lexer.Integer) {
		r, _ := p.integer()
		s.Right = &r
	}
	return s, nil
}

func (p *Parser) consume(t lexer.TokenType, e error) (lexer.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	v, err := p.peek()
	if err != nil {
		err := &Error{"EOL", v.Buffer, len(*v.Buffer), e}
		return lexer.Token{}, err
	}
	err = &Error{v.Lexeme(), v.Buffer, v.Start, e}
	return lexer.Token{}, err
}

func (p *Parser) match(tt ...lexer.TokenType) bool {
	for _, t := range tt {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	other, err := p.peek()
	if err != nil {
		return false
	}
	return other.Type == t
}

func (p *Parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	if p.current > len(p.buffer)-1 {
		return true
	}
	return false
}

func (p *Parser) previous() lexer.Token {
	return p.buffer[p.current-1]
}

func (p *Parser) peek() (lexer.Token, error) {
	if p.isAtEnd() {
		return p.previous(), ErrParserBufferOutOfRange
	}
	return p.buffer[p.current], nil
}
