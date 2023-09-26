package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type PackingOrder struct {
	ID           int64     `orm:"column(id)" json:"id"`
	Code         string    `orm:"column(code)" json:"code"`
	SiteID       int64     `orm:"column(site_id)" json:"site_id"`
	SiteIDGP     string    `orm:"column(site_id_gp)" json:"site_id_gp"`
	RegionID     int64     `orm:"column(region_id)" json:"region_id"`
	RegionIDGP   string    `orm:"column(region_id_gp)" json:"region_id_gp"`
	DeliveryDate time.Time `orm:"column(delivery_date)" json:"delivery_date"`
	Note         string    `orm:"column(note)" json:"note"`
	Item         string    `orm:"column(item)" json:"item"`
	CreatedAt    time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt    time.Time `orm:"column(updated_at)" json:"updated_at"`
	Status       int8      `orm:"column(status)" json:"status"`
}

func init() {
	orm.RegisterModel(new(PackingOrder))
}

func (m *PackingOrder) TableName() string {
	return "packing_order"
}
