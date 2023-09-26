// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

type DeliveryRunReturnItemResponse struct {
	ID             int64   `json:"id"`
	ReceiveQty     float64 `json:"receive_qty"`
	ReturnReason   int8    `json:"return_reason"`
	ReturnEvidence string  `json:"return_evidence"`
	Subtotal       float64 `json:"subtotal"`

	DeliveryRunReturnID int64  `json:"delivery_run_return_id"`
	DeliveryOrderItemID string `json:"delivery_order_item_id"`
}

type DeliveryRunReturnItemGetRequest struct {
	Offset                  int
	Limit                   int
	OrderBy                 string
	ArrDeliveryRunReturnIDs []int64
	ArrDeliveryOrderItemIDs []string
}
