// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"strings"

	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetPaymentMethods : function to get data from database based on parameters
func GetPaymentMethods(rq *orm.RequestQuery) (m []*model.PaymentMethod, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PaymentMethod))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PaymentMethod
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetPaymentMethod find a single data payment method using field and value condition.
func GetPaymentMethod(field string, values ...interface{}) (*model.PaymentMethod, error) {
	m := new(model.PaymentMethod)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetFilterPaymentMethods : function to get data from database based on parameters with filtered permission
func GetFilterPaymentMethods(rq *orm.RequestQuery) (m []*model.PaymentMethod, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PaymentMethod))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PaymentMethod
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidPaymentMethod : function to check if id is valid in database
func ValidPaymentMethod(id int64) (paymentmethod *model.PaymentMethod, e error) {
	paymentmethod = &model.PaymentMethod{ID: id}
	e = paymentmethod.Read("ID")

	return
}

// CheckPaymentMethodData : function to check data based on filter and exclude parameters
func CheckPaymentMethodData(filter, exclude map[string]interface{}) (pm []*model.PaymentMethod, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.PaymentMethod))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err := o.All(&pm); err == nil {
		return pm, total, nil
	}

	return pm, 0, err
}

// CheckPaymentMethodHasPayChan : function to check if payment method has payment channel
func CheckPaymentMethodHasPayChan(paymentMethodID int64) (totalPaymentChannel float64, e error) {
	o := orm.NewOrm()
	o.Using("read_only")

	if e = o.Raw("SELECT COUNT(id) FROM payment_channel pc"+
		" where pc.payment_method_id = ? and pc.publish_fva = 1", paymentMethodID).QueryRow(&totalPaymentChannel); e != nil {
		return totalPaymentChannel, e
	}

	return
}

// GetPaymentMethodsFieldPurchaser : function to get data from database based on parameters
func GetPaymentMethodsFieldPurchaser(rq *orm.RequestQuery) (m []*model.PaymentMethod, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PaymentMethod))
	o := orm.NewOrm()
	o.Using("read_only")
	config_app := new(model.ConfigApp)

	o.QueryTable(config_app).Filter("field", "Field Purchaser Payment Method").Filter("attribute", "payment_method_id").Limit(1).One(config_app)

	configAppID := strings.Split(config_app.Value, ",")

	if total, err = q.Exclude("status", 3).Filter("id__in", configAppID).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PaymentMethod
	if _, err = q.Exclude("status", 3).Filter("id__in", configAppID).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}
