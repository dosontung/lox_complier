package tokenize

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type TokenType string

type LogError struct {
	sb strings.Builder
}

func (l *LogError) Len() int {
	return l.sb.Len()
}
func (l *LogError) String() string {
	return l.sb.String()
}
func (l *LogError) writeError(lineIdx int, message string) {
	// Construct the error message
	l.sb.WriteString("[line ")
	l.sb.WriteString(fmt.Sprintf("%d] Error: ", lineIdx))
	l.sb.WriteString(message)
	l.sb.WriteString("\n")
}

const (
	// Single-character tokens.
	LEFT_PAREN    TokenType = "LEFT_PAREN"
	RIGHT_PAREN   TokenType = "RIGHT_PAREN"
	LEFT_BRACE    TokenType = "LEFT_BRACE"
	RIGHT_BRACE   TokenType = "RIGHT_BRACE"
	COMMA         TokenType = "COMMA"
	DOT           TokenType = "DOT"
	MINUS         TokenType = "MINUS"
	PLUS          TokenType = "PLUS"
	SEMICOLON     TokenType = "SEMICOLON"
	SLASH         TokenType = "SLASH"
	STAR          TokenType = "STAR"
	BANG          TokenType = "BANG"
	BANG_EQUAL    TokenType = "BANG_EQUAL"
	EQUAL         TokenType = "EQUAL"
	EQUAL_EQUAL   TokenType = "EQUAL_EQUAL"
	GREATER       TokenType = "GREATER"
	GREATER_EQUAL TokenType = "GREATER_EQUAL"
	LESS          TokenType = "LESS"
	LESS_EQUAL    TokenType = "LESS_EQUAL"

	// Literals
	IDENTIFIER TokenType = "IDENTIFIER"
	STRING     TokenType = "STRING"
	NUMBER     TokenType = "NUMBER"

	// Keywords
	AND    TokenType = "AND"
	CLASS  TokenType = "CLASS"
	ELSE   TokenType = "ELSE"
	FALSE  TokenType = "FALSE"
	FUN    TokenType = "FUN"
	FOR    TokenType = "FOR"
	IF     TokenType = "IF"
	NIL    TokenType = "NIL"
	OR     TokenType = "OR"
	PRINT  TokenType = "PRINT"
	RETURN TokenType = "RETURN"
	SUPER  TokenType = "SUPER"
	THIS   TokenType = "THIS"
	TRUE   TokenType = "TRUE"
	VAR    TokenType = "VAR"
	WHILE  TokenType = "WHILE"

	// End of file
	EOF TokenType = "EOF"
)

var reservedWords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Token struct {
	Type    TokenType   // The type of token (e.g., IDENTIFIER, STRING, etc.)
	Lexeme  string      // The actual string value of the token
	Literal interface{} // The literal value of the token (e.g., a number or string value)
	Line    int         // The line number where the token appears
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}
func formatLiteral(literal interface{}) string {
	if literal == nil {
		return "null"
	}
	if number, ok := literal.(float64); ok && number == math.Trunc(number) {
		return fmt.Sprintf("%v.0", number)
	}
	return fmt.Sprintf("%v", literal)
}

// String returns a string representation of the token.
func (t *Token) String() string {
	return fmt.Sprintf("%s %s %s", t.Type, t.Lexeme, formatLiteral(t.Literal))
}

type Tokennizer struct {
	tokens   []*Token
	logError *LogError
}

func NewTokennizer() *Tokennizer {
	return &Tokennizer{
		tokens:   make([]*Token, 0),
		logError: &LogError{},
	}
}

func (t *Tokennizer) AddToken(token *Token) {
	t.tokens = append(t.tokens, token)
}

func (t *Tokennizer) Tokens() []*Token {
	return t.tokens
}

