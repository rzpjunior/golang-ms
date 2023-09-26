// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// ValidPurchaseDeliverSignature : function to check if id is valid in database
func ValidPurchaseDeliverSignature(id int64) (purchaseDeliverSignature *model.PurchaseDeliverSignature, e error) {
	purchaseDeliverSignature = &model.PurchaseDeliverSignature{ID: id}
	e = purchaseDeliverSignature.Read("ID")

	return
}

// isRoleAlreadySigned : check if role already sign surat jalan
func IsRoleAlreadySigned(purchaseDeliverID int64, role string) (bool, error) {
	purchaseDeliverSignature := new(model.PurchaseDeliverSignature)
	o := orm.NewOrm()
	o.Using("read_only")

	if total, err := o.QueryTable(purchaseDeliverSignature).Filter("purchase_deliver_id", purchaseDeliverID).Filter("role", role).Count(); err != nil || total != 0 {
		return false, err
	}

	return true, nil
}

// GetFilterPurchaseDeliverSignature : function to check data based on filter and exclude parameters
func GetFilterPurchaseDeliverSignature(filter, exclude map[string]interface{}) (m []*model.PurchaseDeliverSignature, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.PurchaseDeliverSignature))
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
