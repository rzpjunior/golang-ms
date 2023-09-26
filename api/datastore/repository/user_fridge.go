package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetUserFridge find a single data UserFridge using field and value condition.
func GetUserFridge(field string, values ...interface{}) (*model.UserFridge, error) {
	m := new(model.UserFridge)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetUserFridges : function to get data from database based on parameters
func GetUserFridges(rq *orm.RequestQuery) (m []*model.UserFridge, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.UserFridge))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.UserFridge
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		for _, userFridge := range mx {
			userFridge.Branch.Read()
			userFridge.Warehouse.Read()
		}
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterUserFridges : function to get data from database based on parameters with filtered permission
func GetFilterUserFridges(rq *orm.RequestQuery) (m []*model.UserFridge, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.UserFridge))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.UserFridge
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidUserFridge : function to check if id is valid in database
func ValidUserFridge(id int64) (box *model.UserFridge, e error) {
	box = &model.UserFridge{ID: id}
	e = box.Read("ID")

	return
}
