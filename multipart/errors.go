package multipart

import (
	"errors"
)

var (
	ErrZeroValueInput = errors.New("multipart: input value equals zero or nil")
	ErrInvalidKind    = errors.New("multipart: non struct typed input")
)
