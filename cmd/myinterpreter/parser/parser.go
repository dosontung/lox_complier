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
	return parser.expression()
}

func (parser *Parser) nextToken() *tokenize.Token {
	if parser.current < len(parser.tokens) {
		parser.current++
	}
	return parser.tokens[parser.current-1]
}

func (parser *Parser) currentToken() *tokenize.Token {
	return parser.tokens[parser.current]
}

func (parser *Parser) isEnd() bool {
	return parser.current >= len(parser.tokens) || parser.tokens[parser.current].Type == tokenize.EOF
}

func (parser *Parser) expression() Expression {
	return parser.equality()
}

func (parser *Parser) equality() Expression {
	return parser.comparison()
}

func (parser *Parser) comparison() Expression {
	return parser.term()
}

func (parser *Parser) term() Expression {
	return parser.factor()
}
func (parser *Parser) primary() Expression {
	token := parser.nextToken()
	switch token.Type {
	case tokenize.TRUE:
		return &LiteralExpression{"true"}
	case tokenize.FALSE:
		return &LiteralExpression{"false"}
	case tokenize.NUMBER, tokenize.STRING:
		return &LiteralExpression{token.Literal}
	case tokenize.NIL:
		return &LiteralExpression{"nil"}
	case tokenize.LEFT_PAREN:
		expr := &GroupExpression{Expr: parser.expression()}
		if parser.nextToken().Type == tokenize.RIGHT_PAREN {
			return expr
		}

	}
	return &LiteralExpression{"null"}
}

func (parser *Parser) unary() Expression {
	token := parser.currentToken()
	switch token.Type {
	case tokenize.BANG, tokenize.MINUS:
		parser.nextToken()
		return &UnaryExpression{Operator: token, Right: parser.unary()}
	}
	return parser.primary()
}

func (parser *Parser) factor() Expression {
	expr := parser.unary()
	for !parser.isEnd() {
		token := parser.currentToken()
		switch token.Type {
		case tokenize.STAR, tokenize.SLASH:
			parser.nextToken()
			expr = &BinaryExpression{Left: expr, Right: parser.unary(), Operator: token}
		default:
			return expr
		}
	}
	return expr

}
