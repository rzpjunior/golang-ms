// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sales_person

import (
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
)

// Archive : function to change staff status into archive
func Archive(r archiveRequest) (u *model.Staff, e error) {
	r.Staff.Status = int8(2)

	if e = r.Staff.Save("Status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, r.Staff.ID, "staff", "archive", "")
	}

	return u, e
}

// Unarchive : function to change staff status into active
func UnArchive(r unarchiveRequest) (u *model.Staff, e error) {
	r.Staff.Status = int8(1)

	if e = r.Staff.Save("Status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, r.Staff.ID, "staff", "unarchive", "")
	}

	return u, e
}
