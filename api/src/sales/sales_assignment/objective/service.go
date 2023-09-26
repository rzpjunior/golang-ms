// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package objective

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (u *model.SalesAssignmentObjective, e error) {
	r.Code, e = util.GenerateCode(r.Code, "sales_assignment_objective")
	o := orm.NewOrm()
	o.Begin()

	if e == nil {
		u = &model.SalesAssignmentObjective{
			Code:       r.Code,
			Name:       r.Name,
			Objective:  r.Objective,
			SurveyLink: r.SurveyLink,
			Status:     int8(1),
			CreatedAt:  time.Now(),
			CreatedBy:  r.Session.Staff,
		}

		if _, e = o.Insert(u); e != nil {
			o.Rollback()
			return nil, e
		}
		if e = log.AuditLogByUser(r.Session.Staff, u.ID, "sales_assignment_objective", "create", ""); e != nil {
			o.Rollback()
			return nil, e
		}
	}
	o.Commit()

	return u, e
}

func Update(r updateRequest) (u *model.SalesAssignmentObjective, e error) {
	o := orm.NewOrm()
	o.Begin()

	u = &model.SalesAssignmentObjective{
		ID:         r.ID,
		Name:       r.Name,
		Objective:  r.Objective,
		SurveyLink: r.SurveyLink,
	}

	if e = u.Save("Name", "Objective", "SurveyLink"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, u.ID, "sales_assignment_objective", "update", ""); e != nil {
		o.Rollback()
		return nil, e
	}
	o.Commit()

	return
}

func Archive(r archiveRequest) (u *model.SalesAssignmentObjective, e error) {
	o := orm.NewOrm()
	o.Begin()

	u = &model.SalesAssignmentObjective{
		ID:     r.ID,
		Status: int8(2),
	}

	if e = u.Save("id", "status"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, u.ID, "sales_assignment_objective", "archive", ""); e != nil {
		o.Rollback()
		return nil, e
	}
	o.Commit()

	return u, e
}

func UnArchive(r unarchiveRequest) (u *model.SalesAssignmentObjective, e error) {
	o := orm.NewOrm()
	o.Begin()

	u = &model.SalesAssignmentObjective{
		ID:     r.ID,
		Status: int8(1),
	}

	if e = u.Save("id", "status"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, u.ID, "sales_assignment_objective", "unarchive", ""); e != nil {
		o.Rollback()
		return nil, e
	}
	o.Commit()

	return u, e
}
