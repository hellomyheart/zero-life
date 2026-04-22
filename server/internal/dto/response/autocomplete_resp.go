package response

type AutocompleteItemResp struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type,omitempty"`
	Symbol string `json:"symbol,omitempty"`
	Code   string `json:"code,omitempty"`
	Color  string `json:"color,omitempty"`
}
