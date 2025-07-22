package main

import (
	"log"
	"os"
	"path/filepath"
	"fmt"
)

func main() {

	// Initialize the service
	service := NewTransactionReconciliationService()

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}// get the paths for the CSV files
	

	sourceFile := filepath.Join(workingDir, "assets", "data", "csvs", "source_transactions.csv")
	systemFile := filepath.Join(workingDir, "assets", "data", "csvs", "system_transactions.csv")

	// Check if files exist
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		log.Fatalf("Source transactions file not found: %s", sourceFile)
	}
	if _, err := os.Stat(systemFile); os.IsNotExist(err) {
		log.Fatalf("System transactions file not found: %s", systemFile)
	}

	fmt.Println("🔄 Starting Transaction Reconciliation Service")
	fmt.Println("================================================")

	// Process the reconciliation
	result, err := service.ProcessReconciliation(sourceFile, systemFile)
	if err != nil {
		log.Fatalf("Reconciliation failed: %v", err)
	}

	// Print summary
	service.PrintSummary(result)


	// Output the detailed JSON result to a file
	fmt.Println("\n📊 DETAILED RECONCILIATION REPORT:")
	fmt.Println("==================================")
	err = service.OutputReconciliationResult(result)
	if err != nil {
		log.Fatalf("Failed to output reconciliation result: %v", err)
	}

	fmt.Println("\n✅ Reconciliation completed successfully!")
}