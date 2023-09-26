// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package wrt

import (
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (u *model.Wrt, e error) {
	r.Code, e = util.GenerateCode(r.Code, "wrt")
	if e == nil {
		u = &model.Wrt{
			Code:   r.Code,
			Name:   r.Name,
			Note:   r.Note,
			Status: int8(1),
			Area:   r.Area,
			Type:   r.Type,
		}

		if e = u.Save(); e == nil {
			e = log.AuditLogByUser(r.Session.Staff, u.ID, "wrt", "create", "")
		}
	}

	return u, e
}

func Archive(r archiveRequest) (u *model.Wrt, e error) {
	u = &model.Wrt{
		ID:     r.ID,
		Status: int8(2),
	}

	if e = u.Save("id", "status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "wrt", "archive", "")
	}

	return u, e
}

func Unarchive(r unarchiveRequest) (u *model.Wrt, e error) {
	u = &model.Wrt{
		ID:     r.ID,
		Status: int8(1),
	}

	if e = u.Save("id", "status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "wrt", "unarchive", "")
	}

	return u, e
}
