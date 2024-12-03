package parser

import (
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"
)

type Parser struct {
	tokens  []*tokenize.Token
	current int
}

func NewParser(tokens []*tokenize.Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

func (parser *Parser) Parse() Expression {
	return parser.primary()
}

func (parser *Parser) currentToken() *tokenize.Token {
	if parser.current < len(parser.tokens) {
		parser.current++
	}
	return parser.tokens[parser.current-1]
}

func (parser *Parser) isEnd() bool {
	return parser.current >= len(parser.tokens) || parser.tokens[parser.current].Type == tokenize.EOF
}

func (parser *Parser) primary() Expression {
	token := parser.currentToken()
	switch token.Type {
	case tokenize.TRUE:
		return &LiteralExpression{"true"}
	case tokenize.FALSE:
		return &LiteralExpression{"false"}
	case tokenize.NUMBER:
		return &LiteralExpression{token.Literal}
	case tokenize.STRING:
		return &LiteralExpression{token.Literal}
	case tokenize.NIL:
		return &LiteralExpression{"nil"}

	}
	return &LiteralExpression{"null"}
}