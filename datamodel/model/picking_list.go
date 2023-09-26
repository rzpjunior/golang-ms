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
	orm.RegisterModel(new(PickingList))
}

// PickingList model for picking order table.
type PickingList struct {
	ID                     int64     `orm:"column(id);auto" json:"-"`
	Code                   string    `orm:"column(code);size(50);null" json:"code"`
	DeliveryDate           time.Time `orm:"column(delivery_date)" json:"delivery_date"`
	Status                 int8      `orm:"column(status);null" json:"status"`
	Note                   string    `orm:"column(note)" json:"note"`
	RoutingNote            string    `orm:"column(routing_note)" json:"routing_note"`
	StatusConvert          string    `orm:"-" json:"status_convert"`
	TotalWeightPickingList float64   `orm:"-" json:"total_weight_picking_list"`
	TotalItemPickingList   int8      `orm:"-" json:"total_item_picking_list"`
	TotalSalesOrder        int       `orm:"-" json:"total_sales_order"`
	PickerName             string    `orm:"-" json:"picker_name"`

	PickingRouting int8          `json:"picking_routing"`
	Pickers        []*Staff      `orm:"-" json:"staff,omitempty"`
	Warehouse      *Warehouse    `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse,omitempty"`
	SalesOrder     []*SalesOrder `orm:"-" json:"sales_order"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *PickingList) MarshalJSON() ([]byte, error) {
	type Alias PickingList

	alias := &struct {
		ID            string `json:"id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		StatusConvert: util.ConvertStatusPickingList(m.Status),
		Alias:         (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PickingList) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PickingList) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
