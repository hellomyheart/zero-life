package request

type ImportParseReq struct {
	FileID    string            `json:"file_id" binding:"required"`
	Mapping   map[string]string `json:"mapping" binding:"required"`
}

type ImportExecuteReq struct {
	FileID    string            `json:"file_id" binding:"required"`
	Mapping   map[string]string `json:"mapping" binding:"required"`
}
