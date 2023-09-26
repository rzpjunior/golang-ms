package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetUom find a single data uom using field and value condition.
func GetUom(field string, values ...interface{}) (*model.Uom, error) {
	m := new(model.Uom)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetUoms : function to get data from database based on parameters
func GetUoms(rq *orm.RequestQuery) (m []*model.Uom, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Uom))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Uom
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterUoms : function to get data from database based on parameters with filtered permission
func GetFilterUoms(rq *orm.RequestQuery) (m []*model.Uom, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Uom))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Uom
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

func ValidUom(id int64) (uom *model.Uom, e error) {
	uom = &model.Uom{ID: id}
	e = uom.Read("ID")

	return
}

// CheckUomData : function to check data based on filter and exclude parameters
func CheckUomData(filter, exclude map[string]interface{}) (uom []*model.Uom, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.Uom))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err := o.All(&uom); err == nil {
		return uom, total, nil
	}

	return nil, 0, err
}
