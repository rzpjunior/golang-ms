package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetProvince find a single data price set using field and value condition.
func GetProvince(field string, values ...interface{}) (*model.Province, error) {
	m := new(model.Province)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetProvincies : function to get data from database based on parameters
func GetProvincies(rq *orm.RequestQuery) (m []*model.Province, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Province))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Province
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterProvincies : function to get data from database based on parameters with filtered permission
func GetFilterProvincies(rq *orm.RequestQuery) (m []*model.Province, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Province))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Province
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidProvince : function to check if id is valid in database
func ValidProvince(id int64) (province *model.Province, e error) {
	province = &model.Province{ID: id}
	e = province.Read("ID")

	return
}
