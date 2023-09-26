// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetSubDistrict find a single data sub district using field and value condition.
func GetSubDistrict(field string, values ...interface{}) (*model.SubDistrict, error) {
	m := new(model.SubDistrict)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetSubDistricts get all data sub district that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetSubDistricts(rq *orm.RequestQuery) (m []*model.SubDistrict, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.SubDistrict))

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.SubDistrict
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// GetFilterSubDistricts : function to get data from database based on parameters with filtered permission
func GetFilterSubDistricts(rq *orm.RequestQuery) (m []*model.SubDistrict, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SubDistrict))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SubDistrict
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

func ValidSubDistrict(id int64) (subdistrict *model.SubDistrict, e error) {
	subdistrict = &model.SubDistrict{ID: id}
	e = subdistrict.Read("ID")

	return
}
