package transformer

import (
	"fmt"
	"strings"
)

// GenericTransformer is an interface that defines a method to transform input data based on given rules.
type GenericTransformer interface {
	// Transform takes an input data map and a set of transformation rules,
	// and returns a new map with transformed data or an error if no matching fields are found.
	Transform(inputData map[string]interface{}, rules map[string]interface{}) (map[string]interface{}, error)
}

// KeywordTransformer is a concrete implementation of the GenericTransformer interface.
type KeywordTransformer struct{}

// NewKeywordTransformer returns a new instance of KeywordTransformer as a GenericTransformer.
func NewKeywordTransformer() GenericTransformer {
	return &KeywordTransformer{}
}

// Transform applies the transformation rules to the inputData and returns the transformed result.
// It iterates over the rules and extracts values from inputData based on source paths.
func (kt *KeywordTransformer) Transform(inputData map[string]interface{}, rules map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Loop over each rule entry where the key is the target field name and the value is the source path or nested mapping.
	for targetKey, sourcePath := range rules {
		switch path := sourcePath.(type) {
		case string: // If the rule is a string, treat it as a dot-separated path to a value in the input data.
			if val, found := extractValue(inputData, path); found {
				result[targetKey] = val
			}
		case map[string]interface{}: // If the rule is a nested map, process each nested rule.
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

// extractValue traverses the input data map using a dot-separated path to find the target value.
// It returns the value and a boolean indicating whether the value was found.
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
