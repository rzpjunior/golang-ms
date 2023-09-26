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
	orm.RegisterModel(new(MatchedArea))
}

type MatchedArea struct {
	ID              int64          `orm:"column(id)" json:"id"`
	DashboardArea   *DashboardArea `orm:"column(dashboard_area_id);null;rel(fk)" json:"dashboard_area"`
	PublicDataArea1 *PublicArea1   `orm:"column(public_data_area_1_id);null;rel(fk)" json:"public_data_area_1"`
	PublicDataArea2 *PublicArea2   `orm:"column(public_data_area_2_id);null;rel(fk)" json:"public_data_area_2"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *MatchedArea) MarshalJSON() ([]byte, error) {
	type Alias MatchedArea

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *MatchedArea) Save(fields ...string) (err error) {
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
func (m *MatchedArea) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("scrape")

	return o.Read(m, fields...)
}
