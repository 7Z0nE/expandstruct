package main

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	DefaultSeparator = "."
)

func fieldByPath(v reflect.Value, fieldPath string) (reflect.Value, error) {
	names := strings.Split(fieldPath, ".")
	res := v
	// fmt.Printf("res: %v, type: %v, settable: %v\n", res.Interface(), res.Type(), res.CanSet())
	for _, name := range names {
		if v.Kind() == reflect.Struct {
			res = res.FieldByName(name)
			// fmt.Printf("res: %v, type: %v, settable: %v\n", res.Interface(), res.Type(), res.CanSet())
		} else {
			return v, fmt.Errorf("v is not a struct of the given fieldPath: %s does not exists on the struct.", fieldPath)
		}
	}

	return res, nil
}

func ExpandToStruct(m map[string]interface{}, s interface{}) error {
	structVal := reflect.Indirect(reflect.ValueOf(s))
	// fmt.Printf("structVal: %v, %v, %v\n", structVal.Interface(), structVal.Kind(), structVal.CanSet())
	for k, v := range m {
		// fmt.Printf("k: %v, v: %v\n", k, v)
		mapVal := reflect.ValueOf(v)
		fieldVal, err := fieldByPath(structVal, k)
		// fmt.Printf("fieldVal: %v, settable: %v\n", fieldVal, fieldVal.CanSet())
		if err != nil {
			return err
		}

		switch fieldVal.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if mapVal.Kind() != reflect.Int && mapVal.Kind() != reflect.Int8 && mapVal.Kind() != reflect.Int16 && mapVal.Kind() != reflect.Int32 && mapVal.Kind() != reflect.Int64 {
				return fmt.Errorf("Kind in map at key %s should be some int but was %v", k, mapVal.Kind())
			} else {
				fieldVal.SetInt(mapVal.Int())
			}
		case reflect.String:
			fieldVal.SetString(mapVal.String())
		case reflect.Float32, reflect.Float64:
			if mapVal.Kind() != reflect.Float32 && mapVal.Kind() != reflect.Float64 {
				return fmt.Errorf("Kind in map at key %s should be some float but was %v", k, mapVal.Kind())
			} else {
				fieldVal.SetFloat(mapVal.Float())
			}
		case reflect.Bool:
			fieldVal.SetBool(mapVal.Bool())
		default:
			return fmt.Errorf("Unsupported value of kind %v in field. Only Int, Float, Bool, String allowed.", fieldVal.Kind())
		}
	}
	return nil
}
