package response

type ImportUploadResp struct {
	FileID string `json:"file_id"`
}

type ImportPreviewResp struct {
	Total   int              `json:"total"`
	Valid   int              `json:"valid"`
	Invalid int              `json:"invalid"`
	Rows    []ImportRowResp  `json:"rows"`
}

type ImportRowResp struct {
	Index   int              `json:"index"`
	Data    map[string]string `json:"data"`
	IsValid bool             `json:"is_valid"`
	Errors  []string         `json:"errors,omitempty"`
}

type ImportResultResp struct {
	Total   int `json:"total"`
	Success int `json:"success"`
	Failed  int `json:"failed"`
	Skipped int `json:"skipped"`
}
