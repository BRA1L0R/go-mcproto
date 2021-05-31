package serialization

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"

	"github.com/BRA1L0R/go-mcprot/varint"
	"github.com/Tnze/go-mc/nbt"
)

func SerializeFields(inter interface{}, databuf *bytes.Buffer) error {
	t := reflect.ValueOf(inter)

	if t.Kind() != reflect.Struct {
		return errors.New("inter should be an interface not a pointer to one")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		typeField := t.Type().Field(i)

		if !checkDependency(t, typeField) {
			continue
		}

		switch typeField.Tag.Get("type") {
		case "varint":
			encodedData, _ := varint.EncodeVarInt(field.Interface().(int))
			databuf.Write(encodedData)
		case "string":
			data, _ := field.Interface().(string)
			strLenEncoded, _ := varint.EncodeVarInt(len(data))

			databuf.Write(strLenEncoded)
			databuf.Write([]byte(data))
		case "inherit":
			inheritedData := field.Interface()

			if err := binary.Write(databuf, binary.BigEndian, inheritedData); err != nil {
				return err
			}
		case "ignore":
			ignoreLength, err := getLength(t, typeField)
			if err != nil {
				return err
			}

			ignoreBuf := make([]byte, ignoreLength)

			databuf.Write(ignoreBuf)
		case "byte":
			data, _ := field.Interface().(byte)
			databuf.WriteByte(data)
		case "nbt":
			err := nbt.Marshal(databuf, field.Interface())
			if err != nil {
				return err
			}
		case "varintarr":
			length, _ := getLength(t, typeField)
			fieldArr := field.Addr().Interface().([]int)

			fieldArrSize := len(fieldArr)

			for i := 0; i < fieldArrSize; i++ {
				encoded, _ := varint.EncodeVarInt(fieldArr[i])
				databuf.Write(encoded)
			}

			if length > fieldArrSize {
				filler := make([]byte, length-fieldArrSize)
				databuf.Write(filler)
			}
		}
	}

	return nil
}
