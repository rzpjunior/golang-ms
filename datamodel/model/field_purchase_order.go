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
	orm.RegisterModel(new(FieldPurchaseOrder))
}

// FieldPurchaseOrder model for field_purchase_order table.
type FieldPurchaseOrder struct {
	ID            int64          `orm:"column(id);auto" json:"-"`
	Code          string         `orm:"column(code);size(50);null" json:"code,omitempty"`
	PurchaseOrder int64          `orm:"column(purchase_order_id)" json:"-"`
	Stall         *Stall         `orm:"column(stall_id);null;rel(fk)" json:"stall,omitempty"`
	TotalPrice    float64        `orm:"column(total_price)" json:"total_price"`
	TotalItem     int8           `orm:"column(total_item)" json:"total_item"`
	PaymentMethod *PaymentMethod `orm:"column(payment_method_id);null;rel(fk)" json:"payment_method,omitempty"`
	Latitude      float64        `orm:"column(latitude)" json:"latitude"`
	Longitude     float64        `orm:"column(longitude)" json:"longitude"`
	CreatedAt     time.Time      `orm:"column(created_at)" json:"created_at"`
	CreatedBy     int64          `orm:"column(created_by)" json:"-"`

	FieldPurchaseOrderItems []*FieldPurchaseOrderItem `orm:"reverse(many)" json:"field_purchase_order_items,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *FieldPurchaseOrder) MarshalJSON() ([]byte, error) {
	type Alias FieldPurchaseOrder

	alias := &struct {
		ID            string `json:"id"`
		PurchaseOrder string `json:"purchase_order"`
		CreatedBy     string `json:"created_by"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		PurchaseOrder: common.Encrypt(m.PurchaseOrder),
		CreatedBy:     common.Encrypt(m.CreatedBy),
		Alias:         (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *FieldPurchaseOrder) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *FieldPurchaseOrder) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
