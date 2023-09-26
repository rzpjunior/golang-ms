// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import "time"

type MerchantDeliveryLogResponse struct {
	Id        int64     `json:"id"`
	Latitude  *float64  `json:"latitude"`
	Longitude *float64  `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`

	DeliveryRunSheetItemId int64 `json:"delivery_run_sheet_item_id"`
}
