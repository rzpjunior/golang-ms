// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetDivision find a single data division using field and value condition.
func GetDivision(field string, values ...interface{}) (*model.Division, error) {
	m := new(model.Division)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetDivisions get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetDivisions(rq *orm.RequestQuery) (m []*model.Division, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Division))

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Division
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// GetDivisions get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetFilterDivision(rq *orm.RequestQuery) (m []*model.Division, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Division))

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Division
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func ValidDivision(id int64) (div *model.Division, e error) {
	div = &model.Division{ID: id}
	e = div.Read("ID")

	return
}
