// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetPaymentChannels : function to get data from database based on parameters
func GetPaymentChannels(rq *orm.RequestQuery) (m []*model.PaymentChannel, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PaymentChannel))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PaymentChannel
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetPaymentChannel find a single data payment method using field and value condition.
func GetPaymentChannel(field string, values ...interface{}) (*model.PaymentChannel, error) {
	m := new(model.PaymentChannel)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetFilterPaymentChannels : function to get data from database based on parameters with filtered permission
func GetFilterPaymentChannels(rq *orm.RequestQuery) (m []*model.PaymentChannel, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PaymentChannel))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PaymentChannel
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidPaymentChannel : function to check if id is valid in database
func ValidPaymentChannel(id int64) (paymentmethod *model.PaymentChannel, e error) {
	paymentmethod = &model.PaymentChannel{ID: id}
	e = paymentmethod.Read("ID")

	return
}

// CheckPaymentChannelData : function to check data based on filter and exclude parameters
func CheckPaymentChannelData(filter, exclude map[string]interface{}) (pc []*model.PaymentChannel, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.PaymentChannel))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err = o.All(&pc); err == nil {
		return pc, total, err
	}

	return nil, 0, err
}
