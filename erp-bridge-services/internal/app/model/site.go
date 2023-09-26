package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Site struct {
	ID          int64     `orm:"column(id)" json:"id"`
	Code        string    `orm:"column(code)" json:"code"`
	Description string    `orm:"column(description)" json:"description"`
	RegionId    int64     `orm:"column(region_id)" json:"region_id"`
	Status      int8      `orm:"column(status)" json:"status"`
	CreatedAt   time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt   time.Time `orm:"column(updated_at)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(Site))
}

func (m *Site) TableName() string {
	return "site"
}
