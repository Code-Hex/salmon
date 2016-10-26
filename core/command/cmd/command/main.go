package main

import (
	"../../../command"
	"github.com/k0kubun/pp"
)

func main() {
	src := `ping              hello
	 hello1 

	 world`
	parser := command.NewParser(src)

	result, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	pp.Print(result)
}
