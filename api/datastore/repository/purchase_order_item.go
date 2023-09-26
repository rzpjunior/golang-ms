// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetPurchaseOrderItem find a single data price set using field and value condition.
func GetPurchaseOrderItem(field string, values ...interface{}) (*model.PurchaseOrderItem, error) {
	m := new(model.PurchaseOrderItem)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel(2).Limit(1).One(m); err != nil {
		return nil, err
	}
	o.LoadRelated(m, "FieldPurchaseOrderItems", 2)
	return m, nil
}

// GetPurchaseOrderItems : function to get data from database based on parameters
func GetPurchaseOrderItems(rq *orm.RequestQuery) (m []*model.PurchaseOrderItem, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PurchaseOrderItem))

	if total, err = q.Exclude("status", 4).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PurchaseOrderItem
	if _, err = q.Exclude("status", 4).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidPurchaseOrderItem : function to check if id is valid in database
func ValidPurchaseOrderItem(id int64) (poi *model.PurchaseOrderItem, e error) {
	poi = &model.PurchaseOrderItem{ID: id}
	e = poi.Read("ID")

	return
}

// CheckPurchaseOrderItemData : function to get all purchase order item data based on filter and exclude parameters
func CheckPurchaseOrderItemData(filter, exclude map[string]interface{}) (m []*model.PurchaseOrderItem, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.PurchaseOrderItem))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err = o.All(&m); err != nil {
		return nil, 0, err
	}

	return m, total, nil
}

// GetPurchaseOrderItemByProduct : function to get purchase order item data based on Product ID and PO ID
func GetPurchaseOrderItemByProduct(purchaseOrderID, productID int64) (poi *model.PurchaseOrderItem, e error) {
	o := orm.NewOrm()
	o.Using("read_only")

	if e = o.Raw("SELECT * FROM purchase_order_item WHERE purchase_order_id =? AND product_id =?", purchaseOrderID, productID).QueryRow(&poi); e != nil {
		return nil, e
	}

	return
}
