package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type SalesPayment struct {
	ID               int64     `orm:"column(id);auto" json:"-"`
	Code             string    `orm:"column(code);size(50);null" json:"code"`
	Status           int8      `orm:"column(status)" json:"status"`
	RecognitionDate  time.Time `orm:"column(recognition_date)" json:"recognition_date"`
	Amount           float64   `orm:"column(amount)" json:"amount"`
	BankReceiveNum   string    `orm:"column(bank_receive_num)" json:"bank_receive_num"`
	PaidOff          int8      `orm:"column(paid_off)" json:"paid_off"`
	ImageUrl         string    `orm:"column(image_url);null" json:"image_url"`
	Note             string    `orm:"column(note)" json:"note"`
	CreatedAt        time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy        int64     `orm:"column(created_by)" json:"created_by"`
	CancellationNote string    `orm:"-" json:"cancellation_note,omitempty"`
	ReceivedDate     time.Time `orm:"column(received_date);type(date);null" json:"received_date"`
}

func init() {
	orm.RegisterModel(new(SalesPayment))
}

func (m *SalesPayment) TableName() string {
	return "sales_payment"
}
