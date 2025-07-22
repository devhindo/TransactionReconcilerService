package main

import (
	"fmt"
	"log"
	"strings"
	"encoding/json"
	"os"
)

// TransactionReconciliationService is the main service that generates the reconciliation report
type TransactionReconciliationService struct {
	csvReader  *CSVReader // it requires a CSVReader
	reconciler *TransactionReconciler // a reconciler to handle the logic
}

// Constructor: NewTransactionReconciliationService creates a new service instance
func NewTransactionReconciliationService() *TransactionReconciliationService {
	return &TransactionReconciliationService{
		csvReader:  NewCSVReader(),
		reconciler: NewTransactionReconciler(),
	}
}

// ProcessReconciliation reads transactions from CSV files, reconciles them, and returns the result
func (s *TransactionReconciliationService) ProcessReconciliation(sourceFilePath, systemFilePath string) (*ReconciliationResult, error) {
	// Read source transactions
	log.Printf("Reading source transactions from: %s", sourceFilePath)
	sourceTransactions, err := s.csvReader.ReadSourceTransactions(sourceFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read source transactions: %w", err)
	}
	log.Printf("Successfully read %d source transactions", len(sourceTransactions))

	// Read system transactions
	log.Printf("Reading system transactions from: %s", systemFilePath)
	systemTransactions, err := s.csvReader.ReadSystemTransactions(systemFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read system transactions: %w", err)
	}
	log.Printf("Successfully read %d system transactions", len(systemTransactions))

	// Perform reconciliation
	log.Println("Starting reconciliation process...")
	result := s.reconciler.Reconcile(sourceTransactions, systemTransactions)
	log.Println("Reconciliation completed")

	return result, nil
}

// OutputReconciliationResult outputs the reconciliation result in JSON format
func (s *TransactionReconciliationService) OutputReconciliationResult(result *ReconciliationResult) error {
	// Create output in the required format
	output := map[string]interface{}{
		"missing_in_internal":      result.MissingInInternal,
		"missing_in_source":        result.MissingInSource,
		"mismatched_transactions":  result.MismatchedTransactions,
		"summary":                  result.Summary,
	}

	// Convert to JSON with pretty printing
	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal result to JSON: %w", err)
	}

	// Print to console
	fmt.Println(string(jsonOutput))

	// Also save to file
	outputFile := "reconciliation_report.json"
	err = os.WriteFile(outputFile, jsonOutput, 0644)
	if err != nil {
		log.Printf("Warning: Could not save report to file %s: %v", outputFile, err)
	} else {
		log.Printf("Reconciliation report saved to: %s", outputFile)
	}

	return nil
}

// PrintSummary prints a human-readable summary of the reconciliation results
func (s *TransactionReconciliationService) PrintSummary(result *ReconciliationResult) {
	separator := strings.Repeat("=", 60)
	fmt.Println("\n" + separator)
	fmt.Println("TRANSACTION RECONCILIATION SUMMARY")
	fmt.Println(separator)
	fmt.Printf("Total Source Transactions:      %d\n", result.Summary.TotalSourceTransactions)
	fmt.Printf("Total System Transactions:      %d\n", result.Summary.TotalSystemTransactions)
	fmt.Printf("Successfully Matched:           %d\n", result.Summary.SuccessfullyMatchedCount)
	fmt.Printf("Missing in Internal System:     %d\n", result.Summary.MissingInInternalCount)
	fmt.Printf("Missing in Source:              %d\n", result.Summary.MissingInSourceCount)
	fmt.Printf("Mismatched Transactions:        %d\n", result.Summary.MismatchedTransactionsCount)
	fmt.Println(separator)

	// Calculate reconciliation rate
	totalPossibleMatches := result.Summary.TotalSourceTransactions
	if result.Summary.TotalSystemTransactions < totalPossibleMatches {
		totalPossibleMatches = result.Summary.TotalSystemTransactions
	}
	
	if totalPossibleMatches > 0 {
		matchRate := float64(result.Summary.SuccessfullyMatchedCount) / float64(totalPossibleMatches) * 100
		fmt.Printf("Reconciliation Rate:            %.2f%%\n", matchRate)
	}
	fmt.Println(separator)
}
