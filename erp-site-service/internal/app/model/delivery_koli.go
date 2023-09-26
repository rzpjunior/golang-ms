package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type DeliveryKoli struct {
	Id        int64   `orm:"column(id)" json:"id"`
	SopNumber string  `orm:"column(sales_order_code)" json:"sop_number"`
	KoliId    int64   `orm:"column(koli_id)" json:"koli_id"`
	Quantity  float64 `orm:"column(quantity)" json:"quantity"`
	Note      string  `orm:"column(note)" json:"note"`
}

func init() {
	orm.RegisterModel(new(DeliveryKoli))
}

func (m *DeliveryKoli) TableName() string {
	return "delivery_koli"
}
