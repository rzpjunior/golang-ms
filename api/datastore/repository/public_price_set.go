// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetPublicPriceSet : find a single data using field and value condition.
func GetPublicPriceSet(field string, values ...interface{}) (*model.PublicPriceSet, error) {
	m := new(model.PublicPriceSet)
	o := orm.NewOrm()
	o.Using("scrape")

	if err := o.QueryTable(m).Filter(field, values...).Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetPublicPriceSets : function to get public product 1 data from database based on parameters
func GetPublicPriceSets(rq *orm.RequestQuery) (pps []*model.PublicPriceSet, total int64, err error) {
	q, _ := rq.QueryScrape(new(model.PublicPriceSet))

	if total, err = q.All(&pps, rq.Fields...); err == nil {
		return pps, total, nil
	}

	return nil, 0, err
}

// ValidPublicPriceSet : function to check if id is valid in database
func ValidPublicPriceSet(id int64) (pps *model.PublicPriceSet, e error) {
	pps = &model.PublicPriceSet{ID: id}
	e = pps.Read("ID")

	return
}
