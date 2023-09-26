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
	orm.RegisterModel(new(PurchaseOrderImage))
}

// PurchaseOrderImage: struct to hold image data for purchase_order
type PurchaseOrderImage struct {
	ID            int64          `orm:"column(id);auto" json:"-"`
	PurchaseOrder *PurchaseOrder `orm:"column(purchase_order_id);null;rel(fk)" json:"purchase_order,omitempty"`
	ImageURL      string         `orm:"column(image_url)" json:"image_url"`
	CreatedAt     time.Time      `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy     int64          `orm:"column(created_by)" json:"-"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PurchaseOrderImage) MarshalJSON() ([]byte, error) {
	type Alias PurchaseOrderImage

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
func (m *PurchaseOrderImage) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PurchaseOrderImage) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
