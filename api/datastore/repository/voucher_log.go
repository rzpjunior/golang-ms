// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetVoucherLog find a single data voucher using field and value condition.
func GetVoucherLog(field string, values ...interface{}) (*model.VoucherLog, error) {
	m := new(model.VoucherLog)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetAreas get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetVoucherLogs(rq *orm.RequestQuery) (m []*model.VoucherLog, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.VoucherLog))

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.VoucherLog
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func ValidVoucherLog(id int64) (v *model.VoucherLog, e error) {
	v = &model.VoucherLog{ID: id}
	e = v.Read("ID")

	return
}

// GetAreas get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetFilterVoucherLog(rq *orm.RequestQuery) (m []*model.VoucherLog, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.VoucherLog))

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.VoucherLog
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// CheckVoucherLogData : function to check whether data with particular data based on parameter is exist or not
func CheckVoucherLogData(filter, exclude map[string]interface{}) (voucherLog []*model.VoucherLog, countResult int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.VoucherLog))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if countResult, err := o.All(&voucherLog); err == nil {
		return voucherLog, countResult, nil
	}

	return nil, 0, err
}
