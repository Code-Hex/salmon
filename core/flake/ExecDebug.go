package flake

import (
	"fmt"
	"os"
	"strings"

	"github.com/Code-Hex/salmon/core/command"
	"github.com/fatih/color"
	"github.com/k0kubun/pp"
)

func ExecDebug(args []string) error {
	src := strings.Join(args, " ")
	parser := command.NewParser(src)

	exec, err := parser.Parse()
	if err != nil {
		return err
	}

	fmt.Println(color.MagentaString("[Command parse]"))
	pp.Println(exec)
	fmt.Println(color.BlueString("[Command execute]"))

	out, err := command.Execute(src)
	if err != nil {
		return err
	}
	os.Stdout.WriteString(out)

	return nil
}
