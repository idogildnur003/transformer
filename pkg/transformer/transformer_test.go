package transformer_test

import (
	"pathid_assignment/pkg/transformer"
	"pathid_assignment/pkg/unmarshaller"
	"testing"

	"encoding/json"
	"os"
)

func TestKeywordTransformer_JSON(t *testing.T) {
	rules := `{
		"id": "id",
		"external_id": "userPrincipalName",
		"mail": "mail",
		"type": "userType",
		"location": "usageLocation",
		"is_enabled": "accountEnabled",
		"first_name": "givenName",
		"last_name": "surname",
		"sign_in_activity": {
			"lastSignInDateTime": "signInActivity.lastSignInDateTime",
			"lastSignInRequestId": "signInActivity.lastSignInRequestId"
		}
	}`

	rawData := []byte(`{
		"id": "123",
		"userPrincipalName": "user@example.com",
		"mail": "user@example.com",
		"userType": "Member",
		"usageLocation": "US",
		"accountEnabled": true,
		"givenName": "John",
		"surname": "Doe",
		"signInActivity": {
			"lastSignInDateTime": "1995-04-28T13:22:24",
			"lastSignInRequestId": "0e685562-4a32-4728-a7ec-4d288ed7d3d4",
			"lastNonInteractiveSignInDateTime": "2013-03-21T07:49:58",
			"lastNonInteractiveSignInRequestId": "f543e91a-2862-4aad-aced-99c77e303c57",
			"lastSuccessfulSignInDateTime": "1989-06-04T05:16:13",
			"lastSuccessfulSignInRequestId": "1edc6e1e-4037-47e3-a5df-3c914a055146"
		}
	}`)

	rulesPath := "test_rules.json"
	outputPath := "test_output.json"

	if err := os.WriteFile(rulesPath, []byte(rules), 0644); err != nil {
		t.Fatalf("Failed to write rules file: %v", err)
	}
	defer os.Remove(rulesPath)

	transformer := transformer.NewKeywordTransformer()
	unmarshaller := unmarshaller.NewJSONUnmarshaller()

	var rulesMap map[string]interface{}
	if err := json.Unmarshal([]byte(rules), &rulesMap); err != nil {
		t.Fatalf("Failed to unmarshal rules: %v", err)
	}

	data, err := unmarshaller.Unmarshal(rawData, rulesMap)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// **Initialize results slice**
	var results []map[string]interface{}

	// Transform and store the result
	for _, obj := range data {
		result, err := transformer.Transform(obj, rulesMap)
		if err != nil {
			t.Fatalf("Transform failed: %v", err)
		}

		if result["id"] != "123" {
			t.Errorf("Expected id '123', got '%v'", result["id"])
		}

		results = append(results, result)
	}

	// Save transformation result to file
	fileData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal transformed data: %v", err)
	}

	if err := os.WriteFile(outputPath, fileData, 0644); err != nil {
		t.Fatalf("Failed to write output file: %v", err)
	}
	defer os.Remove(outputPath)

	t.Logf("Transformation result written to: %s", outputPath)
}

func TestKeywordTransformer_YAML(t *testing.T) {
	yamlInput := []byte(`
	id: "123"
	userPrincipalName: "user@example.com"
	mail: "user@example.com"
	`)

	rules := map[string]interface{}{
		"id":          "id",
		"external_id": "userPrincipalName",
	}

	unmarshaller := unmarshaller.NewYAMLUnmarshaller()
	data, err := unmarshaller.Unmarshal(yamlInput, rules)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	transformer := transformer.NewKeywordTransformer()

	for _, obj := range data {
		result, err := transformer.Transform(obj, rules)
		if err != nil {
			t.Fatalf("Transform failed: %v", err)
		}

		if result["id"] != "123" {
			t.Errorf("Expected id '123', got '%v'", result["id"])
		}
	}

}

func TestKeywordTransformer_XML(t *testing.T) {
	xmlInput := []byte(`
	<User>
		<id>123</id>
		<userPrincipalName>user@example.com</userPrincipalName>
	</User>
	`)

	rules := map[string]interface{}{
		"id":          "id",
		"external_id": "userPrincipalName",
	}

	unmarshaller := unmarshaller.NewXMLUnmarshaller()
	transformer := transformer.NewKeywordTransformer()

	data, err := unmarshaller.Unmarshal(xmlInput, rules)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	for _, obj := range data {
		result, err := transformer.Transform(obj, rules)
		if err != nil {
			t.Fatalf("Transform failed: %v", err)
		}

		if result["external_id"] != "user@example.com" {
			t.Errorf("Expected external_id 'user@example.com', got '%v'", result["external_id"])
		}
	}
}
