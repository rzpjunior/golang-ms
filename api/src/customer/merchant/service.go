// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package merchant

import (
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/log"
)

// Update : function to update data of merchant
func Update(r updateRequest) (u *model.Merchant, e error) {
	u = &model.Merchant{
		ID:                         r.ID,
		PicName:                    r.PicName,
		AltPhoneNumber:             r.AltPhoneNumber,
		Email:                      r.Email,
		Note:                       r.Note,
		InvoiceTerm:                r.InvoiceTerm,
		PaymentTerm:                r.PaymentTerm,
		PaymentGroup:               r.PaymentGroup,
		BillingAddress:             r.BillingAddress,
		LastUpdatedAt:              time.Now(),
		LastUpdatedBy:              r.Session.Staff.ID,
		BusinessTypeCreditLimit:    r.BusinessTypeCreditLimit,
		CustomCreditLimit:          r.CustomCreditLimit,
		CreditLimitAmount:          r.CreditLimitAmount,
		RemainingCreditLimitAmount: r.Merchant.RemainingCreditLimitAmount,
		KTPPhotosUrl:               r.KTPPhotosStr,
		MerchantPhotosUrl:          r.MerchantPhotosStr,
	}

	if e = u.Save("PicName", "AltPhoneNumber", "Email", "Note", "InvoiceTerm", "PaymentTerm", "PaymentGroup", "BillingAddress", "LastUpdatedAt", "LastUpdatedBy", "BusinessTypeCreditLimit", "CustomCreditLimit", "CreditLimitAmount", "RemainingCreditLimitAmount", "KTPPhotosUrl", "MerchantPhotosUrl"); e != nil {
		return nil, e
	}

	if r.IsCreateCreditLimitLog == 1 {
		if e = log.CreditLimitLogByStaff(u, u.ID, "merchant", r.CreditLimitBefore, r.CreditLimitAfter, r.Session.Staff.ID, "update merchant"); e != nil {
			return nil, e
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, u.ID, "merchant", "update", ""); e != nil {
		return nil, e
	}

	return
}

// Archive : function to update status data into archive
func Archive(r archiveRequest) (u *model.Merchant, e error) {
	u = &model.Merchant{
		ID:     r.ID,
		Status: int8(2),
	}

	if e = u.Save("id", "status"); e != nil {
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, u.ID, "merchant", "archive", ""); e != nil {
		return nil, e
	}

	return u, e
}

// Unarchive : function to update status data into active
func Unarchive(r unarchiveRequest) (u *model.Merchant, e error) {
	o := orm.NewOrm()
	o.Begin()

	u = &model.Merchant{
		ID:     r.ID,
		Status: int8(1),
	}

	if _, e = o.Update(u, "status"); e == nil {
		userMerchant := orm.Params{
			"status": int8(1),
		}

		if _, e = o.QueryTable(new(model.UserMerchant)).Filter("id", r.Merchant.UserMerchant.ID).Update(userMerchant); e == nil {
			if branchs, _, e := repository.GetBranchsByMerchantId(u.ID); e == nil {
				for _, _ = range branchs {
					branch := orm.Params{
						"status": int8(1),
					}

					if _, e = o.QueryTable(new(model.Branch)).Filter("merchant_id", r.ID).Update(branch); e != nil {
						o.Rollback()
						return nil, e
					}
				}

				if e = log.AuditLogByUser(r.Session.Staff, u.ID, "merchant", "unarchive", ""); e != nil {
					o.Rollback()
					return nil, e
				}
			}
		}
	} else {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return
}

// UpdateTag : function to update merchant tag
func UpdateTag(r updateTagRequest) (u *model.Merchant, e error) {
	tagCustomer := r.Merchant.TagCustomer
	r.Merchant.TagCustomer = r.CustomerTagStr
	r.Merchant.LastUpdatedAt = time.Now()
	r.Merchant.LastUpdatedBy = r.Session.Staff.ID

	if e = r.Merchant.Save("TagCustomer", "LastUpdatedAt", "LastUpdatedBy"); e != nil {
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.Merchant.ID, "merchant", "update_customer_tag", "previous : "+tagCustomer); e != nil {
		return nil, e
	}

	return u, e
}

// UpdatePhoneNumber : function to update phone number of merchant
func UpdatePhoneNumber(r updatePhoneRequest) (u *model.Merchant, e error) {
	phoneNumber := r.Merchant.PhoneNumber
	r.Merchant.PhoneNumber = strings.TrimPrefix(r.PhoneNumber, "0")
	r.Merchant.LastUpdatedAt = time.Now()
	r.Merchant.LastUpdatedBy = r.Session.Staff.ID

	if e = r.Merchant.Save("PhoneNumber", "LastUpdatedAt", "LastUpdatedBy"); e != nil {
		return nil, e
	}

	r.Merchant.UserMerchant.Verification = 1
	r.Merchant.UserMerchant.ForceLogout = 1

	if e = r.Merchant.UserMerchant.Save("Verification", "ForceLogout"); e != nil {
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.Merchant.ID, "merchant", "update_phone_number", "previous : "+phoneNumber); e != nil {
		return nil, e
	}

	return u, e
}

// Suspension : function to update data suspended on merchant
func Suspension(r suspensionRequest) (sus string, e error) {
	m := &model.Merchant{
		ID:        r.Merchant.ID,
		Suspended: r.Merchant.Suspended,
	}

	if e = m.Save("Suspended"); e != nil {
		return r.SuspendOrUnSuspend, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, m.ID, "merchant", r.SuspendOrUnSuspend, ""); e != nil {
		return r.SuspendOrUnSuspend, e
	}

	return r.SuspendOrUnSuspend, e
}
