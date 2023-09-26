// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import "time"

type PickingOrderAssignGetRequest struct {
	Offset           int
	Limit            int
	OrderBy          string
	Status           []int
	PickingOrderId   []int64
	SopNumber        []string
	SiteID           []string
	WrtId            []string
	DeliveryDateFrom time.Time
	DeliveryDateTo   time.Time
	CheckerId        []string
}
