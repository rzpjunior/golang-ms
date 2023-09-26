// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

// DeliveryRunReturn model for delivery_run_return table.
type DeliveryRunReturn struct {
	ID          int64     `orm:"column(id)" json:"id"`
	Code        string    `orm:"column(code)" json:"code"`
	TotalPrice  float64   `orm:"column(total_price)" json:"total_price,omitempty"`
	TotalCharge float64   `orm:"column(total_charge)" json:"total_charge,omitempty"`
	CreatedAt   time.Time `orm:"column(created_at)" json:"created_at"`

	DeliveryRunSheetItemID int64 `orm:"column(delivery_run_sheet_item_id)" json:"delivery_run_sheet_item_id"`
}

func init() {
	orm.RegisterModel(new(DeliveryRunReturn))
}

func (m *DeliveryRunReturn) TableName() string {
	return "delivery_run_return"
}
