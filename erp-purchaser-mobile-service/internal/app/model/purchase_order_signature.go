package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type PurchaseOrderSignature struct {
	ID                int64     `orm:"column(id)" json:"id"`
	PurchaseOrderID   int64     `orm:"column(purchase_order_id)" json:"purchase_order_id"`
	PurchaseOrderIDGP string    `orm:"column(purchase_order_id_gp)" json:"purchase_order_id_gp"`
	JobFunction       string    `orm:"column(job_function)" json:"job_function"`
	Name              string    `orm:"column(name)" json:"name"`
	SignatureURL      string    `orm:"column(signature_url)" json:"signature_url"`
	CreatedAt         time.Time `orm:"column(created_at)" json:"created_at"`
	CreatedBy         int64     `orm:"column(created_by)" json:"created_by"`
}

func init() {
	orm.RegisterModel(new(PurchaseOrderSignature))
}

func (m *PurchaseOrderSignature) TableName() string {
	return "purchase_order_signature"
}
