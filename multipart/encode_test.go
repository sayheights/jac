package multipart

import (
	"bytes"
	_ "embed"
	"io"
	"net/textproto"
	"reflect"
	"strings"
	"testing"
)

type testFormWriter struct {
	fieldname   string
	value       string
	contentType string
	buf         *bytes.Buffer
}

func (t *testFormWriter) Write(p []byte) (n int, err error) {
	n, err = t.buf.Write(p)
	t.value += t.buf.String()

	return
}

func (t *testFormWriter) WriteField(fieldName string, value string) error {
	t.fieldname, t.value = fieldName, value
	return nil
}

func (t *testFormWriter) CreatePart(header textproto.MIMEHeader) (io.Writer, error) {
	mime := header.Get("Content-Type")
	t.contentType = mime
	t.fieldname = strings.Split(strings.Split(header.Values("Content-Disposition")[0], "=")[1], "\"")[1]

	return t, nil
}

func newTestFormWriter() *testFormWriter {
	t := &testFormWriter{buf: &bytes.Buffer{}}
	return t
}

type testStructPart struct {
	ID    string  `json:"id" csv:"ID"`
	Name  string  `json:"name" csv:"NAME"`
	Count int     `json:"count" csv:"COUNT AMOUNT"`
	Price float64 `json:"price" csv:"TOTAL PRICE"`
}

var structPart = testStructPart{
	ID:    "ID",
	Name:  "NAME",
	Count: 10,
	Price: 1.2,
}

var structPartArray = []testStructPart{structPart}

//go:embed testdata/struct_part.json
var structPartJSON []byte

//go:embed testdata/struct_part.csv
var structPartCSV []byte

func Test_encoder_encode(t *testing.T) {
	tests := []struct {
		name     string
		f        *formPart
		fw       *testFormWriter
		wantName string
		wantVal  string
		wantErr  bool
	}{
		{
			fw: newTestFormWriter(),
			f: &formPart{
				name: "uploadFile",
				mime: ApplicationJSON,
				typ:  reflect.TypeOf(structPart),
				val:  reflect.ValueOf(structPart),
			},
			wantName: "uploadFile",
			wantVal:  string(structPartJSON),
		},
		{
			fw: newTestFormWriter(),
			f: &formPart{
				name: "uploadFile",
				mime: TextCSV,
				typ:  reflect.TypeOf(structPartArray),
				val:  reflect.ValueOf(structPartArray),
			},
			wantName: "uploadFile",
			wantVal:  string(structPartCSV),
		},
		{
			fw: newTestFormWriter(),
			f: &formPart{
				name: "uploadFile",
				mime: ApplicationJSON,
				typ:  reflect.TypeOf(""),
				val:  reflect.ValueOf("testdata/struct_part.json"),
			},
			wantName: "uploadFile",
			wantVal:  string(structPartJSON),
		},
		{
			fw: newTestFormWriter(),
			f: &formPart{
				name: "str_field",
				typ:  reflect.TypeOf(""),
				val:  reflect.ValueOf("str is good"),
			},
			wantName: "str_field",
			wantVal:  "str is good",
		},
		{
			fw: newTestFormWriter(),
			f: &formPart{
				name: "int_slice",
				typ:  reflect.TypeOf([]int{1, 2}),
				val:  reflect.ValueOf([]int{1}),
			},
			wantName: "int_slice[]",
			wantVal:  "1",
		},
		{
			fw: newTestFormWriter(),
			f: &formPart{
				name: "float_slice",
				typ:  reflect.TypeOf([]float64{1.3}),
				val:  reflect.ValueOf([]float64{1.3}),
			},
			wantName: "float_slice[]",
			wantVal:  "1.3",
		},
		{
			fw: newTestFormWriter(),
			f: &formPart{
				name: "bool_field",
				typ:  reflect.TypeOf(true),
				val:  reflect.ValueOf(true),
			},
			wantName: "bool_field",
			wantVal:  "true",
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &encoder{
				fw: tt.fw,
			}
			if err := e.encode(tt.f); (err != nil) != tt.wantErr {
				t.Errorf("encoder.encode() error = %v, wantErr %v", err, tt.wantErr)
			}
			check := tt.wantName == tt.fw.fieldname
			if !check {
				t.Errorf("%d: encoder.encode() Fieldname = want %v, got %v", i, tt.wantName, tt.fw.fieldname)
			}
			check = tt.wantVal == tt.fw.value
			if !check {
				t.Errorf("%d: encoder.encode() Value = want %v, got %v", i, tt.wantVal, tt.fw.value)
			}
			check = tt.f.mime == toMIME[tt.fw.contentType]
			if !check {
				t.Errorf("%d: encoder.encode() MIME = want %v, got %v", i, tt.f.mime, toMIME[tt.fw.contentType])
			}
		})
	}
}
