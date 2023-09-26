// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package category

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (u *model.Category, e error) {
	o := orm.NewOrm()

	u = &model.Category{
		Code:          r.Code,
		Name:          r.Name,
		GrandParentID: r.GrandParentIDDecrypt,
		ParentID:      r.ParentIDDecrypt,
		Note:          r.Note,
		Status:        int8(1),
	}

	if _, e = o.Insert(u); e != nil {
		return nil, e
	}
	if e = log.AuditLogByUser(r.Session.Staff, u.ID, "category", "create", ""); e != nil {
		return nil, e
	}

	return u, e
}
