package dto

type KoliResponse struct {
	Id     int64  `json:"id,omitempty"`
	Code   string `json:"code,omitempty"`
	Value  string `json:"value,omitempty"`
	Name   string `json:"name,omitempty"`
	Note   string `json:"note,omitempty"`
	Status int8   `json:"status,omitempty"`
}

type GetKoliRequest struct {
	Offset  int
	Limit   int
	OrderBy string
	Status  int
}
