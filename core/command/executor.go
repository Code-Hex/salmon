package command

import (
	"bytes"
	"errors"
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
		return runUsage()
	}

	cmd, ok := run[exec.command]
	if ok {
		return cmd(exec.args...)
	}

	return "", errors.New("Command does not exist")
}

func runUsage() (string, error) {
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
	"convert": plugin.RunConvert,
	"echo":    plugin.RunEcho,
	"ping":    plugin.RunPing,
	"update":  plugin.RunUpdate,
}

var usage = map[string]string{
	"convert": plugin.DetailConvert,
	"echo":    plugin.DetailEcho,
	"ping":    plugin.DetailPing,
	"update":  plugin.DetailUpdate,
}
