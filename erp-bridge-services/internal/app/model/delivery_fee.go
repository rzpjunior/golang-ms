package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type DeliveryFee struct {
	ID            int64   `orm:"column(id)" json:"id"`
	Code          string  `orm:"column(code)" json:"code"`
	Name          string  `orm:"column(name)" json:"name,omitempty"`
	Note          string  `orm:"column(note)" json:"note,omitempty"`
	Status        int32   `orm:"column(status)" json:"status,omitempty"`
	MinimumOrder  float64 `orm:"column(minimum_order)" json:"minimum_order,omitempty"`
	DeliveryFee   float64 `orm:"column(delivery_fee)" json:"delivery_fee,omitempty"`
	RegionId      int64   `orm:"column(region_id)" json:"region_id,omitempty"`
	CutomerTypeId string  `orm:"column(customer_type_id)" json:"customer_type_id,omitempty"`
}

func init() {
	orm.RegisterModel(new(DeliveryFee))
}

func (m *DeliveryFee) TableName() string {
	return "delivery_fee"
}
