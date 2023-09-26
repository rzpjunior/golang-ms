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
	orm.RegisterModel(new(ProductSectionItem))
}

// ProductSectionItem : struct to hold model data for database
type ProductSectionItem struct {
	ID   int64  `orm:"column(id);auto" json:"-"`
	Code string `orm:"column(code)" json:"code"`
	Name string `orm:"column(name)" json:"name"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *ProductSectionItem) MarshalJSON() ([]byte, error) {
	type Alias ProductSectionItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}
