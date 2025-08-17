package main

import (
	"math"
	"strings"
	"time"
)

// TransactionReconciler handles the reconciliation logic
type TransactionReconciler struct{}

// NewTransactionReconciler creates a new reconciler instance
func NewTransactionReconciler() *TransactionReconciler {
	return &TransactionReconciler{}
}

// Reconcile performs the reconciliation between source and system transactions
func (tr *TransactionReconciler) Reconcile(sourceTransactions []SourceTransaction, systemTransactions []SystemTransaction) *ReconciliationResult {
	// Create maps for efficient lookup
	sourceMap := make(map[string]SourceTransaction)
	systemMap := make(map[string]SystemTransaction)

	// Index source transactions by their ID
	for _, txn := range sourceTransactions {
		sourceMap[txn.ProviderTransactionID] = txn
	}

	// Index system transactions by their ID
	for _, txn := range systemTransactions {
		systemMap[txn.TransactionID] = txn
	}

	var missingInInternal []SourceTransaction
	var missingInSource []SystemTransaction
	var mismatchedTransactions []MismatchedTransaction
	matchedCount := 0

	// Find transactions missing in internal system and mismatched transactions
	for id, sourceTxn := range sourceMap {
		if systemTxn, exists := systemMap[id]; exists {
			// Transaction exists in both systems, check for discrepancies
			discrepancies := tr.findDiscrepancies(sourceTxn, systemTxn)
			if len(discrepancies) > 0 {
				mismatchedTransactions = append(mismatchedTransactions, MismatchedTransaction{
					TransactionID: id,
					Discrepancies: discrepancies,
				})
			} else {
				matchedCount++
			}
		} else {
			// Transaction exists in source but not in internal system
			missingInInternal = append(missingInInternal, sourceTxn)
		}
	}

	// Find transactions missing in source
	for id, systemTxn := range systemMap {
		if _, exists := sourceMap[id]; !exists {
			// Transaction exists in system but not in source
			missingInSource = append(missingInSource, systemTxn)
		}
	}

	// Create summary
	summary := ReconciliationSummary{
		TotalSourceTransactions:     len(sourceTransactions),
		TotalSystemTransactions:     len(systemTransactions),
		MissingInInternalCount:      len(missingInInternal),
		MissingInSourceCount:        len(missingInSource),
		MismatchedTransactionsCount: len(mismatchedTransactions),
		SuccessfullyMatchedCount:    matchedCount,
	}

	return &ReconciliationResult{
		MissingInInternal:      missingInInternal,
		MissingInSource:        missingInSource,
		MismatchedTransactions: mismatchedTransactions,
		Summary:                summary,
	}
}

// findDiscrepancies compares a source transaction with a system transaction and returns discrepancies
func (tr *TransactionReconciler) findDiscrepancies(source SourceTransaction, system SystemTransaction) map[string]Discrepancy {
	discrepancies := make(map[string]Discrepancy)

	// Compare User ID
	if source.UserID != system.UserID {
		discrepancies["userId"] = Discrepancy{
			Source: source.UserID,
			System: system.UserID,
		}
	}

	// Compare amounts with tolerance for floating point precision
	if !tr.isAmountEqual(source.Amount, system.Amount) {
		discrepancies["amount"] = Discrepancy{
			Source: source.Amount,
			System: system.Amount,
		}
	}

	// Compare currency
	if source.Currency != system.Currency {
		discrepancies["currency"] = Discrepancy{
			Source: source.Currency,
			System: system.Currency,
		}
	}

	// Compare statuses (normalize before comparison)
	normalizedSourceStatus := tr.normalizeStatus(source.Status)
	normalizedSystemStatus := tr.normalizeStatus(system.Status)
	if normalizedSourceStatus != normalizedSystemStatus {
		discrepancies["status"] = Discrepancy{
			Source: source.Status,
			System: system.Status,
		}
	}

	// Compare payment method
	if source.PaymentMethod != system.PaymentMethod {
		discrepancies["paymentMethod"] = Discrepancy{
			Source: source.PaymentMethod,
			System: system.PaymentMethod,
		}
	}

	// Compare created timestamps (allow small tolerance for time differences)
	if !tr.isTimeEqual(source.CreatedAt, system.CreatedAt) {
		discrepancies["createdAt"] = Discrepancy{
			Source: source.CreatedAt.Format(time.RFC3339),
			System: system.CreatedAt.Format(time.RFC3339),
		}
	}

	// Compare updated timestamps (allow small tolerance for time differences)
	if !tr.isTimeEqual(source.UpdatedAt, system.UpdatedAt) {
		discrepancies["updatedAt"] = Discrepancy{
			Source: source.UpdatedAt.Format(time.RFC3339),
			System: system.UpdatedAt.Format(time.RFC3339),
		}
	}

	// Compare provider reference with system reference ID
	if source.ProviderReference != system.ReferenceID {
		discrepancies["referenceId"] = Discrepancy{
			Source: source.ProviderReference,
			System: system.ReferenceID,
		}
	}

	return discrepancies
}

// isAmountEqual compares two amounts with a small tolerance for floating point precision
func (tr *TransactionReconciler) isAmountEqual(amount1, amount2 float64) bool {
	tolerance := 0.01 // 1 cent tolerance
	return math.Abs(amount1-amount2) < tolerance
}

// normalizeStatus standardizes status values from different systems to a common format
// This ensures accurate matching despite different naming conventions between source and system data
func (tr *TransactionReconciler) normalizeStatus(status string) string {
	// Convert to uppercase for case-insensitive comparison
	normalizedStatus := strings.ToUpper(strings.TrimSpace(status))

	// Handle SUCCEEDED and COMPLETED as the same thing
	if normalizedStatus == "SUCCEEDED" || normalizedStatus == "COMPLETED" {
		return "COMPLETED"
	}

	return normalizedStatus
}

// isTimeEqual compares two timestamps with a tolerance for small time differences
// This handles cases where timestamps might be slightly different due to processing delays
func (tr *TransactionReconciler) isTimeEqual(time1, time2 time.Time) bool {
	tolerance := 5 * time.Second // 5 second tolerance
	return time1.Sub(time2).Abs() <= tolerance
}
