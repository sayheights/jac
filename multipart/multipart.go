package multipart

import (
	"bytes"
	"mime/multipart"
)

// Form holds a multipart/form-data encoded body
// along with the appropriate Content-Type header matching
// the boundary value set in the payload body.
//
// Form is created as a result of the Marshal method and
// can be used to create jac.Request implementations.
type Form struct {
	Body        []byte
	ContentType string
}

// Marshal returns the multipart/form-data encoding of i.
//
// Marshal traverses the fields of i and consumes fields
// tagged with `multi` key, ignoring the untagged fields.
//
// Fields that have their Content-Type specified encode as file parts.
// String values in file parts are considered to be the path to the file
// to attach to the form.
// Struct and slice valued file parts are encoded based on the specified Content-Type.
//
// Fields with no specified Content-Type encode as text field parts.
func Marshal(i interface{}) (Form, error) {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	encoder := &encoder{fw: writer}
	scanner := &scanner{}

	err := scanner.init(i)
	var part *formPart

	for {
		if err != nil {
			return Form{}, err
		}

		part, err = scanner.nextPart()

		if err == eof || part == endPart {
			break
		}

		if err != nil {
			return Form{}, err
		}

		err = encoder.encode(part)

	}

	_ = writer.Close()

	form := Form{}
	form.Body = buf.Bytes()
	form.ContentType = writer.FormDataContentType()

	return form, nil
}
