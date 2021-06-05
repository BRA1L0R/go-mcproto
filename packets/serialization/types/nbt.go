package types

import (
	"bytes"
	"reflect"

	"github.com/Tnze/go-mc/nbt"
)

func SerializeNbt(field reflect.Value, databuf *bytes.Buffer) error {
	if field.Type().Kind() == reflect.Slice {
		length := field.Len()

		for i := 0; i < length; i++ {
			arrField := field.Index(i)

			err := nbt.Marshal(databuf, arrField.Interface())
			if err != nil {
				return err
			}
		}
	} else {
		return nbt.Marshal(databuf, field.Interface())
	}

	return nil
}

func DeserializeNbt(field reflect.Value, length int, databuf *bytes.Buffer) error {
	decoder := nbt.NewDecoder(databuf)
	typeField := field.Type()

	if typeField.Kind() == reflect.Slice {
		if length < 0 {
			return ErrMissingLen
		}

		field.Set(reflect.MakeSlice(typeField, length, length))

		for i := 0; i < length; i++ {
			err := decoder.Decode(field.Index(i).Addr().Interface())
			if err != nil {
				return err
			}
		}
	} else {
		err := decoder.Decode(field.Addr().Interface())
		if err != nil {
			return err
		}
	}

	return nil
}
