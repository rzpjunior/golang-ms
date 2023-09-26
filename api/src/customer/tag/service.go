// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package tag

import (
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/api/log"
)

// Save : function to save data requested into database
func Save(r createRequest) (u *model.TagCustomer, e error) {
	r.Code, e = util.GenerateCode(r.Code, "tag_customer")
	if e == nil {
		u = &model.TagCustomer{
			Code:           r.Code,
			Name:           r.Name,
			Note:           r.Note,
			Status:         int8(1),
		}

		if e = u.Save(); e == nil {
			e = log.AuditLogByUser(r.Session.Staff, u.ID, "tag_customer", "create", "")
		}
	}

	return u, e
}

func Archive(r archiveRequest) (u *model.TagCustomer, e error) {
	u = &model.TagCustomer{
		ID:     r.ID,
		Status: int8(2),
	}

	if e = u.Save("id", "status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "tag_customer", "archive", "")
	}

	return u, e
}

func UnArchive(r unarchiveRequest) (u *model.TagCustomer, e error) {
	u = &model.TagCustomer{
		ID:     r.ID,
		Status: int8(1),
	}

	if e = u.Save("id", "status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "tag_customer", "unarchive", "")
	}

	return u, e
}
