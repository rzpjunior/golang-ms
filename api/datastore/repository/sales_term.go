// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetSalesTerm find a single data sales term using field and value condition.
func GetSalesTerm(field string, values ...interface{}) (*model.SalesTerm, error) {
	m := new(model.SalesTerm)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetSalesTerms : function to get data from database based on parameters
func GetSalesTerms(rq *orm.RequestQuery) (m []*model.SalesTerm, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SalesTerm))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SalesTerm
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterSalesTerms : function to get data from database based on parameters with filtered permission
func GetFilterSalesTerms(rq *orm.RequestQuery, exclude map[string]string) (m []*model.SalesTerm, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SalesTerm))

	if len(exclude) > 0 {
		for i, v := range exclude {
			q = q.Exclude(i, v)
		}
	}

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SalesTerm
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidSalesTerm : function to check if id is valid in database
func ValidSalesTerm(id int64) (salesTerm *model.SalesTerm, e error) {
	salesTerm = &model.SalesTerm{ID: id}
	e = salesTerm.Read("ID")

	return
}

// CheckSalesTermData : function to check data based on filter and exclude parameters
func CheckSalesTermData(filter, exclude map[string]interface{}) (st []*model.SalesTerm, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.SalesTerm))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err := o.All(&st); err == nil {
		return st, total, nil
	}

	return nil, 0, err
}
