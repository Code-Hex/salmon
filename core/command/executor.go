package command

import "./plugin"

func Execute(src string) (string, error) {
	parser := NewParser(src)

	exec, err := parser.Parse()
	if err != nil {
		return "", err
	}

	return run[exec.command](exec.args...)
}

var run = map[string]func(...string) (string, error){
	"ping": plugin.RunPing,
}
