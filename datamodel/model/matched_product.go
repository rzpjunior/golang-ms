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
	orm.RegisterModel(new(MatchedProduct))
}

type MatchedProduct struct {
	ID               int64             `orm:"column(id)" json:"id"`
	DashboardProduct *DashboardProduct `orm:"column(dashboard_product_id);null;rel(fk)" json:"dashboard_product"`
	PublicProduct1   *PublicProduct1   `orm:"column(public_product_1_id);null;rel(fk)" json:"public_product_1"`
	PublicProduct2   *PublicProduct2   `orm:"column(public_product_2_id);null;rel(fk)" json:"public_product_2"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *MatchedProduct) MarshalJSON() ([]byte, error) {
	type Alias MatchedProduct

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *MatchedProduct) Save(fields ...string) (err error) {
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
func (m *MatchedProduct) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("scrape")

	return o.Read(m, fields...)
}
