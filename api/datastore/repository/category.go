package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetCategory find a single data division using field and value condition.
func GetCategory(field string, values ...interface{}) (*model.Category, error) {
	m := new(model.Category)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetCategories : function to get data from database based on parameters
func GetCategories(rq *orm.RequestQuery) (m []*model.Category, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Category))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Category
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterCategories : function to get data from database based on parameters with filtered permission
func GetFilterCategories(rq *orm.RequestQuery) (m []*model.Category, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Category))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Category
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err != nil {
		return nil, total, err
	}

	for _, v := range mx {
		v.GrandParent = &model.Category{ID: v.GrandParentID}
		v.GrandParent.Read("ID")
		v.Parent = &model.Category{ID: v.ParentID}
		v.Parent.Read("ID")
	}

	return mx, total, nil
}

func ValidCategory(id int64) (category *model.Category, e error) {
	category = &model.Category{ID: id}
	e = category.Read("ID")

	return
}

// CheckCategoryData : function to check data based on filter and exclude parameters
func CheckCategoryData(filter, exclude map[string]interface{}) (category []*model.Category, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.Category))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err := o.All(&category); err == nil {
		return category, total, nil
	}

	return nil, 0, err
}
