package main

import (
	"bufio"
	"fmt"
	"os"
)

var hasError = false
var hasRuntimeError = false

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

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

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// if len(fileContents) > 0 {
	// 	panic("Scanner not implemented")
	// } else {
	// 	fmt.Println("EOF  null") // Placeholder, remove this line when implementing the scanner
	// }
	source := string(fileContents)
	scanner := Scanner{source: source, tokens: []Token{}, start: 0, current: 0, line: 1}
	tokens := scanner.scanTokens()

	print(tokens)
	if hasError {
		os.Exit(65)
	}
}
func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		if text == "exit\n" {
			os.Exit(0)
		}
		run(text)
		hasError = false
	}
}
func runFile(contents []byte) {
	run(string(contents))
	if hasError {
		os.Exit(0)
	}
	if hasRuntimeError {
		os.Exit(70)
	}
}
func print(tokens []Token) {
	for _, token := range tokens {
		fmt.Printf("%s %s %s\n", tokenNames[token.Type], token.Lexeme, "null")
	}
}
func run(source string) {
	scanner := Scanner{source: source, tokens: []Token{}, start: 0, current: 0, line: 1}
	scanner.scanTokens()

	parser := Parser{tokens: scanner.tokens, current: 0}
	expr := parser.parse()
	if hasError {
		return
	}
	interpreter := Interpreter{}
	interpreter.interpret(expr)
}
func error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	hasError = true
}
