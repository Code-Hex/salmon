package command

type Token int

var eof = rune(0)

const (
	EOF Token = iota
	EOL
	SPACE // white space
	IDENT // identify

	// command keywords
	PING
)

var CommandNames = map[string]Token{
	"ping": PING,
}
