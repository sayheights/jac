package multipart

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	tp "net/textproto"
	"os"
	"reflect"

	"github.com/darrae/monk"

	"github.com/gocarina/gocsv"
)

func structEncoder(p *formPart, pw formWriter) error {
	enc := structEncoderByMIME[p.mime]
	buf, err := enc(p.val)
	if err != nil {
		return err
	}
	ext := fileExtensions[p.mime.String()]

	filename, err := saveFile(buf, ext)
	if err != nil {
		return err
	}
	del := func() error {
		_ = createFileEncoder(filename)(p, pw)
		return os.Remove(filename)
	}

	return del()
}

func createFileEncoder(filename string) partWriter {
	return func(p *formPart, fw formWriter) error {
		s, err := readFile(filename)
		if err != nil {
			return err
		}
		mime := mimeHeader(p.name, filename, p.mime)

		w, err := fw.CreatePart(mime)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, bytes.NewReader(s))

		return err
	}
}

type encodeStructFunc func(val reflect.Value) ([]byte, error)

var structEncoderByMIME = map[ContentType]encodeStructFunc{
	TextCSV:         csvEncoder,
	ApplicationJSON: jsonEncoder,
}

func jsonEncoder(val reflect.Value) ([]byte, error) {
	buf, err := json.Marshal(val.Interface())
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func csvEncoder(val reflect.Value) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := gocsv.Marshal(val.Interface(), buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func mimeHeader(name, filename string, ct ContentType) tp.MIMEHeader {
	mimeHeader := tp.MIMEHeader{}
	mimeHeader.Set("Content-Disposition", fmt.Sprintf(
		`form-data; name="%s"; filename="%s"`, name, filename))
	mimeHeader.Add("Content-Type", ct.String())

	return mimeHeader
}

func saveFile(b []byte, suffix string) (string, error) {
	name := "tmp" + "." + suffix
	file, err := os.Create(name)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(file, bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	return name, nil
}

func readFile(s string) ([]byte, error) {
	return monk.Read(s)
}
