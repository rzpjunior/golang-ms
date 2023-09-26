// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package policy

import (
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
)

func Update(r updateRequest) (u *model.AreaPolicy, e error) {
	u = &model.AreaPolicy{
		ID:                  r.ID,
		MinOrder:            r.MinOrder,
		DeliveryFee:         r.DeliveryFee,
		OrderTimeLimit:      r.OrderTimeLimit,
		Area:                r.Area,
		DefaultPriceSet:     r.DefaultPriceSet,
		MaxDayDeliveryDate:  r.MaxDayDeliveryDate,
		WeeklyDayOff:        r.WeeklyDayOff,
		DraftOrderTimeLimit: r.DraftOrderTimeLimit,
	}

	if e = u.Save(); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "area_policy", "update", "")
	}

	return
}
