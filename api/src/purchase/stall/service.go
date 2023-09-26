// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stall

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
)

// Save : function to save data requested into database
func Save(r createRequest) (m *model.Stall, err error) {
	o := orm.NewOrm()
	o.Begin()

	r.Code, err = util.GenerateCode(r.Code, "supplier")

	if err != nil {
		o.Rollback()
		return nil, err
	}

	m = &model.Stall{
		Code:        r.Code,
		Name:        r.Name,
		PhoneNumber: r.PhoneNumber,
		CreatedAt:   time.Now(),
		CreatedBy:   r.Session.Staff.ID,
	}

	if _, err = o.Insert(m); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, m.ID, "stall", "create", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()
	return m, nil
}
