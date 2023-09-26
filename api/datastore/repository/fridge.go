package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetFridge find a single data fridge using field and value condition.
func GetFridge(field string, values ...interface{}) (*model.Fridge, error) {
	m := new(model.Fridge)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetFridges : function to get data from database based on parameters
func GetFridges(rq *orm.RequestQuery) (m []*model.Fridge, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Fridge))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Fridge
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterFridges : function to get data from database based on parameters with filtered permission
func GetFilterFridges(rq *orm.RequestQuery) (m []*model.Fridge, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Fridge))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Fridge
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidFridge : function to check if id is valid in database
func ValidFridge(id int64) (box *model.Fridge, e error) {
	box = &model.Fridge{ID: id}
	e = box.Read("ID")

	return
}
