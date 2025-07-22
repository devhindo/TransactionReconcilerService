# TransactionReconcilerService

This project is a Go-based service for reconciling transactions. It is designed to be a simple and efficient solution for managing transaction data. It's part of the initial step towards joining the Extreme Solution team as a backend developer.

## Project Creation

Using `gh` cli tool:

![Project Creation](./assets/imgs/projecCreation.png)

## Project Configuration


- Initializing Go project:

![Initializing Go project](./assets/imgs/initGoProject.png)

## Adding the CSVs files

![Adding the CSVs files](./assets/imgs/CSVsFiles.png)

## Created the necessary data structures used across the service

seen in [models.go](./models.go)



## created an struct and constructor for the service

```Go
// TransactionReconciliationService is the main service that orchestrates the reconciliation process
type TransactionReconciliationService struct {
	csvReader  *CSVReader
	reconciler *TransactionReconciler
}

// NewTransactionReconciliationService creates a new service instance
func NewTransactionReconciliationService() *TransactionReconciliationService {
	return &TransactionReconciliationService{
		csvReader:  NewCSVReader(),
		reconciler: NewTransactionReconciler(),
	}
}
```
## Creating the CSVReader struct and constructor

can be seen in [csv_reader.go](./csv_reader.go)

```Go
type CSVReader struct{}

// NewCSVReader constructor to create a new CSV reader instance
func NewCSVReader() *CSVReader {
	return &CSVReader{}
}

// ReadSourceTransactions reads source transactions from a CSV file
func (r *CSVReader) ReadSourceTransactions(filePath string) ([]SourceTransaction, error) {}

// ReadSystemTransactions reads system transactions from a CSV file
func (r *CSVReader) ReadSystemTransactions(filePath string) ([]SystemTransaction, error) {}
```
## Creating reconciler struct to implement the reconciliation logic

can be seen in [reconciler.go](./reconciler.go)

```Go

// TransactionReconciler handles the reconciliation logic
type TransactionReconciler struct{}

// NewTransactionReconciler creates a new reconciler instance
func NewTransactionReconciler() *TransactionReconciler {
	return &TransactionReconciler{}
}

// Reconcile performs the reconciliation between source and system transactions
func (tr *TransactionReconciler) Reconcile(sourceTransactions []SourceTransaction, systemTransactions []SystemTransaction) *ReconciliationResult {}

// findDiscrepancies compares a source transaction with a system transaction and returns discrepancies
func (tr *TransactionReconciler) findDiscrepancies(source SourceTransaction, system SystemTransaction) map[string]Discrepancy {
	discrepancies := make(map[string]Discrepancy) {}

// the logic to compare amounts, I added a tolerance for floating point percision
func (tr *TransactionReconciler) isAmountEqual(amount1, amount2 float64) bool {}

// normalize the status values for comparison
func (tr *TransactionReconciler) normalizeStatus(status string) string
```
