// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetPublicArea1 : find a single data using field and value condition.
func GetPublicArea1(field string, values ...interface{}) (*model.PublicArea1, error) {
	m := new(model.PublicArea1)
	o := orm.NewOrm()
	o.Using("scrape")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetPublicArea1s : function to get data from database based on parameters
func GetPublicArea1s(rq *orm.RequestQuery) (pa []*model.PublicArea1, total int64, err error) {
	q, _ := rq.QueryScrape(new(model.PublicArea1))

	if total, err = q.All(&pa, rq.Fields...); err == nil {
		return pa, total, nil
	}

	return nil, 0, err
}

// ValidPublicArea1 : function to check if id is valid in database
func ValidPublicArea1(id int64) (pa *model.PublicArea1, e error) {
	pa = &model.PublicArea1{ID: id}
	e = pa.Read("ID")

	return
}
