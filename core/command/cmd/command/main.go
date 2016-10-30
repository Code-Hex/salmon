package main

import (
	"fmt"

	"github.com/Code-Hex/salmon/core/command"
)

func main() {

	src := `help`
	/*
		parser := command.NewParser(src)

		result, err := parser.Parse()
		if err != nil {
			panic(err)
		}
	*/
	result, err := command.Execute(src)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
