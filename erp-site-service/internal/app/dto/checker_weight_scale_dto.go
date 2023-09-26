// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

type CheckerWeightScaleGetRequest struct {
	PickingOrderItemId int64
}

type CheckerWeightScaleGetResponse struct {
	PickingOrderItemId int64   `json:"picking_order_item_id"`
	ProductPicture     string  `json:"product_picture"`
	ProductName        string  `json:"product_name"`
	ProductId          string  `json:"product_id"`
	OrderQty           float64 `json:"order_qty"`
	OrderMinQty        float64 `json:"order_min_qty"`
}

type CheckerWeightScaleUpdateRequest struct {
	PickingOrderItemId int64   `json:"-"`
	CheckQty           float64 `json:"check_qty"`
}

// response
type CheckerSuccessResponse struct {
	Success bool `json:"success"`
}
