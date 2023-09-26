// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

type DeliveryKoliGetRequest struct {
	SopNumber string
}

type DeliveryKoliResponse struct {
	Id        int64   `json:"id,omitempty"`
	SopNumber string  `json:"sop_number,omitempty"`
	KoliId    int64   `json:"koli_id,omitempty"`
	Quantity  float64 `json:"quantity,omitempty"`
	Note      string  `json:"note,omitempty"`
}
