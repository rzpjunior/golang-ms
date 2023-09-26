package dto

import "time"

type PickingListResponse struct {
	Id              string    `json:"id,omitempty"`
	DocDate         time.Time `json:"doc_date,omitempty"`
	SiteId          string    `json:"site_id,omitempty"`
	RequestShipDate time.Time `json:"request_ship_date,omitempty"`
	Status          int8      `json:"status,omitempty"`
}

// Get PickingList
type GetPickingListRequest struct {
	Limit           int
	Offset          int
	Status          []int
	Search          string
	OrderBy         string
	SiteId          string
	RequestShipDate time.Time
}
