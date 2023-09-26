// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package schedule

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
	"time"
)

// Save : function to save data requested into database
func Save(r createRequest) (u *model.PriceSchedule, e error) {
	o := orm.NewOrm()
	o.Begin()

	if e == nil {
		u = &model.PriceSchedule{
			PriceSet:     r.PriceSet,
			Status:       int8(1),
			ScheduleDate: r.ScheduleDateStr,
			ScheduleTime: r.ScheduleTimeStr,
			CreatedAt:    time.Now(),
			CreatedBy:    r.Session.Staff,
		}

		if _, e = o.Insert(u); e == nil {

			var arrPsd []*model.PriceScheduleDump

			for _, v := range r.InsertProductPrice {
				psd := &model.PriceScheduleDump{
					PriceSchedule: u,
					Product:       v.Product,
					PriceSet:      r.PriceSet,
					UnitPrice:     float64(v.UnitPrice),
				}
				arrPsd = append(arrPsd, psd)
			}
			if _, e = o.InsertMulti(100, &arrPsd); e != nil {
				o.Rollback()
				return
			}

		} else {
			o.Rollback()
			return
		}
	}
	o.Commit()

	return u, e
}

func Cancel(r cancelRequest) (u *model.PriceSchedule, e error) {
	o := orm.NewOrm()
	o.Begin()

	u = &model.PriceSchedule{
		ID:     r.ID,
		Status: 3,
		Note:   r.CancellationNote,
	}

	if _, e := o.Update(u, "Status", "Note"); e == nil {
		o.Raw("DELETE FROM price_schedule_dump WHERE price_schedule_id = ?", r.ID).Exec()
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "price_set", "cancel", r.CancellationNote)
	} else {
		o.Rollback()
	}

	o.Commit()

	return
}
