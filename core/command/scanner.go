package command

import (
	"bytes"
	"fmt"
)

type Scanner struct {
	src    []rune
	line   int
	offset int
}

func NewScanner(src string) *Scanner {
	return &Scanner{src: []rune(src)}
}

func (sc *Scanner) Scan() (tok Token, literal string, err error) {

	sc.skipWhiteSpace()
	ch := sc.peek()
	fmt.Println(ch)
	if isLetter(ch) {
		literal = sc.scanIdentifier()
		if token, ok := CommandNames[literal]; ok {
			tok = token
		} else {
			tok = IDENT
		}
		return
	}

	if sc.isEOF() {
		tok = EOF
		return
	}

	err = fmt.Errorf(`syntax error "%s"`, string(ch))
	return
}

func (sc *Scanner) skipWhiteSpace() {
	for isWhitespace(sc.peek()) {
		sc.next()
	}
}

func (sc *Scanner) scanIdentifier() string {
	var buf bytes.Buffer

	for {
		if !isLetter(sc.peek()) && !isDigit(sc.peek()) {
			break
		}
		buf.WriteRune(sc.peek())
		sc.next()
	}

	return buf.String()

}

// 1 文字ずつ読んでいく.
// error があると eof を返す
func (sc *Scanner) read() rune {
	ch := sc.peek()
	sc.next()
	return ch
}

func (sc *Scanner) peek() rune {
	if sc.isEOF() {
		return eof
	}
	return sc.src[sc.offset]
}

func (sc *Scanner) next() {
	if !sc.isEOF() {
		if sc.peek() == '\n' {
			sc.line++
		}
		sc.offset++
	}
}

// 1 文字戻る
func (sc *Scanner) back() {
	sc.offset--
}

func isWhitespace(ch rune) bool {
	return ch == ' '
}

func (sc *Scanner) isEOF() bool {
	return len(sc.src) <= sc.offset
}

func (sc *Scanner) isEOL(ch rune) bool {
	return ch == '\n' || ch == eof
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}
