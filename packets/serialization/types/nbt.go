package types

import (
	"bytes"
	"reflect"

	"github.com/Tnze/go-mc/nbt"
)

func SerializeNbt(field reflect.Value, databuf *bytes.Buffer) error {
	return nbt.Marshal(databuf, field.Interface())
}

func DeserializeNbt(field reflect.Value, typeField reflect.StructField, length int, databuf *bytes.Buffer) error {
	decoder := nbt.NewDecoder(databuf)

	if typeField.Type.Kind() == reflect.Slice {
		if length < 0 {
			return nil
		}
		field.Set(reflect.MakeSlice(typeField.Type, length, length))

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
