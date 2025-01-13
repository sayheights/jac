package multipart

import (
	"errors"
	"strings"
)

var (
	eot                   = errors.New("end of tag")
	errInvalidContentType = errors.New("invalid content type")
)

// optionDelim is seperator used in seperating
// tag option values.
const (
	optionDelim = ','
	keyDelim    = "multi:"
	quote       = "\""
)

const keyOffset = len(keyDelim) + 1

type tag struct {
	fieldName   string
	contentType ContentType
}

type tagParser struct {
	in      string
	out     tag
	isValue bool

	scan func(t *tagParser) error
}

func (t *tagParser) parse() (tag, error) {
	var err error
	t.scan = scanMulti
	if t.isValue {
		t.scan = scanFieldName
	}
	for err = t.scan(t); err == nil; {
		err = t.scan(t)
	}

	if err == eot {
		err = nil
	}

	return t.out, err
}

func parseTag(in string) (tag, error) {
	t := &tagParser{in: in, out: tag{}}

	return t.parse()
}

func parseTagValue(in string) (tag, error) {
	t := &tagParser{in: in, out: tag{}, isValue: true}

	return t.parse()
}

func scanMulti(t *tagParser) error {
	t.in = strings.Trim(t.in[keyOffset:], quote)
	t.scan = scanFieldName
	return nil
}

func scanFieldName(t *tagParser) error {
	if x := strings.IndexRune(t.in, optionDelim); x >= 0 {
		t.out.fieldName = t.in[:x]
		t.in = t.in[x+1:]
		t.scan = scanContentType
		return nil
	}

	t.out.fieldName = t.in
	return eot
}

func scanContentType(t *tagParser) error {
	contentType, ok := toMIME[strings.TrimSpace(t.in)]
	if !ok {
		return errInvalidContentType
	}

	t.out.contentType = contentType
	return eot
}
