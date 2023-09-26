package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Bank struct {
	ID              int64     `orm:"column(id)" json:"id"`
	Code            string    `orm:"column(code)" json:"code"`
	Description     string    `orm:"column(description)" json:"description"`
	Value           string    `orm:"column(value)" json:"value"`
	ImageUrl        string    `orm:"column(image_url)" json:"image_url"`
	PaymentGuideUrl string    `orm:"column(payment_guide_url)" json:"payment_guide_url"`
	PublishIVA      int8      `orm:"column(publish_iva)" json:"publish_iva"`
	PublishFVA      int8      `orm:"column(publish_fva)" json:"publish_fva"`
	Status          int8      `orm:"column(status)" json:"status"`
	CreatedAt       time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt       time.Time `orm:"column(updated_at)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(Bank))
}

func (m *Bank) TableName() string {
	return "bank"
}
