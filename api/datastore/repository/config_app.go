// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetConfigApp find a single data app_config using field and value condition.
func GetConfigApp(field string, values ...interface{}) (*model.ConfigApp, error) {
	m := new(model.ConfigApp)
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(m)

	if err := o.Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetConfigApps : function to get data from database based on parameters
func GetConfigApps(rq *orm.RequestQuery) (m []*model.ConfigApp, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.ConfigApp))

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.ConfigApp
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterConfigApps : function to get data from database based on parameters with filtered permission
func GetFilterConfigApps(rq *orm.RequestQuery) (m []*model.ConfigApp, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.ConfigApp))

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.ConfigApp
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

func GetConfigAppsByAttribute(field string, values ...interface{}) (ca []*model.ConfigApp, total int64, err error) {
	m := new(model.ConfigApp)
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(m)
	if total, err = o.Filter(field, values...).All(&ca); err != nil {
		return nil, 0, err
	}
	return ca, total, nil
}
