package dto

type WrtRequest struct {
	Platform string     `json:"platform" valid:"required"`
	Data     dataGetWrt `json:"data" valid:"required"`
}

type dataGetWrt struct {
	RegionID     string `json:"region_id" valid:"required"`
	Type         string `json:"type" valid:"required"`
	DataResponse WrtResponse
}

type WrtResponse struct {
	ID     string `orm:"column(id);auto" json:"id"`
	Code   string `orm:"column(code)" json:"code"`
	Name   string `orm:"column(name)" json:"name"`
	Note   string `orm:"column(note)" json:"note"`
	Status string `orm:"column(status)" json:"status"`
	Type   string `orm:"column(type)" json:"type"`

	Region *RegionResponse `json:"region,omitempty"`
}
