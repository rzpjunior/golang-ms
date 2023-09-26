// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetSalesAssignments : function to get data from database based on parameters
func GetSalesAssignments(rq *orm.RequestQuery) (m []*model.SalesAssignment, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SalesAssignment))

	if total, err = q.Exclude("status", 4).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SalesAssignment
	if _, err = q.Exclude("status", 4).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidSalesAssignment : function to check if id is valid in database
func ValidSalesAssignment(id int64) (SalesAssignment *model.SalesAssignment, e error) {
	SalesAssignment = &model.SalesAssignment{ID: id}
	e = SalesAssignment.Read("ID")

	return
}
