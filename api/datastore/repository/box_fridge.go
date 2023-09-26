package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetBoxFridge find a single data box fridge using field and value condition.
func GetBoxFridge(field string, values ...interface{}) (*model.BoxFridge, error) {
	m := new(model.BoxFridge)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetBoxFridges : function to get data from database based on parameters
func GetBoxFridges(rq *orm.RequestQuery) (m []*model.BoxFridge, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.BoxFridge))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.BoxFridge
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterBoxFridges : function to get data from database based on parameters with filtered permission
func GetFilterBoxFridges(rq *orm.RequestQuery) (m []*model.BoxFridge, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.BoxFridge))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.BoxFridge
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidBoxFridge : function to check if id is valid in database
func ValidBoxFridge(id int64) (box *model.BoxFridge, e error) {
	box = &model.BoxFridge{ID: id}
	e = box.Read("ID")

	return
}
