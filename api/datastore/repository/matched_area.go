// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetMatchedArea : find a single data using field and value condition.
func GetMatchedArea(field string, values ...interface{}) (*model.MatchedArea, error) {
	m := new(model.MatchedArea)
	o := orm.NewOrm()
	o.Using("scrape")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetMatchedAreas : function to get data from database based on parameters
func GetMatchedAreas(rq *orm.RequestQuery, publicDataID ...string) (ma []*model.MatchedArea, total int64, err error) {
	q, _ := rq.QueryScrape(new(model.MatchedArea))

	if len(publicDataID) > 0 {
		q = q.Filter("public_data_area_"+publicDataID[0]+"_id__isnull", false)
	}

	if total, err = q.RelatedSel(1).All(&ma, rq.Fields...); err == nil {
		return ma, total, nil
	}

	return nil, 0, err
}

// ValidMatchedArea : function to check if id is valid in database
func ValidMatchedArea(id int64) (ma *model.MatchedArea, e error) {
	ma = &model.MatchedArea{ID: id}
	e = ma.Read("ID")

	return
}
