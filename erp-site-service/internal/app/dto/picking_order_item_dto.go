// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

type PickingOrderItemGetRequest struct {
	OrderBy              string
	Status               []int
	PickingOrderAssignId []int64
	ItemNumber           []string
}

type PickingOrderItemGetDetailRequest struct {
	Id                   int64
	PickingOrderAssignId int64
	ItemNumber           string
}
