package utils

import "encoding/json"

// GetJSONKeys extracts top-level keys from JSON data
func GetJSONKeys(data []byte) []string {
	var jsonObject map[string]interface{}
	var keys []string
	if err := json.Unmarshal(data, &jsonObject); err == nil {
		for key := range jsonObject {
			keys = append(keys, key)
		}
	}
	return keys
}

// FindMappedKey finds the correct key mapping from a string keys map.
func FindMappedKey(targetKey string, rulesMap map[string]interface{}) string {
	for k, v := range rulesMap {
		if mappedKey, ok := v.(string); ok && mappedKey == targetKey {
			return k
		}
	}
	return targetKey // Default fallback
}
