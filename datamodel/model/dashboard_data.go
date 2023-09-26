// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

type DataDashboard struct {
	ProductID   int64   `orm:"column(product_id)" json:"product_id"`
	ProductCode string  `orm:"column(product_code)" json:"product_code"`
	ProductName string  `orm:"column(product_name)" json:"product_name"`
	UOM         string  `orm:"column(uom)" json:"uom"`
	Area        string  `orm:"column(area)" json:"area"`
	Price       float64 `orm:"column(price)" json:"price"`
}
