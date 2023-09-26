// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetAreaBusinessPolicy find a single data reservation using field and value condition.
func GetAreaBusinessPolicy(field string, values ...interface{}) (*model.AreaBusinessPolicy, error) {
	m := new(model.AreaBusinessPolicy)
	o := orm.NewOrm()
	o.Using("read_only")
	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetAreaBusinessPolicies : function to get data from database based on parameters
func GetAreaBusinessPolicies(rq *orm.RequestQuery) (m []*model.AreaBusinessPolicy, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.AreaBusinessPolicy))

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.AreaBusinessPolicy
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterAreaBusinessPolicies : function to get data from database based on parameters with filtered permission
func GetFilterAreaBusinessPolicies(rq *orm.RequestQuery) (m []*model.AreaBusinessPolicy, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.AreaBusinessPolicy))

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.AreaBusinessPolicy
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidAreaBusinessPolicy : function to check if id is valid in database
func ValidAreaBusinessPolicy(id int64) (m *model.AreaBusinessPolicy, e error) {
	m = &model.AreaBusinessPolicy{ID: id}
	e = m.Read("ID")

	return
}

// GetAreaBusinessPolicyDelivery : function to get delivery fee by Area ID and Business Type ID
func GetAreaBusinessPolicyDelivery(areaID, businessTypeID int64) (m *model.AreaBusinessPolicy, e error) {
	area := &model.Area{ID: areaID}
	businessType := &model.BusinessType{ID: businessTypeID}
	ap := &model.AreaPolicy{Area: area}
	m = &model.AreaBusinessPolicy{Area: area, BusinessType: businessType}

	if e = m.Read("Area", "BusinessType"); e != nil {
		// If didn't have configuration area business policy, get delivery fee and min order on area policy
		if e = ap.Read("Area"); e != nil {
			return nil, e
		}
		m.DeliveryFee = ap.DeliveryFee
		m.MinOrder = ap.MinOrder
	}

	return m, nil
}
