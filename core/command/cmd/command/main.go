package main

import (
	"fmt"

	"../../../command"
)

func main() {

	src := `ping hello
		 hello1

		 world`
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
