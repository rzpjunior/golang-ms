// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetProductImage find a single data using field and value condition.
func GetProductImage(field string, values ...interface{}) (*model.ProductImage, error) {
	m := new(model.ProductImage)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetProductImages : function to get data from database based on parameters
func GetProductImages(rq *orm.RequestQuery, tagProductImage string) (m []*model.ProductImage, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.ProductImage))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	if _, err = q.Exclude("status", 3).RelatedSel().All(&m, rq.Fields...); err == nil {
		return m, total, nil
	}

	return m, total, err
}

// GetFilterProductImages : function to get data from database based on parameters with filtered permission
func GetFilterProductImages(rq *orm.RequestQuery) (m []*model.ProductImage, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.ProductImage))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, 0, err
	}

	if _, err = q.Exclude("status", 3).All(&m, rq.Fields...); err == nil {
		return m, total, nil
	}

	return nil, total, err
}

// ValidProductImage : function to check if id is valid in database
func ValidProductImage(id int64) (product *model.ProductImage, e error) {
	product = &model.ProductImage{ID: id}
	e = product.Read("ID")

	return
}
