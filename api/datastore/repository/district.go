package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetDistrict find a single data price set using field and value condition.
func GetDistrict(field string, values ...interface{}) (*model.District, error) {
	m := new(model.District)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetDistricts : function to get data from database based on parameters
func GetDistricts(rq *orm.RequestQuery) (m []*model.District, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.District))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.District
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterDistricts : function to get data from database based on parameters with filtered permission
func GetFilterDistricts(rq *orm.RequestQuery) (m []*model.District, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.District))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.District
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidDistrict : function to check if id is valid in database
func ValidDistrict(id int64) (district *model.District, e error) {
	district = &model.District{ID: id}
	e = district.Read("ID")

	return
}
