// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"

	"git.edenfarm.id/project-version2/datamodel/util"
)

func init() {
	orm.RegisterModel(new(SalesPayment))
}

// Purchase Order: struct to hold model data for database
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

	TxnXendit      *TxnXendit      `orm:"column(txn_xendit_id);null;rel(fk)" json:"txn_xendit"`
	SalesInvoice   *SalesInvoice   `orm:"column(sales_invoice_id);null;rel(fk)" json:"sales_invoice"`
	PaymentMethod  *PaymentMethod  `orm:"column(payment_method_id);null;rel(fk)" json:"payment_method"`
	PaymentChannel *PaymentChannel `orm:"column(payment_channel_id);null;rel(fk)" json:"payment_channel"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *SalesPayment) MarshalJSON() ([]byte, error) {
	type Alias SalesPayment

	var recognitionDateStr string
	var receivedDateStr string

	timeISOLayout := time.RFC3339

	if m.RecognitionDate.Format("2006-01-02 15:04:05") == "0001-01-01 00:00:00" {
		recognitionDateStr = ""
	} else {
		recognitionDateStr = m.RecognitionDate.Format(timeISOLayout)
	}

	if m.ReceivedDate.Format("2006-01-02 15:04:05") == "0001-01-01 00:00:00" {
		receivedDateStr = ""
	} else {
		receivedDateStr = m.ReceivedDate.Format(timeISOLayout)
	}

	return json.Marshal(&struct {
		ID              string `json:"id"`
		StatusConvert   string `json:"status_convert"`
		RecognitionDate string `json:"recognition_date"`
		ReceivedDate    string `json:"received_date"`
		*Alias
	}{
		ID:              common.Encrypt(m.ID),
		StatusConvert:   util.ConvertStatusDoc(m.Status),
		RecognitionDate: recognitionDateStr,
		ReceivedDate:    receivedDateStr,
		Alias:           (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *SalesPayment) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *SalesPayment) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
