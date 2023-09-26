// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uom

import (
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
)

// Save : function to save data requested into database
func Save(r createRequest) (u *model.Uom, e error) {
	r.Code, e = util.GenerateCode(r.Code, "uom")
	if e == nil {
		u = &model.Uom{
			Code:           r.Code,
			Name:           r.Name,
			DecimalEnabled: r.DecimalEnabled,
			Note:           r.Note,
			Status:         int8(1),
		}

		if e = u.Save(); e == nil {
			e = log.AuditLogByUser(r.Session.Staff, u.ID, "uom", "create", "")
		}
	}

	return u, e
}
