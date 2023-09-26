package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type PurchaseOrderImage struct {
	ID                int64     `orm:"column(id)" json:"id"`
	PurchaseOrderID   int64     `orm:"column(purchase_order_id)" json:"purchase_order_id"`
	PurchaseOrderIDGP string    `orm:"column(purchase_order_id_gp)" json:"purchase_order_id_gp"`
	ImageURL          string    `orm:"column(image_url)" json:"image_url"`
	CreatedAt         time.Time `orm:"column(created_at)" json:"created_at"`
}

func init() {
	orm.RegisterModel(new(PurchaseOrderImage))
}

func (m *PurchaseOrderImage) TableName() string {
	return "purchase_order_image"
}
