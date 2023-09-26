// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetSupplierBadge find a single data supplier type using field and value condition.
func GetSupplierBadge(field string, values ...interface{}) (*model.SupplierBadge, error) {
	m := new(model.SupplierBadge)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetSupplierBadges : function to get data from database based on parameters
func GetSupplierBadges(rq *orm.RequestQuery, supplierCommodityId int64) (m []*model.SupplierBadge, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SupplierBadge))
	o := orm.NewOrm()
	o.Using("read_only")

	var listSupplierBadgeIds []int64

	cond := q.GetCond()

	if supplierCommodityId != 0 {

		_, err = o.Raw("SELECT DISTINCT(supplier_badge_id) FROM supplier_group WHERE supplier_commodity_id = ?", supplierCommodityId).QueryRows(&listSupplierBadgeIds)

		if err != nil {
			return nil, total, err
		}

		cond1 := orm.NewCondition()
		cond1 = cond1.And("id__in", listSupplierBadgeIds)
		cond = cond.AndCond(cond1)
	}

	if total, err = q.SetCond(cond).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SupplierBadge
	if _, err = q.SetCond(cond).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterSupplierBadges : function to get data from database based on parameters with filtered permission
func GetFilterSupplierBadges(rq *orm.RequestQuery, supplierCommodityId int64) (m []*model.SupplierBadge, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SupplierBadge))
	o := orm.NewOrm()
	o.Using("read_only")

	var listSupplierBadgeIds []int64
	cond := q.GetCond()

	if supplierCommodityId != 0 {

		_, err = o.Raw("SELECT DISTINCT(supplier_badge_id) FROM supplier_group WHERE supplier_commodity_id = ?", supplierCommodityId).QueryRows(&listSupplierBadgeIds)

		if err != nil {
			return nil, total, err
		}

		cond1 := orm.NewCondition()
		cond1 = cond1.And("id__in", listSupplierBadgeIds)
		cond = cond.AndCond(cond1)
	}

	

	if total, err = q.SetCond(cond).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SupplierBadge
	if _, err = q.SetCond(cond).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidSupplierBadge : function to check if supplier badge data is in table based on supplier badge id
func ValidSupplierBadge(id int64) (SupplierBadge *model.SupplierBadge, e error) {
	SupplierBadge = &model.SupplierBadge{ID: id}
	e = SupplierBadge.Read("ID")

	return
}

func CheckValidSupplierBadge(filter, exclude map[string]interface{}) (tc []*model.SupplierBadge, total int64, err error) {
	m := new(model.SupplierBadge)
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

func IsExistsSupplierBadge(name string) (exists bool, err error) {
	m := new(model.SupplierBadge)
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
