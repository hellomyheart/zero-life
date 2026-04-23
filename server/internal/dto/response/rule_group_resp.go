package response

import "time"

type RuleGroupResp struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Name      string    `json:"name"`
	Order     int       `json:"order"`
	IsActive  bool      `json:"is_active"`
	RuleCount int64     `json:"rule_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RuleGroupExecuteResultResp struct {
	GroupID      int      `json:"group_id"`
	MatchedCount int      `json:"matched_count"`
	SuccessCount int      `json:"success_count"`
	FailCount    int      `json:"fail_count"`
	Errors       []string `json:"errors,omitempty"`
}
