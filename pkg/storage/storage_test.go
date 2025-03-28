package storage_test

import (
	"encoding/json"
	"os"
	"pathid_assignment/pkg/storage"
	"testing"
)

const OutputDir = "output"

func TestStorage(t *testing.T) {
	err := os.Mkdir(OutputDir, 0755)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer os.RemoveAll(OutputDir)

	store := storage.NewStorage()

	// Creating test user data
	testUser := []map[string]interface{}{
		{
			"id":          "user-123",
			"external_id": "ext-001",
			"mail":        "user@example.com",
			"type":        "Member",
			"location":    "BS",
			"is_enabled":  true,
			"first_name":  "Brian",
			"last_name":   "Reid",
		},
	}

	// Creating test sign-in activity data
	testSignInActivity := []map[string]interface{}{
		{
			"id": "user-123",
			"sign_in_activity": map[string]interface{}{
				"lastSignInDateTime":                "2025-03-15T08:00:00Z",
				"lastSignInRequestId":               "abcd-1234",
				"lastNonInteractiveSignInDateTime":  "2025-03-14T10:00:00Z",
				"lastNonInteractiveSignInRequestId": "efgh-5678",
				"lastSuccessfulSignInDateTime":      "2025-03-14T07:00:00Z",
				"lastSuccessfulSignInRequestId":     "ijkl-9012",
			},
		},
	}

	// Save user data
	err = store.SaveUsers(testUser, OutputDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Save sign-in activity
	err = store.SaveSignInActivities(testSignInActivity, OutputDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate users file
	fileData, err := os.ReadFile(storage.GenerateFilePath(OutputDir, "users"))
	if err != nil {
		t.Fatalf("Could not read users file: %v", err)
	}

	var users []map[string]interface{}
	err = json.Unmarshal(fileData, &users)
	if err != nil {
		t.Fatalf("Error unmarshalling users JSON: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}

	if users[0]["mail"] != "user@example.com" {
		t.Errorf("Expected mail 'user@example.com', got %s", users[0]["mail"])
	}

	// Validate sign-in activity file
	signInFileData, err := os.ReadFile(storage.GenerateFilePath(OutputDir, "signInActivity"))
	if err != nil {
		t.Fatalf("Could not read sign-in file: %v", err)
	}

	var signInData []map[string]interface{}
	err = json.Unmarshal(signInFileData, &signInData)
	if err != nil {
		t.Fatalf("Error unmarshalling sign-in JSON: %v", err)
	}

	if len(signInData) == 0 {
		t.Errorf("Expected sign-in records, but found none")
	}
}
