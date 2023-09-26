// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package application

import (
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
)

func Update(r updateRequest) (u *model.ConfigApp, e error) {
	u = &model.ConfigApp{
		ID:    r.ID,
		Value: r.Value,
	}

	if e = u.Save("id", "value"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "config_app", "update", "")
	}

	return
}
