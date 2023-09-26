// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/project-version2/datamodel/model"
	"strings"

	"git.edenfarm.id/cuxs/orm"
)

// GetArea find a single data area using field and value condition.
func GetArea(field string, values ...interface{}) (*model.Area, error) {
	m := new(model.Area)
	o := orm.NewOrm()
	o.Using("read_only")
	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetAreas get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetAreas(rq *orm.RequestQuery) (m []*model.Area, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Area))

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Area
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func ValidArea(id int64) (area *model.Area, e error) {
	area = &model.Area{ID: id}
	e = area.Read("ID")

	return
}

// GetAreas get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetFilterArea(rq *orm.RequestQuery,areaId string) (m []*model.Area, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Area))

	cond := q.GetCond()
	var areaArrDecrypt []int64
	if areaId != "" {
		areaArr := strings.Split(areaId,",")
		for _,v := range areaArr{
			areaID,_ := common.Decrypt(v)
			areaArrDecrypt = append(areaArrDecrypt,areaID)
		}
		cond1 := orm.NewCondition()
		cond1 = cond1.And("id__in", areaArrDecrypt)

		cond = cond.AndCond(cond1)
	}

	q = q.SetCond(cond)

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Area
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}
