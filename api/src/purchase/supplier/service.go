// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package supplier

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (u *model.Supplier, e error) {
	r.Code, e = util.GenerateCode(r.Code, "supplier")
	if e == nil {
		u = &model.Supplier{
			SupplierType:         r.SupplierType,
			PaymentTerm:          r.TermPaymentPur,
			PaymentMethod:        r.PaymentMethod,
			SubDistrict:          r.SubDistrict,
			Code:                 r.Code,
			Name:                 r.Name,
			Email:                r.Email,
			PhoneNumber:          r.PhoneNumber,
			AltPhoneNumber:       r.AltPhoneNumber,
			PicName:              r.PicName,
			Address:              r.Address,
			Note:                 r.Note,
			Status:               int8(1),
			SupplierBadge:        r.SupplierBadge,
			SupplierCommodity:    r.SupplierCommodity,
			SupplierOrganization: r.SupplierOrganization,
			BlockNumber:          r.BlockNumber,
			Rejectable:           r.Rejectable,
			Returnable:           r.Returnable,
			CreatedAt:            time.Now(),
			CreatedBy:            r.Session.Staff,
		}

		if e = u.Save(); e == nil {
			if r.ProspectSupplier != nil {
				r.ProspectSupplier.RegStatus = 2

				if e = r.ProspectSupplier.Save("RegStatus"); e == nil {
					e = log.AuditLogByUser(r.Session.Staff, r.ProspectSupplier.ID, "prospect_supplier", "register", "")
				}
			}
			e = log.AuditLogByUser(r.Session.Staff, u.ID, "supplier", "create", "")
		}
	}

	return u, e
}

// Update : function to update data
func Update(r updateRequest) (u *model.Supplier, e error) {
	u = &model.Supplier{
		ID:                   r.ID,
		PaymentTerm:          r.TermPaymentPur,
		PaymentMethod:        r.PaymentMethod,
		Name:                 r.Name,
		Email:                r.Email,
		PhoneNumber:          r.PhoneNumber,
		AltPhoneNumber:       r.AltPhoneNumber,
		PicName:              r.PicName,
		Address:              r.Address,
		Note:                 r.Note,
		SubDistrict:          r.SubDistrict,
		SupplierType:         r.SupplierType,
		SupplierBadge:        r.SupplierBadge,
		SupplierCommodity:    r.SupplierCommodity,
		SupplierOrganization: r.SupplierOrganization,
		BlockNumber:          r.BlockNumber,
		Rejectable:           r.Rejectable,
		Returnable:           r.Returnable,
	}

	if e = u.Save("PaymentTerm", "PaymentMethod", "Name", "Email", "PhoneNumber", "AltPhoneNumber", "PicName", "Address", "Note", "SubDistrict", "SupplierType", "SupplierBadge", "SupplierCommodity", "SupplierOrganization", "BlockNumber", "Rejectable", "Returnable"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "supplier", "update", "")
	}

	return
}

// Archive : function to update status data into archive
func Archive(r archiveRequest) (u *model.Supplier, e error) {
	u = &model.Supplier{
		ID:     r.ID,
		Status: int8(2),
	}

	if e = u.Save("Status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "supplier", "archive", "")
	}

	return u, e
}

// Unarchive : function to update status data into active
func UnArchive(r unarchiveRequest) (u *model.Supplier, e error) {
	u = &model.Supplier{
		ID:     r.ID,
		Status: int8(1),
	}

	if e = u.Save("Status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "supplier", "unarchive", "")
	}

	return u, e
}

// SaveFieldPurchaser : function to save data requested into database
func SaveFieldPurchaser(r createFieldPurchaserRequest) (m *model.Supplier, err error) {
	o := orm.NewOrm()
	o.Begin()
	r.Code, err = util.GenerateCode(r.Code, "supplier")
	if err != nil {
		o.Rollback()
		return nil, err
	}

	m = &model.Supplier{
		Code:                 r.Code,
		Name:                 r.Name,
		PhoneNumber:          r.PhoneNumber,
		PicName:              r.PicName,
		Address:              r.Address,
		Note:                 r.Note,
		Status:               1,
		Returnable:           r.Returnable,
		Rejectable:           r.Rejectable,
		BlockNumber:          r.BlockNumber,
		CreatedAt:            time.Now(),
		CreatedBy:            r.Session.Staff,
		SupplierType:         r.SupplierOrganization.SupplierType,
		PaymentTerm:          r.SupplierOrganization.TermPaymentPur,
		PaymentMethod:        r.PaymentMethod,
		SubDistrict:          r.SupplierOrganization.SubDistrict,
		SupplierBadge:        r.SupplierOrganization.SupplierBadge,
		SupplierCommodity:    r.SupplierOrganization.SupplierCommodity,
		SupplierOrganization: r.SupplierOrganization,
	}

	if _, err = o.Insert(m); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, m.ID, "supplier", "create", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return m, err
}
