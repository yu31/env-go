package loader

import (
	"encoding"
	"reflect"
)

// Setter is implemented by types can self-deserialize values.
// Any type that implements flag.Value also implements Setter.
type Setter interface {
	Set(value string) error
}

type SetterFunc func(value string) error

func (f SetterFunc) Set(value string) error {
	return f(value)
}

// getSetters return all eligible setter instances
func getSetters(field reflect.Value) []Setter {
	var setters []Setter
	if r := isSetter(field); r != nil {
		setters = append(setters, r)
	}
	if t := isTextUnmarshaler(field); t != nil {
		setters = append(setters, SetterFunc(func(value string) error {
			return t.UnmarshalText([]byte(value))
		}))
	}
	if b := isBinaryUnmarshaler(field); b != nil {
		setters = append(setters, SetterFunc(func(value string) error {
			return b.UnmarshalBinary([]byte(value))
		}))
	}
	return setters
}

func isSetter(field reflect.Value) Setter {
	if !field.CanInterface() {
		return nil
	}
	v := field.Interface()
	r, ok := v.(Setter)
	if !ok && field.CanAddr() {
		v := field.Addr().Interface()
		if r, ok = v.(Setter); !ok {
			return nil
		}
	}
	return r
}

func isTextUnmarshaler(field reflect.Value) encoding.TextUnmarshaler {
	if !field.CanInterface() {
		return nil
	}
	v := field.Interface()
	t, ok := v.(encoding.TextUnmarshaler)
	if !ok && field.CanAddr() {
		v := field.Addr().Interface()
		if t, ok = v.(encoding.TextUnmarshaler); !ok {
			return nil
		}
	}
	return t
}

func isBinaryUnmarshaler(field reflect.Value) encoding.BinaryUnmarshaler {
	if !field.CanInterface() {
		return nil
	}
	v := field.Interface()
	b, ok := v.(encoding.BinaryUnmarshaler)
	if !ok && field.CanAddr() {
		v := field.Addr().Interface()
		if b, ok = v.(encoding.BinaryUnmarshaler); !ok {
			return nil
		}
	}
	return b
}
