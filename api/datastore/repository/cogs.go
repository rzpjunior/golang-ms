// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetCogs : find a single data using field and value condition.
func GetCogs(field string, values ...interface{}) (*model.Cogs, error) {
	m := new(model.Cogs)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetCogses : function to get data from database based on parameters
func GetCogses(rq *orm.RequestQuery) (c []*model.Cogs, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Cogs))

	if total, err = q.All(&c, rq.Fields...); err == nil {
		return c, total, nil
	}

	return nil, total, err
}

// GetFilterCogss : function to get data from database based on parameters with filtered permission
func GetFilterCogss(rq *orm.RequestQuery) (c []*model.Cogs, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Cogs))

	if total, err = q.All(&c, rq.Fields...); err == nil {
		return c, total, nil
	}

	return nil, total, err
}

// ValidCogs : function to check if id is valid in database
func ValidCogs(id int64) (Cogs *model.Cogs, e error) {
	Cogs = &model.Cogs{ID: id}
	e = Cogs.Read("ID")

	return
}
