package serialization

import "errors"

var (
	ErrIncorrectFieldType = errors.New(
		"the target field type does not correspond to the one specified in the type tag",
	)
)
