package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetTag find a single data tag using field and value condition.
func GetTag(field string, values ...interface{}) (*model.TagCustomer, error) {
	m := new(model.TagCustomer)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetTags : function to get data from database based on parameters
func GetTags(rq *orm.RequestQuery) (m []*model.TagCustomer, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.TagCustomer))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.TagCustomer
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterTags : function to get data from database based on parameters with filtered permission
func GetFilterTags(rq *orm.RequestQuery) (m []*model.TagCustomer, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.TagCustomer))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.TagCustomer
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

func ValidCustomerTag(id int64) (tag *model.TagCustomer, e error) {
	tag = &model.TagCustomer{ID: id}
	e = tag.Read("ID")

	return
}

// CheckCustomerTagData : function to check data based on filter and exclude parameters
func CheckCustomerTagData(filter, exclude map[string]interface{}) (tc []*model.TagCustomer, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.TagCustomer))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err := o.All(&tc); err == nil {
		return tc, total, nil
	}

	return nil, 0, err
}
