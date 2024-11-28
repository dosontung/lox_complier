package main

import (
	"fmt"
	"os"
	"strings"
)

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
	line_idx := 1
	errCode := 0
	var builder strings.Builder
	if len(fileContents) > 0 {
		for idx := 0; idx < len(fileContents); idx++ {
			charByte := fileContents[idx]
			switch charByte {
			case '(':
				builder.WriteString("LEFT_PAREN ( null\n")
			case ')':
				builder.WriteString("RIGHT_PAREN ) null\n")
			case '{':
				builder.WriteString("LEFT_BRACE { null\n")
			case '}':
				builder.WriteString("RIGHT_BRACE } null\n")
			case '*':
				builder.WriteString("STAR * null\n")
			case '.':
				builder.WriteString("DOT . null\n")
			case ',':
				builder.WriteString("COMMA , null\n")
			case '+':
				builder.WriteString("PLUS + null\n")
			case '-':
				builder.WriteString("MINUS - null\n")
			case ';':
				builder.WriteString("SEMICOLON ; null\n")
			case '/':
				builder.WriteString("SLASH / null\n")
			case '=':
				if idx+1 < len(fileContents) && fileContents[idx+1] == '=' {
					builder.WriteString("EQUAL_EQUAL == null\n")
					idx++
				} else {
					builder.WriteString("EQUAL = null\n")
				}
			case '!':
				if idx+1 < len(fileContents) && fileContents[idx+1] == '=' {
					builder.WriteString("BANG_EQUAL != null\n")
					idx++
				} else {
					builder.WriteString("BANG ! null\n")
				}
			case '<':
				if idx+1 < len(fileContents) && fileContents[idx+1] == '=' {
					builder.WriteString("LESS_EQUAL <= null\n")
					idx++
				} else {
					builder.WriteString("LESS < null\n")
				}
			case '>':
				if idx+1 < len(fileContents) && fileContents[idx+1] == '=' {
					builder.WriteString("GREATER_EQUAL >= null\n")
					idx++
				} else {
					builder.WriteString("GREATER > null\n")
				}
			case '\n':
				line_idx++
			//case '$':
			//	fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: $\n", line_idx)
			//	errCode = 65
			//case '#':
			//	fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: #\n", line_idx)
			//	errCode = 65
			//case '@':
			//	fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: @\n", line_idx)
			//	errCode = 65
			default:
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", line_idx, charByte)
				errCode = 65
			}
		}
		builder.WriteString("EOF  null\n")
		fmt.Print(builder.String())
		os.Exit(errCode)
	} else {

		fmt.Println("EOF  null") // Placeholder, remove this line when implementing the scanner
	}
}
