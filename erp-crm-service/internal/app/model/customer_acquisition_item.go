package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type CustomerAcquisitionItem struct {
	ID                    int64     `orm:"column(id)" json:"id"`
	CustomerAcquisitionID int64     `orm:"column(customer_acquisition_id)" json:"customer_acquisition_id"`
	ItemID                int64     `orm:"column(item_id)" json:"item_id"`
	ItemIDGP              string    `orm:"column(item_id_gp)" json:"item_id_gp"`
	IsTop                 int8      `orm:"column(is_top)" json:"is_top"`
	CreatedAt             time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt             time.Time `orm:"column(updated_at)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(CustomerAcquisitionItem))
}

func (m *CustomerAcquisitionItem) TableName() string {
	return "customer_acquisition_item"
}
