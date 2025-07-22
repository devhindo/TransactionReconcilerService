package models

import (
	"time"
)

// SourceTransaction represents a transaction from the external provider like Stripe, to be parsed from source_transactions.csv
type SourceTransaction struct {
	ProviderTransactionID string    `json:"providerTransactionId"`
	Email                 string    `json:"email"`
	UserID                string    `json:"userId"`
	Provider              string    `json:"provider"`
	Amount                float64   `json:"amount"`
	Currency              string    `json:"currency"`
	Status                string    `json:"status"`
	TransactionType       string    `json:"transactionType"`
	PaymentMethod         string    `json:"paymentMethod"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
	ProviderReference     string    `json:"providerReference"`
	FraudRisk             string    `json:"fraudRisk"`
	DetailsInvoiceID      string    `json:"details_invoiceId"`
	DetailsCustomerName   string    `json:"details_customerName"`
	DetailsDescription    string    `json:"details_description"`
}

// SystemTransaction represents an internal system transaction, to be parsed from system_transactions.csv
type SystemTransaction struct {
	TransactionID       string    `json:"transactionId"`
	UserID              string    `json:"userId"`
	Amount              float64   `json:"amount"`
	Currency            string    `json:"currency"`
	Status              string    `json:"status"`
	PaymentMethod       string    `json:"paymentMethod"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
	ReferenceID         string    `json:"referenceId"`
	MetadataOrderID     string    `json:"metadata_orderId"`
	MetadataDescription string    `json:"metadata_description"`
}


// ReconciliationReport represents the final reconciliation result to be reported at the end of the process
type ReconciliationReport struct {
	MissingInInternal      []SourceTransaction     `json:"missing_in_internal"`
	MissingInSource        []SystemTransaction     `json:"missing_in_source"`
	MismatchedTransactions []MismatchedTransaction `json:"mismatched_transactions"`
}

// MismatchedTransaction represents transactions with the same ID but different amounts/statuses
type MismatchedTransaction struct {
	TransactionID string                 `json:"transactionId"`
	Discrepancies map[string]interface{} `json:"discrepancies"`
}

// Discrepancy represents a field mismatch between source and system
type Discrepancy struct {
	Source interface{} `json:"source"`
	System interface{} `json:"system"`
}
