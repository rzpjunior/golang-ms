package dto

import "time"

type PickingOrderDetailResponse struct {
	Id         string    `json:"id,omitempty"`
	DocDate    time.Time `json:"doc_date,omitempty"`
	SopNumber  string    `json:"sop_number,omitempty"`
	ItemNumber string    `json:"item_number,omitempty"`
	Status     int8      `json:"status,omitempty"`
}

// Get PickingOrderDetail
type GetPickingOrderDetailRequest struct {
	Limit      int
	Offset     int
	Status     []int
	Search     string
	OrderBy    string
	SopNumber  string
	ItemNumber string
}
