// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetUser find a single data reservation using field and value condition.
func GetUser(field string, values ...interface{}) (*model.User, error) {
	m := new(model.User)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetUsers get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetUsers(rq *orm.RequestQuery) (m []*model.User, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.User))

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.User
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// GetFilterUser get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetFilterUser(rq *orm.RequestQuery) (m []*model.User, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.User))

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.User
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func ValidUser(id int64) (usr *model.User, e error) {
	usr = &model.User{ID: id}
	e = usr.Read("ID")

	return
}
