// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetPaymentGroups : function to get data from database based on parameters
func GetPaymentGroups(rq *orm.RequestQuery) (m []*model.PaymentGroup, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PaymentGroup))

	if total, err = q.Filter("status", 1).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PaymentGroup
	if _, err = q.Filter("status", 1).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetPaymentGroup find a single data payment group using field and value condition.
func GetPaymentGroup(field string, values ...interface{}) (*model.PaymentGroup, error) {
	m := new(model.PaymentGroup)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetFilterPaymentGroups : function to get data from database based on parameters with filtered permission
func GetFilterPaymentGroups(rq *orm.RequestQuery) (m []*model.PaymentGroup, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PaymentGroup))

	if total, err = q.Filter("status", 1).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PaymentGroup
	if _, err = q.Filter("status", 1).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidWarehouse : function to check if id is valid in database
func ValidPaymentGroup(id int64) (pg *model.PaymentGroup, e error) {
	pg = &model.PaymentGroup{ID: id}
	e = pg.Read("ID")

	return
}
