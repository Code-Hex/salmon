package command

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/Code-Hex/salmon/core/command/plugin"
)

func Execute(src string) (string, error) {
	parser := NewParser(src)

	exec, err := parser.Parse()
	if err != nil {
		return "", err
	}

	if exec.command == "help" {
		return RunUsage()
	}

	return run[exec.command](exec.args...)
}

func RunUsage() (string, error) {
	var keys []string
	var buf bytes.Buffer

	for k := range usage {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	buf.WriteString("Available commands:\n")
	for _, k := range keys {
		buf.WriteString(fmt.Sprintf("  %-9s %s\n", k, usage[k]))
	}

	return buf.String(), nil
}

var run = map[string]func(...string) (string, error){
	"echo":   plugin.RunEcho,
	"ping":   plugin.RunPing,
	"update": plugin.RunUpdate,
}

var usage = map[string]string{
	"echo":   plugin.DetailEcho,
	"ping":   plugin.DetailPing,
	"update": plugin.DetailUpdate,
}
