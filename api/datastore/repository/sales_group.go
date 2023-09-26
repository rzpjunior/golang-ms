// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"strings"

	"git.edenfarm.id/cuxs/orm"
)

// GetSalesGroup find a single data sales term using field and value condition.
func GetSalesGroup(field string, values ...interface{}) (*model.SalesGroup, error) {
	m := new(model.SalesGroup)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel(2).Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "SalesGroupItems", 1)
	for _, p := range m.SalesGroupItems {
		o.LoadRelated(p.SubDistrict, "District", 1)
	}

	return m, nil
}

// GetSalesGroups : function to get data from database based on parameters
func GetSalesGroups(rq *orm.RequestQuery, city string) (m []*model.SalesGroup, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SalesGroup))
	o := orm.NewOrm()
	o.Using("read_only")

	cond := q.GetCond()
	if city != "0" {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("city__icontains", city)
		cond = cond.AndCond(cond1)
	}

	q = q.SetCond(cond)

	if total, err = q.Exclude("status", 2).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SalesGroup
	if _, err = q.Exclude("status", 2).All(&mx, rq.Fields...); err == nil {

		// Get city string to show on list sales group
		for _, v := range mx {
			citiesArr := strings.Split(v.City, ",")
			var qMark string
			for range citiesArr {
				qMark += "?,"
			}
			qMark = strings.TrimSuffix(qMark, ",")
			if err := o.Raw("SELECT group_concat(name) from city where id in ("+qMark+") order by id asc", citiesArr).QueryRow(&v.CityStr); err != nil {
				return nil, total, err
			}
			if err := o.Raw("SELECT COUNT(id) from staff where sales_group_id = ?", v.ID).QueryRow(&v.SalespersonTotal); err != nil {
				return nil, total, err
			}
		}

		return mx, total, nil
	}

	return nil, total, err
}

// CheckSalesGroupData : function to check data based on filter and exclude parameters
func CheckSalesGroupData(filter, exclude map[string]interface{}) (p []*model.SalesGroup, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.SalesGroup))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err := o.All(&p); err == nil {
		return p, total, nil
	}

	return nil, 0, err
}

// ValidSalesGroup : function to check if id is valid in database
func ValidSalesGroup(id int64) (SalesGroup *model.SalesGroup, e error) {
	SalesGroup = &model.SalesGroup{ID: id}
	e = SalesGroup.Read("ID")

	return
}
