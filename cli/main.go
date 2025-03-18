package main

import (
	"pathid_assignment/pkg/processor"
	"pathid_assignment/pkg/transformer"
	"pathid_assignment/pkg/unmarshaller"

	"fmt"
	"log"

	"github.com/spf13/cobra"
)

const defaultRulesPath = "configs/default_rules.json"

func main() {
	var inputPath, outputPath, rulesPath string

	// Define CLI command
	var rootCmd = &cobra.Command{
		Use:   "read",
		Short: "Read and process input file",
		Run: func(cmd *cobra.Command, args []string) {
			// Validate required flags
			if inputPath == "" || outputPath == "" {
				log.Fatal("Error: Input and output paths must be specified")
			}

			// Use default rules file if none provided
			if rulesPath == "" {
				fmt.Println("No rules file specified. Using default rules file:", defaultRulesPath)
				rulesPath = defaultRulesPath
			}

			fmt.Println("Starting processing...")

			proc := processor.NewProcessor(
				transformer.NewKeywordTransformer(),
				unmarshaller.NewJSONUnmarshaller(),
			)

			proc.Process([]string{inputPath}, rulesPath, outputPath)
			fmt.Println("Processing completed successfully!")
		},
	}

	// Define command flags
	rootCmd.Flags().StringVarP(&inputPath, "input", "i", "", "Path to input file (required)")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to output file (required)")
	rootCmd.Flags().StringVarP(&rulesPath, "rules", "r", "", "Path to rules file (optional, defaults to configs/default_rules.json)")

	// Execute CLI command
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
