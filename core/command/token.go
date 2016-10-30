package command

type Token int

var eof = rune(0)

const (
	EOF Token = iota
	EOL
	SPACE // white space
	IDENT // identify
	HELP

	// Command keywords
	CONVERT
	ECHO
	PING
	UPDATE
)

var CommandNames = map[string]Token{
	"help":    HELP,
	"convert": CONVERT,
	"echo":    ECHO,
	"ping":    PING,
	"update":  UPDATE,
}
