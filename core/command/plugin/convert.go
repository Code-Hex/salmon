package plugin

func isBit(ch rune) bool {
	return ch == '0' || ch == '1'
}

func isHex(ch rune) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}
