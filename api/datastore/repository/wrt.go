package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetWrt find a single data using field and value condition.
func GetWrt(field string, values ...interface{}) (*model.Wrt, error) {
	m := new(model.Wrt)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetWrts : function to get data from database based on parameters
func GetWrts(rq *orm.RequestQuery) (m []*model.Wrt, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Wrt))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Wrt
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterWrts : function to get data from database based on parameters with filtered permission
func GetFilterWrts(rq *orm.RequestQuery) (m []*model.Wrt, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Wrt))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Wrt
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidWrt : function to check if id is valid in database
func ValidWrt(id int64) (wrt *model.Wrt, e error) {
	wrt = &model.Wrt{ID: id}
	e = wrt.Read("ID")

	return
}

// CheckWrtData : function to check data based on filter and exclude parameters
func CheckWrtData(filter, exclude map[string]interface{}) (wrt []*model.Wrt, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.Wrt))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err := o.All(&wrt); err == nil {
		return wrt, total, nil
	}

	return nil, 0, err
}
