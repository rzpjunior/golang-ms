package dto

type GlossaryResponse struct {
	ID        int64  `json:"id,omitempty"`
	Table     string `json:"table,omitempty"`
	Attribute string `json:"attribute,omitempty"`
	ValueInt  int8   `json:"value_int"`
	ValueName string `json:"value_name"`
	Note      string `json:"note"`
}
