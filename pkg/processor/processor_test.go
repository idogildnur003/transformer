package processor_test

import (
	"os"
	"testing"

	"pathid_assignment/pkg/processor"
	"pathid_assignment/pkg/storage"
	"pathid_assignment/pkg/transformer"
	"pathid_assignment/pkg/unmarshaller"
)

func TestProcessor(t *testing.T) {
	defer os.RemoveAll("test_output")

	// Mocking Rules content
	rules := `
		{
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
				"lastSignInRequestId": "signInActivity.lastSignInRequestId",
				"lastNonInteractiveSignInDateTime": "signInActivity.lastNonInteractiveSignInDateTime",
				"lastNonInteractiveSignInRequestId": "signInActivity.lastNonInteractiveSignInRequestId",
				"lastSuccessfulSignInDateTime": "signInActivity.lastSuccessfulSignInDateTime",
   				"lastSuccessfulSignInRequestId": "signInActivity.lastSuccessfulSignInRequestId"
			}
		}`

	// Setup test paths
	rulesPath := "test_rules.json"
	inputPath := "../../data/input"
	outputPath := "test_output"

	if err := os.WriteFile(rulesPath, []byte(rules), 0644); err != nil {
		t.Fatalf("Failed to write rules file: %v", err)
	}
	defer os.Remove(rulesPath)

	// Ensure output directory exists
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Initialize processor
	proc := processor.NewProcessor(transformer.NewKeywordTransformer(), unmarshaller.NewJSONUnmarshaller(), storage.NewStorage())
	proc.Process([]string{inputPath}, rulesPath, outputPath)
}
