// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetFieldPurchaseOrder find a single data field_purchase_order set using field and value condition.
func GetFieldPurchaseOrder(field string, values ...interface{}) (*model.FieldPurchaseOrder, error) {
	m := new(model.FieldPurchaseOrder)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "FieldPurchaseOrderItems", 2)

	return m, nil
}

// GetFieldPurchaseOrders get all data field_purchase_order that matched with query request parameters.
// returning slices of field_purchase_order, total data without limit and error.
func GetFieldPurchaseOrders(rq *orm.RequestQuery) (grs []*model.FieldPurchaseOrder, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.FieldPurchaseOrder))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.FieldPurchaseOrder
	if _, err = q.All(&mx, rq.Fields...); err != nil {
		return nil, total, err
	}

	return mx, total, nil
}

// ValidFieldPurchaseOrder : function to check if id is valid in database
func ValidFieldPurchaseOrder(id int64) (fieldPurchaseOrder *model.FieldPurchaseOrder, e error) {
	fieldPurchaseOrder = &model.FieldPurchaseOrder{ID: id}
	e = fieldPurchaseOrder.Read("ID")

	return
}
