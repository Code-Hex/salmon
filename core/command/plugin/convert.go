package plugin

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

// DetailConvert is description convert command
var DetailConvert = "convert bin, hex, decimal"

// RunConvert is root function of convert command
func RunConvert(args ...string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("Arguments does not exist.")
	}
	return convert(args[0])
}

type NScanner struct {
	src    []rune
	offset int
	eof    bool
	buf    bytes.Buffer
}

func convert(s string) (string, error) {
	sc := &NScanner{src: []rune(s), eof: false}
	switch sc.scan(s) {
	case 'b':
		return sc.binToConvert()
	case 'h':
		return sc.hexToConvert()
	case 'd':
		return sc.decToConvert()
	}

	return "", errors.New("Invalid argument")
}

func (sc *NScanner) binToConvert() (string, error) {
	i64, err := strconv.ParseInt(sc.output(), 2, 64)
	if err != nil {
		return "", err
	}
	return getResult(i64), nil
}

func (sc *NScanner) hexToConvert() (string, error) {
	i64, err := strconv.ParseInt(sc.output(), 16, 64)
	if err != nil {
		return "", err
	}
	return getResult(i64), nil
}

func (sc *NScanner) decToConvert() (string, error) {
	i64, err := strconv.ParseInt(sc.output(), 10, 64)
	if err != nil {
		return "", err
	}
	return getResult(i64), nil
}

func (sc *NScanner) output() string {
	return sc.buf.String()
}

func (sc *NScanner) scan(s string) rune {
	c := sc.peek()
	sc.next()

	if c == '0' && sc.peek() == 'x' {
		sc.next()
		for isHex(sc.peek()) {
			sc.Write()
			sc.next()
		}
		return 'h'
	}

	if c == '0' && sc.peek() == 'b' {
		sc.next()
		for isBin(sc.peek()) {
			sc.Write()
			sc.next()
		}
		return 'b'
	}

	if isNumber(c) {
		sc.buf.WriteRune(c)
		for isNumber(sc.peek()) {
			sc.Write()
			sc.next()
		}
		sc.back()
		return 'd'
	}

	return 0x00
}

func (s *NScanner) Write() {
	s.buf.WriteRune(s.src[s.offset])
}

func (s *NScanner) back() {
	s.offset--
}

func (s *NScanner) next() {
	s.offset++
}

func (s *NScanner) peek() rune {
	if s.offset >= len(s.src) {
		s.eof = true
		return 0x00
	}
	return s.src[s.offset]
}

func getResult(i64 int64) string {
	return fmt.Sprintf("decimal: %d\nhex    : 0x%x\nbinary : 0b%b", i64, i64, i64)
}

func isNumber(c rune) bool {
	return c >= '0' && c <= '9'
}

func isBin(ch rune) bool {
	return ch == '0' || ch == '1'
}

func isHex(ch rune) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}
