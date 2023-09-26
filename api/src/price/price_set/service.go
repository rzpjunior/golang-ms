// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package price_set

import (
	"strconv"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
)

// Save : function to save data requested into database
func Save(r createRequest) (u *model.PriceSet, e error) {
	r.Code, e = util.GenerateCode(r.Code, "price_set")
	o := orm.NewOrm()
	o.Begin()

	if e == nil {
		u = &model.PriceSet{
			Code:   r.Code,
			Name:   r.Name,
			Note:   r.Note,
			Status: int8(1),
		}

		if _, e = o.Insert(u); e == nil {
			pricesetid := strconv.FormatInt(u.ID, 10)

			if _, e = o.Raw("INSERT INTO price (product_id, price_set_id, unit_price, shadow_price, shadow_price_pct) select id, ?, 0, 0, 0 from product where status in (1,2)", pricesetid).Exec(); e == nil {
				e = log.AuditLogByUser(r.Session.Staff, u.ID, "price_set", "create", "")
			}
		}
	}
	o.Commit()

	return u, e
}

func Update(r updateRequest) (u *model.PriceSet, e error) {
	u = &model.PriceSet{
		ID:   r.ID,
		Name: r.Name,
		Note: r.Note,
	}

	if e = u.Save("Name", "Note"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "price_set", "update", "")
	}

	return
}

func Archive(r archiveRequest) (u *model.PriceSet, e error) {
	u = &model.PriceSet{
		ID:     r.ID,
		Status: int8(2),
	}

	if e = u.Save("id", "status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "price_set", "archive", "")
	}

	return u, e
}

func UnArchive(r unarchiveRequest) (u *model.PriceSet, e error) {
	u = &model.PriceSet{
		ID:     r.ID,
		Status: int8(1),
	}

	if e = u.Save("id", "status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "price_set", "unarchive", "")
	}

	return u, e
}
