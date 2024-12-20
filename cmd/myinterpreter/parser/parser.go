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
declaration    → funDecl
               | varDecl
               | statement ;
funDecl        → "fun" function ;
function       → IDENTIFIER "(" parameters? ")" block ;
varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
statement      → exprStmt
               | forStmt
               | ifStmt
               | printStmt
               | whileStmt
               | block ;

forStmt        → "for" "(" ( varDecl | exprStmt | ";" )
                 expression? ";"
                 expression? ")" statement ;

whileStmt      → "while" "(" expression ")" statement ;

ifStmt         → "if" "(" expression ")" statement
               ( "else" statement )? ;

block          → "{" declaration* "}" ;
expression     → assignment ;
assignment     → IDENTIFIER "=" assignment
               | logic_or ;
logic_or       → logic_and ( "or" logic_and )* ;
logic_and      → equality ( "and" equality )* ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary | call ;
call           → primary ( "(" arguments? ")" )* ;
               |    primary;
arguments      → expression ( "," expression )* ;
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

func (parser *Parser) isEnd() bool {
	return parser.current >= len(parser.tokens) || parser.tokens[parser.current].Type == tokenize.EOF
}

func (parser *Parser) currentToken() *tokenize.Token {
	if parser.isEnd() {
		return parser.tokens[len(parser.tokens)-1]
	}
	return parser.tokens[parser.current]
}

func (parser *Parser) previousToken() *tokenize.Token {
	if parser.current > 0 {
		return parser.tokens[parser.current-1]
	}
	return parser.tokens[0]
}

func (parser *Parser) previous() {
	if parser.current > 0 {
		parser.current = parser.current - 1
	}
}

func (parser *Parser) next() {
	if !parser.isEnd() {
		parser.current++
	}
}

func (parser *Parser) match(tokenTypes ...tokenize.TokenType) bool {
	current := parser.currentToken()
	for _, tokenType := range tokenTypes {
		if tokenType == current.Type {
			parser.next()
			return true
		}
	}
	return false
}

func (parser *Parser) mustMatch(tokenType tokenize.TokenType, err string) bool {
	current := parser.currentToken()
	if tokenType == current.Type {
		parser.next()
		return true
	}
	parser.raiseError(err)
	return false
}
func (parser *Parser) check(tokenTypes ...tokenize.TokenType) bool {
	current := parser.currentToken()
	for _, tokenType := range tokenTypes {
		if tokenType == current.Type {
			return true
		}
	}
	return false
}
func (parser *Parser) raiseError(err string) {
	parser.err = errors.New(err)
	fmt.Fprintln(os.Stderr, err)
	os.Exit(65)
}
func (parser *Parser) Error() error {
	return parser.err
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
	switch {
	case parser.match(tokenize.VAR):
		stmt := &core.VarDeclarationStatement{Expr: &core.LiteralExpression{Value: "nil"}}
		if parser.match(tokenize.IDENTIFIER) {
			stmt.Name = parser.previousToken()
			if parser.match(tokenize.EQUAL) {
				stmt.Expr = parser.expression()
			}
		} else {
			//TODO: Error at here ????
			os.Exit(-1)
		}
		if parser.match(tokenize.SEMICOLON) {
			return stmt
		}
		//TODO:  Error at here ????
		return nil
	case parser.match(tokenize.FUN):
		return parser.function()
	default:
		return parser.statement()
	}
	return nil
}

func (parser *Parser) function() core.Statement {
	parser.mustMatch(tokenize.IDENTIFIER, "Expected IDENTIFIER.")
	funcName := parser.previousToken()
	parser.mustMatch(tokenize.LEFT_PAREN, "Expected '('.")
	params := make([]*tokenize.Token, 0)
	for !parser.check(tokenize.RIGHT_PAREN) {
		if parser.mustMatch(tokenize.IDENTIFIER, "Expected IDENTIFIER.") {
			params = append(params, parser.previousToken())
		}
	}
	parser.mustMatch(tokenize.RIGHT_PAREN, "Expected '('.")
	return &core.FuncStatement{Name: funcName, Params: params, Body: parser.block().(*core.BlockStatement).Statements}
}

