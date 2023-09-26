package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetBox find a single data box using field and value condition.
func GetBox(field string, values ...interface{}) (*model.Box, error) {
	m := new(model.Box)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetBoxes : function to get data from database based on parameters
func GetBoxes(rq *orm.RequestQuery) (m []*model.Box, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Box))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Box
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterBoxes : function to get data from database based on parameters with filtered permission
func GetFilterBoxes(rq *orm.RequestQuery) (m []*model.Box, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Box))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Box
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidBox : function to check if id is valid in database
func ValidBox(id int64) (box *model.Box, e error) {
	box = &model.Box{ID: id}
	e = box.Read("ID")

	return
}

// ValidRfid : function to check if rfid is valid in database
func ValidRfid(rfid string) (box *model.Box, e error) {
	box = &model.Box{Rfid: rfid}
	e = box.Read("Rfid")

	return
}
