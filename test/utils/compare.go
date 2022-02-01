package utils

// Compare compares two interface{} values
// whether they are arrays or maps
// and returns true if they are equal
func Compare(obj1 interface{}, obj2 interface{}) bool {
	var matching bool = true
	switch obj1.(type) {
	case []interface{}:
		if len(obj2.([]interface{})) != len(obj1.([]interface{})) {
			matching = false
		} else {
			for i, v := range obj2.([]interface{}) {
				if !MatchMaps(
					v.(map[string]interface{}),
					obj1.([]interface{})[i].(map[string]interface{}),
				) {
					matching = false
					break
				}
			}
		}

	case map[string]interface{}:
		if !MatchMaps(
			obj1.(map[string]interface{}),
			obj2.(map[string]interface{}),
		) {
			matching = false
		}
	}

	return matching
}
