package es

import "reflect"

// Name returns the type's name within its package for a defined type.
func parseTypeName(v any) string {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
