package protobuf

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

func ToByteSliceE(i interface{}) ([]byte, error) {
	var a []byte
	switch v := i.(type) {
	case []interface{}:
		for _, u := range v {
			b, ok := u.(byte)
			if !ok {
				return a, fmt.Errorf("Unable to Cast %#v to []byte", i)
			}
			a = append(a, b)
		}
		return a, nil
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		return a, fmt.Errorf("Unable to Cast %#v to []byte", i)
	}
}

func ToInt64SliceE(i interface{}) ([]int64, error) {
	if i == nil {
		return []int64{}, fmt.Errorf("Unable to Cast %#v to []int", i)
	}

	switch v := i.(type) {
	case []int64:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]int64, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToInt64E(s.Index(j).Interface())
			if err != nil {
				return []int64{}, fmt.Errorf("Unable to Cast %#v to []int", i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []int64{}, fmt.Errorf("Unable to Cast %#v to []int", i)
	}
}
