// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// ValidPurchaseOrderImage : function to check if id is valid in database
func ValidPurchaseOrderImage(id int64) (purchaseOrderImage *model.PurchaseOrderImage, e error) {
	purchaseOrderImage = &model.PurchaseOrderImage{ID: id}
	e = purchaseOrderImage.Read("ID")

	return
}

// GetFilterPurchaseOrderImage : function to check data based on filter and exclude parameters
func GetFilterPurchaseOrderImage(filter, exclude map[string]interface{}) (m []*model.PurchaseOrderImage, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.PurchaseOrderImage))
	orm := orm.NewOrm()
	orm.Using("read_only")

	for i, v := range filter {
		o = o.Filter(i, v)
	}

	for i, v := range exclude {
		o = o.Exclude(i, v)
	}

	total, err = o.All(&m)
	if err != nil {
		return nil, 0, err
	}

	return m, total, nil
}
