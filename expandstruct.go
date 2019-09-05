package expandstruct

import (
	"fmt"
	"reflect"
	"strings"
)

// fieldByPath gets the reflect.Value from a field in a nested struct
// the fieldPath is a dot seperated string of field names as one one
// would use to access the string in go code. All field names are being
// capitalized, as accessing an unexported field does not make sense.
func fieldByPath(v reflect.Value, fieldPath string) (reflect.Value, error) {
	names := strings.Split(fieldPath, ".")
	res := v
	for _, name := range names {
		if v.Kind() == reflect.Struct {
			res = res.FieldByName(strings.Title(name))
		} else {
			return v, fmt.Errorf("v is not a struct of the given fieldPath: %s does not exists on the struct.", fieldPath)
		}
	}

	return res, nil
}

// ExpandToStruct expands the fieldPath value pairs in the map m into
// the given nested struct s. FieldPaths have to be a dot separated
// string of field names, as used when accessing the struct in go.
func ExpandToStruct(m map[string]interface{}, s interface{}) error {
	structVal := reflect.Indirect(reflect.ValueOf(s))
	for k, v := range m {
		mapVal := reflect.ValueOf(v)
		fieldVal, err := fieldByPath(structVal, k)
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
