package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func main() {

	fmt.Println("Hello Extreme Solution!")
}