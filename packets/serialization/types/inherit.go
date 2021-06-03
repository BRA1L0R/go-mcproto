package types

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

func SerializeInherit(field reflect.Value, databuf *bytes.Buffer) error {
	return binary.Write(databuf, binary.BigEndian, field.Interface())
}

func DeserializeInherit(field reflect.Value, databuf *bytes.Buffer) error {
	return binary.Read(databuf, binary.BigEndian, field.Addr().Interface())
}
