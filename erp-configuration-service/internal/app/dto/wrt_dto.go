package dto

type WrtResponse struct {
	ID       int64           `json:"id"`
	RegionID string          `json:"region_id"`
	Code     string          `json:"code"`
	Name     string          `json:"name"`
	Type     int8            `json:"type"`
	Note     string          `json:"note"`
	Region   *RegionResponse `json:"region"`
}

type WrtRequestUpdate struct {
	Type int8   `json:"type" valid:"required"`
	Note string `json:"note"`
}
