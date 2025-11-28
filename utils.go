package cache

import "reflect"

func IsNil[T any](v T) bool {
	if any(v) == nil {
		return true
	}

	value := reflect.ValueOf(v)

	if value.IsValid() {
		switch value.Kind() {
		case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan:
			return value.IsNil()
		default:
			return false
		}
	}

	return false
}
