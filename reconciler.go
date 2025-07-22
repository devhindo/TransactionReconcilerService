package main

import (
	"math"
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

	// Compare amounts with tolerance for floating point precision
	if !tr.isAmountEqual(source.Amount, system.Amount) {
		discrepancies["amount"] = Discrepancy{
			Source: source.Amount,
			System: system.Amount,
		}
	}

	// Compare statuses (normalize status values for comparison)
	sourceStatus := tr.normalizeStatus(source.Status)
	systemStatus := tr.normalizeStatus(system.Status)
	if sourceStatus != systemStatus {
		discrepancies["status"] = Discrepancy{
			Source: source.Status,
			System: system.Status,
		}
	}

	// Compare currencies
	if source.Currency != system.Currency {
		discrepancies["currency"] = Discrepancy{
			Source: source.Currency,
			System: system.Currency,
		}
	}

	// Compare user IDs if they exist
	if source.UserID != "" && system.UserID != "" && source.UserID != system.UserID {
		discrepancies["userId"] = Discrepancy{
			Source: source.UserID,
			System: system.UserID,
		}
	}

	return discrepancies
}

// isAmountEqual compares two amounts with a small tolerance for floating point precision
func (tr *TransactionReconciler) isAmountEqual(amount1, amount2 float64) bool {
	tolerance := 0.01 // 1 cent tolerance
	return math.Abs(amount1-amount2) < tolerance
}

// normalizeStatus normalizes status values for comparison
// Maps similar statuses to common values for better matching
func (tr *TransactionReconciler) normalizeStatus(status string) string {
	statusMap := map[string]string{
		"succeeded": "completed",
		"success":   "completed",
		"completed": "completed",
		"pending":   "pending",
		"failed":    "failed",
		"refunded":  "refunded",
		"disputed":  "disputed",
	}

	if normalized, exists := statusMap[status]; exists {
		return normalized
	}
	return status
}
