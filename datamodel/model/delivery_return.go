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
	orm.RegisterModel(new(DeliveryReturn))
}

// DeliveryReturn: struct to hold model data for database
type DeliveryReturn struct {
	ID              int64     `orm:"column(id);auto" json:"-"`
	Code            string    `orm:"column(code);size(50);null" json:"code"`
	RecognitionDate time.Time `orm:"column(recognition_date)" json:"recognition_date"`
	Note            string    `orm:"column(note)" json:"note"`
	Status          int8      `orm:"column(status)" json:"status"`
	CreatedAt       time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy       int64     `orm:"column(created_by)" json:"created_by"`
	ConfirmedAt     time.Time `orm:"column(confirmed_at);type(timestamp);null" json:"confirmed_at"`
	ConfirmedBy     int64     `orm:"column(confirmed_by)" json:"confirmed_by"`

	Warehouse     *Warehouse     `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse"`
	DeliveryOrder *DeliveryOrder `orm:"column(delivery_order_id);null;rel(fk)" json:"delivery_order"`

	DeliveryReturnItems []*DeliveryReturnItem `orm:"reverse(many)" json:"delivery_return_item,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *DeliveryReturn) MarshalJSON() ([]byte, error) {
	type Alias DeliveryReturn

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *DeliveryReturn) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *DeliveryReturn) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
