package ast

// Visitor declares the interface for the AST visitor class. It declares
// signatures invoked by respective AST expression nodes.
type Visitor interface {
	VisitRoot(Expr)
	VisitQuery(Expr)
	VisitFilter(Expr)
	VisitIdentity(Expr)
	VisitSelector(Expr)
	VisitIterator(Expr)
	VisitSpan(Expr)
	VisitString(Expr)
	VisitInteger(Expr)
}