func (t *Tokennizer) LogError() *LogError {
	return t.logError
}
func (t *Tokennizer) Scan(fileContents []byte) {
	lineIdx := 1
	isComment := false

	if len(fileContents) > 0 {
		for idx := 0; idx < len(fileContents); idx++ {
			charByte := fileContents[idx]
			if isComment {
				if charByte == '\n' {
					lineIdx++
					isComment = false
				}
				continue
			}
			switch {
			case charByte == '(':
				t.AddToken(NewToken(LEFT_PAREN, "(", nil, lineIdx))
			case charByte == ')':
				t.AddToken(NewToken(RIGHT_PAREN, ")", nil, lineIdx))
			case charByte == '{':
				t.AddToken(NewToken(LEFT_BRACE, "{", nil, lineIdx))
			case charByte == '}':
				t.AddToken(NewToken(RIGHT_BRACE, "}", nil, lineIdx))
			case charByte == '*':
				t.AddToken(NewToken(STAR, "*", nil, lineIdx))
			case charByte == '.':
				t.AddToken(NewToken(DOT, ".", nil, lineIdx))
			case charByte == ',':
				t.AddToken(NewToken(COMMA, ",", nil, lineIdx))
			case charByte == '+':
				t.AddToken(NewToken(PLUS, "+", nil, lineIdx))
			case charByte == '-':
				t.AddToken(NewToken(MINUS, "-", nil, lineIdx))
			case charByte == ';':
				t.AddToken(NewToken(SEMICOLON, ";", nil, lineIdx))
			case charByte == '/':
				if idx+1 < len(fileContents) && fileContents[idx+1] == '/' {
					isComment = true
					idx++
				} else {
					t.AddToken(NewToken(SLASH, "/", nil, lineIdx))
				}

			case charByte == '=':
				if idx+1 < len(fileContents) && fileContents[idx+1] == '=' {
					t.AddToken(NewToken(EQUAL_EQUAL, "==", nil, lineIdx))
					idx++
				} else {
					t.AddToken(NewToken(EQUAL, "=", nil, lineIdx))
				}
			case charByte == '!':
				if idx+1 < len(fileContents) && fileContents[idx+1] == '=' {
					t.AddToken(NewToken(BANG_EQUAL, "!=", nil, lineIdx))
					idx++
				} else {
					t.AddToken(NewToken(BANG, "!", nil, lineIdx))
				}
			case charByte == '<':
				if idx+1 < len(fileContents) && fileContents[idx+1] == '=' {
					t.AddToken(NewToken(LESS_EQUAL, "<=", nil, lineIdx))
					idx++
				} else {
					t.AddToken(NewToken(LESS, "<", nil, lineIdx))
				}
			case charByte == '>':
				if idx+1 < len(fileContents) && fileContents[idx+1] == '=' {
					t.AddToken(NewToken(GREATER_EQUAL, ">=", nil, lineIdx))
					idx++
				} else {
					t.AddToken(NewToken(GREATER, ">", nil, lineIdx))
				}
			case charByte == ' ':
				continue
			case charByte == '\t':
				continue
			case charByte == '"':
				err, str, newIdx := getString(idx+1, fileContents)
				idx = newIdx
				if err != nil {
					t.logError.writeError(lineIdx, "Unterminated string.")
					idx--
				} else {
					t.AddToken(NewToken(STRING, fmt.Sprintf("\"%s\"", str), str, lineIdx))

				}
			case charByte >= '0' && charByte <= '9':
				err, number, precision, newIdx := getNumber(idx, fileContents)
				idx = newIdx - 1
				if err == nil {
					if !isInteger(number) {
						t.AddToken(NewToken(NUMBER, fmt.Sprintf("%.*f", precision, number), number, lineIdx))
					} else {
						t.AddToken(NewToken(NUMBER, fmt.Sprintf("%.*f", precision, number), number, lineIdx))
					}
				}
			case (charByte >= 'a' && charByte <= 'z') || (charByte >= 'A' && charByte <= 'Z') || charByte == '_':
				identifier, newIdx := getIdentifier(idx, fileContents)
				if value, ok := reservedWords[identifier]; ok {
					t.AddToken(NewToken(value, identifier, nil, lineIdx))
				} else {
					t.AddToken(NewToken(IDENTIFIER, identifier, nil, lineIdx))
				}
				idx = newIdx - 1
			case charByte == '\n':
				lineIdx++
			default:
				t.logError.writeError(lineIdx, fmt.Sprintf("Unexpected character: %c", charByte))
			}
		}
		t.AddToken(NewToken(EOF, "", nil, lineIdx))

	}
	if len(t.tokens) == 0 {
		t.AddToken(NewToken(EOF, "", nil, lineIdx))
	}

}

func getString(idx int, fileContents []byte) (error, string, int) {
	var sb strings.Builder
	hasError := true
	i := idx
	for i = idx; i < len(fileContents); i++ {
		if fileContents[i] == '\n' {
			break
		}
		if fileContents[i] == '"' {
			hasError = false
			break
		}
		sb.WriteByte(fileContents[i])
	}
	if hasError {
		return errors.New("Error: Unterminated string."), "", i
	}
	return nil, sb.String(), i
}

func getNumber(idx int, fileContents []byte) (error, float64, int, int) {
	i := idx
	precision := 0
	for i = idx; i < len(fileContents); i++ {
		charByte := fileContents[i]
		if charByte == '.' {
			precision = i
		}
		if (charByte >= '0' && charByte <= '9') || charByte == '.' {
			continue
		}
		break
	}
	floatValue, err := strconv.ParseFloat(string(fileContents[idx:i]), 64)
	if precision != 0 {
		precision = i - precision - 1
	}
	if err != nil {
		return err, floatValue, precision, i
	}
	return nil, floatValue, precision, i

}

func getIdentifier(idx int, fileContents []byte) (string, int) {
	var i int
	for i = idx; i < len(fileContents); i++ {
		charByte := fileContents[i]
		if (charByte >= 'a' && charByte <= 'z') || (charByte >= 'A' && charByte <= 'Z') || charByte == '_' || (charByte >= '0' && charByte <= '9') {
			continue
		}
		break
	}
	return string(fileContents[idx:i]), i
}

func isInteger(f float64) bool {
	return f == math.Trunc(f)
}
