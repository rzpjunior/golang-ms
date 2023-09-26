// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"strings"

	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetRole find a single data reservation using field and value condition.
func GetRole(field string, values ...interface{}) (*model.Role, error) {
	m := new(model.Role)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetRoles get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetRoles(rq *orm.RequestQuery) (m []*model.Role, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Role))

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Role
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// GetRoles get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetFilterRole(rq *orm.RequestQuery) (m []*model.Role, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Role))

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Role
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func ValidRole(id int64) (role *model.Role, e error) {
	role = &model.Role{ID: id}
	e = role.Read("ID")

	return
}

// Validate Role is Allowed in Field Purchaser Mobile App
func IsRoleFieldPurchaser(roleID string) (valid bool, err error) {
	configApp := new(model.ConfigApp)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(configApp).Filter("attribute", "fieldpurchaser_role_id").Limit(1).One(configApp); err != nil {
		return false, err
	}

	alloweRoleIDs := strings.Split(configApp.Value, ",")
	for _, alloweRoleID := range alloweRoleIDs {
		if roleID == alloweRoleID {
			return true, nil
		}
	}

	return false, err
}
