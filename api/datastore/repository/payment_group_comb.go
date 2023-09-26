// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetPaymentGroupCombs : function to get data from database based on parameters
func GetPaymentGroupCombs(rq *orm.RequestQuery) (m []*model.PaymentGroupComb, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PaymentGroupComb))

	if total, err = q.Distinct().Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PaymentGroupComb
	if _, err = q.Distinct().All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetPaymentGroupComb find a single data payment group using field and value condition.
func GetPaymentGroupComb(field string, values ...interface{}) (*model.PaymentGroupComb, error) {
	m := new(model.PaymentGroupComb)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetFilterPaymentGroupCombs : function to get data from database based on parameters with filtered permission
func GetFilterPaymentGroupCombs(rq *orm.RequestQuery) (m []*model.PaymentGroupComb, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PaymentGroupComb))

	if total, err = q.Distinct().Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PaymentGroupComb
	if _, err = q.Distinct().All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidPaymentGroupComb : function to check if id is valid in database
func ValidPaymentGroupComb(id int64) (pgc *model.PaymentGroupComb, e error) {
	pgc = &model.PaymentGroupComb{ID: id}
	e = pgc.Read("ID")

	return
}

// CheckPaymentGroupCombData : function to check data based on filter and exclude parameters
func CheckPaymentGroupCombData(filter, exclude map[string]interface{}) (pgc []*model.PaymentGroupComb, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.PaymentGroupComb))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err := o.All(&pgc); err == nil {
		return pgc, total, nil
	}

	return nil, 0, err
}
