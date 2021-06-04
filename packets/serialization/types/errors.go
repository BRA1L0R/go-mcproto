package types

import "errors"

var (
	ErrNotSlice           = errors.New("mcproto: struct field is not a slice")
	ErrIgnoreLenUnknown   = errors.New("mcproto: len tag missing where absolutely necessary")
	ErrIncorrectFieldType = errors.New(
		"mcproto: the target field type does not correspond to the one specified in the type tag",
	)
)
