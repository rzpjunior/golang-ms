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
	orm.RegisterModel(new(SkuDiscountItemTier))
}

// SkuDiscountItemTier : struct to hold model data for database
type SkuDiscountItemTier struct {
	ID         int64   `orm:"column(id);auto" json:"-"`
	TierLevel  int8    `orm:"column(tier_level)" json:"tier_level"`
	MinimumQty float64 `orm:"column(minimum_qty)" json:"minimum_qty"`
	DiscAmount float64 `orm:"column(disc_amount)" json:"disc_amount"`

	SkuDiscountItem *SkuDiscountItem `orm:"column(sku_discount_item_id);null;rel(fk)" json:"sku_discount_item"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *SkuDiscountItemTier) MarshalJSON() ([]byte, error) {
	type Alias SkuDiscountItemTier

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *SkuDiscountItemTier) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *SkuDiscountItemTier) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
