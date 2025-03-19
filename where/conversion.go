package where

import "reflect"

// convertToInterfaceArray
// @param values
// @return []interface{}
func convertToInterfaceArray(values interface{}) []interface{} {
	if values == nil {
		return nil
	}
	v := reflect.ValueOf(values)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		for v.Kind() == reflect.Pointer {
			newValue := getValue(values)
			if newValue == nil {
				return nil
			}
			return []interface{}{newValue}
		}
		return []interface{}{values}
	}
	if v.Len() == 0 {
		return nil
	}
	result := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i] = v.Index(i).Interface()
	}
	return result
}

// getValue
// @param value
// @return interface{}
func getValue(value interface{}) interface{} {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr {
		return value
	}
	if v.IsNil() {
		return nil
	}

	return getValue(reflect.Indirect(v).Interface())
}
