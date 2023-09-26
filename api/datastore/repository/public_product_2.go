// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetPublicProduct2 : find a single data using field and value condition.
func GetPublicProduct2(field string, values ...interface{}) (*model.PublicProduct2, error) {
	m := new(model.PublicProduct2)
	o := orm.NewOrm()
	o.Using("scrape")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetPublicProduct2s : function to get public product 1 data from database based on parameters
func GetPublicProduct2s(rq *orm.RequestQuery) (pp []*model.PublicProduct2, total int64, err error) {
	q, _ := rq.QueryScrape(new(model.PublicProduct2))

	if total, err = q.All(&pp, rq.Fields...); err == nil {
		return pp, total, nil
	}

	return nil, 0, err
}

// ValidPublicProduct2 : function to check if id is valid in database
func ValidPublicProduct2(id int64) (pp *model.PublicProduct2, e error) {
	pp = &model.PublicProduct2{ID: id}
	e = pp.Read("ID")

	return
}

// GetPublicProduct2ForXls : function to get all public product 2 data to be exported to xls
func GetPublicProduct2ForXls() (pp []*model.PublicProductForXls, err error) {
	o := orm.NewOrm()
	o.Using("scrape")

	query := "select pro.name product_name, pro.uom, pro.product_images, dpro.name dashboard_product_name " +
		"from public_product_2 pro " +
		"left join matched_product mp on pro.id = mp.public_product_2_id " +
		"left join dashboard_product dpro on mp.dashboard_product_id = dpro.id " +
		"order by pro.name"
	if _, err = o.Raw(query).QueryRows(&pp); err == nil {
		return pp, err
	}

	return nil, err
}
