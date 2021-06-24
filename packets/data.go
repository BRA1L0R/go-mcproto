package packets

import (
	"bytes"
	"errors"
	"reflect"

	"github.com/BRA1L0R/go-mcproto/packets/serialization"
)

func getReflectValue(inter interface{}) (val reflect.Value, err error) {
	v := reflect.ValueOf(inter)

	if v.Kind() != reflect.Ptr {
		err = errors.New("mcproto: inter should be a pointer to an interface")
		return
	}

	val = v.Elem()
	return
}

// SerializeData serializes all the fields in the struct with a type struct tag
// into the data buffer
func (p *MinecraftPacket) SerializeData(inter interface{}) error {
	val, err := getReflectValue(inter)
	if err != nil {
		return err
	}

	databuf := bytes.NewBuffer(p.Data)
	if err := serialization.SerializeFields(val, databuf); err != nil {
		return err
	}

	p.Data = databuf.Bytes()
	return nil
}

// DeserializeData deserializes the data buffer into the struct fields which
// have a type struct tag
func (p *MinecraftPacket) DeserializeData(inter interface{}) error {
	val, err := getReflectValue(inter)
	if err != nil {
		return err
	}

	databuf := bytes.NewBuffer(p.Data)
	return serialization.DeserializeFields(val, databuf)
}
