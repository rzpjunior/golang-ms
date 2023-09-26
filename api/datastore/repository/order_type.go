// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetOrderType find a single data using field and value condition.
func GetOrderType(field string, values ...interface{}) (*model.OrderType, error) {
	m := new(model.OrderType)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetOrderTypes : function to get data from database based on parameters
func GetOrderTypes(rq *orm.RequestQuery) (m []*model.OrderType, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.OrderType))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	if _, err = q.Exclude("status", 3).RelatedSel().All(&m, rq.Fields...); err == nil {
		return m, total, nil
	}

	return m, total, err
}

// GetFilterOrderTypes : function to get data from database based on parameters with filtered permission
func GetFilterOrderTypes(rq *orm.RequestQuery) (m []*model.OrderType, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.OrderType))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, 0, err
	}

	if _, err = q.Exclude("status", 3).All(&m, rq.Fields...); err == nil {
		return m, total, nil
	}

	return nil, total, err
}

// ValidOrderType : function to check if id is valid in database
func ValidOrderType(id int64) (orderType *model.OrderType, e error) {
	orderType = &model.OrderType{ID: id}
	e = orderType.Read("ID")

	return
}

// CheckOrderTypeData : function to check data based on filter and exclude parameters
func CheckOrderTypeData(filter, exclude map[string]interface{}) (ot []*model.OrderType, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.OrderType))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err := o.All(&ot); err == nil {
		return ot, total, nil
	}

	return nil, 0, err
}
