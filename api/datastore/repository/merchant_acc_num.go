// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetMerchantAccNum find a single data division using field and value condition.
func GetMerchantAccNum(field string, values ...interface{}) (*model.MerchantAccNum, error) {
	m := new(model.MerchantAccNum)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetMerchantAccNums : function to get data from database based on parameters
func GetMerchantAccNums(rq *orm.RequestQuery) (man []*model.MerchantAccNum, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.MerchantAccNum))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.MerchantAccNum
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterMerchantAccNums : function to get data from database based on parameters with filtered permission
func GetFilterMerchantAccNums(rq *orm.RequestQuery) (man []*model.MerchantAccNum, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.MerchantAccNum))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.MerchantAccNum
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidMerchantAccNum : function to check if id is valid in database
func ValidMerchantAccNum(id int64) (merchant *model.MerchantAccNum, e error) {
	merchant = &model.MerchantAccNum{ID: id}
	e = merchant.Read("ID")

	return
}

// CheckMerchantAccNumData : function to check data based on filter and exclude parameters
func CheckMerchantAccNumData(filter, exclude map[string]interface{}) (man []*model.MerchantAccNum, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.MerchantAccNum))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err = o.All(&man); err == nil {
		return man, total, err
	}

	return nil, 0, err
}
