// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

type PrintLabelGetRequest struct {
	TypePrint string
	Condition string
}

type PrintLabelGetResponse struct {
	Data string `json:"data"`
}

type LabelPickingRequest struct {
	SalesOrder struct {
		Code   string `json:"code"`
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
		Wrt struct {
			Name string `json:"name"`
		} `json:"wrt"`
		OrderType struct {
			Value string `json:"value"`
		} `json:"order_type"`
	} `json:"sales_order"`
	TotalKoli int64 `json:"total_koli"`
	Helper    struct {
		Code string `json:"code"`
	} `json:"helper"`
}

type RePrintLabelGetRequest struct {
	Increments     []int64 `json:"increment_prints"`
	SalesOrderCode string  `json:"sales_order_code"`
}

type LabelPickingReprintRequest struct {
	SalesOrder struct {
		Code   string `json:"code"`
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
		Wrt struct {
			Name string `json:"name"`
		} `json:"wrt"`
		OrderType struct {
			Value string `json:"value"`
		} `json:"order_type"`
	} `json:"sales_order"`
	TotalKoli  int64 `json:"total_koli"`
	Increments int64 `json:"increment"`
	Helper     struct {
		Code string `json:"code"`
	} `json:"helper"`
}

type DeliveryKoli struct {
	Id        int64   `json:"id"`
	SopNumber string  `json:"sop_number"`
	KoliId    int64   `json:"koli_id"`
	Quantity  float64 `json:"quantity"`
	Note      string  `json:"note"`
	Increment int64   `json:"increment"`
}
type DeliveryKoliRes struct {
	Total        int64           `json:"total"`
	DeliveryKoli []*DeliveryKoli `json:"delivery_koli"`
}
