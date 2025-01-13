package httptest

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var logger = log.New(os.Stderr, "jactest: ", 0)

//go:generate stringer -type=diffType -linecomment
type diffType int

const (
	missingKey      diffType = iota + 1 // PARAMETER MISSING
	missingValue                        // MISSING VALUE
	unkownParameter                     // UNEXPECTED PARAMETER
	unknownValue                        // UNKNOWN VALUE
	valueMismatch                       // VALUE MISMATCH
)

type diff struct {
	typ diffType
	key string
	val string
	in  string
}

func (d *diff) report() string {
	if d == nil {
		return "NO DIFF"
	}
	report := fmt.Sprintf("[\033[34m%s\033[0m] - \033[31m%s\033[0m \033[33m%s: %s",
		strings.ToUpper(d.in), d.typ, d.key, d.val)
	logger.Println(report)

	return report
}

func diffParams(expected, actual []parameter, allowExtra bool) *diff {
	for _, param := range expected {
		diff := param.diff(actual)
		if diff != nil {
			return diff
		}
	}
	if allowExtra {
		return nil
	}
	for _, p := range actual {
		if !p.matched {
			return &diff{
				in:  p.in,
				key: p.key,
				typ: unkownParameter,
				val: p.valString(),
			}
		}
	}
	return nil
}

func diffSet(e, a map[string]bool) (string, bool) {
	for k := range e {
		if ok := a[k]; !ok {
			return k, true
		}
	}

	return "", false
}

func resolveDiffType(a, o map[string]bool) diffType {
	typ := missingValue
	if len(a) == 1 {
		typ = unknownValue
		if len(o) == 1 {
			typ = valueMismatch
		}
	}

	return typ
}
