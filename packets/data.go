package packets

import (
	"errors"
	"reflect"
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
