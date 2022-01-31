package utils

import "reflect"

// MatchMaps checks if the given maps are equal or not
// except for matching "Id" and "Token" field as they are
// generated randomly
func MatchMaps(m1 map[string]interface{}, m2 map[string]interface{}) bool {
	m1copy := CopyMap(m1)
	m2copy := CopyMap(m2)

	m1copy["Id"] = ""
	m2copy["Id"] = ""

	m1copy["Token"] = ""
	m2copy["Token"] = ""

	return reflect.DeepEqual(m1copy, m2copy)
}

// CopyMap copies the given map and returns a new map
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

// AreKeysEqual checks if the keys of m1 are available in m2
func AreKeysSame(m1 map[string]interface{}, m2 map[string]interface{}) bool {
	for k := range m1 {
		if _, ok := m2[k]; !ok {
			return false
		}
	}

	return true
}

// CheckFieldsExists checks if the given fields exist in map[string]interface{} or not
func CheckFieldsExist(fields map[string]interface{}, fieldNames []string) bool {
	for _, fieldName := range fieldNames {
		if _, ok := fields[fieldName]; !ok {
			return false
		}
	}

	return true
}
