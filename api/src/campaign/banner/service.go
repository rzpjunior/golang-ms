// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package banner

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (b *model.Banner, err error) {
	if r.Code, err = util.GenerateCode(r.Code, "banner"); err != nil {
		return nil, err
	}

	o := orm.NewOrm()
	o.Begin()

	b = &model.Banner{
		Code:           r.Code,
		Name:           r.Name,
		Area:           r.AreaStr,
		Archetype:      r.ArchetypeStr,
		StartDate:      r.StartTimestamp,
		EndDate:        r.EndTimestamp,
		NavigationType: r.NavigationType,
		ImageUrl:       r.ImageUrl,
		Queue:          r.Queue,
		Status:         1,
		Note:           r.Note,
		CreatedAt:      time.Now(),
		CreatedBy:      r.Session.Staff.ID,
	}

	if r.NavigationType == 1 {
		b.NavigationUrl = r.NavigationUrl
	} else if r.NavigationType == 2 {
		b.TagProduct = r.TagProduct
	} else if r.NavigationType == 3 {
		b.Product = r.Product
	} else if r.NavigationType == 6 {
		b.ProductSection = r.ProductSection
	}

	if _, err = o.Insert(b); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, b.ID, "banner", "create", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	b.Area = util.EncIdInStr(b.Area)
	b.Archetype = util.EncIdInStr(b.Archetype)

	return
}

// Archive : function to update status data into archive
func Archive(r archiveRequest) (b *model.Banner, err error) {
	o := orm.NewOrm()
	o.Begin()

	if err = r.Banner.Save("Status"); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, r.Banner.ID, "banner", "archive", r.Note); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()
	return r.Banner, nil
}
