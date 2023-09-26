// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetBusinessTypes : function to get data from database based on parameters
func GetBusinessTypes(rq *orm.RequestQuery) (m []*model.BusinessType, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.BusinessType))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.BusinessType
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetBusinessType find a single data payment method using field and value condition.
func GetBusinessType(field string, values ...interface{}) (*model.BusinessType, error) {
	m := new(model.BusinessType)
	o := orm.NewOrm()
	o.Using("read_only")
	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetFilterBusinessTypes : function to get data from database based on parameters with filtered permission
func GetFilterBusinessTypes(rq *orm.RequestQuery) (m []*model.BusinessType, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.BusinessType))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.BusinessType
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidBusinessType : function to check if id is valid in database
func ValidBusinessType(id int64) (businessType *model.BusinessType, e error) {
	businessType = &model.BusinessType{ID: id}
	e = businessType.Read("ID")

	return
}
