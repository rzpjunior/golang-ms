// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import "time"

type DeliveryRunSheetResponse struct {
	ID                int64     `json:"id"`
	Code              string    `json:"code"`
	DeliveryDate      time.Time `json:"delivery_date"`
	StartedAt         time.Time `json:"started_at"`
	FinishedAt        time.Time `json:"finished_at"`
	StartingLatitude  *float64  `json:"starting_latitude"`
	StartingLongitude *float64  `json:"starting_longitude"`
	FinishedLatitude  *float64  `json:"finished_latitude"`
	FinishedLongitude *float64  `json:"finished_longitude"`
	Status            int8      `json:"status"`

	CourierID string `json:"courier_id"`
}

type DeliveryRunSheetGetRequest struct {
	Offset        int
	Limit         int
	OrderBy       string
	GroupBy       string
	Status        []int
	ArrCourierIDs []string
}
