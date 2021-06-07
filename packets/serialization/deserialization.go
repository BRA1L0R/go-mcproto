package serialization

import (
	"bytes"
	"reflect"

	"github.com/BRA1L0R/go-mcproto/packets/serialization/tagutils"
	"github.com/BRA1L0R/go-mcproto/packets/serialization/types"
)

func DeserializeFields(t reflect.Value, databuf *bytes.Buffer) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		typeField := t.Type().Field(i)

		lengthTag, err := tagutils.GetLength(t, typeField)
		if err != nil {
			return err
		}

		if !tagutils.CheckDependency(t, typeField) {
			continue
		}

		switch typeField.Tag.Get("mc") {
		case "varint":
			err = types.DeserializeVarInt(field, databuf)
		case "varlong":
			err = types.DeserializeVarLong(field, databuf)
		case "string":
			err = types.DeserializeString(field, databuf)
		case "inherit":
			err = types.DeserializeInherit(field, databuf)
		case "ignore":
			err = types.DeserializeIgnore(lengthTag, databuf)
		case "bytes":
			err = types.DeserializeBytes(field, lengthTag, databuf)
		case "nbt":
			err = types.DeserializeNbt(field, lengthTag, databuf)
		case "array":
			err = types.DeserializeArray(field, lengthTag, databuf, DeserializeFields)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
