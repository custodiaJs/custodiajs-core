package utils

import "reflect"

func IsByteSlice(value interface{}) bool {
	v := reflect.ValueOf(value)
	return v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Uint8
}
