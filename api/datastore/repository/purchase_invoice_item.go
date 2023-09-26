// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"

	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetPurchaseInvoiceItem find a single data price set using field and value condition.
func GetPurchaseInvoiceItem(field string, values ...interface{}) (*model.PurchaseInvoiceItem, error) {
	m := new(model.PurchaseInvoiceItem)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel(2).Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetPurchaseInvoiceItems : function to get data from database based on parameters
func GetPurchaseInvoiceItems(rq *orm.RequestQuery) (m []*model.PurchaseInvoiceItem, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PurchaseInvoiceItem))

	if total, err = q.Exclude("status", 4).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PurchaseInvoiceItem
	if _, err = q.Exclude("status", 4).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidPurchaseInvoiceItem : function to check if id is valid in database
func ValidPurchaseInvoiceItem(id int64) (poi *model.PurchaseInvoiceItem, e error) {
	poi = &model.PurchaseInvoiceItem{ID: id}
	e = poi.Read("ID")

	return
}
