package repository

import (
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetEdenPointCampaign find a single data price set using field and value condition.
func GetEdenPointCampaign(field string, values ...interface{}) (*model.EdenPointCampaign, error) {
	m := new(model.EdenPointCampaign)
	o := orm.NewOrm()
	o.Using("read_only")
	currentTime := time.Now()

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

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

	customerTagArr := strings.Split(m.TagCustomer, ",")
	qMark = ""
	if m.TagCustomer != "" {
		for range customerTagArr {
			qMark += "?,"
		}
		qMark = qMark[:len(qMark)-1]
		o.Raw("select group_concat(name) from tag_customer where id in ("+qMark+")", customerTagArr).QueryRow(&m.TagCustomerName)
		m.TagCustomerNameArr = strings.Split(m.TagCustomerName, ",")
	}

	m.Area = util.EncIdInStr(m.Area)
	m.Archetype = util.EncIdInStr(m.Archetype)
	m.TagCustomer = util.EncIdInStr(m.TagCustomer)

	return m, nil
}

// GetEdenPointCampaigns : function to get data from database based on parameters
func GetEdenPointCampaigns(rq *orm.RequestQuery, areaID, archetypeID, customerTagID int64, period []string) (m []*model.EdenPointCampaign, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")
	q, _ := rq.QueryReadOnly(new(model.EdenPointCampaign))
	cond := q.GetCond()
	currentTime := time.Now()

	if areaID > 0 {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("area__icontains", ","+strconv.Itoa(int(areaID))+",").Or("area__istartswith", strconv.Itoa(int(areaID))+",").Or("area__iendswith", ","+strconv.Itoa(int(areaID))).Or("area", areaID)

		cond = cond.AndCond(cond1)
	}

	if archetypeID > 0 {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("archetype__icontains", ","+strconv.Itoa(int(archetypeID))+",").Or("archetype__istartswith", strconv.Itoa(int(archetypeID))+",").Or("archetype__iendswith", ","+strconv.Itoa(int(archetypeID))).Or("archetype", archetypeID)

		cond = cond.AndCond(cond1)
	}

	if customerTagID > 0 {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("tag_customer__icontains", ","+strconv.Itoa(int(customerTagID))+",").Or("tag_customer__istartswith", strconv.Itoa(int(customerTagID))+",").Or("tag_customer__iendswith", ","+strconv.Itoa(int(customerTagID))).Or("tag_customer", customerTagID)

		cond = cond.AndCond(cond1)
	}

	if len(period) > 0 {
		cond1 := orm.NewCondition()
		cond1 = cond1.Or("start_date__between", period).Or("end_date__between", period)

		cond2 := orm.NewCondition()
		cond2 = cond2.And("start_date__lte", period[0]).And("end_date__gte", period[1])

		cond1 = cond1.OrCond(cond2)

		cond = cond.AndCond(cond1)
	}

	q = q.SetCond(cond)

	if total, err = q.All(&m, rq.Fields...); err != nil {
		return nil, 0, err
	}

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

		customerTagArr := strings.Split(v.TagCustomer, ",")
		qMark = ""
		if v.TagCustomer != "" {
			for range customerTagArr {
				qMark += "?,"
			}
			qMark = qMark[:len(qMark)-1]
			o.Raw("select group_concat(name) from tag_customer where id in ("+qMark+")", customerTagArr).QueryRow(&v.TagCustomerName)
			v.TagCustomerNameArr = strings.Split(v.TagCustomerName, ",")
		}

		v.Area = util.EncIdInStr(v.Area)
		v.Archetype = util.EncIdInStr(v.Archetype)
		v.TagCustomer = util.EncIdInStr(v.TagCustomer)
	}

	return m, total, nil
}

// ValidEdenPointCampaign : function to check if id is valid in database
func ValidEdenPointCampaign(id int64) (epc *model.EdenPointCampaign, e error) {
	epc = &model.EdenPointCampaign{ID: id}
	e = epc.Read("ID")

	return
}

// GetEdenPointCampaignData : function to get data based on filter and exclude parameters
func GetEdenPointCampaignData(filter, exclude map[string]interface{}, isGetCountOnly int8) (edc []*model.EdenPointCampaign, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.EdenPointCampaign))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if isGetCountOnly == 1 {
		if total, err = o.Distinct().Count(); err != nil {
			return nil, 0, nil
		}

		return nil, total, nil
	}

	if total, err = o.All(&edc); err != nil {
		return nil, 0, err
	}

	return edc, total, nil
}
