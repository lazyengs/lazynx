package utils

import (
	"reflect"
	"strings"

	"github.com/charmbracelet/bubbles/v2/key"
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

func RemoveDuplicateBindings(bindings []key.Binding) []key.Binding {
	seen := make(map[string]struct{})
	result := make([]key.Binding, 0, len(bindings))

	// Process bindings in reverse order
	for i := len(bindings) - 1; i >= 0; i-- {
		b := bindings[i]
		k := strings.Join(b.Keys(), " ")
		if _, ok := seen[k]; ok {
			// duplicate, skip
			continue
		}
		seen[k] = struct{}{}
		// Add to the beginning of result to maintain original order
		result = append([]key.Binding{b}, result...)
	}

	return result
}
