// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// ValidFieldPurchaseOrderItem : function to check if id is valid in database
func ValidFieldPurchaseOrderItem(id int64) (fieldPurchaseOrderItem *model.FieldPurchaseOrderItem, e error) {
	fieldPurchaseOrderItem = &model.FieldPurchaseOrderItem{ID: id}
	e = fieldPurchaseOrderItem.Read("ID")

	return
}

// GetFilterFieldPurchaseOrderItems : function to check data based on filter and exclude parameters
func GetFilterFieldPurchaseOrderItems(filter, exclude map[string]interface{}) (m []*model.FieldPurchaseOrderItem, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.FieldPurchaseOrderItem))
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

	for _, v := range m {
		orm.LoadRelated(v, "FieldPurchaseOrder", 2)
	}

	return m, total, nil
}
