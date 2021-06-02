package serialization

import (
	"bytes"
	"encoding/binary"
	"reflect"

	"github.com/BRA1L0R/go-mcprot/varint"
	"github.com/Tnze/go-mc/nbt"
)

func SerializeFields(t reflect.Value, databuf *bytes.Buffer) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		typeField := t.Type().Field(i)

		if !checkDependency(t, typeField) {
			continue
		}

		switch typeField.Tag.Get("type") {
		case "varint":
			encodedData, _ := varint.EncodeVarInt(field.Interface().(int32))
			databuf.Write(encodedData)
		case "string":
			data, _ := field.Interface().(string)
			strLenEncoded, _ := varint.EncodeVarInt(int32(len(data)))

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
		case "nbt":
			err := nbt.Marshal(databuf, field.Interface())
			if err != nil {
				return err
			}
		case "varintarr":
			length, _ := getLength(t, typeField)
			fieldArr := field.Interface().([]int32)

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
