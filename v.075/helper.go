package crowd075

import (
	"reflect"
)

func indirect(o interface{}) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(o))
}
