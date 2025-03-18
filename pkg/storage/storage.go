package storage

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// Storage holds mutexes for thread-safe access to file operations for users and sign-in activities.
type Storage struct {
	userMutex   sync.Mutex
	signInMutex sync.Mutex
}

// NewStorage initializes a Storage instance with file paths.
func NewStorage() *Storage {
	return &Storage{}
}

// SaveUsers stores users map into given file path.
func (s *Storage) SaveUsers(users []map[string]interface{}, usersFilePath string) error {
	s.userMutex.Lock()
	defer s.userMutex.Unlock()

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(generateFilePath(usersFilePath, "users"), data, 0644)
}

// SaveSignInActivities serializes sign-in activities into JSON format and writes it to a file.
// It first transforms the raw activity data into a structured format before writing.
func (s *Storage) SaveSignInActivities(activities []map[string]interface{}, signInFilePath string) error {
	s.signInMutex.Lock()
	defer s.signInMutex.Unlock()

	var structuredData []map[string]interface{}

	for _, activity := range activities {
		// Ensure the activity has an ID
		userID, userExists := activity["id"].(string)
		if !userExists {
			log.Println("Skipping sign-in activity due to missing user ID")
			return nil
		}

		// Convert sign_in_activity to map[string]string
		signInActivity, ok := activity["sign_in_activity"]
		if !ok {
			log.Printf("Unexpected sign-in activity format for user %s: %+v", userID, signInActivity)
			return nil
		}

		// Assert that signInActivity is a map.
		signInMap, ok := signInActivity.(map[string]interface{})
		if !ok {
			log.Printf("Unexpected sign-in activity format for user %s: %+v", userID, signInActivity)
			return nil
		}

		// Define mappings between timetamp and the corresponding requestId key.
		// That is because 'userID' has already been discovered and 'type' will be discovered after crossing keys by the mentioned two keys and signInMap.
		signInMappings := map[string]string{
			"lastSignInDateTime":               "lastSignInRequestId",
			"lastNonInteractiveSignInDateTime": "lastNonInteractiveSignInRequestId",
			"lastSuccessfulSignInDateTime":     "lastSuccessfulSignInRequestId",
		}

		// Iterate over each mapping to extract valid timestamp and request ID pairs.
		for timeStampKey, requestIdKey := range signInMappings {
			timeStamp, timeStampFound := signInMap[timeStampKey]
			requestID, requestFound := signInMap[requestIdKey]

			if timeStampFound && requestFound {
				structuredData = append(structuredData, map[string]interface{}{
					"userId":    userID,
					"timeStamp": timeStamp,
					"requestId": requestID,
					"type":      timeStampKey,
				})
			}
		}
	}

	data, err := json.MarshalIndent(structuredData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(generateFilePath(signInFilePath, "signInActivity"), data, 0644)
}

// generateFilePath constructs a valid file path by combining a base directory with predefined file names.
func generateFilePath(baseDir, fileType string) string {
	// Define hardcoded file names based on type
	fileNames := map[string]string{
		"users":          "users.json",
		"signInActivity": "signin.json",
	}

	// Retrieve file name or default to "output.json"
	fileName, exists := fileNames[fileType]
	if !exists {
		fileName = "output.json"
	}

	// Concatenate base directory with file name to generate the full path
	return filepath.Join(baseDir, fileName)
}
