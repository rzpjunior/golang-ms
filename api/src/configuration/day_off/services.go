package day_off

import (
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
)

// Save : function to save data requested into database
func Save(r createRequest) (u *model.DayOff, e error) {

	u = &model.DayOff{
		OffDate: r.OffDate,
		Note:    r.Note,
		Status:  1,
	}

	if e = u.Save(); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "day_off", "create", "")
	}

	return u, e
}
