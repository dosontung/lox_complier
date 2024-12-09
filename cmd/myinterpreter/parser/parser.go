package parser

import (
	"errors"
	"fmt"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/core"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"
	"os"
)

/*
program        → declaration * EOF ;
declaration    → varDecl | statement ;
varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
statement      → exprStmt | printStmt ;
expression     → assignment ;
assignment     → IDENTIFIER "=" assignment
               | equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               |    primary;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")"
               | IDENTIFIER ;
*/

type Parser struct {
	tokens []*tokenize.Token
	//env     *Environment
	current int
	err     error
}

func NewParser(tokens []*tokenize.Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

func (parser *Parser) nextToken() *tokenize.Token {
	if parser.current < len(parser.tokens) {
		parser.current++
	}
	return parser.tokens[parser.current-1]
}

func (parser *Parser) currentToken() *tokenize.Token {
	if parser.current >= len(parser.tokens) {
		return parser.tokens[len(parser.tokens)-1]
	}
	return parser.tokens[parser.current]
}

func (parser *Parser) peekToken() *tokenize.Token {
	if parser.current < len(parser.tokens) {
		return parser.tokens[parser.current+1]
	}
	return parser.tokens[len(parser.tokens)-1]
}

func (parser *Parser) previous() {
	if parser.current > 0 {
		parser.current = parser.current - 1
	}

}
func (parser *Parser) Error() error {
	return parser.err
}
func (parser *Parser) isEnd() bool {
	return parser.current >= len(parser.tokens) || parser.tokens[parser.current].Type == tokenize.EOF
}

func (parser *Parser) Parse() core.Expression {
	return parser.expression()
}

func (parser *Parser) ParseStmt() []core.Statement {
	out := make([]core.Statement, 0)
	for !parser.isEnd() {
		out = append(out, parser.declaration())
	}
	return out
}

func (parser *Parser) declaration() core.Statement {
	token := parser.currentToken()
	if token.Type == tokenize.VAR {
		parser.nextToken()
		nameToken := parser.currentToken()
		parser.nextToken()
		exprToken := parser.currentToken()
		var stmt core.Statement
		if nameToken.Type == tokenize.IDENTIFIER && exprToken.Type == tokenize.EQUAL {
			parser.nextToken()
			stmt = &core.VarDeclarationStatement{Name: nameToken, Expr: parser.expression()}
		} else if nameToken.Type == tokenize.IDENTIFIER {
			stmt = &core.VarDeclarationStatement{Name: nameToken}
		}
		stmtTrailing := parser.currentToken()
		if stmtTrailing.Type == tokenize.SEMICOLON {
			parser.nextToken()
			return stmt
		}
		//TODO: How to fix this
		os.Exit(-1)
	}

	return parser.printStatement()
}

func (parser *Parser) printStatement() core.Statement {
	token := parser.currentToken()
	if token.Type == tokenize.PRINT {
		parser.nextToken()
		stmt := &core.PrintStatement{Expr: parser.expression()}
		stmtTrailing := parser.nextToken()
		if stmtTrailing.Type == tokenize.SEMICOLON {
			return stmt
		}
	}
	expr := &core.ExpressionStatement{Expr: parser.expression()}
	stmtTrailing := parser.nextToken()
	if stmtTrailing.Type == tokenize.SEMICOLON {
		return expr
	}
	return nil

}

func (parser *Parser) expression() core.Expression {
	return parser.assignment()
}

func (parser *Parser) assignment() core.Expression {
	token := parser.currentToken()
	if token.Type == tokenize.IDENTIFIER && parser.peekToken().Type == tokenize.EQUAL {
		parser.nextToken()
		parser.nextToken()
		return &core.AssignExpression{Name: token, Expr: parser.assignment()}
	}
	return parser.equality()
}
func (parser *Parser) equality() core.Expression {
	expr := parser.comparison()
	for !parser.isEnd() {
		token := parser.currentToken()
		switch token.Type {
		case tokenize.BANG_EQUAL, tokenize.EQUAL_EQUAL:
			parser.nextToken()
			expr = &core.BinaryExpression{Left: expr, Right: parser.comparison(), Operator: token}
		default:
			return expr
		}
	}
	return expr
}

func (parser *Parser) comparison() core.Expression {
	expr := parser.term()
	for !parser.isEnd() {
		token := parser.currentToken()
		switch token.Type {
		case tokenize.GREATER, tokenize.GREATER_EQUAL, tokenize.LESS_EQUAL, tokenize.LESS:
			parser.nextToken()
			expr = &core.BinaryExpression{Left: expr, Right: parser.term(), Operator: token}
		default:
			return expr
		}
	}
	return expr
}

func (parser *Parser) term() core.Expression {
	expr := parser.factor()
	for !parser.isEnd() {
		token := parser.currentToken()
		switch token.Type {
		case tokenize.PLUS, tokenize.MINUS:
			parser.nextToken()
			expr = &core.BinaryExpression{Left: expr, Right: parser.factor(), Operator: token}
		default:
			return expr
		}
	}
	return expr
}

func (parser *Parser) factor() core.Expression {
	expr := parser.unary()
	for !parser.isEnd() {
		token := parser.currentToken()
		switch token.Type {
		case tokenize.STAR, tokenize.SLASH:
			parser.nextToken()
			expr = &core.BinaryExpression{Left: expr, Right: parser.unary(), Operator: token}
		default:
			return expr
		}
	}
	return expr

}

func (parser *Parser) unary() core.Expression {
	token := parser.currentToken()
	switch token.Type {
	case tokenize.BANG, tokenize.MINUS:
		parser.nextToken()
		return &core.UnaryExpression{Operator: token, Right: parser.unary()}
	}
	return parser.primary()
}

func (parser *Parser) primary() core.Expression {
	token := parser.nextToken()
	switch token.Type {
	case tokenize.TRUE:
		return &core.LiteralExpression{Value: true}
	case tokenize.FALSE:
		return &core.LiteralExpression{Value: false}
	case tokenize.NUMBER, tokenize.STRING:
		return &core.LiteralExpression{Value: token.Literal}
	case tokenize.NIL:
		return &core.LiteralExpression{Value: "nil"}
	case tokenize.IDENTIFIER:
		return &core.VarExpression{Name: token}
	case tokenize.LEFT_PAREN:
		expr := &core.GroupExpression{Expr: parser.expression()}
		if parser.nextToken().Type == tokenize.RIGHT_PAREN {
			return expr
		}

	}
	parser.err = errors.New(fmt.Sprintf("[line %d] Error at '%s': Expect expression.", token.Line, token.Lexeme))
	return nil
}
