package plugin

import "strings"

// DetailEcho is description echo command
var DetailEcho = "Returns a string that has been your input."

// RunEcho is root function of echo command
func RunEcho(args ...string) (string, error) {
	return strings.Join(args, " "), nil
}
