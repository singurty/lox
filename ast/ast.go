package ast

import (
	"fmt"
	"strings"

	"github.com/singurty/lox/token"
)

type Expr interface {
	String() string
}

type Ternary struct {
	Condition Expr
	Then Expr
	Else Expr
}

func (t *Ternary) String() string {
	var sb strings.Builder
	sb.WriteString("if (")
	sb.WriteString(t.Condition.String())
	sb.WriteString(") then (")
	sb.WriteString(t.Then.String())
	sb.WriteString(") else (")
	sb.WriteString(t.Else.String())
	sb.WriteString(")")
	return sb.String()
}

type Binary struct {
	Left Expr
	Operator token.Token
	Right Expr
}

// pretty print for binary
func (b *Binary) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(b.Operator.Lexeme)
	sb.WriteString(" ")
	sb.WriteString(b.Left.String())
	sb.WriteString(" ")
	sb.WriteString(b.Right.String())
	sb.WriteString(")")
	return sb.String()
}

// for parenthesized expressions
type Grouping struct {
	Expression Expr
}

func (g *Grouping) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(g.Expression.String())
	sb.WriteString(")")
	return sb.String()
}

type Literal struct {
	Value interface{}
}

func (l *Literal) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%v", l.Value))
	return sb.String()
}

type Unary struct {
	Operator token.Token
	Right Expr
}

func (u *Unary) String() string {
	var sb strings.Builder
	sb.WriteString(u.Operator.Lexeme)
	sb.WriteString(u.Right.String())
	return sb.String()
}

type Stmt interface {
}

type ExprStmt struct {
	Expression Expr
}

type PrintStmt struct {
	Expression Expr
}