func (parser *Parser) statement() core.Statement {
	var stmt core.Statement
	switch {
	case parser.match(tokenize.PRINT):
		stmt = &core.PrintStatement{Expr: parser.expression()}
	case parser.check(tokenize.LEFT_BRACE):
		//fmt.Println("GOOOG")
		blockStmt := parser.block()
		return blockStmt
	case parser.match(tokenize.IF):
		return parser.ifStatement()
	case parser.match(tokenize.WHILE):
		return parser.while()
	case parser.match(tokenize.FOR):
		return parser.forstmt()
	default:
		stmt = &core.ExpressionStatement{Expr: parser.expression()}
	}
	if parser.match(tokenize.SEMICOLON) {
		return stmt
	}
	return nil

}

func (parser *Parser) forstmt() core.Statement {
	parser.mustMatch(tokenize.LEFT_PAREN, "Expected \"(\".")
	var var_statement core.Statement
	var expr1, expr2 core.Expression
	if parser.match(tokenize.SEMICOLON) {
		// Nothing
	} else {
		var_statement = parser.declaration()
		if var_statement.Type() != core.VAR_DECLARATION && var_statement.Type() != core.EXPRESSION {
			parser.raiseError("Unexpected statement")
		}
	}
	if parser.match(tokenize.RIGHT_PAREN) {
		return &core.ForStatement{VarStatment: var_statement, Body: parser.statement()}
	}

	expr1 = parser.expression()
	parser.mustMatch(tokenize.SEMICOLON, "Expected \";\".")
	if parser.match(tokenize.RIGHT_PAREN) {
		return &core.ForStatement{VarStatment: var_statement, Expr1: expr1, Expr2: expr2, Body: parser.statement()}
	}
	expr2 = parser.expression()
	if parser.match(tokenize.RIGHT_PAREN) {
		return &core.ForStatement{VarStatment: var_statement, Expr1: expr1, Expr2: expr2, Body: parser.statement()}
	}
	parser.raiseError("For statement expected")
	return nil
}

func (parser *Parser) while() core.Statement {
	parser.mustMatch(tokenize.LEFT_PAREN, "Expected \"(\".")
	expr := parser.expression()
	parser.mustMatch(tokenize.RIGHT_PAREN, "Expected \")\".")
	stmt := parser.statement()
	return &core.WhileStatement{Expr: expr, Body: stmt}
}

func (parser *Parser) ifStatement() core.Statement {
	parser.mustMatch(tokenize.LEFT_PAREN, "Expected \"(\".")
	expr := parser.expression()
	parser.mustMatch(tokenize.RIGHT_PAREN, "Expected \")\".")
	thenBranch := parser.statement()
	var elseBranch core.Statement
	if parser.match(tokenize.ELSE) {
		elseBranch = parser.statement()
	}
	return &core.IFElseStatement{Expr: expr, ThenBranch: thenBranch, ElseBranch: elseBranch}

}

func (parser *Parser) block() core.Statement {
	blocks := make([]core.Statement, 0)
	parser.mustMatch(tokenize.LEFT_BRACE, "Expected '{'.")
	for !parser.check(tokenize.RIGHT_BRACE) && !parser.isEnd() {
		declaration := parser.declaration()
		blocks = append(blocks, declaration)

	}
	if parser.match(tokenize.RIGHT_BRACE) {
		return &core.BlockStatement{Statements: blocks}
	}
	parser.err = errors.New("Expect '}' after block.")
	return nil

}

func (parser *Parser) expression() core.Expression {
	return parser.assignment()
}

func (parser *Parser) assignment() core.Expression {
	expr := parser.logicOr()
	if parser.match(tokenize.EQUAL) {
		if expr.Type() == core.VARIABLE {
			return &core.AssignExpression{Name: expr.(*core.VarExpression).Name, Expr: parser.assignment()}
		}
	}
	return expr
}

