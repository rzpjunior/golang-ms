package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetCity find a single data price set using field and value condition.
func GetCity(field string, values ...interface{}) (*model.City, error) {
	m := new(model.City)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetCities : function to get data from database based on parameters
func GetCities(rq *orm.RequestQuery) (m []*model.City, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.City))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.City
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterCities : function to get data from database based on parameters with filtered permission
func GetFilterCities(rq *orm.RequestQuery) (m []*model.City, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.City))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.City
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidCity : function to check if id is valid in database
func ValidCity(id int64) (city *model.City, e error) {
	city = &model.City{ID: id}
	e = city.Read("ID")

	return
}
