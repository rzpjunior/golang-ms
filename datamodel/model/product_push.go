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
	orm.RegisterModel(new(ProductPush))
}

// ProductPush : struct to hold model data for database
type ProductPush struct {
	ID        int64     `orm:"column(id);auto" json:"-"`
	StartDate time.Time `orm:"column(start_date)" json:"start_date"`
	Status    int8      `orm:"column(status)" json:"status"`

	Product   *Product   `orm:"column(product_id);null;rel(fk)" json:"product"`
	Area      *Area      `orm:"column(area_id);null;rel(fk)" json:"area"`
	Archetype *Archetype `orm:"column(archetype_id);null;rel(fk)" json:"archetype"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *ProductPush) MarshalJSON() ([]byte, error) {
	type Alias ProductPush

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *ProductPush) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *ProductPush) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
