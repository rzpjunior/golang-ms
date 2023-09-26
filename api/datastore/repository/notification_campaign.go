// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"strconv"
	"strings"

	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

// GetNotificationCampaign find a single data sales term using field and value condition.
func GetNotificationCampaign(field string, values ...interface{}) (*model.NotificationCampaign, error) {
	m := new(model.NotificationCampaign)
	o := orm.NewOrm()

	// support for emoji
	o.Raw("SET NAMES 'utf8mb4'").Exec()

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
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

	// checking redirect_to to glossary
	var redirectToGlossary *model.Glossary
	redirectToGlossary, err := GetGlossaryMultipleValue("table", "notification_campaign", "attribute", "redirect_to", "value_int", m.RedirectTo)
	if err != nil {
		return nil, err
	}

	// set redirect_to_name
	m.RedirectToName = redirectToGlossary.ValueName

	// set value name by redirect_to
	switch redirectToGlossary.ValueName {
	case "Product":
		id, _ := strconv.Atoi(m.RedirectValue)
		product := &model.Product{
			ID: orm.ToInt64(id),
		}
		err = product.Read("ID")
		if err != nil {
			return nil, err
		}

		idDec := common.Encrypt(product.ID)
		m.RedirectValue = idDec
		m.RedirectValueName = product.Code + " - " + product.Name

	case "Product Tag":
		id, _ := strconv.Atoi(m.RedirectValue)
		productTag := &model.TagProduct{
			ID: orm.ToInt64(id),
		}
		err := productTag.Read("ID")
		if err != nil {
			return nil, err
		}

		idDec := common.Encrypt(productTag.ID)
		m.RedirectValue = idDec
		m.RedirectValueName = productTag.Name
	case "URL":
		m.RedirectValueName = "URL"
	case "Cart":
		m.RedirectValueName = "Cart"
	case "Promo":
		m.RedirectValueName = "Promo"
	default:
		m.RedirectValueName = "Home"
	}

	// checking push_now
	var pushNowGlossary *model.Glossary
	if pushNowGlossary, err = GetGlossaryMultipleValue("table", "notification_campaign", "attribute", "push_now", "value_int", m.PushNow); err != nil {
		return nil, err
	}
	// checking push now
	if pushNowGlossary.ValueName == "yes" {
		m.PushNowStatus = true
	} else {
		m.PushNowStatus = false
	}

	return m, nil
}

// GetNotificationCampaigns : function to get data from database based on parameters
func GetNotificationCampaigns(rq *orm.RequestQuery, area string) (m []*model.NotificationCampaign, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.NotificationCampaign))
	o := orm.NewOrm()
	o.Using("read_only")
	cond := q.GetCond()

	if area != "" {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("area__icontains", ","+area+",").Or("area__istartswith", area+",").Or("area__iendswith", ","+area).Or("area", area)

		cond = cond.AndCond(cond1)
	}

	q = q.SetCond(cond)

	if total, err = q.All(&m, rq.Fields...); err != nil {
		return nil, 0, err
	}

	for _, r := range m {
		// get area
		areaArr := strings.Split(r.Area, ",")
		qMark := ""
		if len(r.Area) != 0 {
			for range areaArr {
				qMark += "?,"
			}
			qMark = qMark[:len(qMark)-1]
			o.Raw("select group_concat(name) from area where id in ("+qMark+")", areaArr).QueryRow(&r.AreaName)
			r.AreaNameArr = strings.Split(r.AreaName, ",")
		}

		// get archetype
		archetypeArr := strings.Split(r.Archetype, ",")
		qMark = ""
		if len(r.Archetype) != 0 {
			for range archetypeArr {
				qMark += "?,"
			}
			qMark = qMark[:len(qMark)-1]
			o.Raw("select group_concat(name) from archetype where id in ("+qMark+")", archetypeArr).QueryRow(&r.ArchetypeName)
			r.ArchetypeNameArr = strings.Split(r.ArchetypeName, ",")
		}

		// checking redirect_to to glossary
		var redirectToGlossary *model.Glossary
		redirectToGlossary, err = GetGlossaryMultipleValue("table", "notification_campaign", "attribute", "redirect_to", "value_int", r.RedirectTo)
		if err != nil {
			return nil, 0, err
		}

		// set redirect_to_name
		r.RedirectToName = redirectToGlossary.ValueName

		// set value name by redirect_to
		switch redirectToGlossary.ValueName {
		case "Product":
			id, _ := strconv.Atoi(r.RedirectValue)
			product := &model.Product{
				ID: orm.ToInt64(id),
			}
			err = product.Read("ID")
			if err != nil {
				return nil, 0, err
			}

			idDec := common.Encrypt(product.ID)
			r.RedirectValue = idDec
			r.RedirectValueName = product.Code + " - " + product.Name

		case "Product Tag":
			id, _ := strconv.Atoi(r.RedirectValue)
			productTag := &model.TagProduct{
				ID: orm.ToInt64(id),
			}
			err = productTag.Read("ID")
			if err != nil {
				return nil, 0, err
			}

			idDec := common.Encrypt(productTag.ID)
			r.RedirectValue = idDec
			r.RedirectValueName = productTag.Name
		case "URL":
			r.RedirectValueName = "URL"
		case "Cart":
			r.RedirectValueName = "Cart"
		case "Promo":
			r.RedirectValueName = "Promo"
		default:
			r.RedirectValueName = "Home"
		}

		// checking push_now
		var pushNowGlossary *model.Glossary
		if pushNowGlossary, err = GetGlossaryMultipleValue("table", "notification_campaign", "attribute", "push_now", "value_int", r.PushNow); err != nil {
			return nil, 0, err
		}
		// checking push now
		if pushNowGlossary.ValueName == "yes" {
			r.PushNowStatus = true
		} else {
			r.PushNowStatus = false
		}

	}
	return m, total, nil
}

// ValidNotificationCampaign : function to check if id is valid in database
func ValidNotificationCampaign(id int64) (stock *model.Stock, e error) {
	stock = &model.Stock{ID: id}
	e = stock.Read("ID")

	return
}
