// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

// DeliveryRunReturnItem model for delivery_run_return_item table.
type DeliveryRunReturnItem struct {
	ID             int64   `orm:"column(id)" json:"id"`
	ReceiveQty     float64 `orm:"column(receive_qty)" json:"receive_qty"`
	ReturnReason   int8    `orm:"column(return_reason)" json:"return_reason"`
	ReturnEvidence string  `orm:"column(return_evidence)" json:"return_evidence"`
	Subtotal       float64 `orm:"column(subtotal)" json:"subtotal"`

	DeliveryRunReturnID int64  `orm:"column(delivery_run_return_id)" json:"delivery_run_return_id"`
	DeliveryOrderItemID string `orm:"column(delivery_order_item_id)" json:"delivery_order_item_id"`
}

func init() {
	orm.RegisterModel(new(DeliveryRunReturnItem))
}

func (m *DeliveryRunReturnItem) TableName() string {
	return "delivery_run_return_item"
}
