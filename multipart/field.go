package multipart

import (
	"reflect"
	"strconv"
)

type FieldKindError struct {
	kind string
}

func (f *FieldKindError) Error() string {
	return "multipart: unhandled field kind " + f.kind
}

// basicFieldEncoder writes a field with a basic kind to a multipart form.
func basicFieldEncoder(p *formPart, fw formWriter) error {
	enc := valueEncoderByKind[p.typ.Kind()]
	err := fw.WriteField(p.name, enc(p.val))

	return err
}

// arrayFieldEncoder writes a field with array type to a multipart form.
func arrayFieldEncoder(p *formPart, pw formWriter) error {
	name := p.name + "[]"
	for i := 0; i < p.val.Len(); i++ {
		elem := p.val.Index(i)
		if !isBasicKind(elem.Kind()) {
			return &FieldKindError{"undhandled error"}
		}
		enc := valueEncoderByKind[elem.Type().Kind()]
		err := pw.WriteField(name, enc(p.val.Index(i)))
		if err != nil {
			return err
		}
	}
	return nil
}

type valueEncoder func(v reflect.Value) string

// intEncoder encodes an integer value suitable for
// writing it into a multipart-form.
func intEncoder(v reflect.Value) string {
	return strconv.Itoa(int(v.Int()))
}

func floatEncoder(v reflect.Value) string {
	return strconv.FormatFloat(v.Float(), 'f', -1, 64)
}

func boolEncoder(v reflect.Value) string {
	return strconv.FormatBool(v.Bool())
}

func stringEncoder(v reflect.Value) string {
	return v.String()
}

var valueEncoderByKind = map[reflect.Kind]valueEncoder{
	reflect.Bool:    boolEncoder,
	reflect.Int:     intEncoder,
	reflect.Int8:    intEncoder,
	reflect.Int16:   intEncoder,
	reflect.Int32:   intEncoder,
	reflect.Int64:   intEncoder,
	reflect.Float32: floatEncoder,
	reflect.Float64: floatEncoder,
	reflect.String:  stringEncoder,
}
