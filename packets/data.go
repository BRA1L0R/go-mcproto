package packets

import (
	"bytes"
	"errors"
	"reflect"

	"github.com/BRA1L0R/go-mcproto/packets/serialization"
)

type standardPacket struct {
	Data   *bytes.Buffer
	Length int32

	PacketID int32
}

func (p *standardPacket) SerializeData(inter interface{}) error {
	val, err := getReflectValue(inter)
	if err != nil {
		return err
	}

	return serialization.SerializeFields(val, p.Data)
}

func (p *standardPacket) DeserializeData(inter interface{}) error {
	val, err := getReflectValue(inter)
	if err != nil {
		return err
	}

	return serialization.DeserializeFields(val, p.Data)
}

func getReflectValue(inter interface{}) (val reflect.Value, err error) {
	v := reflect.ValueOf(inter)

	if v.Kind() != reflect.Ptr {
		err = errors.New("inter should be a pointer to an interface")
		return
	}

	val = v.Elem()
	return
}

func (p *standardPacket) InitializePacket() {
	p.Data = bytes.NewBuffer([]byte{})
}
