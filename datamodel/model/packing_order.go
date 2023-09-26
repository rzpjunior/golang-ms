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
	orm.RegisterModel(new(PackingOrder))
}

// PackingOrder model for packing order table.
type PackingOrder struct {
	ID            int64     `orm:"column(id);auto" json:"-"`
	Code          string    `orm:"column(code);size(50);null" json:"code"`
	Note          string    `orm:"column(note);size(250);null" json:"note"`
	DeliveryDate  time.Time `orm:"column(delivery_date)" json:"delivery_date"`
	Status        int8      `orm:"column(status);null" json:"status"`
	StatusConvert string    `orm:"-" json:"status_convert"`

	Warehouse             *Warehouse               `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse,omitempty"`
	Area                  *Area                    `orm:"column(area_id);null;rel(fk)" json:"area,omitempty"`
	PackingOrderItems     []*PackingOrderItem      `orm:"reverse(many)" json:"packing_order_items,omitempty"`
	PackingRecommendation []*PackingRecommendation `orm:"-" json:"packing_recommendation,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *PackingOrder) MarshalJSON() ([]byte, error) {
	type Alias PackingOrder

	alias := &struct {
		ID            string `json:"id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		StatusConvert: util.ConvertStatusMaster(m.Status),
		Alias:         (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PackingOrder) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PackingOrder) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
