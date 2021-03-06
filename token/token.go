package token

import "fmt"

type Type int

const (
	// Single-character tokens
	LEFT_PAREN = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR
	QUESTION_MARK
	COLON

	// One or two character tokens
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals
	IDENTIFIER
	STRING
	NUMBER

	// Keywords
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NULL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
	BREAK
	CONTINUE

	EOF
)

type Token struct {
	Type Type
	Lexeme string
	Literal interface{}
	Line int
}

func (token *Token) String() string {
	return fmt.Sprintf("%v %v %v", token.Type, token.Lexeme, token.Literal)
}
