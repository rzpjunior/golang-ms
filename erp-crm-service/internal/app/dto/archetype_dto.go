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

type ArchetypeGetListRequest struct {
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
	Search         string `json:"search"`
	Status         int8   `json:"status"`
	CustomerTypeID string `json:"customer_type_id"`
}
