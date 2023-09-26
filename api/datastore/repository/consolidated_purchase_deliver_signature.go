// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// ValidConsolidatedPurchaseDeliverSignature : function to check if id is valid in database
func ValidConsolidatedPurchaseDeliverSignature(id int64) (consolidatedPurchaseDeliverSignature *model.ConsolidatedPurchaseDeliverSignature, e error) {
	consolidatedPurchaseDeliverSignature = &model.ConsolidatedPurchaseDeliverSignature{ID: id}
	e = consolidatedPurchaseDeliverSignature.Read("ID")

	return
}

// isRoleCPAlreadySigned : check if role already sign consolidated surat jalan
func IsRoleCPAlreadySigned(consolidatedPurchaseDeliverID int64, role string) (bool, error) {
	consolidatedPurchaseDeliverSignature := new(model.ConsolidatedPurchaseDeliverSignature)
	o := orm.NewOrm()
	o.Using("read_only")

	if total, err := o.QueryTable(consolidatedPurchaseDeliverSignature).Filter("consolidated_purchase_deliver_id", consolidatedPurchaseDeliverID).Filter("role", role).Count(); err != nil || total != 0 {
		return false, err
	}

	return true, nil
}

// GetFilterConsolidatedPurchaseDeliverSignature : function to check data based on filter and exclude parameters
func GetFilterConsolidatedPurchaseDeliverSignature(filter, exclude map[string]interface{}) (m []*model.ConsolidatedPurchaseDeliverSignature, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.ConsolidatedPurchaseDeliverSignature))
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
