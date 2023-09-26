// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import (
	"time"
)

type DeliveryRunReturnResponse struct {
	ID          int64     `json:"id"`
	Code        string    `json:"code"`
	TotalPrice  float64   `json:"total_price"`
	TotalCharge float64   `json:"total_charge"`
	CreatedAt   time.Time `json:"created_at"`

	DeliveryRunSheetItemID int64 `json:"delivery_run_sheet_item_id"`
}

type DeliveryRunReturnGetRequest struct {
	Offset                     int
	Limit                      int
	OrderBy                    string
	ArrDeliveryRunSheetItemIDs []int64
}
