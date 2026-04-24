package response

import "time"

type RecurringTransactionResp struct {
	ID             uint64     `json:"id"`
	Description    string     `json:"description"`
	Amount         string     `json:"amount"`
	SourceID       uint64     `json:"source_id"`
	DestinationID  *uint64    `json:"destination_id"`
	CategoryID     *uint64    `json:"category_id"`
	Notes          string     `json:"notes"`
	RecurrenceType string     `json:"recurrence_type"`
	RepeatEvery    int        `json:"repeat_every"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	NextOccurrence time.Time  `json:"next_occurrence"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type RecurringTransactionLogResp struct {
	ID                  uint64    `json:"id"`
	RecurringTransactionID uint64 `json:"recurring_transaction_id"`
	TransactionID       uint64    `json:"transaction_id"`
	OccurrenceDate      time.Time `json:"occurrence_date"`
	CreatedAt           time.Time `json:"created_at"`
}