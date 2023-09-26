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
	orm.RegisterModel(new(SkuDiscount))
}

// SkuDiscount : struct to hold model data for database
type SkuDiscount struct {
	ID                int64     `orm:"column(id);auto" json:"-"`
	Code              string    `orm:"column(code)" json:"code"`
	Name              string    `orm:"column(name)" json:"name"`
	Note              string    `orm:"column(note)" json:"note"`
	Status            int8      `orm:"column(status)" json:"status"`
	StartTimestamp    time.Time `orm:"column(start_timestamp)" json:"start_timestamp"`
	EndTimestamp      time.Time `orm:"column(end_timestamp)" json:"end_timestamp"`
	OrderChannels     string    `orm:"column(order_channel)" json:"order_channel"`
	OrderChannelsName string    `orm:"-" json:"order_channel_name"`
	PriceSets         string    `orm:"column(price_set)" json:"price_set"`
	PriceSetsName     string    `orm:"-" json:"price_set_name"`

	Division         *Division          `orm:"column(division_id);null;rel(fk)" json:"division"`
	SkuDiscountItems []*SkuDiscountItem `orm:"reverse(many)" json:"sku_discount_item,omitempty"`

	// data log
	CreatedAt  time.Time `orm:"column(created_at)" json:"created_at"`
	CreatedBy  int64     `orm:"column(created_by)" json:"created_by"`
	ArchivedAt time.Time `orm:"column(archived_at)" json:"archived_at"`
	ArchivedBy int64     `orm:"column(archived_by)" json:"archived_by"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *SkuDiscount) MarshalJSON() ([]byte, error) {
	type Alias SkuDiscount

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
		StatusConvert string `json:"status_convert"`
	}{
		ID:            common.Encrypt(m.ID),
		Alias:         (*Alias)(m),
		StatusConvert: util.ConvertStatusDoc(m.Status),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *SkuDiscount) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *SkuDiscount) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
