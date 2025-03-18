package transformer

import (
	"fmt"
	"strings"
)

type GenericTransformer interface {
	Transform(inputData map[string]interface{}, rules map[string]interface{}) (map[string]interface{}, error)
}

type KeywordTransformer struct{}

func NewKeywordTransformer() GenericTransformer {
	return &KeywordTransformer{}
}

func (kt *KeywordTransformer) Transform(inputData map[string]interface{}, rules map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for targetKey, sourcePath := range rules {
		switch path := sourcePath.(type) {
		case string:
			if val, found := extractValue(inputData, path); found {
				result[targetKey] = val
			}
		case map[string]interface{}:
			nestedMap := make(map[string]interface{})
			for nestedTargetKey, nestedSourcePath := range path {
				if nestedPath, ok := nestedSourcePath.(string); ok {
					if val, found := extractValue(inputData, nestedPath); found {
						nestedMap[nestedTargetKey] = val
					}
				}
			}
			if len(nestedMap) > 0 {
				result[targetKey] = nestedMap
			}
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no matching fields found")
	}

	return result, nil
}

func extractValue(data map[string]interface{}, path string) (interface{}, bool) {
	keys := strings.Split(path, ".")
	var value interface{} = data

	for _, key := range keys {
		if m, ok := value.(map[string]interface{}); ok {
			if val, exists := m[key]; exists {
				value = val
			} else {
				return nil, false
			}
		} else {
			return nil, false
		}
	}

	return value, true
}
