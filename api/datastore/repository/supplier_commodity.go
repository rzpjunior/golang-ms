// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetSupplierCommodity find a single data supplier commodity using field and value condition.
func GetSupplierCommodity(field string, values ...interface{}) (*model.SupplierCommodity, error) {
	m := new(model.SupplierCommodity)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetSupplierCommodities : function to get data from database based on parameters
func GetSupplierCommodities(rq *orm.RequestQuery) (m []*model.SupplierCommodity, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SupplierCommodity))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SupplierCommodity
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

func GetSupplierCommoditiesFilter(rq *orm.RequestQuery) (m []*model.SupplierCommodity, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SupplierCommodity))

	o := orm.NewOrm()
	o.Using("read_only")

	var listSupplierCommodityIds []int64

	_, err = o.Raw("SELECT DISTINCT(supplier_commodity_id) FROM `supplier_group`").QueryRows(&listSupplierCommodityIds)

	if err != nil {
		return nil, total, err
	}

	if total, err = q.Exclude("status", 3).Filter("id__in", listSupplierCommodityIds).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SupplierCommodity
	if _, err = q.Exclude("status", 3).Filter("id__in", listSupplierCommodityIds).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidSupplierCommodity : function to check if supplier commodity data is in table based on supplier badge id
func ValidSupplierCommodity(id int64) (SupplierCommodity *model.SupplierCommodity, e error) {
	SupplierCommodity = &model.SupplierCommodity{ID: id}
	e = SupplierCommodity.Read("ID")

	return
}

func IsExistsSupplierCommodity(name string) (exists bool, err error) {
	m := new(model.SupplierCommodity)
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
