package main

import (
	"os"
	"pathid_assignment/pkg/processor"
	"pathid_assignment/pkg/storage"
	"pathid_assignment/pkg/transformer"
	"pathid_assignment/pkg/unmarshaller"

	"fmt"
	"log"

	"github.com/spf13/cobra"
)

const defaultRulesPath = "configs/default_mapping_config.json"
const defaultOutputPath = "data/output"

func main() {
	var inputPath, outputPath, rulesPath string

	// Define CLI command
	var rootCmd = &cobra.Command{
		Use:   "read",
		Short: "Read and process input file",
		Run: func(cmd *cobra.Command, args []string) {
			// Validate required input flag
			if inputPath == "" {
				log.Fatal("Error: Input and output paths must be specified")
			}

			// Use default rules file if none provided
			if rulesPath == "" {
				fmt.Println("No rules file specified. Using default rules file:", defaultRulesPath)
				rulesPath = defaultRulesPath
			}

			// Use default output directory path if none provided
			if outputPath == "" {
				fmt.Println("No output directory path specified. Using default path:", defaultOutputPath)
				outputPath = defaultOutputPath

				// Check if output directory is empty and clear it if necessary
				fmt.Println("Clearing output directory...")
				clearOutputDirectory(outputPath)
			}

			fmt.Println("Starting processing...")

			proc := processor.NewProcessor(
				transformer.NewKeywordTransformer(),
				unmarshaller.NewJSONUnmarshaller(),
				storage.NewStorage(),
			)

			proc.Process([]string{inputPath}, rulesPath, outputPath)
			fmt.Println("Processing completed successfully!")
		},
	}

	// Define command flags
	rootCmd.Flags().StringVarP(&inputPath, "input", "i", "", "Path to input file (required)")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to output directory (optional), default to data/output")
	rootCmd.Flags().StringVarP(&rulesPath, "rules", "r", "", "Path to rules file (optional, defaults to configs/default_rules.json)")

	// Execute CLI command
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

func clearOutputDirectory(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Error reading output directory: %v", err)
	}

	for _, entry := range entries {
		filePath := fmt.Sprintf("%s/%s", dir, entry.Name())
		if err := os.RemoveAll(filePath); err != nil {
			log.Fatalf("Error clearing output file %s: %v", filePath, err)
		}
	}
}
