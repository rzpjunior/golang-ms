// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetDayOff find a single data DayOff using field and value condition.
func GetDayOff(field string, values ...interface{}) (*model.DayOff, error) {
	m := new(model.DayOff)
	o := orm.NewOrm()
	o.Using("read_only")
	if err := o.QueryTable(m).Filter(field, values...).Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetDayOffs get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetDayOffs(rq *orm.RequestQuery) (m []*model.DayOff, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.DayOff))

	// get data requested
	var mx []*model.DayOff
	if total, err = q.All(&mx, rq.Fields...); err == nil && total > 0 {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, 0, err
}
