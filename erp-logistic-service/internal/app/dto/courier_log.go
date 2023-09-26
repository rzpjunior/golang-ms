// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import "time"

type CourierLogResponse struct {
	ID        int64     `json:"id"`
	Latitude  *float64  `json:"latitude"`
	Longitude *float64  `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`

	CourierID    string `json:"courier_id"`
	SalesOrderID string `json:"sales_order_id"`
}

type CourierLog struct {
	ID           int64     `json:"id"`
	Latitude     *float64  `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude    *float64  `json:"longtude,omitempty" bson:"longtude,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	CourierID    string    `json:"courier_id,omitempty" bson:"courier_id,omitempty"`
	SalesOrderID string    `json:"sales_order_id,omitempty" bson:"sales_order_id,omitempty"`
}
