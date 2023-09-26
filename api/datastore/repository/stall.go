// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetStall find a single data stall using field and value condition.
func GetStall(field string, values ...interface{}) (*model.Stall, error) {
	m := new(model.Stall)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetStalls get all data stall that matched with query request parameters.
// returning slices of stall, total data without limit and error.
func GetStalls(rq *orm.RequestQuery) (m []*model.Stall, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Stall))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Stall
	if _, err = q.All(&mx, rq.Fields...); err != nil {
		return nil, total, err
	}

	return mx, total, nil
}

func ValidStall(id int64) (stall *model.Stall, e error) {
	stall = &model.Stall{ID: id}
	e = stall.Read("ID")

	return
}
