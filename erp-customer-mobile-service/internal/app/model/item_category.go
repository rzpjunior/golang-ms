package model

type ItemCategory struct {
	ID       int64  `orm:"column(id);auto" json:"-"`
	RegionID string `orm:"column(region_id)" json:"area"`
	Name     string `orm:"column(name)" json:"name"`
	ImageUrl string `orm:"column(image_url)" json:"image_url"`
	Note     string `orm:"column(note)" json:"note"`
	Status   int8   `orm:"column(status)" json:"status"`
}
