package utils

import (
	"reflect"

	"github.com/charmbracelet/bubbles/key"
)

// keyMapToSlice uses reflection to extract all key bindings from a struct
func KeyMapToSlice(keymap any) []key.Binding {
	var bindings []key.Binding
	typ := reflect.TypeOf(keymap)
	if typ.Kind() != reflect.Struct {
		return bindings
	}

	val := reflect.ValueOf(keymap)
	for i := range typ.NumField() {
		field := val.Field(i)
		if field.Type() == reflect.TypeOf(key.Binding{}) {
			bindings = append(bindings, field.Interface().(key.Binding))
		}
	}
	return bindings
}
