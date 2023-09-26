// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetCreditLimits : function to get data from database based on parameters
func GetCreditLimits(rq *orm.RequestQuery) (m []*model.CreditLimit, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.CreditLimit))

	if total, err = q.Distinct().Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.CreditLimit
	if _, err = q.Distinct().All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetCreditLimit find a single data payment group using field and value condition.
func GetCreditLimit(field string, values ...interface{}) (*model.CreditLimit, error) {
	m := new(model.CreditLimit)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetFilterCreditLimits : function to get data from database based on parameters with filtered permission
func GetFilterCreditLimits(rq *orm.RequestQuery) (m []*model.CreditLimit, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.CreditLimit))

	if total, err = q.Distinct().Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.CreditLimit
	if _, err = q.Distinct().All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidCreditLimit : function to check if id is valid in database
func ValidCreditLimit(id int64) (cl *model.CreditLimit, e error) {
	cl = &model.CreditLimit{ID: id}
	e = cl.Read("ID")

	return
}

// CheckCreditLimitData : function to check data based on filter and exclude parameters
func CheckCreditLimitData(filter, exclude map[string]interface{}) (cl []*model.CreditLimit, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.CreditLimit))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err := o.All(&cl); err == nil {
		return cl, total, nil
	}

	return nil, 0, err
}

// CheckSingleCreditLimitData : function to check one data based on filter and exclude parameters
func CheckSingleCreditLimitData(businessTypeId int64, paymentTermId int64, businessTypeCreditLimit int8) (*model.CreditLimit, error) {
	m := new(model.CreditLimit)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter("business_type_id", businessTypeId).Filter("business_type_credit_limit", businessTypeCreditLimit).Filter("term_payment_sls_id", paymentTermId).RelatedSel().Limit(1).One(m); err != nil {
		return nil, nil
	}

	return m, nil
}

// CheckSingleCreditLimitDataByMerchant : function to check one data based on merchant
func CheckSingleCreditLimitDataByMerchant(merchant *model.Merchant) (*model.CreditLimit, error) {
	m := new(model.CreditLimit)
	o := orm.NewOrm()
	o.Using("read_only")

	if merchant.BusinessTypeCreditLimit == 0 {
		merchant.BusinessTypeCreditLimit = 2
	}

	if err := o.QueryTable(m).Filter("business_type_id", merchant.BusinessType.ID).Filter("business_type_credit_limit", merchant.BusinessTypeCreditLimit).Filter("term_payment_sls_id", merchant.PaymentTerm.ID).RelatedSel().Limit(1).One(m); err != nil {
		return nil, nil
	}

	return m, nil
}
