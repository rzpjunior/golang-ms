// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetPublicArea2 : find a single data using field and value condition.
func GetPublicArea2(field string, values ...interface{}) (*model.PublicArea2, error) {
	m := new(model.PublicArea2)
	o := orm.NewOrm()
	o.Using("scrape")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetPublicArea2s : function to get data from database based on parameters
func GetPublicArea2s(rq *orm.RequestQuery) (pa []*model.PublicArea2, total int64, err error) {
	q, _ := rq.QueryScrape(new(model.PublicArea2))

	if total, err = q.RelatedSel(1).All(&pa, rq.Fields...); err == nil {
		return pa, total, nil
	}

	return nil, 0, err
}

// ValidPublicArea2 : function to check if id is valid in database
func ValidPublicArea2(id int64) (pa *model.PublicArea2, e error) {
	pa = &model.PublicArea2{ID: id}
	e = pa.Read("ID")

	return
}
