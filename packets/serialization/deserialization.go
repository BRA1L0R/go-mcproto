package serialization

import (
	"bytes"
	"encoding/binary"
	"reflect"

	"github.com/BRA1L0R/go-mcproto/varint"
	"github.com/Tnze/go-mc/nbt"
)

func DeserializeFields(t reflect.Value, databuf *bytes.Buffer) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		typeField := t.Type().Field(i)

		lengthTag, err := getLength(t, typeField)
		if err != nil {
			return err
		}

		if !checkDependency(t, typeField) {
			continue
		}

		switch typeField.Tag.Get("type") {
		case "varint":
			decodedVal, _, err := varint.DecodeReaderVarInt(databuf)
			if err != nil {
				return err
			}

			field.SetInt(int64(decodedVal))
		case "string":
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
		case "inherit":
			err := binary.Read(databuf, binary.BigEndian, field.Addr().Interface())
			if err != nil {
				return err
			}
		case "ignore":
			ignoreBuf := make([]byte, lengthTag)

			_, err := databuf.Read(ignoreBuf)
			if err != nil {
				return err
			}

		// Instead of byte you can use inherit
		// case "byte":
		// 	value, err := databuf.ReadByte()
		// 	if err != nil {
		// 		return err
		// 	}

		// 	field.SetUint(uint64(value))
		case "bytes":
			buf := make([]byte, lengthTag)
			_, err = databuf.Read(buf)
			if err != nil {
				return err
			}

			field.SetBytes(buf)
		case "nbt":
			decoder := nbt.NewDecoder(databuf)

			if typeField.Type.Kind() == reflect.Slice {
				if lengthTag <= 0 {
					continue
				}
				field.Set(reflect.MakeSlice(typeField.Type, lengthTag, lengthTag))

				for i := 0; i < lengthTag; i++ {
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
		case "varintarr":
			length, err := getLength(t, typeField)
			// fmt.Printf("%v\n\n\n", length)
			if err != nil {
				return err
			}

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
		}
	}

	return nil
}
