package dto

// type WrtResponse struct {
// 	ID        int64           `json:"id"`
// 	RegionID  int64           `json:"region_id"`
// 	Code      string          `json:"code"`
// 	StartTime string          `json:"start_time"`
// 	EndTime   string          `json:"end_time"`
// 	Region    *RegionResponse `json:"region"`
// }

type WrtResponse struct {
	ID     string `orm:"column(id);auto" json:"id"`
	Code   string `orm:"column(code)" json:"code"`
	Name   string `orm:"column(name)" json:"name"`
	Note   string `orm:"column(note)" json:"note"`
	Status string `orm:"column(status)" json:"status"`
	Type   string `orm:"column(type)" json:"type"`

	Region *RegionResponse `json:"region,omitempty"`
}

type GetWrtListRequest struct {
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
	Search   string `json:"search"`
	RegionId string `json:"region_id"`
	Type     int32  `json:"type"`
}

type GetWrtDetailRequest struct {
	Id int32 `json:"id"`
}

type WrtGP struct {
	GnL_Region string `json:"gnL_Region,omitempty"`
	GnL_WRT_ID string `json:"gnL_WRT_ID,omitempty"`
	Strttime   string `json:"strttime"`
	Endtime    string `json:"endtime"`
	Inactive   *int32 `json:"inactive,omitempty"`
}
