// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)


// GetSupplierGroup find a single data supplier group using field and value condition.
func GetSupplierGroup(field string, values ...interface{}) (*model.SupplierGroup, error) {
	m := new(model.SupplierGroup)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel(1).Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetSupplierGroups : function to get data from database based on parameters
func GetSupplierGroups(rq *orm.RequestQuery, supplierCommodityId int64, supplierBadgeId int64, supplierTypeId int64) (m []*model.SupplierGroup, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SupplierGroup))

	var mx []*model.SupplierGroup

	cond := q.GetCond()

	o := orm.NewOrm()
	o.Using("read_only")

	if supplierCommodityId != 0 {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("supplier_commodity_id", supplierCommodityId)
		cond = cond.AndCond(cond1)
	}

	if supplierBadgeId != 0 {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("supplier_badge_id", supplierBadgeId)
		cond = cond.AndCond(cond1)
	}

	if supplierTypeId != 0 {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("supplier_type_id", supplierTypeId)
		cond = cond.AndCond(cond1)
	}

	if total, err = q.SetCond(cond).RelatedSel(1).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

func CheckValidSupplierGroup(filter, exclude map[string]interface{}) (tc []*model.SupplierGroup, total int64, err error) {
	m := new(model.SupplierGroup)
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(m)

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err := o.All(&tc); err == nil {
		return tc, total, nil
	}

	return nil, 0, err
}

func IsExistsSupplierGroup(supplierCommodityId int64, supplierBadgeId int64, supplierTypeId int64) (exists bool, err error) {
	m := new(model.SupplierGroup)
	rq := orm.RequestQuery{}
	q, _ := rq.QueryReadOnly(m)

	if supplierCommodityId == 0 || supplierBadgeId == 0 || supplierTypeId == 0 {
		return false, err
	}

	cond := orm.NewCondition()
	cond = cond.And("supplier_commodity_id", supplierCommodityId).And("supplier_badge_id", supplierBadgeId).And("supplier_type_id", supplierTypeId)

	total, err := q.SetCond(cond).Count()

	if err != nil {
		return false, err
	}

	if total > 0 {
		return true, err
	}

	return false, err
}