// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetMatchedProduct : find a single data using field and value condition.
func GetMatchedProduct(field string, values ...interface{}) (*model.MatchedProduct, error) {
	m := new(model.MatchedProduct)
	o := orm.NewOrm()
	o.Using("scrape")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetMatchedProducts : function to get data from database based on parameters
func GetMatchedProducts(rq *orm.RequestQuery) (mp []*model.MatchedProduct, total int64, err error) {
	q, _ := rq.QueryScrape(new(model.MatchedProduct))

	if total, err = q.RelatedSel(1).All(&mp, rq.Fields...); err == nil {
		return mp, total, nil
	}

	return nil, 0, err
}

// ValidMatchedProduct : function to check if id is valid in database
func ValidMatchedProduct(id int64) (mp *model.MatchedProduct, e error) {
	mp = &model.MatchedProduct{ID: id}
	e = mp.Read("ID")

	return
}

// GetMatchingProduct : function to get list of product to matching with public product
func GetMatchingProduct() (pmt []*model.ProductMatchingTemplate, e error) {
	o := orm.NewOrm()
	o.Using("scrape")

	query := "select dp.code dashboard_product_code, dp.name dashboard_product_name, pp1.name public_product_1_name, pp2.name public_product_2_name " +
		"from dashboard_product dp  " +
		"left join matched_product mp on dp.id = mp.dashboard_product_id " +
		"left join public_product_1 pp1 on mp.public_product_1_id = pp1.id " +
		"left join public_product_2 pp2 on mp.public_product_2_id = pp2.id " +
		"order by dp.id"
	if _, e = o.Raw(query).QueryRows(&pmt); e == nil {
		return pmt, nil
	}

	return nil, e
}
