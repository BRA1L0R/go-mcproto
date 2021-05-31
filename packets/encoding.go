package packets

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
	"strconv"

	"github.com/BRA1L0R/go-mcprot/varint"
)

func encodeFields(inter interface{}, databuf *bytes.Buffer) error {
	t := reflect.ValueOf(inter)

	if t.Kind() != reflect.Struct {
		return errors.New("inter should be an interface not a pointer to one")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		typeField := t.Type().Field(i)

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
			ignoreLength, _ := strconv.Atoi(typeField.Tag.Get("len"))
			ignoreBuf := make([]byte, ignoreLength)

			databuf.Write(ignoreBuf)
		case "byte":
			data, _ := field.Interface().(byte)
			databuf.WriteByte(data)
		}
	}

	return nil
}

func decodeFields(interPtr interface{}, databuf *bytes.Buffer) error {
	v := reflect.ValueOf(interPtr)

	if v.Kind() != reflect.Ptr {
		return errors.New("interPtr should be a pointer to an interface")
	}

	t := v.Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		typeField := t.Type().Field(i)

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
			databuf.Read(dataString)

			field.SetString(string(dataString))
		case "inherit":
			err := binary.Read(databuf, binary.BigEndian, field.Addr().Interface())
			if err != nil {
				return err
			}
		case "ignore":
			ignoreLength, _ := strconv.Atoi(typeField.Tag.Get("len"))
			ignoreBuf := make([]byte, ignoreLength)

			databuf.Read(ignoreBuf)
		case "byte":
			value, err := databuf.ReadByte()
			if err != nil {
				return err
			}

			field.SetUint(uint64(value))
		}
	}

	return nil
}
