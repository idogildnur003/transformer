package processor

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"pathid_assignment/pkg/models"
	"pathid_assignment/pkg/storage"
	"pathid_assignment/pkg/transformer"
	"pathid_assignment/pkg/unmarshaller"
)

// Processor struct manages the transformation and storage process.
type Processor struct {
	Transformer  transformer.GenericTransformer
	Storage      *storage.Storage
	Unmarshaller unmarshaller.Unmarshaller
}

// NewProcessor initializes a new Processor with given Transformer and Unmarshaller.
func NewProcessor(transformer transformer.GenericTransformer, unmarshaller unmarshaller.Unmarshaller, storage *storage.Storage) *Processor {
	return &Processor{
		Transformer:  transformer,
		Storage:      storage,
		Unmarshaller: unmarshaller,
	}
}

// Process reads input files, transforms their contents, and stores the results - Runs the main workflow.
func (p *Processor) Process(inputPaths []string, rulesPath string, outputPath string) {
	var wg sync.WaitGroup
	var users models.UserModel
	workerCount := runtime.NumCPU()

	// Semaphore controls the max number of concurrent goroutines.
	semaphore := make(chan struct{}, workerCount)

	// Loading the rules file into the memory.
	rulesData, err := os.ReadFile(rulesPath)
	if err != nil {
		log.Fatalf("Error reading rules file: %v", err)
	}

	// Parsing the rules file.
	rules := make(map[string]interface{})
	err = json.Unmarshal(rulesData, &rules)
	if err != nil {
		log.Fatalf("Error unmarshalling rules: %v", err)
	}

	// Process both files and directories
	allFiles := []string{}
	for _, path := range inputPaths {
		info, err := os.Stat(path)
		if err != nil {
			log.Printf("Error accessing path %s: %v", path, err)
			continue
		}

		if info.IsDir() {
			files, err := filepath.Glob(filepath.Join(path, "*.json")) // Assuming JSON input files
			if err != nil {
				log.Printf("Error reading directory %s: %v", path, err)
				continue
			}
			allFiles = append(allFiles, files...)
		} else {
			allFiles = append(allFiles, path)
		}
	}

	// Process each input file concurrently, while respecting the semaphore limits.
	for _, inputFilepath := range allFiles {
		wg.Add(1)

		go func(filePath string, rulesMap map[string]interface{}) {
			defer wg.Done()

			semaphore <- struct{}{}        // Adding empty struct to semphore as registering new goroutine.
			defer func() { <-semaphore }() // Release semaphore slot as releasing goroutine.

			log.Printf("Processing file: %s", filePath)
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Error reading file %s: %v", filePath, err)
				return
			}

			// Unmarshaling input data using the configured Unmarshaller.
			objs, err := p.Unmarshaller.UnmarshalByProperty(fileData, rules, "value")
			if err != nil {
				log.Fatalf("Error unmarshalling rules: %v", err)
				return
			}

			// Transform and store each object concurrently, while respecting the semaphore limits.
			for _, obj := range objs {
				wg.Add(1)

				go func(obj map[string]interface{}) {
					defer wg.Done()

					semaphore <- struct{}{}        // Adding empty struct to semphore as registering new goroutine.
					defer func() { <-semaphore }() // Release semaphore slot as releasing goroutine.

					// Transforming the object using the Transformer.
					data, err := p.Transformer.Transform(obj, rulesMap)
					if err != nil {
						log.Println("Error Transforming: " + err.Error())
						return
					}

					users.UserMutex.Lock()

					if signInActivity, exists := data["sign_in_activity"]; exists {
						users.Activities = append(users.Activities, map[string]interface{}{
							"id":               data["id"], // Maintain reference to user ID
							"sign_in_activity": signInActivity,
						})
					}

					// Remove sign_in_activity from user before appending to users list
					delete(data, "sign_in_activity")

					users.Users = append(users.Users, data)
					users.UserMutex.Unlock()
				}(obj)
			}
		}(inputFilepath, rules)

	}

	// Wait for all processing to finish before exiting.
	wg.Wait()

	// Stores all users in output path, in designated json file.
	err = p.Storage.SaveUsers(users.Users, outputPath)
	if err != nil {
		log.Println("Error saving users in file: " + err.Error())
		return
	}

	// Stores all users sign-in activities in output path, in designated json file.
	err = p.Storage.SaveSignInActivities(users.Activities, outputPath)
	if err != nil {
		log.Println("Error saving sign in activities in file: " + err.Error())
		return
	}
}
