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


