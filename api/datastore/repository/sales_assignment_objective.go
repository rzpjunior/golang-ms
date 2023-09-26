package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetSalesAssignmentObjective find a single data sales assignment objective using field and value condition.
func GetSalesAssignmentObjective(field string, values ...interface{}) (*model.SalesAssignmentObjective, error) {
	m := new(model.SalesAssignmentObjective)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetSalesAssignmentObjectives : function to get data from database based on parameters
func GetSalesAssignmentObjectives(rq *orm.RequestQuery) (m []*model.SalesAssignmentObjective, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SalesAssignmentObjective))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SalesAssignmentObjective
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterSalesAssignmentObjectives : function to get data from database based on parameters with filtered permission
func GetFilterSalesAssignmentObjectives(rq *orm.RequestQuery) (m []*model.SalesAssignmentObjective, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SalesAssignmentObjective))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SalesAssignmentObjective
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidSalesAssignmentObjective : function to check if id is valid in database
func ValidSalesAssignmentObjective(id int64) (obj *model.SalesAssignmentObjective, e error) {
	obj = &model.SalesAssignmentObjective{ID: id}
	e = obj.Read("ID")

	return
}

// CheckSalesAssignmentObjective : function to check data based on filter and exclude parameters
func CheckSalesAssignmentObjective(filter, exclude map[string]interface{}) (ps []*model.SalesAssignmentObjective, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.SalesAssignmentObjective))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}
	if total, err := o.All(&ps); err == nil {
		return ps, total, nil
	}

	return nil, 0, err
}

// GetAllSalesAssignmentObjectives : function to get all task assignment objective data
func GetAllSalesAssignmentObjectives() (m []*model.SalesAssignmentObjective, err error) {
	w := new(*model.SalesAssignmentObjective)
	o := orm.NewOrm()
	o.Using("read_only")

	if _, err := o.QueryTable(w).Exclude("status", 3).All(&m); err != nil {
		return nil, err
	}

	return m, nil
}
