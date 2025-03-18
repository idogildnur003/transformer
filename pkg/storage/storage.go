package storage

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Storage struct {
	userMutex      sync.Mutex
	signInMutex    sync.Mutex
	usersFilePath  string
	signInFilePath string
}

// NewStorage initializes a Storage instance with file paths.
func NewStorage(usersFilePath, signInFilePath string) *Storage {
	return &Storage{
		usersFilePath:  usersFilePath,
		signInFilePath: signInFilePath,
	}
}

func (s *Storage) SaveUsers(users []map[string]interface{}) error {
	s.userMutex.Lock()
	defer s.userMutex.Unlock()

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.usersFilePath, data, 0644)
}

func (s *Storage) SaveSignInActivities(activities []map[string]interface{}) error {
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

		signInMap, ok := signInActivity.(map[string]interface{})
		if !ok {
			log.Printf("Unexpected sign-in activity format for user %s: %+v", userID, signInActivity)
			return nil
		}

		signInMappings := map[string]string{
			"lastSignInDateTime":               "lastSignInRequestId",
			"lastNonInteractiveSignInDateTime": "lastNonInteractiveSignInRequestId",
			"lastSuccessfulSignInDateTime":     "lastSuccessfulSignInRequestId",
		}

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

	return os.WriteFile(s.signInFilePath, data, 0644)
}

func (s *Storage) SaveSignInActivity(activity map[string]interface{}) error {
	s.signInMutex.Lock()
	defer s.signInMutex.Unlock()

	var existingData []map[string]interface{}

	// Read existing sign-in file
	file, err := os.ReadFile(s.signInFilePath)
	if err == nil {
		json.Unmarshal(file, &existingData)
	}

	// Debug log: Check structure before processing
	log.Printf("Processing sign-in activity for user: %+v", activity)

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

	signInMap, ok := signInActivity.(map[string]interface{})
	if !ok {
		log.Printf("Unexpected sign-in activity format for user %s: %+v", userID, signInActivity)
		return nil
	}

	signInMappings := map[string]string{
		"lastSignInDateTime":               "lastSignInRequestId",
		"lastNonInteractiveSignInDateTime": "lastNonInteractiveSignInRequestId",
		"lastSuccessfulSignInDateTime":     "lastSuccessfulSignInRequestId",
	}

	for timeStampKey, requestIdKey := range signInMappings {
		timeStamp, timeStampFound := signInMap[timeStampKey]
		requestID, requestFound := signInMap[requestIdKey]

		if timeStampFound && requestFound {
			existingData = append(existingData, map[string]interface{}{
				"userId":    userID,
				"timeStamp": timeStamp,
				"requestId": requestID,
				"type":      timeStampKey,
			})
		}
	}

	// Save back to file
	data, err := json.MarshalIndent(existingData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.signInFilePath, data, 0644)
}
