// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetBanners : function to get data from database based on parameters
func GetBanners(rq *orm.RequestQuery, area, archetype string) (m []*model.Banner, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Banner))
	o := orm.NewOrm()
	o.Using("read_only")

	cond := q.GetCond()

	if area != "" {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("area__icontains", ","+area+",").Or("area__istartswith", area+",").Or("area__iendswith", ","+area).Or("area", area)

		cond = cond.AndCond(cond1)
	}

	if archetype != "" {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("archetype__icontains", ","+archetype+",").Or("archetype__istartswith", archetype+",").Or("archetype__iendswith", ","+archetype).Or("archetype", archetype)

		cond = cond.AndCond(cond1)
	}

	q = q.SetCond(cond)

	if total, err = q.All(&m, rq.Fields...); err != nil {
		return nil, 0, err
	}

	currentTime := time.Now()
	for _, v := range m {
		if v.Status == 1 {
			if currentTime.Before(v.StartDate) {
				v.Status = 5
			} else if currentTime.After(v.EndDate) {
				v.Status = 2
			}
		}

		areaArr := strings.Split(v.Area, ",")
		qMark := ""
		if v.Area != "" {
			for range areaArr {
				qMark += "?,"
			}
			qMark = qMark[:len(qMark)-1]
			o.Raw("select group_concat(name) from area where id in ("+qMark+")", areaArr).QueryRow(&v.AreaName)
			v.AreaNameArr = strings.Split(v.AreaName, ",")
		}

		archetypeArr := strings.Split(v.Archetype, ",")
		qMark = ""
		if v.Archetype != "" {
			for range archetypeArr {
				qMark += "?,"
			}
			qMark = qMark[:len(qMark)-1]
			o.Raw("select group_concat(name) from archetype where id in ("+qMark+")", archetypeArr).QueryRow(&v.ArchetypeName)
			v.ArchetypeNameArr = strings.Split(v.ArchetypeName, ",")
		}

		v.Area = util.EncIdInStr(v.Area)
		v.Archetype = util.EncIdInStr(v.Archetype)

		if v.NavigationType == 2 {
			v.TagProduct.Read("ID")
			v.Product = nil
		} else if v.NavigationType == 3 {
			v.Product.Read("ID")
			v.TagProduct = nil
		} else {
			v.TagProduct = nil
			v.Product = nil
		}

		o.Raw("select value_name from glossary where `table` = 'banner' and attribute = 'navigate_type' and value_int = ?", v.NavigationType).QueryRow(&v.NavigationTypeName)
	}

	return m, total, nil
}

// GetBanner find a single data payment method using field and value condition.
func GetBanner(field string, values ...interface{}) (*model.Banner, error) {
	m := new(model.Banner)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	currentTime := time.Now()
	if m.Status == 1 {
		if currentTime.Before(m.StartDate) {
			m.Status = 5
		} else if currentTime.After(m.EndDate) {
			m.Status = 2
		}
	}

	areaArr := strings.Split(m.Area, ",")
	qMark := ""
	if m.Area != "" {
		for range areaArr {
			qMark += "?,"
		}
		qMark = qMark[:len(qMark)-1]
		o.Raw("select group_concat(name) from area where id in ("+qMark+")", areaArr).QueryRow(&m.AreaName)
		m.AreaNameArr = strings.Split(m.AreaName, ",")
	}

	archetypeArr := strings.Split(m.Archetype, ",")
	qMark = ""
	if m.Archetype != "" {
		for range archetypeArr {
			qMark += "?,"
		}
		qMark = qMark[:len(qMark)-1]
		o.Raw("select group_concat(name) from archetype where id in ("+qMark+")", archetypeArr).QueryRow(&m.ArchetypeName)
		m.ArchetypeNameArr = strings.Split(m.ArchetypeName, ",")
	}

	m.Area = util.EncIdInStr(m.Area)
	m.Archetype = util.EncIdInStr(m.Archetype)

	if m.TagProduct == nil || m.TagProduct.ID == 0 {
		m.TagProduct = nil
	}

	if m.Product == nil || m.Product.ID == 0 {
		m.Product = nil
	}

	o.Raw("select value_name from glossary where `table` = 'banner' and attribute = 'navigate_type' and value_int = ?", m.NavigationType).QueryRow(&m.NavigationTypeName)

	return m, nil
}

// ValidBanner : function to check if id is valid in database
func ValidBanner(id int64) (m *model.Banner, e error) {
	m = &model.Banner{ID: id}
	e = m.Read("ID")

	return
}
