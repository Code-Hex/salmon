package plugin

// DetailPing is description ping command
var DetailPing = "Return 'pong'."

func RunPing(args ...string) (string, error) {
	return "pong", nil
}
