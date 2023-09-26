// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package tag_product

import (
	"strings"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
)

// Unarchive: function for update status(1)
func Unarchive(r unarchiveRequest) (u *model.TagProduct, e error) {
	o := orm.NewOrm()
	o.Begin()

	u = &model.TagProduct{
		ID:     r.ID,
		Status: 1,
	}

	if _, err := o.Update(u, "status"); err != nil {
		o.Rollback()
		return nil, err
	}
	err := log.AuditLogByUser(r.Session.Staff, u.ID, "tag_product", "unarchive", "Unarchive Product Tag")
	if err != nil {
		o.Rollback()
		return nil, err
	}
	o.Commit()
	return
}

// Archive: function for update status(2)
func Archive(r archiveRequest) (u *model.TagProduct, e error) {
	o := orm.NewOrm()
	o.Begin()

	u = &model.TagProduct{
		ID:     r.ID,
		Status: 2,
	}

	if _, err := o.Update(u, "status"); err != nil {
		o.Rollback()
		return nil, err
	}
	err := log.AuditLogByUser(r.Session.Staff, u.ID, "tag_product", "archive", "Archive Product Tag")
	if err != nil {
		o.Rollback()
		return nil, err
	}
	o.Commit()
	return
}

// Update: function for update tag product
func Update(r updateRequest) (u *model.TagProduct, e error) {
	o := orm.NewOrm()
	o.Begin()

	u = &model.TagProduct{
		ID:       r.ID,
		Name:     r.Name,
		Area:     strings.ToLower(strings.Join(r.Area, ",")),
		ImageUrl: r.Image,
		Note:     r.Note,
	}

	if _, err := o.Update(u, "name", "image_url", "note", "area"); err != nil {
		o.Rollback()
		return nil, err
	}
	err := log.AuditLogByUser(r.Session.Staff, u.ID, "tag_product", "update", "Update Product Tag")
	if err != nil {
		o.Rollback()
		return nil, err
	}
	o.Commit()
	return
}

// Save: function for insert tag product
func Save(r createRequest) (tp *model.TagProduct, e error) {
	r.Code, e = util.GenerateCode(r.Code, "tag_product")

	o := orm.NewOrm()
	o.Begin()

	tp = &model.TagProduct{
		Code:     r.Code,
		Name:     r.Name,
		Area:     strings.ToLower(strings.Join(r.Area, ",")),
		Status:   1,
		ImageUrl: r.Image,
		Note:     r.Note,
		Value:    r.Value,
	}

	if _, e := o.Insert(tp); e != nil {
		o.Rollback()
		return nil, e
	}

	err := log.AuditLogByUser(r.Session.Staff, tp.ID, "tag_product", "create", "Create Product Tag")
	if err != nil {
		o.Rollback()
		return nil, err
	}
	o.Commit()

	return tp, e
}
