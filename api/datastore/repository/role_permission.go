// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetRolePermission find a single data role permission using field and value condition.
func GetRolePermission(field string, values ...interface{}) (*model.RolePermission, error) {
	m := new(model.RolePermission)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetRolePermissions get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetRolePermissions(rq *orm.RequestQuery) (m []*model.RolePermission, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.RolePermission))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.RolePermission
	if _, err = q.RelatedSel(1).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func ValidRolePermission(id int64) (rp *model.RolePermission, e error) {
	rp = &model.RolePermission{ID: id}
	e = rp.Read("ID")

	return
}
