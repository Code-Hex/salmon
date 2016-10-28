package command

import "fmt"

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

func (p *Parser) Parse() (*Exec, error) {
	exec := &Exec{}

LOOP:
	for {
		tok, lit, err := p.sc.Scan()
		if err != nil {
			return exec, err
		}

		switch tok {
		case EOL:
			continue LOOP
		case EOF:
			break LOOP
		case IDENT:
			exec.args = append(exec.args, lit)
		default:
			if exec.command != "" {
				return exec, fmt.Errorf(`syntax error: command name already parsed "%s"`, lit)
			}
			exec.command = lit
		}
	}

	return exec, nil
}
