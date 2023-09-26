// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

// MerchantDeliveryLog model for merchant_delivery_log table IN MONGO.
// this table hold the location of courier's delivery to the associated merchant
type MerchantDeliveryLog struct {
	Id        int64    `json:"id" bson:"id"`
	Latitude  *float64 `json:"latitude" bson:"latitude"`
	Longitude *float64 `json:"longitude" bson:"longitude"`
	CreatedAt int64    `json:"created_at" bson:"created_at"`

	DeliveryRunSheetItemId int64 `json:"delivery_run_sheet_item_id,omitempty" bson:"delivery_run_sheet_item_id"`
}
