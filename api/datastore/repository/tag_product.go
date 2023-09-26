// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"strings"

	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetProductTag find a single data supplier using field and value condition.
func GetProductTag(field string, values ...interface{}) (*model.TagProduct, error) {
	var arrValueArea []interface{}
	var where string

	m := new(model.TagProduct)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	//to create dynamic question mark for arrValueArea
	arrArea := strings.Split(m.Area, ",")
	for _, i := range arrArea {
		where = where + "? ,"
		arrValueArea = append(arrValueArea, i)
	}
	where = strings.TrimSuffix(where, " ,")

	if err := o.Raw("select group_concat(name) from area a where a.value in ("+where+")", arrValueArea).QueryRow(&m.Area); err != nil {
		return nil, err
	}

	m.Area = strings.ReplaceAll(m.Area, ",", ", ")
	return m, nil
}

// GetProductTags get all data with query request parameters.
// returning slices of User, total data without limit and error.
func GetProductTags(rq *orm.RequestQuery) (m []*model.TagProduct, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.TagProduct))
	o := orm.NewOrm()
	o.Using("read_only")
	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.TagProduct
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			var arrValueArea []interface{}
			var where string
			//to create dynamic question mark for arrValueArea
			arrArea := strings.Split(v.Area, ",")
			for _, i := range arrArea {
				where = where + "? ,"
				arrValueArea = append(arrValueArea, i)
			}
			where = strings.TrimSuffix(where, " ,")

			if err := o.Raw("select group_concat(name) from area a where a.value in ("+where+")", arrValueArea).QueryRow(&v.Area); err != nil {
				return nil, total, err
			}
			v.Area = strings.ReplaceAll(v.Area, ",", ", ")
		}
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// GetFilterProductTags : function to get data from database based on parameters with filtered permission
func GetFilterProductTags(rq *orm.RequestQuery, area string) (m []*model.TagProduct, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.TagProduct))
	o := orm.NewOrm()
	o.Using("read_only")
	if area != "" {

		// set condition for given querystring
		cond := q.GetCond()

		c := orm.NewCondition()
		cond1 := c.And("area__startswith", area+",").
			Or("area__endswith", ","+area).
			Or("area__contains", ","+area+",").
			Or("area", area)

		conditions := cond.AndCond(cond1)

		// get total data
		if total, err = q.SetCond(conditions).Count(); err != nil || total == 0 {
			return nil, total, err
		}

		// get data requested
		var mx []*model.TagProduct
		if _, err = q.SetCond(conditions).All(&mx, rq.Fields...); err == nil {
			for _, v := range mx {
				var arrValueArea []interface{}
				var where string
				//to create dynamic question mark for arrValueArea
				arrArea := strings.Split(v.Area, ",")
				for _, i := range arrArea {
					where = where + "? ,"
					arrValueArea = append(arrValueArea, i)
				}
				where = strings.TrimSuffix(where, " ,")

				if err := o.Raw("select group_concat(name) from area a where a.value in ("+where+")", arrValueArea).QueryRow(&v.Area); err != nil {
					return nil, total, err
				}
				v.Area = strings.ReplaceAll(v.Area, ",", ", ")
			}
			return mx, total, nil
		}

		// return error some thing went wrong
		return nil, total, err
	} else {
		// get total data
		if total, err = q.Count(); err != nil || total == 0 {
			return nil, total, err
		}

		// get data requested
		var mx []*model.TagProduct
		if _, err = q.All(&mx, rq.Fields...); err == nil {
			for _, v := range mx {
				var arrValueArea []interface{}
				var where string
				//to create dynamic question mark for arrValueArea
				arrArea := strings.Split(v.Area, ",")
				for _, i := range arrArea {
					where = where + "? ,"
					arrValueArea = append(arrValueArea, i)
				}
				where = strings.TrimSuffix(where, " ,")

				if err := o.Raw("select group_concat(name) from area a where a.value in ("+where+")", arrValueArea).QueryRow(&v.Area); err != nil {
					return nil, total, err
				}
			}
			return mx, total, nil
		}

		// return error some thing went wrong
		return nil, total, err
	}
}

// ValidTagProduct : function to check if id is valid in database
func ValidTagProduct(id int64) (tp *model.TagProduct, e error) {
	tp = &model.TagProduct{ID: id}
	e = tp.Read("ID")

	return
}
