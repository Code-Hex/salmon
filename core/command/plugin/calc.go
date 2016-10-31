package plugin

import (
	"strconv"
	"strings"

	"github.com/alfredxing/calc/compute"
)

// DetailCalc is description calc command
var DetailCalc = "calculator. see calc help."
var help = `Operators:
  +, -, *, /, ^, %
Functions: 
  sin, cos, tan, cot, sec, csc, asin, acos, atan, acot, asec, acsc, sqrt, log, lg, ln, abs
Constants:
  e, pi, π
`

// RunCalc is root function of calc command
func RunCalc(args ...string) (string, error) {
	src := strings.Join(args, " ")

	if strings.Contains(src, "help") {
		return help, nil
	}

	res, err := compute.Evaluate()
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(res, 'G', -1, 64), nil
}
