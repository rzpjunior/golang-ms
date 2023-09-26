// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package koli

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
)

func Save(r createRequest) (k *model.Koli, e error) {
	o := orm.NewOrm()
	e = o.Begin()
	r.Code, e = util.GenerateCode(r.Code, "koli")

	if e == nil {
		k = &model.Koli{
			Code:   r.Code,
			Name:   r.Name,
			Value:  r.Value,
			Note:   r.Note,
			Status: int8(1),
		}

	} else {
		e = o.Rollback()
		return nil, e
	}
	e = log.AuditLogByUser(r.Session.Staff, k.ID, "koli", "create", "")

	o.Commit()
	return k, nil
}
