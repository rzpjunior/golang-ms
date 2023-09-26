package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

func init() {
	orm.RegisterModel(new(PaymentChannel))
}

// PaymentChannel : struct to hold payment term model data for database
type PaymentChannel struct {
	ID              int64  `orm:"column(id);auto" json:"-"`
	Code            string `orm:"column(code);size(50);null" json:"code"`
	Value           string `orm:"column(value);size(50);null" json:"value"`
	Name            string `orm:"column(name);size(100);null" json:"name"`
	ImageUrl        string `orm:"column(image_url);size(300);null" json:"image_url"`
	Note            string `orm:"column(note);size(255)" json:"note"`
	Status          int8   `orm:"column(status)" json:"status"`
	PublishIva      int8   `orm:"column(publish_iva)" json:"publish_iva"`
	PublishFva      int8   `orm:"column(publish_fva)" json:"publish_fva"`
	PaymentGuideURL string `orm:"column(payment_guide_url);" json:"payment_guide_url"`
	PaymentMethodID int64  `orm:"column(payment_method_id);null;" json:"payment_method"`
}

// TableName : set table name used by model
func (PaymentChannel) TableName() string {
	return "payment_channel"
}
