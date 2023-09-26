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
	orm.RegisterModel(new(AreaPolicy))
}

// AreaPolicy : struct to hold model data for database
type AreaPolicy struct {
	ID                  int64     `orm:"column(id);auto" json:"-"`
	MinOrder            float64   `orm:"column(min_order);null;digits(20);decimals(2)" json:"min_order"`
	DeliveryFee         float64   `orm:"column(delivery_fee);null;digits(20);decimals(2)" json:"delivery_fee"`
	OrderTimeLimit      string    `orm:"column(order_time_limit)" json:"order_time_limit"`
	DraftOrderTimeLimit string    `orm:"column(draft_order_time_limit)" json:"draft_order_time_limit"`
	DefaultPriceSet     *PriceSet `orm:"column(default_price_set);rel(fk)" json:"default_price_set"`
	Area                *Area     `orm:"column(area_id);rel(fk)" json:"area,omitempty"`
	MaxDayDeliveryDate  int       `orm:"column(max_day_delivery_date)" json:"max_day_delivery_date,omitempty"`
	WeeklyDayOff        int       `orm:"column(weekly_day_off)" json:"weekly_day_off,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *AreaPolicy) MarshalJSON() ([]byte, error) {
	type Alias AreaPolicy

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *AreaPolicy) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *AreaPolicy) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
