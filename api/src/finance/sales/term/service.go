// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package term

import (
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
)

// Save : function to save data requested into database
func Save(r createRequest) (u *model.SalesTerm, e error) {
	r.Code, e = util.GenerateCode(r.Code, "term_payment_sls")
	if e == nil {
		u = &model.SalesTerm{
			Code:      r.Code,
			Name:      r.Name,
			DaysValue: r.DaysValue,
			Note:      r.Note,
			Status:    int8(1),
		}

		if e = u.Save(); e == nil {
			e = log.AuditLogByUser(r.Session.Staff, u.ID, "sales_term", "create", "")
		}
	}

	return u, e
}

// Archive : function to update status data into archive
func Archive(r archiveRequest) (u *model.SalesTerm, e error) {
	u = &model.SalesTerm{
		ID:     r.ID,
		Status: int8(2),
	}

	if e = u.Save("id", "status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "sales_term", "archive", "")
	}

	return u, e
}

// Unarchive : function to update status data into active
func Unarchive(r unarchiveRequest) (u *model.SalesTerm, e error) {
	u = &model.SalesTerm{
		ID:     r.ID,
		Status: int8(1),
	}

	if e = u.Save("id", "status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "sales_term", "unarchive", "")
	}

	return u, e
}
