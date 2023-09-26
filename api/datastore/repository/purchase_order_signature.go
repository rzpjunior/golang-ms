// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// ValidPurchaseOrderSignature : function to check if id is valid in database
func ValidPurchaseOrderSignature(id int64) (purchaseOrderSignature *model.PurchaseOrderSignature, e error) {
	purchaseOrderSignature = &model.PurchaseOrderSignature{ID: id}
	e = purchaseOrderSignature.Read("ID")

	return
}

// isRoleAlreadySignedPurchaseOrder : check if role already sign purchase order
func IsRoleAlreadySignedPurchaseOrder(purchaseOrderID int64, jobFunction string) (bool, error) {
	purchaseOrderSignature := new(model.PurchaseOrderSignature)
	o := orm.NewOrm()
	o.Using("read_only")

	if total, err := o.QueryTable(purchaseOrderSignature).Filter("purchase_order_id", purchaseOrderID).Filter("job_function", jobFunction).Count(); err != nil || total != 0 {
		return false, err
	}

	return true, nil
}

// GetFilterPurchaseOrderSignature : function to check data based on filter and exclude parameters
func GetFilterPurchaseOrderSignature(filter, exclude map[string]interface{}) (m []*model.PurchaseOrderSignature, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.PurchaseOrderSignature))
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
