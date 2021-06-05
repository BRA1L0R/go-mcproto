package types

import (
	"bytes"
	"reflect"
)

type serializer func(reflect.Value, *bytes.Buffer) error

func SerializeArray(field reflect.Value, databuf *bytes.Buffer, SerializeFields serializer) error {
	if field.Kind() != reflect.Slice {
		return ErrNotSlice
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

func DeserializeArray(field reflect.Value, length int, databuf *bytes.Buffer, DeserializeFields serializer) error {
	if field.Kind() != reflect.Slice {
		return ErrNotSlice
	}

	if length < 0 {
		return ErrMissingLen
	}

	field.Set(reflect.MakeSlice(field.Type(), length, length))

	for i := 0; i < length; i++ {
		arrField := field.Index(i)

		if err := DeserializeFields(arrField, databuf); err != nil {
			return err
		}
	}

	return nil
}
