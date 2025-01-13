package multipart

import (
	"errors"
	"io"
	"net/textproto"
	refl "reflect"
)

type partWriter func(p *formPart, fw formWriter) error

type formWriter interface {
	WriteField(string, string) error
	CreatePart(textproto.MIMEHeader) (io.Writer, error)
}

type encoder struct {
	fw formWriter
}

func (e *encoder) encode(f *formPart) error {
	var enc partWriter
	var err error
	switch f.mime {
	case 0:
		enc, err = fieldWriter(f)
	default:

		enc, err = fileWriter(f)
	}
	if err != nil {
		return err
	}
	return enc(f, e.fw)
}

func fieldWriter(f *formPart) (partWriter, error) {
	if isBasicKind(f.typ.Kind()) {
		return basicFieldEncoder, nil
	}
	if f.typ.Kind() == refl.Array || f.typ.Kind() == refl.Slice {
		if !isBasicKind(f.typ.Elem().Kind()) {
			return nil, errors.New("multipart: invalid array element type " + f.typ.Kind().String())
		}
		return arrayFieldEncoder, nil
	}

	return nil, errors.New("multipart: invalid field part type")
}

func fileWriter(f *formPart) (partWriter, error) {
	switch f.typ.Kind() {
	case refl.Struct, refl.Slice, refl.Array:
		return structEncoder, nil
	case refl.String:
		return createFileEncoder(f.val.String()), nil
	}

	return nil, errors.New("multipart: invalid type for file field " + f.typ.String())
}
