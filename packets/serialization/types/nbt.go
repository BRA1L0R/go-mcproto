package types

import (
	"bytes"
	"reflect"

	"github.com/Tnze/go-mc/nbt"
)

func SerializeNbt(field reflect.Value, databuf *bytes.Buffer) error {
	return nbt.Marshal(databuf, field.Interface())
}

func DeserializeNbt(field reflect.Value, length int, databuf *bytes.Buffer) error {
	decoder := nbt.NewDecoder(databuf)

	return decoder.Decode(field.Addr().Interface())
}
