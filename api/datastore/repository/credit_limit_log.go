// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetCreditLimitLogLogs : function to get data from database based on parameters
func GetCreditLimitLogLogs(rq *orm.RequestQuery) (m []*model.CreditLimitLog, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.CreditLimitLog))

	if total, err = q.Distinct().Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.CreditLimitLog
	if _, err = q.Distinct().All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetCreditLimitLog find a single data payment group using field and value condition.
func GetCreditLimitLog(field string, values ...interface{}) (*model.CreditLimitLog, error) {
	m := new(model.CreditLimitLog)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetFilterCreditLimitLogs : function to get data from database based on parameters with filtered permission
func GetFilterCreditLimitLogs(rq *orm.RequestQuery) (m []*model.CreditLimitLog, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.CreditLimitLog))

	if total, err = q.Distinct().Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.CreditLimitLog
	if _, err = q.Distinct().All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidCreditLimitLog : function to check if id is valid in database
func ValidCreditLimitLog(id int64) (cl *model.CreditLimitLog, e error) {
	cl = &model.CreditLimitLog{ID: id}
	e = cl.Read("ID")

	return
}

// CheckCreditLimitLogData : function to check data based on filter and exclude parameters
func CheckCreditLimitLogData(filter, exclude map[string]interface{}) (cl []*model.CreditLimitLog, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.CreditLimitLog))

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
