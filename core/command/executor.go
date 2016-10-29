import "github.com/Code-Hex/salmon/core/command/plugin"

func Execute(src string) (string, error) {
	parser := NewParser(src)

	exec, err := parser.Parse()
	if err != nil {
		return "", err
	}

	return run[exec.command](exec.args...)
}

var run = map[string]func(...string) (string, error){
	"echo": plugin.RunEcho,
	"ping": plugin.RunPing,
}
