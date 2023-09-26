// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetSupplierOrganization find a single data supplier type using field and value condition.
func GetSupplierOrganization(field string, values ...interface{}) (*model.SupplierOrganization, error) {
	m := new(model.SupplierOrganization)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	if m.SubDistrict != nil {
		o.LoadRelated(m.SubDistrict.District.City, "Province", 2)
	}

	return m, nil
}

// GetSupplierOrganizations : function to get data from database based on parameters
func GetSupplierOrganizations(rq *orm.RequestQuery) (m []*model.SupplierOrganization, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SupplierOrganization))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SupplierOrganization
	if _, err = q.Exclude("status", 3).RelatedSel("SupplierType").All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

func ValidSupplierOrganization(id int64) (SupplierOrganization *model.SupplierOrganization, e error) {
	SupplierOrganization = &model.SupplierOrganization{ID: id}
	e = SupplierOrganization.Read("ID")

	return
}

func IsExistsSupplierOrganization(name string) (exists bool, err error) {
	m := new(model.SupplierOrganization)
	o := orm.NewOrm()
	o.Using("read_only")

	total, err := o.QueryTable(m).Filter("name", name).Count()

	if err != nil {
		return false, err
	}

	if total > 0 {
		return true, nil
	}

	return false, nil
}
