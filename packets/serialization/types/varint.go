package types

import (
	"bytes"
	"reflect"

	"github.com/BRA1L0R/go-mcproto/varint"
)

func SerializeVarInt(field reflect.Value, databuf *bytes.Buffer) error {
	if field.Kind() == reflect.Slice {
		fieldArr, ok := field.Interface().([]int32)
		if !ok {
			return ErrIncorrectFieldType
		}

		fieldArrSize := len(fieldArr)

		for i := 0; i < fieldArrSize; i++ {
			encoded, _ := varint.EncodeVarInt(fieldArr[i])
			databuf.Write(encoded)
		}
	} else {
		value, ok := field.Interface().(int32)
		if !ok {
			return ErrIncorrectFieldType
		}

		encodedData, _ := varint.EncodeVarInt(value)
		databuf.Write(encodedData)
	}

	return nil
}

func DeserializeVarint(field reflect.Value, length int, databuf *bytes.Buffer) error {
	if length >= 0 {
		arrayField, ok := field.Addr().Interface().(*[]int32)
		if !ok {
			return ErrIncorrectFieldType
		}

		*arrayField = make([]int32, length)

		for i := 0; i < length; i++ {
			decodedVarint, _, err := varint.DecodeReaderVarInt(databuf)
			if err != nil {
				return err
			}

			(*arrayField)[i] = decodedVarint
		}
	} else {
		decodedVal, _, err := varint.DecodeReaderVarInt(databuf)
		if err != nil {
			return err
		}

		field.SetInt(int64(decodedVal))
	}

	return nil
}
