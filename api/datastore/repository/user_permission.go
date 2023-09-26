// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetWarehouses get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetUserPermissions(rq *orm.RequestQuery) (m []*model.UserPermission, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.UserPermission))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.UserPermission
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}
