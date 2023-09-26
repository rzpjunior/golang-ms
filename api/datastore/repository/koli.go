// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetKoli find a single data koli using field and value condition.
func GetKoli(field string, values ...interface{}) (*model.Koli, error) {
	m := new(model.Koli)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetKolis get all data koli that matched with query request parameters.
func GetKolis(rq *orm.RequestQuery) (m []*model.Koli, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Koli))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Koli
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func ValidKoli(id int64) (koli *model.Koli, e error) {
	koli = &model.Koli{ID: id}
	e = koli.Read("ID")

	return
}
