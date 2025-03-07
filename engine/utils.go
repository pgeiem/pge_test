package engine

import (
	"reflect"
)

// Check if only one field of the passed struct is set
func isOnlyOneFieldSet(s interface{}) bool {
	cpt := 0
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() && !reflect.ValueOf(v.Field(i).Interface()).IsZero() {
			if cpt > 0 {
				return false
			}
			cpt++
		}
	}
	return true
}
