// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

// DeliveryRunSheet model for delivery_run_sheet table.
type DeliveryRunSheet struct {
	ID                int64     `orm:"column(id)" json:"id"`
	Code              string    `orm:"column(code)" json:"code"`
	DeliveryDate      time.Time `orm:"column(delivery_date)" json:"delivery_date"`
	StartedAt         time.Time `orm:"column(started_at);type(timestamp)" json:"started_at"`
	FinishedAt        time.Time `orm:"column(finished_at);type(timestamp)" json:"finished_at"`
	StartingLatitude  *float64  `orm:"column(starting_latitude)" json:"starting_latitude"`
	StartingLongitude *float64  `orm:"column(starting_longitude)" json:"starting_longitude"`
	FinishedLatitude  *float64  `orm:"column(finished_latitude)" json:"finished_latitude"`
	FinishedLongitude *float64  `orm:"column(finished_longitude)" json:"finished_longitude"`
	Status            int8      `orm:"column(status)" json:"status"`

	CourierID string `orm:"column(courier_id)" json:"courier_id"`
}

func init() {
	orm.RegisterModel(new(DeliveryRunSheet))
}

func (m *DeliveryRunSheet) TableName() string {
	return "delivery_run_sheet"
}
