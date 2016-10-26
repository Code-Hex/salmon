package command

type Exec struct {
	command string
	args    []string
}

type Parser struct {
	sc *Scanner
}

func NewParser(src string) *Parser {
	return &Parser{sc: NewScanner(src)}
}

// Parse parses a SQL SELECT statement.
func (p *Parser) Parse() (*Exec, error) {
	exec := &Exec{}

	// Next we should loop over all our comma-delimited fields.
	for {
		// Read a field.
		tok, lit, err := p.sc.Scan()
		if err != nil {
			return exec, err
		}

		if tok == EOF {
			break
		}

		switch tok {
		case EOL:
			continue
		case IDENT:
			exec.args = append(exec.args, lit)
		default:
			exec.command = lit
		}
	}

	return exec, nil
}
