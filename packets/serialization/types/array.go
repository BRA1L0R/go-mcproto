package types

import (
	"bytes"
	"errors"
	"reflect"
)

func SerializeArray(field reflect.Value, databuf *bytes.Buffer, SerializeFields func(f reflect.Value, buf *bytes.Buffer) error) error {
	if field.Kind() != reflect.Slice {
		return errors.New("mcproto: struct field is not a slice")
	}

	length := field.Len()
	for i := 0; i < length; i++ {
		arrField := field.Index(i)

		if err := SerializeFields(arrField, databuf); err != nil {
			return err
		}
	}

	return nil
}
