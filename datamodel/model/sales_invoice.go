// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(SalesInvoice))
}

// Sales Invoice: struct to hold model sales_invoice for database
type SalesInvoice struct {
	ID                 int64     `orm:"column(id);auto" json:"-"`
	Code               string    `orm:"column(code);size(50);null" json:"code"`
	CodeExt            string    `orm:"column(code_ext);size(50);null" json:"code_ext"`
	Status             int8      `orm:"column(status)" json:"status"`
	RecognitionDate    time.Time `orm:"column(recognition_date)" json:"recognition_date"`
	DueDate            time.Time `orm:"column(due_date)" json:"due_date"`
	BillingAddress     string    `orm:"column(billing_address)" json:"billing_address"`
	DeliveryFee        float64   `orm:"column(delivery_fee)" json:"delivery_fee"`
	VouRedeemCode      string    `orm:"column(vou_redeem_code);null" json:"vou_redeem_code"`
	VouDiscAmount      float64   `orm:"column(vou_disc_amount);null" json:"vou_disc_amount"`
	Adjustment         int8      `orm:"column(adjustment)" json:"adjustment"`
	AdjAmount          float64   `orm:"column(adj_amount)" json:"adj_amount"`
	AdjNote            string    `orm:"column(adj_note)" json:"adj_note"`
	TotalPrice         float64   `orm:"column(total_price)" json:"total_price"`
	TotalCharge        float64   `orm:"column(total_charge)" json:"total_charge"`
	DeltaPrint         int64     `orm:"column(delta_print)" json:"delta_print"`
	Note               string    `orm:"column(note)" json:"note"`
	VoucherID          int64     `orm:"column(voucher_id);null" json:"voucher_id"`
	RemainingAmount    float64   `orm:"-" json:"remaining_amount"`
	TotalPaid          float64   `orm:"-" json:"total_paid"`
	VoucherType        int8      `orm:"-" json:"voucher_type"`
	CreatedAt          time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy          int64     `orm:"column(created_by)" json:"created_by"`
	LastUpdatedAt      time.Time `orm:"column(last_updated_at);type(timestamp);null" json:"last_updated_at"`
	LastUpdatedBy      int64     `orm:"column(last_updated_by)" json:"last_updated_by"`
	PointRedeemAmount  float64   `orm:"column(point_redeem_amount)" json:"point_redeem_amount"`
	TotalSkuDiscAmount float64   `orm:"column(total_sku_disc_amount)" json:"total_sku_disc_amount"`

	SalesOrder   *SalesOrder   `orm:"column(sales_order_id);null;rel(fk)" json:"sales_order"`
	PaymentGroup *PaymentGroup `orm:"column(payment_group_sls_id);null;rel(fk)" json:"payment_group_sls"`
	SalesTerm    *SalesTerm    `orm:"column(term_payment_sls_id);null;rel(fk)" json:"term_payment_sls"`
	InvoiceTerm  *InvoiceTerm  `orm:"column(term_invoice_sls_id);null;rel(fk)" json:"term_invoice_sls"`
	//Voucher      	*Voucher      `orm:"column(voucher_id);null;rel(fk)" json:"voucher"`
	SalesInvoiceItems []*SalesInvoiceItem `orm:"reverse(many)" json:"sales_invoice_items,omitempty"`
	SalesPayment      []*SalesPayment     `orm:"reverse(many)" json:"sales_payment,omitempty"`

	MerchantAccNum    []*MerchantAccNum `orm:"-" json:"merchant_acc_num"`
	XenditBCA         string            `orm:"-" json:"xendit_bca"`
	XenditPermata     string            `orm:"-" json:"xendit_permata"`
	DeliveryKoli      []*DeliveryKoli   `orm:"-" json:"delivery_koli"`
	TotalKoli         float64           `orm:"-" json:"total_koli"`
	PaymentPercentage float64           `orm:"-" json:"payment_percentage"`
	StatusDescription string            `orm:"-" json:"status_description"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *SalesInvoice) MarshalJSON() ([]byte, error) {
	type Alias SalesInvoice

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *SalesInvoice) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *SalesInvoice) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
