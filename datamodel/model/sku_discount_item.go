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
	orm.RegisterModel(new(SkuDiscountItem))
}

// SkuDiscountItem : struct to hold model data for database
type SkuDiscountItem struct {
	ID                   int64   `orm:"column(id);auto" json:"-"`
	OverallQuota         int64   `orm:"column(overall_quota)" json:"overall_quota"`
	OverallQuotaPerUser  int64   `orm:"column(overall_quota_per_user)" json:"overall_quota_per_user"`
	DailyQuotaPerUser    int64   `orm:"column(daily_quota_per_user)" json:"daily_quota_per_user"`
	RemOverallQuota      float64 `orm:"column(rem_overall_quota)" json:"rem_overall_quota"`
	Budget               float64 `orm:"column(budget);digits(11);decimals(2)" json:"budget"`
	RemBudget            float64 `orm:"column(rem_budget);digits(11);decimals(2)" json:"rem_budget"`
	RemDailyQuotaPerUser int64   `orm:"-" json:"rem_daily_quota_per_user"`
	RemQuotaPerUser      int64   `orm:"-" json:"rem_quota_per_user"`
	IsUseBudget          int8    `orm:"column(use_budget)" json:"is_use_budget"`

	SkuDiscount          *SkuDiscount           `orm:"column(sku_discount_id);null;rel(fk)" json:"sku_discount"`
	Product              *Product               `orm:"column(product_id);null;rel(fk)" json:"product"`
	SkuDiscountItemTiers []*SkuDiscountItemTier `orm:"reverse(many)" json:"sku_discount_item_tier,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *SkuDiscountItem) MarshalJSON() ([]byte, error) {
	type Alias SkuDiscountItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *SkuDiscountItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *SkuDiscountItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
