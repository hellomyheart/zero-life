package response

import "time"

type RecurrenceResp struct {
	ID             uint64     `json:"id"`
	UserID         uint64     `json:"user_id"`
	Title          string     `json:"title"`
	Type           string     `json:"type"`
	Amount         string     `json:"amount"`
	SourceID       uint64     `json:"source_id"`
	SourceName     string     `json:"source_name"`
	DestinationID  *uint64    `json:"destination_id"`
	DestinationName string    `json:"destination_name,omitempty"`
	CategoryID     *uint64    `json:"category_id"`
	Description    string     `json:"description"`
	Notes          string     `json:"notes"`
	TagNames       string     `json:"tag_names"`
	RepeatFreq     string     `json:"repeat_freq"`
	RepeatInterval int        `json:"repeat_interval"`
	NextDate       time.Time  `json:"next_date"`
	EndDate        *time.Time `json:"end_date"`
	Repetitions    int        `json:"repetitions"`
	MaxRepetitions *int       `json:"max_repetitions"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
