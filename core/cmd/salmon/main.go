package main

import (
	"os"

	"../../../core"
)

func main() {
	os.Exit(core.Generate(os.Stdout).Swim())
}
