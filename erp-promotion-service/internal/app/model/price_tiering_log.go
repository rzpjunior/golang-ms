package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type PriceTieringLog struct {
	ID               int64     `orm:"column(id);auto" json:"id"`
	PriceTieringIDGP string    `orm:"column(price_tiering_id_gp)" json:"price_tiering_id_gp"`
	CustomerID       int64     `orm:"column(customer_id)" json:"customer_id"`
	AddressIDGP      string    `orm:"column(address_id_gp)" json:"address_id_gp"`
	SalesOrderIDGP   string    `orm:"column(sales_order_id_gp)" json:"sales_order_id_gp"`
	ItemID           int64     `orm:"column(item_id)" json:"item_id"`
	DiscountQty      float64   `orm:"column(discount_qty)" json:"discount_qty"`
	DiscountAmount   float64   `orm:"column(discount_amount)" json:"discount_amount"`
	CreatedAt        time.Time `orm:"column(created_at)" json:"created_at"`
	Status           int8      `orm:"column(status)" json:"status"`
}

func init() {
	orm.RegisterModel(new(PriceTieringLog))
}

func (m *PriceTieringLog) TableName() string {
	return "price_tiering_log"
}

func (m *PriceTieringLog) MarshalJSON() ([]byte, error) {
	type Alias PriceTieringLog

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
