// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetPurchaseTerm find a single data price set using field and value condition.
func GetPurchaseTerm(field string, values ...interface{}) (*model.PurchaseTerm, error) {
	m := new(model.PurchaseTerm)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetPurchaseTerms : function to get data from database based on parameters
func GetPurchaseTerms(rq *orm.RequestQuery) (m []*model.PurchaseTerm, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PurchaseTerm))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PurchaseTerm
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterPurchaseTerms : function to get data from database based on parameters with filtered permission
func GetFilterPurchaseTerms(rq *orm.RequestQuery) (m []*model.PurchaseTerm, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PurchaseTerm))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PurchaseTerm
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidPurchaseTerm : function to check if id is valid in database
func ValidPurchaseTerm(id int64) (purchaseTerm *model.PurchaseTerm, e error) {
	purchaseTerm = &model.PurchaseTerm{ID: id}
	e = purchaseTerm.Read("ID")

	return
}

// CheckPurchaseTermData : function to check data based on filter and exclude parameters
func CheckPurchaseTermData(filter, exclude map[string]interface{}) (pt []*model.PurchaseTerm, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.PurchaseTerm))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}
	if total, err = o.All(&pt); err == nil {
		return pt, total, nil
	}

	return nil, 0, err
}
