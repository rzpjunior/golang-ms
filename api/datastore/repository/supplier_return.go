// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetSupplierReturn find a single data supplier return using field and value condition.
func GetSupplierReturn(field string, values ...interface{}) (*model.SupplierReturn, error) {
	m := new(model.SupplierReturn)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel(2).Limit(1).One(m); err != nil {
		return nil, err
	}
	o.LoadRelated(m, "SupplierReturnItems", 2)

	return m, nil
}

// GetSupplierReturns : function to get data from database based on parameters
func GetSupplierReturns(rq *orm.RequestQuery) (m []*model.SupplierReturn, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	q, _ := rq.QueryReadOnly(new(model.SupplierReturn))

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SupplierReturn
	if _, err = q.All(&mx, rq.Fields...); err != nil {
		return nil, total, err
	}

	return mx, total, nil
}

// ValidSupplierReturn : function to check if id is valid in database
func ValidSupplierReturn(id int64) (supplierReturn *model.SupplierReturn, e error) {
	supplierReturn = &model.SupplierReturn{ID: id}
	e = supplierReturn.Read("ID")

	return
}
