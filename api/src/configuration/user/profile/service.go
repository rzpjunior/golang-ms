// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package profile

import (
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
)

// Update : function to update staff data
func Update(r updateRequest) (staff *model.Staff, e error) {
	r.Staff.DisplayName = r.DisplayName
	r.Staff.PhoneNumber = r.PhoneNumber

	if e = r.Staff.Save("DisplayName", "PhoneNumber"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, r.Session.Staff.User.ID, "user_profile", "update_profile", "")
	}

	return
}

// UpdatePassword : function to update password of user
func UpdatePassword(r updatePasswordRequest) (user *model.User, e error) {
	user = &model.User{
		ID:       r.Session.Staff.User.ID,
		Password: r.PasswordHash,
	}

	if e = user.Save("Password"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, r.Session.Staff.User.ID, "user_profile", "update_password", "")
	}

	return
}
