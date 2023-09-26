// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"
)

// ValidPurchasePlanItem : function to check if id is valid in database
func ValidPurchasePlanItem(id int64) (ppi *model.PurchasePlanItem, err error) {
	ppi = &model.PurchasePlanItem{ID: id}
	err = ppi.Read("ID")

	return
}
