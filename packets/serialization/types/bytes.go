package types

import (
	"bytes"
	"reflect"
)

func SerializeBytes(field reflect.Value, length int, databuf *bytes.Buffer) error {
	buf := make([]byte, length)
	_, err := databuf.Read(buf)
	if err != nil {
		return err
	}

	field.SetBytes(buf)
	return nil
}
