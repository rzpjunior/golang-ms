package dto

import "time"

type PickingOrderResponse struct {
	Id        string    `json:"id,omitempty"`
	DocDate   time.Time `json:"doc_date,omitempty"`
	SopNumber string    `json:"sop_number,omitempty"`
	PickerId  string    `json:"picker_id,omitempty"`
	Status    int8      `json:"status,omitempty"`
}

// Get PickingOrder
type GetPickingOrderRequest struct {
	Limit     int
	Offset    int
	Status    []int
	Search    string
	OrderBy   string
	SopNumber string
	PickerId  string
}
