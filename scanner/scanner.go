package scanner

import (
	"github.com/singurty/lox/token"
)

// Scanner transforms the source into tokens
type Scanner struct {
	source string
	tokens []token.Token
	start int
	current int
	line int
}

func New(source string) Scanner {
	scanner := Scanner{source: source, tokens: make([]token.Token, 0), start: 0, current: 0, line: 1}
	return scanner
}

func (sc *Scanner) ScanTokens() ([]token.Token) {
	for !sc.isAtEnd() {
		// beginning of the next lexeme
		sc.start = sc.current
		sc.scanToken()
	}
	sc.addToken(token.EOF)
	return sc.tokens
}

func (sc *Scanner) scanToken() (byte) {
	c := sc.source[sc.current]
	sc.current++
	switch (c) {
		// single-character tokens
		case '(':
			sc.addToken(token.LEFT_PAREN)
			break
		case ')':
			sc.addToken(token.RIGHT_PAREN)
			break
		case '{':
			sc.addToken(token.LEFT_BRACE)
			break
		case '}':
			sc.addToken(token.RIGHT_BRACE)
			break
		case ',':
			sc.addToken(token.COMMA)
			break
		case '.':
			sc.addToken(token.DOT)
			break
		case '-':
			sc.addToken(token.MINUS)
			break
		case '+':
			sc.addToken(token.PLUS)
			break
		case ';':
			sc.addToken(token.SEMICOLON)
			break
		case '*':
			sc.addToken(token.STAR)
			break
		default:
	}
	return c
}

func (sc *Scanner) addToken(Type token.Type) {
	sc.tokens = append(sc.tokens, token.Token{Type: Type})
}

func (sc *Scanner) isAtEnd() bool {
	return sc.current >= len(sc.source)
}

func (sc *Scanner) match(expected byte) bool {
	if (sc.isAtEnd()) {
		return false
	}
	if (sc.source[sc.current] != expected) {
		return false
	}
	sc.current++
	return true
}
