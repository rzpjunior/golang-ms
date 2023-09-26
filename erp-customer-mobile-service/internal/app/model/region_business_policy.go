// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

// RegionBusinessPolicy : struct to hold model data for database
type RegionBusinessPolicy struct {
	ID            int64     `orm:"column(id);auto" json:"-"`
	MinOrder      float64   `orm:"column(min_order);null;digits(20);decimals(2)" json:"min_order"`
	DeliveryFee   float64   `orm:"column(delivery_fee);null;digits(20);decimals(2)" json:"delivery_fee"`
	CreatedAt     time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy     string    `orm:"column(created_by)" json:"created_by"`
	LastUpdatedAt time.Time `orm:"column(last_updated_at);type(timestamp);null" json:"last_updated_at"`
	LastUpdatedBy string    `orm:"column(last_updated_by)" json:"last_updated_by"`
}

func init() {
	orm.RegisterModel(new(RegionBusinessPolicy))
}

func (m *RegionBusinessPolicy) TableName() string {
	return "region_business_policy"
}
