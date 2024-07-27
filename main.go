package main

import (
	"pikolang-interpreter/cli"
	"pikolang-interpreter/interpreter"
)

func main() {
	input := cli.GetFileContentsFromArgs()
	parser := interpreter.NewParser(string(input))
	parser.Run()
}
