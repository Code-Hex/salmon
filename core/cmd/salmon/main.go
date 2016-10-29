package main

import (
	"os"

	"github.com/Code-Hex/salmon/core"
)

func main() {
	os.Exit(core.Generate(os.Stdout).Swim())
}
