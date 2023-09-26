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
)

func init() {
	orm.RegisterModel(new(SkuDiscountLog))
}

// SkuDiscountLog model
type SkuDiscountLog struct {
	ID             int64     `orm:"column(id);auto" json:"-"`
	DiscountAmount float64   `orm:"column(discount_amount)" json:"discount_amount"`
	DiscountQty    float64   `orm:"column(discount_qty)" json:"discount_qty"`
	CreatedAt      time.Time `orm:"column(created_at)" json:"created_at"`
	Status         int8      `orm:"column(status)" json:"status"`

	SkuDiscount     *SkuDiscount     `orm:"column(sku_discount_id);null;rel(fk)" json:"sku_discount"`
	SkuDiscountItem *SkuDiscountItem `orm:"column(sku_discount_item_id);null;rel(fk)" json:"sku_discount_item"`
	Merchant        *Merchant        `orm:"column(merchant_id);null;rel(fk)" json:"merchant"`
	Branch          *Branch          `orm:"column(branch_id);null;rel(fk)" json:"branch"`
	SalesOrderItem  *SalesOrderItem  `orm:"column(sales_order_item_id);null;rel(fk)" json:"sales_order_item"`
	Product         *Product         `orm:"column(product_id);null;rel(fk)" json:"product"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *SkuDiscountLog) MarshalJSON() ([]byte, error) {
	type Alias SkuDiscountLog

	alias := &struct {
		ID            string `json:"id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		StatusConvert: util.ConvertStatusMaster(m.Status),
		Alias:         (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating User struct into user table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *SkuDiscountLog) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting user data
// this also will truncated all data from all table
// that have relation with this user.
func (m *SkuDiscountLog) Delete() (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		var i int64
		if i, err = o.Delete(m); i == 0 && err == nil {
			err = orm.ErrNoAffected
		}
		return
	}
	return orm.ErrNoRows
}

// Read execute select based on data struct that already
// assigned.
func (m *SkuDiscountLog) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
