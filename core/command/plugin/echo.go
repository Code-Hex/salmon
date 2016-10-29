package plugin

import "strings"

func RunEcho(args ...string) (string, error) {
	return strings.Join(args, ", "), nil
}
