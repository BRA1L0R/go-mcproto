package serialization

import "errors"

var (
	ErrIncorrectFieldType = errors.New(
		"mcproto: the target field type does not correspond to the one specified in the type tag",
	)
)
