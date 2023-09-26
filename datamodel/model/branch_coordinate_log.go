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
	orm.RegisterModel(new(BranchCoordinateLog))
}

type BranchCoordinateLog struct {
	ID             int64       `orm:"column(id);auto" json:"-"`
	Branch         *Branch     `orm:"column(branch_id);null;rel(fk)" json:"branch"`
	SalesOrder     *SalesOrder `orm:"column(sales_order_id);null;rel(fk)" json:"sales_order"`
	Latitude       float64     `orm:"column(latitude)" json:"latitude,omitempty"`
	Longitude      float64     `orm:"column(longitude)" json:"longitude,omitempty"`
	LogChannelID   int8        `orm:"column(log_channel_id);null" json:"log_channel"`
	MainCoordinate int8        `orm:"column(main_coordinate)" json:"main_coordinate"`
	CreatedAt      time.Time   `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy      int64       `orm:"column(created_by)" json:"created_by"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *BranchCoordinateLog) MarshalJSON() ([]byte, error) {
	type Alias BranchCoordinateLog

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *BranchCoordinateLog) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *BranchCoordinateLog) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
