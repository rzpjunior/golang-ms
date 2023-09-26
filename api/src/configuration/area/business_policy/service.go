// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package business_policy

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Update : function to update area business config
func Update(r updateRequest) (u *model.AreaBusinessPolicy, e error) {
	o := orm.NewOrm()
	o.Begin()

	u = &model.AreaBusinessPolicy{
		ID:            r.ID,
		MinOrder:      r.MinOrder,
		DeliveryFee:   r.DeliveryFee,
		LastUpdatedAt: time.Now(),
		LastUpdatedBy: r.Session.Staff.ID,
	}

	if _, e = o.Update(u, "MinOrder", "DeliveryFee", "LastUpdatedAt", "LastUpdatedBy"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, u.ID, "area_business_type", "update", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return
}
