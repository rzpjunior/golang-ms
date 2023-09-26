// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package product_section

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (productSection *model.ProductSection, err error) {
	o := orm.NewOrm()
	o.Begin()

	if r.Code, err = util.GenerateCode(r.Code, "product_section"); err != nil {
		return nil, err
	}

	productSection = &model.ProductSection{
		Code:            r.Code,
		Name:            r.Name,
		BackgroundImage: r.BackgroundImage,
		Area:            r.AreaStr,
		Archetype:       r.ArchetypeStr,
		StartAt:         r.StartAt,
		EndAt:           r.EndAt,
		Sequence:        r.Sequence,
		Product:         r.ProductStr,
		Status:          1,
		Type:            r.Type,
		CreatedAt:       time.Now(),
	}

	if _, err = o.Insert(productSection); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, productSection.ID, "product_section", "create", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()
	return
}

// Update: function to update data from database
func Update(r updateRequest) (productSection *model.ProductSection, e error) {
	o := orm.NewOrm()
	o.Begin()

	productSection = &model.ProductSection{
		ID:              r.ID,
		Name:            r.Name,
		BackgroundImage: r.BackgroundImage,
		Area:            r.AreaStr,
		Archetype:       r.ArchetypeStr,
		StartAt:         r.StartAt,
		EndAt:           r.EndAt,
		Sequence:        r.Sequence,
		Product:         r.ProductStr,
		Type:            r.Type,
		UpdatedAt:       time.Now(),
	}

	if _, e = o.Update(productSection, "name", "area", "archetype", "background_image", "start_at", "end_at", "sequence", "product", "updated_at"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, productSection.ID, "product_section", "update", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return productSection, e
}

// Archive : function to update status data into archive
func Archive(r archiveRequest) (productSection *model.ProductSection, err error) {
	o := orm.NewOrm()
	o.Begin()

	if err = r.ProductSection.Save("Status"); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, r.ProductSection.ID, "product_section", "archive", r.Note); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()
	return r.ProductSection, nil
}
