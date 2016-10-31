package plugin

import (
	"errors"
	"strconv"
	"strings"

	"github.com/alfredxing/calc/compute"
)

// DetailCalc is description calc command
var DetailCalc = `calculator. invoke "calc usage".`
var help = `Operators:
  +, -, *, /, ^, %
Functions: 
  sin, cos, tan, cot, sec, csc, asin, acos, atan, acot, asec, acsc, sqrt, log, lg, ln, abs
Constants:
  e, pi, π
`

// RunCalc is root function of calc command
func RunCalc(args ...string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("Arguments does not exist.")
	}

	src := strings.Join(args, " ")

	if strings.Contains(src, "usage") {
		return help, nil
	}

	res, err := compute.Evaluate(src)
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(res, 'G', -1, 64), nil
}
