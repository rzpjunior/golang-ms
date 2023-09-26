// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetSupplierType find a single data supplier type using field and value condition.
func GetSupplierType(field string, values ...interface{}) (*model.SupplierType, error) {
	m := new(model.SupplierType)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetSupplierTypes : function to get data from database based on parameters
func GetSupplierTypes(rq *orm.RequestQuery) (m []*model.SupplierType, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SupplierType))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SupplierType
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterSupplierTypes : function to get data from database based on parameters with filtered permission
func GetFilterSupplierTypes(rq *orm.RequestQuery, supplierBadgeId int64, supplierCommodityID int64) (m []*model.SupplierType, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SupplierType))

	cond := q.GetCond()

	if supplierBadgeId != 0 && supplierCommodityID != 0 {
		o := orm.NewOrm()
		o.Using("read_only")

		var listSupplierTypeIds []int64

		_, err = o.Raw("SELECT DISTINCT(supplier_type_id) FROM supplier_group WHERE supplier_badge_id = ? AND supplier_commodity_id = ?", supplierBadgeId, supplierCommodityID).QueryRows(&listSupplierTypeIds)

		if err != nil {
			return nil, total, err
		}

		cond1 := orm.NewCondition()
		cond1 = cond1.And("id__in", listSupplierTypeIds)
		cond = cond.AndCond(cond1)
	}

	if total, err = q.SetCond(cond).Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SupplierType
	if _, err = q.SetCond(cond).Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

func ValidSupplierType(id int64) (suppliertype *model.SupplierType, e error) {
	suppliertype = &model.SupplierType{ID: id}
	e = suppliertype.Read("ID")

	return
}

func CheckValidSupplierType(filter, exclude map[string]interface{}) (tc []*model.SupplierGroup, total int64, err error) {
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

func IsExistsSupplierType(name string) (exists bool, err error) {
	m := new(model.SupplierType)
	o := orm.NewOrm()
	o.Using("read_only")

	total, err := o.QueryTable(m).Filter("name", name).RelatedSel().Count()

	if err != nil {
		return false, err
	}

	if total > 0 {
		return true, err
	}

	return false, err
}