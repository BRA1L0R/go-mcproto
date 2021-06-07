package serialization

import (
	"bytes"
	"reflect"

	"github.com/BRA1L0R/go-mcproto/packets/serialization/tagutils"
	"github.com/BRA1L0R/go-mcproto/packets/serialization/types"
)

func SerializeFields(t reflect.Value, databuf *bytes.Buffer) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		typeField := t.Type().Field(i)

		if !tagutils.CheckDependency(t, typeField) {
			continue
		}

		lengthTag, err := tagutils.GetLength(t, typeField)
		if err != nil {
			return err
		}

		switch typeField.Tag.Get("mc") {
		case "varint":
			err = types.SerializeVarInt(field, databuf)
		case "varlong":
			err = types.SerializeVarLong(field, databuf)
		case "string":
			err = types.SerializeString(field, databuf)
		case "inherit":
			err = types.SerializeInherit(field, databuf)
		case "ignore":
			err = types.SerializeIgnore(lengthTag, databuf)
		case "bytes":
			err = types.SerializeBytes(field, databuf)
		case "nbt":
			err = types.SerializeNbt(field, databuf)
		case "array":
			err = types.SerializeArray(field, databuf, SerializeFields)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
