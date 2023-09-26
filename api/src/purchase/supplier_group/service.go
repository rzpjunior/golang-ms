// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package supplier_group

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
)

// Save : function to save data requested into database
func Save(r createRequest) (supplierGroup *model.SupplierGroup, err error) {
	o := orm.NewOrm()

	o.Begin()

	supplierGroup = &model.SupplierGroup{
		SupplierCommodity: r.SupplierCommodity,
		SupplierBadge:     r.SupplierBadge,
		SupplierType:      r.SupplierType,
	}

	_, err = o.Insert(supplierGroup)

	if err != nil {
		o.Rollback()
		return nil, err
	}

	err = log.AuditLogByUser(r.Session.Staff, supplierGroup.ID, "supplier_group", "create", "")

	if err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return supplierGroup, err
}

func Update(u updateRequest) (supplierGroup *model.SupplierGroup, err error) {
	o := orm.NewOrm()
	o.Begin()

	supplierGroup = &model.SupplierGroup{
		ID:                u.ID,
		SupplierCommodity: u.SupplierCommodity,
		SupplierBadge:     u.SupplierBadge,
		SupplierType:      u.SupplierType,
	}

	err = supplierGroup.Save("SupplierCommodity", "SupplierBadge", "SupplierType")

	if err != nil {
		o.Rollback()
		return nil, err
	}

	err = log.AuditLogByUser(u.Session.Staff, u.ID, "supplier_group", "update", "")

	if err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return supplierGroup, err
}
