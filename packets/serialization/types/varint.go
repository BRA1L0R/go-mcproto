package types

import (
	"bytes"
	"reflect"

	"github.com/BRA1L0R/go-mcproto/varint"
)

func SerializeVarInt(field reflect.Value, databuf *bytes.Buffer) error {
	value, ok := field.Interface().(int32)
	if !ok {
		return ErrIncorrectFieldType
	}

	encodedData, _ := varint.EncodeVarInt(value)
	databuf.Write(encodedData)

	return nil
}

func SerializeVarLong(field reflect.Value, databuf *bytes.Buffer) error {
	value, ok := field.Interface().(int64)
	if !ok {
		return ErrIncorrectFieldType
	}

	encodedData, _ := varint.EncodeVarLong(value)
	databuf.Write(encodedData)

	return nil
}

func DeserializeVarInt(field reflect.Value, databuf *bytes.Buffer) error {
	decodedVal, _, err := varint.DecodeReaderVarInt(databuf)
	if err != nil {
		return err
	}

	field.SetInt(int64(decodedVal))

	return nil
}

func DeserializeVarLong(field reflect.Value, databuf *bytes.Buffer) error {
	decodedVal, _, err := varint.DecodeReaderVarLong(databuf)
	if err != nil {
		return err
	}

	field.SetInt(decodedVal)

	return nil
}
