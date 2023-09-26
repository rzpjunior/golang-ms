// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package prospect_supplier

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
)

// Save : function to insert data requested into database
func Save(r createRequest) (ps *model.ProspectSupplier, e error) {
	r.Code, e = util.GenerateCode(r.Code, "prospect_supplier", 6)

	o := orm.NewOrm()
	o.Begin()

	ps = &model.ProspectSupplier{
		Code:           r.Code,
		Name:           r.Name,
		PhoneNumber:    r.PhoneNumber,
		AltPhoneNumber: r.AltPhoneNumber,
		SubDistrict:    r.SubDistrict,
		StreetAddress:  r.StreetAddress,
		Commodity:      r.CommodityStr,
		PicName:        r.PicName,
		PicPhoneNumber: r.PicPhoneNumber,
		TimeConsent:    r.TimeConsent,
		RegStatus:      1,
		PicAddress:     r.PicAddress,
	}

	if _, e = o.Insert(ps); e == nil {
		e = log.AuditLogByUser(nil, ps.ID, "prospect_supplier", "create", "")
	} else {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return
}

// Register : function to register prospect supplier data requested into supplier table
func Register(r registerRequest) (u *model.ProspectSupplier, e error) {
	r.Code, e = util.GenerateCode(r.Code, "supplier")
	if e == nil {
		f := &model.Supplier{
			SupplierType: r.SupplierType,
			//TermPaymentPur:	r.TermPaymentPur,
			//PaymentMethod:	r.PaymentMethod,
			SubDistrict: r.SubDistrict,
			Code:        r.Code,
			Name:        r.Name,
			//Email:          r.Email,
			PhoneNumber:    r.PhoneNumber,
			AltPhoneNumber: r.AltPhoneNumber,
			Address:        r.Address,
			PicName:        r.PicName,
			Note:           r.Note,
			Status:         int8(1),
		}

		if e = f.Save(); e == nil {
			u = &model.ProspectSupplier{
				ID:        r.ID,
				RegStatus: int8(2),
			}
			if e = u.Save("RegStatus"); e == nil {
				e = log.AuditLogByUser(r.Session.Staff, u.ID, "prospect_supplier", "register", "")
				e = log.AuditLogByUser(r.Session.Staff, u.ID, "supplier", "create", "")
			}
		}
	}

	return u, e
}

func Decline(r declineRequest) (u *model.ProspectSupplier, e error) {
	u = &model.ProspectSupplier{
		ID:        r.ID,
		RegStatus: int8(3),
	}

	if e = u.Save("RegStatus"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "prospect_supplier", "decline", "")
	}

	return u, e
}
