package unmarshaller

import (
	"encoding/json"
	"errors"
	"fmt"
	"pathid_assignment/pkg/utils"
)

// Unmarshaller interface defines methods for different data formats
// ensuring compatibility with various structured data inputs.
type Unmarshaller interface {
	Unmarshal(data []byte, rules map[string]interface{}) ([]map[string]interface{}, error)
	UnmarshalByProperty(data []byte, rules map[string]interface{}, prop string) ([]map[string]interface{}, error)
}

// JSONUnmarshaller implements the Unmarshaller interface for JSON format.
type JSONUnmarshaller struct{}

// NewJSONUnmarshaller creates a new instance of JSONUnmarshaller.
func NewJSONUnmarshaller() Unmarshaller {
	return &JSONUnmarshaller{}
}

// UnmarshalByProperty extracts the specified property from JSON data, if it exists as an array, and processes it accordingly.
func (u *JSONUnmarshaller) UnmarshalByProperty(data []byte, rules map[string]interface{}, prop string) ([]map[string]interface{}, error) {
	var jsonObject map[string]interface{}
	if err := json.Unmarshal(data, &jsonObject); err != nil {
		return nil, err
	}

	// If the specified property exists and is an array, extract and send to Unmarshal
	if value, exists := jsonObject[prop]; exists {
		if extractedArray, ok := value.([]interface{}); ok {
			var extractedData []byte
			extractedData, err := json.Marshal(extractedArray)
			if err != nil {
				return nil, err
			}
			return u.Unmarshal(extractedData, rules)
		} else {
			fmt.Printf("Warning: '%s' field is not an array\n", prop)
		}
	}

	// If no matching property is found, fallback to the default unmarshal logic
	return u.Unmarshal(data, rules)
}

// Unmarshal processes JSON data based on defined rules and extracts relevant objects.
func (u *JSONUnmarshaller) Unmarshal(data []byte, rules map[string]interface{}) ([]map[string]interface{}, error) {
	var jsonObject map[string]interface{}
	if err := json.Unmarshal(data, &jsonObject); err == nil {
		// Try extracting data using rule-based key lookup
		if extracted := extractByRulesIterative(jsonObject, rules); len(extracted) > 0 {
			return extracted, nil
		}

		// If the JSON is a valid object but not in an array, return it directly
		return []map[string]interface{}{jsonObject}, nil
	}

	// If the JSON is already an array of objects
	var jsonArray []map[string]interface{}
	if err := json.Unmarshal(data, &jsonArray); err == nil {
		return jsonArray, nil
	}

	return nil, errors.New("expected structure not found in JSON. Ensure that the JSON format matches the expected rules structure. JSON keys detected: " + fmt.Sprint(utils.GetJSONKeys(data)))
}

// extractByRulesIterative uses an iterative BFS approach to find and extract the desired array or object
func extractByRulesIterative(obj map[string]interface{}, rules map[string]interface{}) []map[string]interface{} {
	queue := []map[string]interface{}{obj} // Initialize queue with root object

	for len(queue) > 0 {
		current := queue[0] // Get first item from queue
		queue = queue[1:]   // Remove first item

		for key := range rules {
			if nestedValue, exists := current[key]; exists {
				switch value := nestedValue.(type) {
				case []interface{}:
					var extracted []map[string]interface{}
					for _, item := range value {
						if obj, ok := item.(map[string]interface{}); ok {
							extracted = append(extracted, obj)
						}
					}
					if len(extracted) > 0 {
						return extracted
					}
				case map[string]interface{}:
					// If the rule key itself is a valid object, return it as a single object
					return []map[string]interface{}{value}
				}
			}
		}
	}
	return nil
}

type XMLUnmarshaller struct{}

// UnmarshalByProperty implements Unmarshaller.
func (u *XMLUnmarshaller) UnmarshalByProperty(data []byte, rules map[string]interface{}, prop string) ([]map[string]interface{}, error) {
	panic("unimplemented")
}

func NewXMLUnmarshaller() Unmarshaller {
	return &XMLUnmarshaller{}
}

func (u *XMLUnmarshaller) Unmarshal(data []byte, rules map[string]interface{}) ([]map[string]interface{}, error) {
	return nil, errors.New("XML unmarshalling not yet implemented")
}

type YAMLUnmarshaller struct{}

// UnmarshalByProperty implements Unmarshaller.
func (u *YAMLUnmarshaller) UnmarshalByProperty(data []byte, rules map[string]interface{}, prop string) ([]map[string]interface{}, error) {
	panic("unimplemented")
}

func NewYAMLUnmarshaller() Unmarshaller {
	return &YAMLUnmarshaller{}
}

func (u *YAMLUnmarshaller) Unmarshal(data []byte, rules map[string]interface{}) ([]map[string]interface{}, error) {
	return nil, errors.New("YAML unmarshalling not yet implemented")
}
