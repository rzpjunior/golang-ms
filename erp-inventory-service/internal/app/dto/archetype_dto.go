package dto

type ArchetypeResponse struct {
	ID             string                `json:"id"`
	Code           string                `json:"code"`
	Description    string                `json:"description"`
	CustomerTypeID string                `json:"customer_type_id"`
	Status         int8                  `json:"status"`
	ConvertStatus  string                `json:"convert_status"`
	CustomerType   *CustomerTypeResponse `json:"customer_type,omitempty"`
}

type CustomerTypeResponse struct {
	ID            string `json:"id"`
	Code          string `json:"code"`
	Description   string `json:"description"`
	CustomerGroup string `json:"customer_group,omitempty"`
	Status        int8   `json:"status,omitempty"`
	ConvertStatus string `json:"convert_status,omitempty"`
}
