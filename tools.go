package frule_module

import (
	"reflect"
)

func inSlice(item string, list []string) bool {
	for _, check := range list {
		if check == item {
			return true
		}
	}
	return false
}

func intersectSlices(left, right []string) []string {
	var result []string
	hash := make(map[string]struct{})

	for _, item := range left {
		hash[item] = struct{}{}
	}
	for _, item := range right {
		if _, ok := hash[item]; ok {
			result = append(result, item)
		}
	}
	return result
}

func getFieldValueByTag(input interface{}, tagValue string) reflect.Value {
	for i := 0; i < reflect.TypeOf(input).NumField(); i++ {
		field := reflect.TypeOf(input).Field(i)
		if field.Tag == reflect.StructTag("sql:\""+tagValue+"\"") {
			return reflect.ValueOf(input).FieldByName(field.Name)
		}
	}
	return reflect.Zero(reflect.TypeOf(input))
}
