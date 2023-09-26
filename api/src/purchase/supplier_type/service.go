// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package supplier_type

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
)

// Save : function to save data requested into database
func Save(r createRequest) (u *model.SupplierType, err error) {
	r.Code, _ = util.GenerateCode(r.Code, "supplier_type")

	o := orm.NewOrm()

	o.Begin()

	supplierType := &model.SupplierType{
		Code:         r.Code,
		Name:         r.Name,
		Abbreviation: r.Abbreviation,
		Note:         r.Note,
		Status:       1,
		CreatedAt:    time.Now(),
		UpdatedBy:    r.Session.Staff,
	}

	_, err = o.Insert(supplierType)

	if err != nil {
		o.Rollback()
		return nil, err
	}

	err = log.AuditLogByUser(r.Session.Staff, supplierType.ID, "supplier type", "create", "")

	if err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return supplierType, err
}

func Update(u updateRequest) (supplierType *model.SupplierType, err error) {
	o := orm.NewOrm()
	o.Begin()

	u.SupplierType.Name = u.Name
	u.SupplierType.Note = u.Note
	u.SupplierType.Abbreviation = u.Abbreviation
	u.SupplierType.UpdatedAt = time.Now()
	u.SupplierType.UpdatedBy = u.Session.Staff

	_, err = o.Update(u.SupplierType, "Name", "Note", "UpdatedAt", "Abbreviation", "UpdatedBy")

	if err != nil {
		o.Rollback()
		return nil, err
	}

	err = log.AuditLogByUser(u.Session.Staff, u.SupplierType.ID, "supplier type", "update", "")

	if err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return u.SupplierType, err
}
