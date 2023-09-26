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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	orm.RegisterModel(new(PickingOrder))
}

// PickingOrder model for picking order table.
type PickingOrder struct {
	ID              int64              `orm:"column(id);auto" json:"-"`
	JobsID          primitive.ObjectID `orm:"-" json:"jobs_id"`
	Code            string             `orm:"column(code);size(50);null" json:"code"`
	RecognitionDate time.Time          `orm:"column(recognition_date)" json:"recognition_date"`
	Status          int8               `orm:"column(status);null" json:"status"`
	Note            string             `orm:"column(note)" json:"note"`
	StatusConvert   string             `orm:"-" json:"status_convert"`

	Warehouse          *Warehouse            `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse,omitempty"`
	PickingOrderAssign []*PickingOrderAssign `orm:"reverse(many)" json:"picking_order_assign,omitempty"`
	Jobs               *Jobs                 `orm:"-" json:"jobs,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *PickingOrder) MarshalJSON() ([]byte, error) {
	type Alias PickingOrder

	alias := &struct {
		ID            string `json:"id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		StatusConvert: util.ConvertStatusPicking(m.Status),
		Alias:         (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PickingOrder) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PickingOrder) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
