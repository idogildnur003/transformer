package unmarshaller_test

import (
	"encoding/json"
	"pathid_assignment/pkg/unmarshaller"
	"testing"
)

func TestJSONUnmarshaller(t *testing.T) {
	rulesJSON := `{
		"users": [
			{
				"id": "id",
				"name": "name"
			}
		]
	}`

	var rules map[string]interface{}
	if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		t.Fatalf("Failed to parse rules JSON: %v", err)
	}

	jsonData := `{
		"users": [
			{"id": "123", "name": "Alice"},
			{"id": "456", "name": "Bob"}
		]
	}`

	unmarshaller := unmarshaller.NewJSONUnmarshaller()

	result, err := unmarshaller.Unmarshal([]byte(jsonData), rules)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 objects, got %d", len(result))
	}

	if result[0]["id"] != "123" {
		t.Errorf("Expected first id '123', got '%v'", result[0]["id"])
	}
}

func TestXMLUnmarshaller(t *testing.T) {
	unmarshaller := unmarshaller.NewXMLUnmarshaller()
	xmlData := []byte(`
		<UserList>
			<User><id>123</id><name>Alice</name></User>
			<User><id>456</id><name>Bob</name></User>
		</UserList>
	`)

	rules := map[string]interface{}{
		"id":   "id",
		"name": "name",
	}

	_, err := unmarshaller.Unmarshal(xmlData, rules)
	if err == nil {
		t.Fatal("Expected XML unmarshal to fail, but it did not")
	}
}

func TestYAMLUnmarshaller(t *testing.T) {
	unmarshaller := unmarshaller.NewYAMLUnmarshaller()
	yamlData := []byte(`
	users:
	  - id: "123"
	    name: "Alice"
	  - id: "456"
	    name: "Bob"
	`)

	rules := map[string]interface{}{
		"id":   "id",
		"name": "name",
	}

	_, err := unmarshaller.Unmarshal(yamlData, rules)
	if err == nil {
		t.Fatal("Expected YAML unmarshal to fail, but it did not")
	}
}
