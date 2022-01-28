// Copyright 2014 The godump Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package godump

import (
	"fmt"
	"reflect"
	"strconv"
)

type variable struct {
	// Output dump string
	Out string

	// Indent counter
	indent int64
}

func (v *variable) dump(val reflect.Value, name string) {
	v.indent++

	if val.IsValid() && val.CanInterface() {
		typ := val.Type()

		switch typ.Kind() {
		case reflect.Array, reflect.Slice:
			v.printType(name, val.Interface())
			l := val.Len()
			for i := 0; i < l; i++ {
				v.dump(val.Index(i), strconv.Itoa(i))
			}
		case reflect.Map:
			v.printType(name, val.Interface())
			//l := val.Len()
			keys := val.MapKeys()
			for _, k := range keys {
				keyStr, err := v.toString(k.Interface())
				if err == nil {
					v.dump(val.MapIndex(k), keyStr)
				} else {
					v.dump(k, "key:")
					v.dump(val.MapIndex(k), "value:")
				}
			}
		case reflect.Ptr:
			v.printType(name, val.Interface())
			v.dump(val.Elem(), name)
		case reflect.Struct:
			v.printType(name, val.Interface())
			for i := 0; i < typ.NumField(); i++ {
				field := typ.Field(i)
				v.dump(val.FieldByIndex([]int{i}), field.Name)
			}
		default:
			v.printValue(name, val.Interface())
		}
	} else {
		v.printValue(name, "")
	}

	v.indent--
}

func (v *variable) toString(value interface{}) (string, error) {
	switch value.(type) {
	case string:
		if value, ok := value.(string); ok {
			return value, nil
		}
	case bool:
		if value, ok := value.(bool); ok {
			return strconv.FormatBool(value), nil
		}
	case uint8:
		if value, ok := value.(uint8); ok {
			return strconv.Itoa(int(value)), nil
		}
	case int:
		if value, ok := value.(int); ok {
			return strconv.Itoa(value), nil
		}
	case int32:
		if value, ok := value.(int32); ok {
			return strconv.Itoa(int(value)), nil
		}
	case int64:
		if value, ok := value.(int64); ok {
			return strconv.FormatInt(value, 10), nil
		}
	case float32:
		if value, ok := value.(float32); ok {
			return strconv.FormatFloat(float64(value), 'f', 10, 64), nil
		}
	case float64:
		if value, ok := value.(float64); ok {
			return strconv.FormatFloat(value, 'f', 10, 64), nil
		}
	}
	return "", fmt.Errorf("value type unknown")
}

func (v *variable) printType(name string, vv interface{}) {
	v.printIndent()
	v.Out = fmt.Sprintf("%s%s(%T)\n", v.Out, name, vv)
}

func (v *variable) printValue(name string, vv interface{}) {
	v.printIndent()
	v.Out = fmt.Sprintf("%s%s(%T) %#v\n", v.Out, name, vv, vv)
}

func (v *variable) printIndent() {
	var i int64
	for i = 0; i < v.indent; i++ {
		v.Out = fmt.Sprintf("%s  ", v.Out)
	}
}

// Print to standard out the value that is passed as the argument with indentation.
// Pointers are dereferenced.
func Dump(v interface{}) {
	val := reflect.ValueOf(v)
	dump := &variable{indent: -1}
	dump.dump(val, "")
	fmt.Printf("%s", dump.Out)
}

// Return the value that is passed as the argument with indentation.
// Pointers are dereferenced.
func Sdump(v interface{}) string {
	val := reflect.ValueOf(v)
	dump := &variable{indent: -1}
	dump.dump(val, "")
	return dump.Out
}
