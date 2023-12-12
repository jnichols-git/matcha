package main

import (
	"fmt"
	"os"
	"slices"
)

var examples = []string{
	"hello",
	"echo",
	"fileserver",
	"middleware",
}

var examplesMap = map[string]func(){
	"hello":      HelloExample,
	"echo":       EchoExample,
	"fileserver": FileServerExample,
	"middleware": MiddlewareExample,
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run ./examples/ example")
		os.Exit(1)
	}
	example := os.Args[1]
	if !slices.Contains(examples, example) {
		fmt.Printf("'%s' is not a valid example.\n", example)
		for _, valid := range examples {
			fmt.Println(" -", valid)
		}
		os.Exit(1)
	}
	examplesMap[example]()
}
