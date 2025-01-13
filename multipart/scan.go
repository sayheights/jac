package multipart

import (
	"errors"
	"reflect"
	refl "reflect"
)

const tagKey = "multi"

type formPart struct {
	name string
	mime ContentType
	typ  refl.Type
	val  refl.Value
}

func (f *formPart) Equal(o *formPart) bool {
	if f == nil || o == nil {
		return f == o
	}
	return (f.name == o.name &&
		f.mime == o.mime &&
		f.typ.String() == o.typ.String() &&
		f.val.String() == o.val.String())
}

type scanner struct {
	input  refl.Value
	fields []refl.StructField
	amount int
	index  int
}

func (s *scanner) init(i interface{}) error {
	val, err := validate(i)
	if err != nil {
		return err
	}
	s.input = val
	s.fields = refl.VisibleFields(val.Type())
	s.amount = len(s.fields)

	return nil
}

func validate(i interface{}) (refl.Value, error) {
	val := refl.Indirect(refl.ValueOf(i))
	if !val.IsValid() {
		return val, ErrZeroValueInput
	}

	if val.Kind() != refl.Struct {
		return val, ErrInvalidKind
	}

	return val, nil
}

var (
	endPart = (*formPart)(nil)
	eof     = errors.New("end of fields")
)

func (s *scanner) nextPart() (*formPart, error) {
	field := s.next()
	if shouldSkip(field) {
		s.index += 1
		return s.nextPart()
	}
	if !s.isDone() {
		part, err := s.buildPart()
		return part, err
	}

	return endPart, nil
}

func (s *scanner) buildPart() (*formPart, error) {
	field := s.fields[s.index]
	part := &formPart{}
	part.val = s.input.FieldByIndex(field.Index)
	part.typ = field.Type

	tag, err := parseTagValue(field.Tag.Get(tagKey))
	if err != nil {
		return nil, err
	}
	part.name, part.mime = tag.fieldName, tag.contentType
	s.index += 1

	return part, nil
}

func (s *scanner) isDone() bool {
	return s.amount <= s.index
}

func shouldSkip(st refl.StructField) bool {
	return skipNonMulti(st) && skipAnonymous(st)
}

func skipNonMulti(st refl.StructField) bool {
	_, ok := st.Tag.Lookup(tagKey)
	return !ok
}

func skipAnonymous(st refl.StructField) bool {
	return st.Anonymous
}

func (s *scanner) next() refl.StructField {
	if s.isDone() {
		return reflect.StructField{}
	}
	field := s.fields[s.index]

	return field
}
