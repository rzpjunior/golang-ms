// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetNotification find a single data division using field and value condition.
func GetNotification(field string, values ...interface{}) (*model.Notification, error) {
	m := new(model.Notification)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetNotifications : function to get data from database based on parameters
func GetNotifications(rq *orm.RequestQuery) (m []*model.Notification, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Notification))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Notification
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterNotifications : function to get data from database based on parameters with filtered permission
func GetFilterNotifications(rq *orm.RequestQuery) (m []*model.Notification, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Notification))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Notification
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidNotification : function to check if id is valid in database
func ValidNotification(id int64) (notification *model.Notification, e error) {
	notification = &model.Notification{ID: id}
	e = notification.Read("ID")

	return
}
