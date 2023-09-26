package dto

type WrtResponse struct {
	ID        string `json:"id"`
	RegionId  string `json:"region_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    int8   `json:"status"`
}

// Get Wrt
type GetWrtRequest struct {
	Limit    int
	Offset   int
	Search   string
	RegionId int
	SiteId   string
}
