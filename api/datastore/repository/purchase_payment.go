// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"

	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetPurchasePayment find a single data price set using field and value condition.
func GetPurchasePayment(field string, values ...interface{}) (*model.PurchasePayment, error) {
	m := new(model.PurchasePayment)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetPurchasePayments : function to get data from database based on parameters
func GetPurchasePayments(rq *orm.RequestQuery) (m []*model.PurchasePayment, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PurchasePayment))

	if total, err = q.Filter("status__in", 2, 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PurchasePayment
	if _, err = q.Filter("status__in", 2, 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterPurchasePayments : function to get data from database based on parameters with filtered permission
func GetFilterPurchasePayments(rq *orm.RequestQuery) (m []*model.PurchasePayment, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PurchasePayment))

	if total, err = q.Filter("status__in", 2, 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PurchasePayment
	if _, err = q.Filter("status__in", 2, 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidPurchasePayment : function to check if id is valid in database
func ValidPurchasePayment(id int64) (purchaseTerm *model.PurchasePayment, e error) {
	purchaseTerm = &model.PurchasePayment{ID: id}
	e = purchaseTerm.Read("ID")

	return
}

// GetDataPurchasePayment : function to get data based on filter and exclude parameters
func GetDataPurchasePayment(filter map[string]interface{}, exclude map[string]interface{}) (pp []*model.PurchasePayment, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	q := o.QueryTable(new(model.PurchasePayment))

	for k, v := range filter {
		q = q.Filter(k, v)
	}

	for k, v := range exclude {
		q = q.Exclude(k, v)
	}

	if total, err := q.All(&pp); err == nil {
		return pp, total, nil
	}

	return nil, 0, err
}
