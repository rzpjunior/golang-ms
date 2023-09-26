// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package archetype

import (
	"strconv"

	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
)

// Save : function to save data requested into database
func Save(r createRequest) (u *model.Archetype, e error) {
	r.Code, e = util.GenerateCode(r.Code, "arc")
	if e == nil {
		customerGroup, _ := strconv.Atoi(r.CustomerGroup)
		u = &model.Archetype{
			Code:          r.Code,
			BusinessType:  r.BusinessType,
			Name:          r.Name,
			Note:          r.Note,
			Status:        1,
			CustomerGroup: int8(customerGroup),
			AuxData:       2,
			// Abbreviation: r.Abbreviation,
		}

		if e = u.Save(); e == nil {
			e = log.AuditLogByUser(r.Session.Staff, u.ID, "archetype", "create", "")
		}
	}

	return u, e
}

// Archive : function to update status data into archive
func Archive(r archiveRequest) (u *model.Archetype, e error) {
	u = &model.Archetype{
		ID:     r.ID,
		Status: int8(2),
	}

	if e = u.Save("id", "status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "archetype", "archive", "")
	}

	return u, e
}

// Unarchive : function to update status data into active
func Unarchive(r unarchiveRequest) (u *model.Archetype, e error) {
	u = &model.Archetype{
		ID:     r.ID,
		Status: int8(1),
	}

	if e = u.Save("id", "status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "archetype", "unarchive", "")
	}

	return u, e
}
