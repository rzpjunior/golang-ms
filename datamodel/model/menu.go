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
	orm.RegisterModel(new(Menu))
}

// Menu : struct to hold model data for database
type Menu struct {
	ID     int64  `orm:"column(id);auto" json:"-"`
	Title  string `orm:"column(title)" json:"title"`
	Url    string `orm:"column(url)" json:"url"`
	Icon   string `orm:"column(icon)" json:"icon"`
	Status int8   `orm:"column(status)" json:"status"`
	Order  int32  `orm:"column(order)" json:"order"`

	Parent    *Menu       `orm:"column(parent_id);null;rel(fk)" json:"parent,omitempty"`
	Privilege *Permission `orm:"column(permission_id);null;rel(fk)" json:"privilege,omitemtpy"`
	Child     []*Menu     `orm:"reverse(many)" json:"child,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *Menu) MarshalJSON() ([]byte, error) {
	type Alias Menu

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *Menu) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *Menu) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
