package cache

import "reflect"

// IsNil check whether the value v is "nil" in the Go sense.
// It correctly handles both typed nil values (e.g. interface{} containing nil)
// and untyped nil, and safely checks reflect-based nil for pointers, maps,
// slices, maps, channels, functions and interfaces without panicking on invalid values.
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
