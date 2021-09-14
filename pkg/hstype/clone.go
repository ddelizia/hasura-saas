package hstype

import "reflect"

func GenerareEmptyInterface(data interface{}) interface{} {
	vp := reflect.New(reflect.TypeOf(data).Elem())
	inter := vp.Interface()
	return inter
}
