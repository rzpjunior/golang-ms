// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(PaymentGroupComb))
}

// PaymentGroupComb : struct to hold payment group comb model data for database
type PaymentGroupComb struct {
	ID           int64         `orm:"column(id);auto" json:"-"`
	PaymentGroup *PaymentGroup `orm:"column(payment_group_sls_id);null;rel(fk)" json:"payment_group"`
	PaymentTerm  *SalesTerm    `orm:"column(term_payment_sls_id);null;rel(fk)" json:"payment_term"`
	InvoiceTerm  *InvoiceTerm  `orm:"column(term_invoice_sls_id);null;rel(fk)" json:"invoice_term"`
}

// TableName : set table name used by model
func (PaymentGroupComb) TableName() string {
	return "payment_group_comb"
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PaymentGroupComb) MarshalJSON() ([]byte, error) {
	type Alias PaymentGroupComb

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Read : function to get data from database
func (m *PaymentGroupComb) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
