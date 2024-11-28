package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type LogError struct {
	sb strings.Builder
}

func (l *LogError) writeError(lineIdx int, message string) {
	// Construct the error message
	l.sb.WriteString("[line ")
	l.sb.WriteString(fmt.Sprintf("%d] Error: ", lineIdx))
	l.sb.WriteString(message)
	l.sb.WriteString("\n")
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	//fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	// Uncomment this block to pass the first stage
	//
	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	lineIdx := 1
	isComment := false
	var builder strings.Builder
	var logError LogError
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
				builder.WriteString("LEFT_PAREN ( null\n")
			case charByte == ')':
				builder.WriteString("RIGHT_PAREN ) null\n")
			case charByte == '{':
				builder.WriteString("LEFT_BRACE { null\n")
			case charByte == '}':
				builder.WriteString("RIGHT_BRACE } null\n")
			case charByte == '*':
				builder.WriteString("STAR * null\n")
			case charByte == '.':
				builder.WriteString("DOT . null\n")
			case charByte == ',':
				builder.WriteString("COMMA , null\n")
			case charByte == '+':
				builder.WriteString("PLUS + null\n")
			case charByte == '-':
				builder.WriteString("MINUS - null\n")
			case charByte == ';':
				builder.WriteString("SEMICOLON ; null\n")
			case charByte == '/':
				if idx+1 < len(fileContents) && fileContents[idx+1] == '/' {
					isComment = true
					idx++
				} else {
					builder.WriteString("SLASH / null\n")
				}

			case charByte == '=':
				if idx+1 < len(fileContents) && fileContents[idx+1] == '=' {
					builder.WriteString("EQUAL_EQUAL == null\n")
					idx++
				} else {
					builder.WriteString("EQUAL = null\n")
				}
			case charByte == '!':
				if idx+1 < len(fileContents) && fileContents[idx+1] == '=' {
					builder.WriteString("BANG_EQUAL != null\n")
					idx++
				} else {
					builder.WriteString("BANG ! null\n")
				}
			case charByte == '<':
				if idx+1 < len(fileContents) && fileContents[idx+1] == '=' {
					builder.WriteString("LESS_EQUAL <= null\n")
					idx++
				} else {
					builder.WriteString("LESS < null\n")
				}
			case charByte == '>':
				if idx+1 < len(fileContents) && fileContents[idx+1] == '=' {
					builder.WriteString("GREATER_EQUAL >= null\n")
					idx++
				} else {
					builder.WriteString("GREATER > null\n")
				}
			case charByte == ' ':
				continue
			case charByte == '\t':
				continue
			case charByte == '"':
				err, str, newIdx := getString(idx+1, fileContents)
				idx = newIdx
				if err != nil {
					logError.writeError(lineIdx, "Unterminated string.")
					idx--
					//fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", line_idx)
				} else {
					builder.WriteString(fmt.Sprintf("STRING \"%s\" %s\n", str, str))
				}
			case charByte >= '0' && charByte <= '9':
				err, number, precision, newIdx := getNumber(idx, fileContents)
				idx = newIdx
				if err == nil {
					if !isInteger(number) {
						builder.WriteString(fmt.Sprintf("NUMBER %.*f %.*f\n", precision, number, precision, number))
					} else {
						builder.WriteString(fmt.Sprintf("NUMBER %.*f %.1f\n", precision, number, number))
					}
				}

			case charByte == '\n':
				lineIdx++
			default:
				logError.writeError(lineIdx, fmt.Sprintf("Unexpected character: %c", charByte))
				//fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", line_idx, charByte)
			}
		}
		builder.WriteString("EOF  null\n")
		fmt.Print(builder.String())

	} else {
		fmt.Println("EOF  null") // Placeholder, remove this line when implementing the scanner
	}
	if logError.sb.Len() > 0 {
		fmt.Fprint(os.Stderr, logError.sb.String())
		os.Exit(65)
	}
}
func isInteger(f float64) bool {
	return f == math.Trunc(f)
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
			continue
		}
		if charByte < '0' || charByte > '9' {
			break
		}
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