func (parser *Parser) logicOr() core.Expression {
	expr := parser.logicAnd()
	for parser.match(tokenize.OR) {
		expr = &core.LogicalExpression{Left: expr, Operator: parser.previousToken(), Right: parser.logicAnd()}
	}
	return expr

}

func (parser *Parser) logicAnd() core.Expression {
	expr := parser.equality()
	for parser.match(tokenize.AND) {
		expr = &core.LogicalExpression{Left: expr, Operator: parser.previousToken(), Right: parser.equality()}
	}

	return expr

}

func (parser *Parser) equality() core.Expression {
	expr := parser.comparison()
	for !parser.isEnd() {
		if parser.match(tokenize.BANG_EQUAL, tokenize.EQUAL_EQUAL) {
			expr = &core.BinaryExpression{Left: expr, Operator: parser.previousToken(), Right: parser.comparison()}
		} else {
			return expr
		}
	}
	return expr
}

func (parser *Parser) comparison() core.Expression {
	expr := parser.term()
	for !parser.isEnd() {
		if parser.match(tokenize.GREATER, tokenize.GREATER_EQUAL, tokenize.LESS_EQUAL, tokenize.LESS) {
			expr = &core.BinaryExpression{Left: expr, Operator: parser.previousToken(), Right: parser.term()}
		} else {
			return expr
		}
	}
	return expr
}

func (parser *Parser) term() core.Expression {
	expr := parser.factor()
	for !parser.isEnd() {
		if parser.match(tokenize.PLUS, tokenize.MINUS) {
			expr = &core.BinaryExpression{Left: expr, Operator: parser.previousToken(), Right: parser.factor()}
		} else {
			return expr
		}
	}
	return expr
}

func (parser *Parser) factor() core.Expression {
	expr := parser.unary()
	for !parser.isEnd() {
		if parser.match(tokenize.STAR, tokenize.SLASH) {
			expr = &core.BinaryExpression{Left: expr, Operator: parser.previousToken(), Right: parser.unary()}
		} else {
			return expr
		}
	}
	return expr

}

func (parser *Parser) unary() core.Expression {
	if parser.match(tokenize.BANG, tokenize.MINUS) {
		return &core.UnaryExpression{Operator: parser.previousToken(), Right: parser.unary()}
	}
	return parser.call()
}

func (parser *Parser) call() core.Expression {
	callee := parser.primary()
	if parser.match(tokenize.LEFT_PAREN) {
		params := parser.arguments()
		parser.mustMatch(tokenize.RIGHT_PAREN, "Expected )")
		return &core.CallExpression{Callee: callee, Params: params}
	}
	return callee
}

func (parser *Parser) arguments() []core.Expression {
	params := make([]core.Expression, 0)
	for !parser.check(tokenize.RIGHT_PAREN) {
		params = append(params, parser.expression())
		if !parser.match(tokenize.COMMA) {
			break
		}
	}
	return params
}
func (parser *Parser) primary() core.Expression {
	var expr core.Expression
	token := parser.currentToken()
	switch {
	case parser.match(tokenize.TRUE):
		expr = &core.LiteralExpression{Value: true}
	case parser.match(tokenize.FALSE):
		expr = &core.LiteralExpression{Value: false}
	case parser.match(tokenize.NUMBER, tokenize.STRING):
		expr = &core.LiteralExpression{Value: token.Literal}
	case parser.match(tokenize.NIL):
		expr = &core.LiteralExpression{Value: nil}
	case parser.match(tokenize.IDENTIFIER):
		expr = &core.VarExpression{Name: token}
	case parser.match(tokenize.LEFT_PAREN):
		expr = &core.GroupExpression{Expr: parser.expression()}
		if parser.match(tokenize.RIGHT_PAREN) {
			return expr
		}
		parser.raiseError(fmt.Sprintf("[line %d] Error at '%s': Expect expression.", parser.currentToken().Line, parser.currentToken().Lexeme))
		return nil
	default:
		parser.raiseError(fmt.Sprintf("[line %d] Error at '%s': Expect expression.", token.Line, token.Lexeme))
		return nil
	}
	return expr

}
