package serialization

import (
	"errors"
	"reflect"
	"strconv"
)

// checkDependency
//
// returns false either if the dependency has not been found
// or the dependency was not met
//
// returns true if the dependency has been found AND it's true
func checkDependency(inter reflect.Value, field reflect.StructField) bool {
	depend := field.Tag.Get("depends_on")
	if depend == "" {
		return true
	}

	val := inter.FieldByName(depend)
	return val.Bool()
}

func getLength(inter reflect.Value, field reflect.StructField) (int, error) {
	lengthTag := field.Tag.Get("len")
	if lengthTag == "" {
		return 0, nil
	} else if length, err := strconv.Atoi(lengthTag); err == nil {
		return length, nil
	} else if lenField := inter.FieldByName(lengthTag); lenField.IsValid() {
		return int(lenField.Int()), nil
	} else {
		return 0, errors.New("no way to decode the length found")
	}
}
