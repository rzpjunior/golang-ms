// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import "time"

type AddressCoordinateLogResponse struct {
	ID             int64     `json:"id"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	LogChannelID   int8      `json:"log_channel"`
	MainCoordinate int8      `json:"main_coordinate"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      int64     `json:"created_by"`

	AddressID    string `json:"address_id"`
	SalesOrderID string `json:"sales_order_id"`
}

type AddressCoordinateLogGetRequest struct {
	OrderBy          string
	GroupBy          string
	ArrAddressIDs    []string
	ArrSalesOrderIDs []string
}
