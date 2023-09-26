// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetAreaPolicy find a single data reservation using field and value condition.
func GetAreaPolicy(field string, values ...interface{}) (*model.AreaPolicy, error) {
	m := new(model.AreaPolicy)
	o := orm.NewOrm()
	o.Using("read_only")
	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetAreaPolicies : function to get data from database based on parameters
func GetAreaPolicies(rq *orm.RequestQuery) (m []*model.AreaPolicy, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.AreaPolicy))

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.AreaPolicy
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterAreaPolicies : function to get data from database based on parameters with filtered permission
func GetFilterAreaPolicies(rq *orm.RequestQuery) (m []*model.AreaPolicy, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.AreaPolicy))

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.AreaPolicy
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidAreaPolicy : function to check if id is valid in database
func ValidAreaPolicy(id int64) (ap *model.AreaPolicy, e error) {
	ap = &model.AreaPolicy{ID: id}
	e = ap.Read("ID")

	return
}
