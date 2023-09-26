// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

func init() {
	orm.RegisterModel(new(SalesOrderItem))
}
func (m *SalesOrderItem) TableName() string {
	return "sales_order_item"
}

// Sales Order Item: struct to hold model data for database
type SalesOrderItem struct {
	ID               int64   `orm:"column(id);auto" json:"-"`
	SalesOrderID     int64   `orm:"column(sales_order_id)" json:"sales_order_id"`
	ItemIDGP         string  `orm:"column(item_id_gp)" json:"item_id_gp"`
	PriceTieringIDGP string  `orm:"column(price_tiering_id_gp);null" json:"price_tiering_id_gp"`
	OrderQty         float64 `orm:"column(order_qty)" json:"order_qty"`
	UnitPrice        float64 `orm:"column(unit_price)" json:"unit_price"`
	UomIDGP          string  `orm:"column(uom_gp)" json:"uom_gp"`
	Subtotal         float64 `orm:"column(subtotal)" json:"subtotal"`
	Weight           float64 `orm:"column(weight)" json:"weight"`
}
