package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type SalesInvoice struct {
	ID                int64     `orm:"column(id);auto" json:"-"`
	Code              string    `orm:"column(code);size(50);null" json:"code"`
	CodeExt           string    `orm:"column(code_ext);size(50);null" json:"code_ext"`
	Status            int8      `orm:"column(status)" json:"status"`
	RecognitionDate   time.Time `orm:"column(recognition_date)" json:"recognition_date"`
	DueDate           time.Time `orm:"column(due_date)" json:"due_date"`
	BillingAddress    string    `orm:"column(billing_address)" json:"billing_address"`
	DeliveryFee       float64   `orm:"column(delivery_fee)" json:"delivery_fee"`
	VouRedeemCode     string    `orm:"column(vou_redeem_code);null" json:"vou_redeem_code"`
	VouDiscAmount     float64   `orm:"column(vou_disc_amount);null" json:"vou_disc_amount"`
	PointRedeemAmount float64   `orm:"column(point_redeem_amount);null" json:"point_redeem_amount"`
	Adjustment        int8      `orm:"column(adjustment)" json:"adjustment"`
	AdjAmount         float64   `orm:"column(adj_amount)" json:"adj_amount"`
	AdjNote           string    `orm:"column(adj_note)" json:"adj_note"`
	TotalPrice        float64   `orm:"column(total_price)" json:"total_price"`
	TotalCharge       float64   `orm:"column(total_charge)" json:"total_charge"`
	DeltaPrint        int64     `orm:"column(delta_print)" json:"delta_print"`
	Note              string    `orm:"column(note)" json:"note"`
	VoucherID         int64     `orm:"column(voucher_id);null" json:"voucher_id"`
	RemainingAmount   float64   `orm:"-" json:"remaining_amount"`
}

func init() {
	orm.RegisterModel(new(SalesInvoice))
}

func (m *SalesInvoice) TableName() string {
	return "sales_invoice"
}
