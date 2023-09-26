// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

// SalesRecap: struct to hold sales recap data for sales recap list
type SalesRecap struct {
	Product                *Product `json:"product"`
	SumSoQty               float64  `json:"sum_so_qty"`
	SumPoQty               float64  `json:"sum_po_qty"`
	AvailableStock         float64  `json:"available_stock"`
	ExpectedRemainingStock float64  `json:"expected_remaining_stock"`
	SpareQty               float64  `json:"spare_qty"`
	Week1Demand            float64  `json:"week_1_demand"`
	Week2Demand            float64  `json:"week_2_demand"`
	WeekAvgDemand          float64  `json:"week_avg_demand"`
}
