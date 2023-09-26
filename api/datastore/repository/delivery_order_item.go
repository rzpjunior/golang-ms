// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"
)

// ValidDeliveryOrder : function to check if id is valid in database
func ValidDeliveryOrderItem(id int64) (deliveryOrderItem *model.DeliveryOrderItem, e error) {
	deliveryOrderItem = &model.DeliveryOrderItem{ID: id}
	e = deliveryOrderItem.Read("ID")

	return
}
