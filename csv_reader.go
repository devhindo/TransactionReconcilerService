package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

// creating CSVReader struct (class) handles reading and parsing CSV files
type CSVReader struct{}

// NewCSVReader constructor to create a new CSV reader instance
func NewCSVReader() *CSVReader {
	return &CSVReader{}
}

// ReadSourceTransactions reads and parses source transactions from CSV file
func (r *CSVReader) ReadSourceTransactions(filePath string) ([]SourceTransaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source transactions file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV records: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	// Skip header row
	records = records[1:]
	transactions := make([]SourceTransaction, 0, len(records))

	for i, record := range records {
		if len(record) < 16 {
			return nil, fmt.Errorf("invalid record at line %d: expected 16 fields, got %d", i+2, len(record))
		}

		amount, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid amount at line %d: %w", i+2, err)
		}

		createdAt, err := time.Parse(time.RFC3339, record[9])
		if err != nil {
			return nil, fmt.Errorf("invalid createdAt at line %d: %w", i+2, err)
		}

		updatedAt, err := time.Parse(time.RFC3339, record[10])
		if err != nil {
			return nil, fmt.Errorf("invalid updatedAt at line %d: %w", i+2, err)
		}

		transaction := SourceTransaction{
			ProviderTransactionID: record[0],
			Email:                 record[1],
			UserID:                record[2],
			Provider:              record[3],
			Amount:                amount,
			Currency:              record[5],
			Status:                record[6],
			TransactionType:       record[7],
			PaymentMethod:         record[8],
			CreatedAt:             createdAt,
			UpdatedAt:             updatedAt,
			ProviderReference:     record[11],
			FraudRisk:             record[12],
			DetailsInvoiceID:      record[13],
			DetailsCustomerName:   record[14],
			DetailsDescription:    record[15],
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// ReadSystemTransactions reads and parses system transactions from CSV file
func (r *CSVReader) ReadSystemTransactions(filePath string) ([]SystemTransaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open system transactions file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV records: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	// Skip header row
	records = records[1:]
	transactions := make([]SystemTransaction, 0, len(records))

	for i, record := range records {
		if len(record) < 11 {
			return nil, fmt.Errorf("invalid record at line %d: expected 11 fields, got %d", i+2, len(record))
		}

		amount, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid amount at line %d: %w", i+2, err)
		}

		createdAt, err := time.Parse(time.RFC3339, record[6])
		if err != nil {
			return nil, fmt.Errorf("invalid createdAt at line %d: %w", i+2, err)
		}

		updatedAt, err := time.Parse(time.RFC3339, record[7])
		if err != nil {
			return nil, fmt.Errorf("invalid updatedAt at line %d: %w", i+2, err)
		}

		transaction := SystemTransaction{
			TransactionID:       record[0],
			UserID:              record[1],
			Amount:              amount,
			Currency:            record[3],
			Status:              record[4],
			PaymentMethod:       record[5],
			CreatedAt:           createdAt,
			UpdatedAt:           updatedAt,
			ReferenceID:         record[8],
			MetadataOrderID:     record[9],
			MetadataDescription: record[10],
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
