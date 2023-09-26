package dto

type WrtResponse struct {
	ID        int64           `json:"id"`
	RegionID  int64           `json:"region_id"`
	Code      string          `json:"code"`
	StartTime string          `json:"start_time"`
	EndTime   string          `json:"end_time"`
	Region    *RegionResponse `json:"region"`
}
