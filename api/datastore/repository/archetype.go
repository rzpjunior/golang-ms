// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetArchetype find a single data reservation using field and value condition.
func GetArchetype(field string, values ...interface{}) (*model.Archetype, error) {
	m := new(model.Archetype)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetArchetypes : function to get data from database based on parameters
func GetArchetypes(rq *orm.RequestQuery) (m []*model.Archetype, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Archetype))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Archetype
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterArchetypes : function to get data from database based on parameters with filtered permission
func GetFilterArchetypes(rq *orm.RequestQuery) (m []*model.Archetype, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Archetype))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Archetype
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidArchetype : function to check if id is valid in database
func ValidArchetype(id int64) (archetype *model.Archetype, e error) {
	archetype = &model.Archetype{ID: id}
	e = archetype.Read("ID")

	return
}

// CheckArchetypeData : function to check data based on filter and exclude parameters
func CheckArchetypeData(filter, exclude map[string]interface{}) (archetype []*model.Archetype, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.Archetype))

	for i, v := range filter {
		o = o.Filter(i, v)
	}

	for i, v := range exclude {
		o = o.Exclude(i, v)
	}

	if total, err := o.All(&archetype); err == nil {
		return archetype, total, nil
	}

	return nil, 0, err
}
