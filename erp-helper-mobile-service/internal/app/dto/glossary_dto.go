package dto

type GlossaryResponse struct {
	ID        int64  `json:"id,omitempty"`
	Table     string `json:"table,omitempty"`
	Attribute string `json:"attribute,omitempty"`
	ValueInt  int8   `json:"value_int,omitempty"`
	ValueName string `json:"value_name,omitempty"`
	Note      string `json:"note,omitempty"`
}

// Get Glossary
type GetGlossaryRequest struct {
	Table     string
	Attribute string
	ValueInt  int
	ValueName string
}
