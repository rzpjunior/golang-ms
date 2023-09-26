package dto

type CheckbookListRequest struct {
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
	Status   int32  `json:"status"`
	RegionID string `json:"region_id"`
	OrderBy  string `json:"order_by"`
}

type CheckbookGP struct {
	ID            string `json:"id"`
	CheckbookCode string `json:"checkbook_code"`
	CheckbookDesc string `json:"checkbook_desc"`
}
