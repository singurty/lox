package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"errors"

	"github.com/singurty/lox/interpreter"
	"github.com/singurty/lox/parser"
	"github.com/singurty/lox/scanner"
	"github.com/singurty/lox/resolver"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Printf("Usage: %v [file]\n", os.Args[0])
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">> ")
		text, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Println("\nExiting..")
			return
		} else if err != nil {
			panic(err)
		}
		err = run(text)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func runFile(file string) {
	content, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err = run(string(content))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func run(source string) error {
	scanner := scanner.New(source)
	tokens := scanner.ScanTokens()
	if scanner.HadError {
		return errors.New("scanner error")
	}
	parser := parser.New(tokens)
	statements := parser.Parse()
	if parser.HadError {
		return errors.New("parser error")
	}
	resolver := resolver.NewResolver()
	err := resolver.Resolve(statements)
	if err != nil {
		return err
	}
	err = interpreter.Interpret(statements, resolver)
	if err != nil {
		return err
	}
	return nil
}
