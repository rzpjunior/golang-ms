// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(PublicProduct1))
}

type PublicProduct1 struct {
	ID            int64     `orm:"column(id)" json:"id"`
	UOM           string    `orm:"column(uom)" json:"uom"`
	Name          string    `orm:"column(name)" json:"name"`
	ProductImages string    `orm:"column(product_images)" json:"product_images"`
	CreatedAt     time.Time `orm:"column(created_at);type(datetime)" json:"created_at"`
	LastUpdatedAt time.Time `orm:"column(last_updated_at);type(datetime)" json:"last_updated_at"`
}

// TableName : set table name used by model
func (PublicProduct1) TableName() string {
	return "public_product_1"
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PublicProduct1) MarshalJSON() ([]byte, error) {
	type Alias PublicProduct1

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PublicProduct1) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	o.Using("scrape")

	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PublicProduct1) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("scrape")

	return o.Read(m, fields...)
}
