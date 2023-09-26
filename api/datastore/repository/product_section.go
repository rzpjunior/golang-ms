// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetProductSections : function to get data from database based on parameters
func GetProductSections(rq *orm.RequestQuery, area string, archetype string, status []int) (m []*model.ProductSection, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.ProductSection))
	o := orm.NewOrm()
	o.Using("read_only")

	var productSections []*model.ProductSection
	currentTime := time.Now()

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

	if _, err = q.All(&m, rq.Fields...); err != nil {
		return nil, 0, err
	}

	for _, v := range m {
		if v.Status == 1 {
			if currentTime.Before(v.StartAt) {
				v.Status = 5
			}
			if currentTime.After(v.EndAt) {
				v.Status = 2
			}
		}

		// check params status
		if len(status) != 0 {
			for _, value := range status {
				if v.Status == int8(value) {
					productSections = append(productSections, v)
				}
			}
		} else {
			productSections = append(productSections, v)
		}

		// get area
		areaArr := strings.Split(v.Area, ",")
		qMark := ""
		if len(v.Area) != 0 {
			for range areaArr {
				qMark += "?,"
			}
			qMark = qMark[:len(qMark)-1]
			o.Raw("select group_concat(name) from area where id in ("+qMark+")", areaArr).QueryRow(&v.AreaName)
			v.AreaNameArr = strings.Split(v.AreaName, ",")
		}

		// get archetype
		archetypeArr := strings.Split(v.Archetype, ",")
		qMark = ""
		if len(v.Archetype) != 0 {
			for range archetypeArr {
				qMark += "?,"
			}
			qMark = qMark[:len(qMark)-1]
			o.Raw("select group_concat(name) from archetype where id in ("+qMark+")", archetypeArr).QueryRow(&v.ArchetypeName)
			v.ArchetypeNameArr = strings.Split(v.ArchetypeName, ",")
		}
	}

	total = int64(len(productSections))
	return productSections, total, nil
}

// GetProductSection find a single data payment method using field and value condition.
func GetProductSection(field string, values ...interface{}) (*model.ProductSection, error) {
	m := new(model.ProductSection)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	// change current status by time.now
	currentTime := time.Now()
	if m.Status == 1 {
		if currentTime.Before(m.StartAt) {
			m.Status = 5
		}
		if currentTime.After(m.EndAt) {
			m.Status = 2
		}
	}

	// get area
	areaArr := strings.Split(m.Area, ",")
	qMark := ""
	if len(m.Area) != 0 {
		for _, idStr := range areaArr {
			id, _ := strconv.Atoi(idStr)
			area := &model.Area{ID: int64(id)}
			if err := area.Read("ID"); err != nil {
				return nil, err
			}
			m.AreaArr = append(m.AreaArr, area)
			qMark += "?,"
		}
		qMark = qMark[:len(qMark)-1]
		o.Raw("select group_concat(name) from area where id in ("+qMark+")", areaArr).QueryRow(&m.AreaName)
		m.AreaNameArr = strings.Split(m.AreaName, ",")
	}

	// get archetype
	archetypeArr := strings.Split(m.Archetype, ",")
	qMark = ""
	if len(m.Archetype) != 0 {
		for _, idStr := range archetypeArr {
			id, _ := strconv.Atoi(idStr)
			archetype := &model.Archetype{ID: int64(id)}
			if err := archetype.Read("ID"); err != nil {
				return nil, err
			}
			m.ArchetypeArr = append(m.ArchetypeArr, archetype)
			qMark += "?,"
		}
		qMark = qMark[:len(qMark)-1]
		o.Raw("select group_concat(name) from archetype where id in ("+qMark+")", archetypeArr).QueryRow(&m.ArchetypeName)
		m.ArchetypeNameArr = strings.Split(m.ArchetypeName, ",")
	}

	// get product section item
	productArr := strings.Split(m.Product, ",")
	qMark = ""
	if len(m.Product) != 0 {
		for range productArr {
			qMark += "?,"
		}
		qMark = qMark[:len(qMark)-1]

		o.Raw("select id, code, name from product where id in ("+qMark+")", productArr).QueryRows(&m.ProductSectionItem)
	}
	return m, nil
}

// CheckIsIntersect : to check if there already exist active product section at based on parameters
func CheckIsIntersect(sectionType int8, startDate, endDate string) (isExist bool, e error) {
	o := orm.NewOrm()
	o.Using("read_only")

	q := "select exists(select id from product_section where status = 1 and type = ? and (" +
		"(? BETWEEN start_at and end_at) or (? BETWEEN start_at and end_at) or (start_at BETWEEN ? and ?) or (end_at BETWEEN ? and ?)" +
		"))"
	if e = o.Raw(q, sectionType, startDate, endDate, startDate, endDate, startDate, endDate).QueryRow(&isExist); e != nil {
		return true, e
	}

	return false, e
}

// ValidProductSection : function to check if id is valid in database
func ValidProductSection(id int64) (m *model.ProductSection, e error) {
	m = &model.ProductSection{ID: id}
	e = m.Read("ID")

	return
}
