package main

import (
	"time"
)

// SourceTransaction represents a transaction from the external provider like Stripe, to be parsed from source_transactions.csv
type SourceTransaction struct {
	ProviderTransactionID string    `csv:"providerTransactionId" json:"providerTransactionId"`
	Email                 string    `csv:"email" json:"email"`
	UserID                string    `csv:"userId" json:"userId"`
	Provider              string    `csv:"provider" json:"provider"`
	Amount                float64   `csv:"amount" json:"amount"`
	Currency              string    `csv:"currency" json:"currency"`
	Status                string    `csv:"status" json:"status"`
	TransactionType       string    `csv:"transactionType" json:"transactionType"`
	PaymentMethod         string    `csv:"paymentMethod" json:"paymentMethod"`
	CreatedAt             time.Time `csv:"createdAt" json:"createdAt"`
	UpdatedAt             time.Time `csv:"updatedAt" json:"updatedAt"`
	ProviderReference     string    `csv:"providerReference" json:"providerReference"`
	FraudRisk             string    `csv:"fraudRisk" json:"fraudRisk"`
	DetailsInvoiceID      string    `csv:"details_invoiceId" json:"details_invoiceId"`
	DetailsCustomerName   string    `csv:"details_customerName" json:"details_customerName"`
	DetailsDescription    string    `csv:"details_description" json:"details_description"`
}

// SystemTransaction represents an internal system transaction, to be parsed from system_transactions.csv
type SystemTransaction struct {
	TransactionID         string    `csv:"transactionId" json:"transactionId"`
	UserID                string    `csv:"userId" json:"userId"`
	Amount                float64   `csv:"amount" json:"amount"`
	Currency              string    `csv:"currency" json:"currency"`
	Status                string    `csv:"status" json:"status"`
	PaymentMethod         string    `csv:"paymentMethod" json:"paymentMethod"`
	CreatedAt             time.Time `csv:"createdAt" json:"createdAt"`
	UpdatedAt             time.Time `csv:"updatedAt" json:"updatedAt"`
	ReferenceID           string    `csv:"referenceId" json:"referenceId"`
	MetadataOrderID       string    `csv:"metadata_orderId" json:"metadata_orderId"`
	MetadataDescription   string    `csv:"metadata_description" json:"metadata_description"`
}

// Discrepancy represents a field mismatch between source and system
type Discrepancy struct {
	Source interface{} `json:"source"`
	System interface{} `json:"system"`
}

// MismatchedTransaction represents transactions with the same ID but different amounts/statuses
type MismatchedTransaction struct {
	TransactionID string                 `json:"transactionId"`
	Discrepancies map[string]Discrepancy `json:"discrepancies"`
}

// ReconciliationResult represents the complete reconciliation report
type ReconciliationResult struct {
	MissingInInternal      []SourceTransaction     `json:"missing_in_internal"`
	MissingInSource        []SystemTransaction     `json:"missing_in_source"`
	MismatchedTransactions []MismatchedTransaction `json:"mismatched_transactions"`
	Summary                ReconciliationSummary   `json:"summary"`
}

// ReconciliationSummary provides statistics about the reconciliation
type ReconciliationSummary struct {
	TotalSourceTransactions      int `json:"total_source_transactions"`
	TotalSystemTransactions      int `json:"total_system_transactions"`
	MissingInInternalCount       int `json:"missing_in_internal_count"`
	MissingInSourceCount         int `json:"missing_in_source_count"`
	MismatchedTransactionsCount  int `json:"mismatched_transactions_count"`
	SuccessfullyMatchedCount     int `json:"successfully_matched_count"`
}
