// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import "time"

// Fulfillment model
type Fulfillment struct {
	TotalSo                 float64   `json:"total_so"`
	TotalSoUnfulfilled      float64   `json:"total_so_unfulfilled"`
	FulfillmentRate         float64   `json:"fulfillment_rate"`
	UnfulfillmentRate       float64   `json:"unfulfillment_rate"`
	TotalCust               float64   `json:"total_cust"`
	TotalCustUnfulfilled    float64   `json:"total_cust_unfulfilled"`
	CustFulfillmentRate     float64   `json:"cust_fulfillment_rate"`
	TotalProductUnfulfilled float64   `json:"total_prod_unfulfilled"`
	LastUpdatedAt           time.Time `json:"last_updated_at"`
}

// Report Fulfillment model
type ReportFulfillment struct {
	WeekNumber      string    `orm:"column(week_number)" json:"week_number"`
	StartDate       time.Time `orm:"column(start_date)" json:"start_date"`
	EndDate         time.Time `orm:"column(end_date)" json:"end_date"`
	FulfillmentRate float64   `orm:"column(fulfillment_rate)" json:"fulfillment_rate"`
}

// Unfulfilled Product model
type UnfulfilledProduct struct {
	ID              int64   `orm:"column(product_id)" json:"product_id"`
	UnfulfilledSO   int64   `orm:"column(count_so)" json:"unfulfilled_so"`
	UnfulfilledCust int64   `orm:"column(count_cust)" json:"unfulfilled_cust"`
	UnfulfilledQty  float64 `orm:"column(unfulfilled_qty)" json:"unfulfilled_qty"`
	Product         string  `orm:"column(product)" json:"product"`
	Uom             string  `orm:"column(uom)" json:"uom"`
}
