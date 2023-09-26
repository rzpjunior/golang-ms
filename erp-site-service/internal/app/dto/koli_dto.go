// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

type KoliGetRequest struct {
	Offset  int
	Limit   int
	OrderBy string
	Status  int
}

type KoliResponse struct {
	Id     int64  `json:"id,omitempty"`
	Code   string `json:"code,omitempty"`
	Value  string `json:"value,omitempty"`
	Name   string `json:"name,omitempty"`
	Note   string `json:"note,omitempty"`
	Status int8   `json:"status,omitempty"`
}
