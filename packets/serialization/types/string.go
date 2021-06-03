package types

import (
	"bytes"
	"reflect"

	"github.com/BRA1L0R/go-mcproto/varint"
)

func SerializeString(field reflect.Value, databuf *bytes.Buffer) error {
	data, ok := field.Interface().(string)
	if !ok {
		return ErrIncorrectFieldType
	}

	strLenEncoded, _ := varint.EncodeVarInt(int32(len(data)))

	databuf.Write(strLenEncoded)
	databuf.Write([]byte(data))

	return nil
}

func DeserializeString(field reflect.Value, databuf *bytes.Buffer) error {
	if field.Kind() != reflect.String {
		return ErrIncorrectFieldType
	}

	strLength, _, err := varint.DecodeReaderVarInt(databuf)
	if err != nil {
		return err
	}

	dataString := make([]byte, strLength)

	_, err = databuf.Read(dataString)
	if err != nil {
		return err
	}

	field.SetString(string(dataString))
	return nil
}
