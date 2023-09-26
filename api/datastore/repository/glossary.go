package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetGlossary : find a single data using field and value condition.
func GetGlossary(field string, values ...interface{}) (*model.Glossary, error) {
	m := new(model.Glossary)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetGlossaryMultipleValue : find a single data using multiple field and value condition.
func GetGlossaryMultipleValue(field1 string, values1 interface{}, field2 string, values2 interface{}, field3 string, values3 interface{}) (*model.Glossary, error) {
	m := new(model.Glossary)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field1, values1).Filter(field2, values2).Filter(field3, values3).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetGlossaries : function to get data from database based on parameters
func GetGlossaries(rq *orm.RequestQuery) (m []*model.Glossary, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Glossary))

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Glossary
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetGlossaries : function to get data from database based on parameters
func GetGlossariesByFilter(filter map[string]interface{}, exclude map[string]interface{}) (g []*model.Glossary, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	q := o.QueryTable(new(model.Glossary))

	for k, v := range filter {
		q = q.Filter(k, v)
	}

	for k, v := range exclude {
		q = q.Exclude(k, v)
	}

	if total, err := q.All(&g); err == nil {
		return g, total, nil
	}

	return nil, 0, err
}
