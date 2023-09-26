// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

// AddressCoordinateLog model for address_coordinate_log table.
type AddressCoordinateLog struct {
	ID             int64     `orm:"column(id)" json:"id"`
	Latitude       float64   `orm:"column(latitude)" json:"latitude"`
	Longitude      float64   `orm:"column(longitude)" json:"longitude"`
	LogChannelID   int8      `orm:"column(log_channel_id)" json:"log_channel"`
	MainCoordinate int8      `orm:"column(main_coordinate)" json:"main_coordinate"`
	CreatedAt      time.Time `orm:"column(created_at)" json:"created_at"`
	CreatedBy      int64     `orm:"column(created_by)" json:"created_by"`

	AddressID    string `orm:"column(address_id)" json:"address_id"`
	SalesOrderID string `orm:"column(sales_order_id)" json:"sales_order_id"`
}

func init() {
	orm.RegisterModel(new(AddressCoordinateLog))
}

func (m *AddressCoordinateLog) TableName() string {
	return "address_coordinate_log"
}
