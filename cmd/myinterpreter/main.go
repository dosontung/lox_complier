package main

import (
	"fmt"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/parser"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/tokenize"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	//fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" && command != "evaluate" && command != "parse" {
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
	tkn := tokenize.NewTokennizer()

	if command == "tokenize" {
		tkn.Scan(fileContents)

		for _, token := range tkn.Tokens() {
			fmt.Println(token.String())
		}

		if tkn.LogError().Len() > 0 {
			fmt.Fprint(os.Stderr, tkn.LogError().String())
			os.Exit(65)
		}
	}
	if command == "parse" {
		visitorImp := &parser.VisitorImpl{}
		tkn.Scan(fileContents)
		parser := parser.NewParser(tkn.Tokens())
		expression := parser.Parse()
		fmt.Println(expression.Accept(visitorImp).(string))

	}

}
