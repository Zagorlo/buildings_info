package consts

import "errors"

const (
	NilPointerReference = "nil pointer reference"
)

var (
	ErrorNilPointerReference = errors.New(NilPointerReference)
)
