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
	orm.RegisterModel(new(PackingHelperLog))
}

// PackingOrder model for packing order table.
type PackingHelperLog struct {
	ID        int64     `orm:"column(id);auto" json:"-"`
	QtyWeight float64   `orm:"column(qty_weight);null" json:"qty_weight"`
	QtyPack   float64   `orm:"column(qty_pack);null" json:"qty_pack"`
	CreatedAt time.Time `orm:"column(created_at);size(30);null" json:"created_at"`

	PackingOrderItem *PackingOrderItem `orm:"column(packing_order_item_id);null;rel(fk)" json:"packing_order_item_id,omitempty"`
	Helper           *Staff            `orm:"column(staff_id);null;rel(fk)" json:"staff_id,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *PackingHelperLog) MarshalJSON() ([]byte, error) {
	type Alias PackingHelperLog

	alias := &struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PackingHelperLog) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PackingHelperLog) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
