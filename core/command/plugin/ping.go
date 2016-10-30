package plugin

// DetailPing is description ping command
var DetailPing = "Return 'pong'."

// RunPing is root function of ping command
func RunPing(args ...string) (string, error) {
	return "pong", nil
}
