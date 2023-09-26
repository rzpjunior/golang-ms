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
	orm.RegisterModel(new(PurchaseDeliver))
}

// PurchaseDeliver: struct to hold purchase_deliver model data for database
type PurchaseDeliver struct {
	ID                          int64                        `orm:"column(id);auto" json:"-"`
	Code                        string                       `orm:"column(code);size(50);null" json:"code"`
	PurchaseOrder               *PurchaseOrder               `orm:"column(purchase_order_id);null;rel(fk)" json:"purchase_order,omitempty"`
	FieldPurchaseOrder          *FieldPurchaseOrder          `orm:"column(field_purchase_order_id);null;rel(fk)" json:"field_purchase_order,omitempty"`
	ConsolidatedPurchaseDeliver *ConsolidatedPurchaseDeliver `orm:"column(consolidated_purchase_deliver_id);null;rel(fk)" json:"consolidated_purchase_deliver,omitempty"`
	Stall                       *Stall                       `orm:"column(stall_id);null;rel(fk)" json:"stall,omitempty"`
	DeltaPrint                  int8                         `orm:"column(delta_print)" json:"delta_print"`
	CreatedAt                   time.Time                    `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy                   int64                        `orm:"column(created_by)" json:"-"`

	PurchaseDeliverSignature []*PurchaseDeliverSignature `orm:"reverse(many)" json:"signature,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PurchaseDeliver) MarshalJSON() ([]byte, error) {
	type Alias PurchaseDeliver

	return json.Marshal(&struct {
		ID        string `json:"id"`
		CreatedBy string `json:"created_by"`
		*Alias
	}{
		ID:        common.Encrypt(m.ID),
		CreatedBy: common.Encrypt(m.CreatedBy),
		Alias:     (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PurchaseDeliver) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PurchaseDeliver) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
