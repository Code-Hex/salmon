package plugin

import "github.com/k0kubun/pp"

func RunPing(args ...string) (string, error) {
	pp.Print(args)
	return "pong", nil
}
