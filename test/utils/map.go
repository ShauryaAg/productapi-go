package utils

import "reflect"

// func to check if the given maps are same or not, except for Id and Token
func MatchMaps(m1 map[string]interface{}, m2 map[string]interface{}) bool {
	m1copy := CopyMap(m1)
	m2copy := CopyMap(m2)

	m1copy["Id"] = ""
	m2copy["Id"] = ""

	m1copy["Token"] = ""
	m2copy["Token"] = ""

	return reflect.DeepEqual(m1copy, m2copy)
}

func CopyMap(m map[string]interface{}) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range m {
		vm, ok := v.(map[string]interface{})
		if ok {
			cp[k] = CopyMap(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}

// func to check if keys in map 1 and map 2 are same
func AreKeysSame(m1 map[string]interface{}, m2 map[string]interface{}) bool {
	for k := range m1 {
		if _, ok := m2[k]; !ok {
			return false
		}
	}

	return true
}

// func to check if the given fields in map[string]interface{} exist or not
func CheckFieldsExist(fields map[string]interface{}, fieldNames []string) bool {
	for _, fieldName := range fieldNames {
		if _, ok := fields[fieldName]; !ok {
			return false
		}
	}

	return true
}
