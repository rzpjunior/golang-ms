// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package agent

import (
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/service/xendit"
	"git.edenfarm.id/project-version2/api/util"
)

// Save : function to save data requested into database
func Save(r createRequest) (agent *model.Merchant, e error) {
	o := orm.NewOrm()
	o.Begin()
	var referred *model.Merchant
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	r.UserMerchantCode, e = util.GenerateCode(r.UserMerchantCode, "user_merchant")
	userMerchant := &model.UserMerchant{
		Code: r.UserMerchantCode,
		//Uid:          "-",
		//Password:     "-",
		//FirebaseID:   "-",
		FirebaseToken: "-",
		LoginToken:    "-",
		Verification:  int8(1),
		Status:        int8(1),
	}

	if _, e = o.Insert(userMerchant); e == nil {
		r.Code, e = util.GenerateCustomerCode(r.Code, "merchant")

		cGroup, _ := strconv.Atoi(r.CustomerGroup)
		if r.Merchant != nil {
			referred = r.Merchant
		}
		agent = &model.Merchant{
			UserMerchant:               userMerchant,
			CustomerGroup:              int8(cGroup),
			InvoiceTerm:                r.InvoiceTerm,
			PaymentTerm:                r.PaymentTerm,
			PaymentGroup:               r.PaymentGroup,
			BusinessTypeCreditLimit:    r.BusinessTypeCreditLimit,
			BusinessType:               r.BusinessType,
			FinanceArea:                r.FinanceArea,
			Code:                       r.Code,
			Name:                       r.Name,
			PicName:                    r.PicName,
			PhoneNumber:                strings.TrimPrefix(r.PhoneNumber, "0"),
			AltPhoneNumber:             r.AltPhoneNumber,
			Email:                      r.Email,
			BillingAddress:             r.BillingAddress,
			Note:                       r.Note,
			TagCustomer:                r.CustomerTagStr,
			Status:                     int8(1),
			ReferralCode:               r.ReferralCode,
			ReferrerCode:               r.ReferrerCode,
			Referrer:                   referred,
			CreatedAt:                  time.Now(),
			CreatedBy:                  r.Session.Staff.ID,
			Suspended:                  2,
			CustomCreditLimit:          2,
			CreditLimitAmount:          r.CreditLimitAmount,
			RemainingCreditLimitAmount: r.CreditLimitAmount,
		}

		if r.ProspectCust != nil {
			agent.ProspectCustomer = r.ProspectCust
			r.ProspectCust.RegStatus = int8(2)
			r.ProspectCust.Save("reg_status")

			log.AuditLogByUser(r.Session.Staff, r.ProspectCust.ID, "prospect_customer", "register", "")

			pr := &model.ProspectCustomer{
				ID:          r.ProspectCust.ID,
				ProcessedAt: time.Now(),
				ProcessedBy: r.Session.Staff.ID,
			}

			if _, e = o.Update(pr, "ProcessedAt", "ProcessedBy"); e != nil {
				o.Rollback()
				return nil, e
			}

		}
		if _, e = o.Insert(agent); e == nil {
			var arrPsa []*model.MerchantPriceSet
			var priceSet *model.PriceSet
			// transaction to merchant price set
			// get default price set id for branch
			for _, v := range r.PriceSetArea {
				mps := &model.MerchantPriceSet{
					PriceSet: v.PriceSet,
					Area:     v.Area,
					Merchant: agent,
				}
				arrPsa = append(arrPsa, mps)

				// area suitable with selected sub district
				if v.Area.ID == r.SubDistrict.Area.ID {
					priceSet = v.PriceSet
				}
			}

			if _, e := o.InsertMulti(100, &arrPsa); e == nil {

				r.BranchCode, e = util.GenerateCustomerCode(r.BranchCode, "branch")
				branch := &model.Branch{
					Merchant:        agent,
					Area:            r.ShippingArea,
					Archetype:       r.Archetype,
					PriceSet:        priceSet,
					Warehouse:       r.WarehouseCoverage.Warehouse,
					Salesperson:     r.Salesperson,
					SubDistrict:     r.SubDistrict,
					Code:            r.BranchCode,
					Name:            r.Name,
					PicName:         r.RecipientName,
					PhoneNumber:     r.RecipientPhoneNumber,
					AltPhoneNumber:  r.RecipientAltPhoneNumber,
					AddressName:     r.RecipientName,
					ShippingAddress: r.ShippingAddress,
					Note:            r.ShippingNote,
					MainBranch:      int8(1),
					Status:          int8(1),
					CreatedAt:       time.Now(),
					CreatedBy:       r.Session.Staff.ID,
				}

				if _, e = o.Insert(branch); e == nil {

					if e = log.AuditLogByUser(r.Session.Staff, agent.ID, "agent", "create", ""); e == nil {
						cr := &model.CodeGeneratorReferral{
							Code:      r.ReferralCode,
							CreatedAt: time.Now(),
						}
						if _, e = o.Insert(cr); e == nil {
							if e = log.AuditLogByUser(r.Session.Staff, branch.ID, "shipping_address", "create shipping address", ""); e != nil {
								o.Rollback()
								return nil, e
							}
						} else {
							o.Rollback()
						}

					} else {
						o.Rollback()
						return nil, e
					}

				} else {
					o.Rollback()
					return nil, e
				}
			} else {
				o.Rollback()
				return nil, e
			}
			if r.IsCreateCreditLimitLog == 1 {
				if e = log.CreditLimitLogByStaff(agent, agent.ID, "merchant", r.CreditLimitBefore, r.CreditLimitAfter, r.Session.Staff.ID, "update agent"); e != nil {
					o.Rollback()
					return nil, e
				}
			}
		} else {
			o.Rollback()
			return nil, e
		}

		if r.ProspectCustomerId != "" && r.ProspectCust.SalespersonID != 0 {
			// notification FS Apps when Register Prospective Customer
			messageNotif := &util.MessageNotification{}

			orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0017'").QueryRow(&messageNotif)
			messageNotif.Message = util.ReplaceNotificationSalesOrder(messageNotif.Message, "#Name#", r.Name)

			mn := &util.ModelNotification{
				SendTo:    r.Staff.User.SalesAppNotifToken,
				Title:     messageNotif.Title,
				Message:   messageNotif.Message,
				Type:      "2",
				RefID:     r.ProspectCust.ID,
				ServerKey: util.FieldSalesServerKeyFireBase,
				StaffID:   r.Staff.ID,
			}
			util.PostModelNotificationFieldSales(mn)
		}

	} else {
		o.Rollback()
		return nil, e
	}
	o.Commit()
	if e == nil {
		xendit.BCAXenditFixedVA(agent)
		xendit.PermataXenditFixedVA(agent)
	}
	return
}

func Archive(r archiveRequest) (agent *model.Merchant, e error) {
	o := orm.NewOrm()
	o.Begin()

	agent = &model.Merchant{
		ID:     r.ID,
		Status: int8(2),
	}

	if _, e = o.Update(agent, "status"); e == nil {
		userMerchant := orm.Params{
			"status": int8(2),
		}

		if _, e = o.QueryTable(new(model.UserMerchant)).Filter("id", r.Merchant.UserMerchant.ID).Update(userMerchant); e == nil {
			if branchs, _, e := repository.GetBranchsByMerchantId(agent.ID); e == nil {
				for _, _ = range branchs {
					branch := orm.Params{
						"status": int8(2),
					}

					if _, e = o.QueryTable(new(model.Branch)).Filter("merchant_id", r.ID).Update(branch); e != nil {
						o.Rollback()
						return nil, e
					}
				}

				if e = log.AuditLogByUser(r.Session.Staff, agent.ID, "agent", "archive", ""); e != nil {
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

func Unarchive(r unarchiveRequest) (agent *model.Merchant, e error) {
	o := orm.NewOrm()
	o.Begin()

	agent = &model.Merchant{
		ID:     r.ID,
		Status: int8(1),
	}

	if _, e = o.Update(agent, "status"); e == nil {
		userMerchant := orm.Params{
			"status": int8(1),
		}

		if _, e = o.QueryTable(new(model.UserMerchant)).Filter("id", r.Merchant.UserMerchant.ID).Update(userMerchant); e == nil {
			if branchs, _, e := repository.GetBranchsByMerchantId(agent.ID); e == nil {
				for _, _ = range branchs {
					branch := orm.Params{
						"status": int8(1),
					}

					if _, e = o.QueryTable(new(model.Branch)).Filter("merchant_id", r.ID).Update(branch); e != nil {
						o.Rollback()
						return nil, e
					}
				}

				if e = log.AuditLogByUser(r.Session.Staff, agent.ID, "agent", "unarchive", ""); e != nil {
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

func Update(r updateRequest) (agent *model.Merchant, e error) {
	o := orm.NewOrm()
	o.Begin()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	pca := make(map[int64]*model.PriceSet)
	agent = &model.Merchant{
		ID:                         r.ID,
		PicName:                    r.PicName,
		AltPhoneNumber:             r.AltPhoneNumber,
		Email:                      r.Email,
		Note:                       r.Note,
		BillingAddress:             r.BillingAddress,
		InvoiceTerm:                r.InvoiceTerm,
		PaymentTerm:                r.PaymentTerm,
		PaymentGroup:               r.PaymentGroup,
		LastUpdatedAt:              time.Now(),
		LastUpdatedBy:              r.Session.Staff.ID,
		CustomCreditLimit:          r.CustomCreditLimit,
		BusinessTypeCreditLimit:    r.BusinessTypeCreditLimit,
		CreditLimitAmount:          r.CreditLimitAmount,
		RemainingCreditLimitAmount: r.Merchant.RemainingCreditLimitAmount,
	}

	if _, e = o.Update(agent, "pic_name", "alt_phone_number", "email", "note", "term_invoice_sls_id", "term_payment_sls_id", "payment_method_id", "billing_address", "payment_group_sls_id", "LastUpdatedAt", "LastUpdatedBy", "business_type_credit_limit", "custom_credit_limit", "credit_limit_amount", "credit_limit_remaining"); e == nil {
		// delete insert for merchant_price_set
		var mpsList []*model.MerchantPriceSet
		orSelect.Raw("select * from merchant_price_set mps where merchant_id = ?", agent.ID).QueryRows(&mpsList)
		for _, v := range mpsList {
			if _, e := o.Delete(v); e != nil {
				o.Rollback()
			}
		}
		var arrPsa []*model.MerchantPriceSet

		for _, v := range r.PriceSetArea {
			mps := &model.MerchantPriceSet{
				PriceSet: v.PriceSet,
				Area:     v.Area,
				Merchant: agent,
			}
			arrPsa = append(arrPsa, mps)
			pca[v.Area.ID] = v.PriceSet
		}

		if _, e := o.InsertMulti(100, &arrPsa); e == nil {
			if branchs, _, e := repository.GetBranchsByMerchantIdforUpdate(agent.ID); e == nil {
				for _, v := range branchs {
					if pc, ok := pca[v.Area.ID]; ok {
						branch := &model.Branch{
							ID:            v.ID,
							Salesperson:   r.Salesperson,
							PriceSet:      pc,
							LastUpdatedAt: time.Now(),
							LastUpdatedBy: r.Session.Staff.ID,
						}
						if _, e = o.Update(branch, "salesperson_id", "PriceSet", "LastUpdatedAt", "LastUpdatedBy"); e != nil {
							o.Rollback()
							return nil, e
						}
					}
				}

				if e = log.AuditLogByUser(r.Session.Staff, agent.ID, "agent", "update", r.NotePriceSetChange); e != nil {
					o.Rollback()
					return nil, e
				}
			}
		}
		if r.IsCreateCreditLimitLog == 1 {
			if e = log.CreditLimitLogByStaff(agent, agent.ID, "merchant", r.CreditLimitBefore, r.CreditLimitAfter, r.Session.Staff.ID, "update agent"); e != nil {
				o.Rollback()
				return nil, e
			}
		}
	} else {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return
}

func UpdateTagCustomer(r updateTagRequest) (agent *model.Merchant, e error) {
	agent = &model.Merchant{
		ID:            r.ID,
		TagCustomer:   r.CustomerTagStr,
		LastUpdatedAt: time.Now(),
		LastUpdatedBy: r.Session.Staff.ID,
	}
	if e = agent.Save("tag_customer", "LastUpdatedAt", "LastUpdatedBy"); e == nil {
		if e = log.AuditLogByUser(r.Session.Staff, agent.ID, "agent", "update tag customer", "previous:"+r.Merchant.TagCustomer); e != nil {
			return nil, e
		}
	}

	agent.TagCustomer = util.EncIdInStr(agent.TagCustomer)

	return
}

func UpdatePhoneNumber(r updatePhoneNumber) (agent *model.Merchant, e error) {

	o := orm.NewOrm()
	o.Begin()

	oldPoneNumber := r.Merchant.PhoneNumber
	r.Merchant.PhoneNumber = strings.TrimPrefix(r.PhoneNumber, "0")
	r.Merchant.LastUpdatedAt = time.Now()
	r.Merchant.LastUpdatedBy = r.Session.Staff.ID

	if _, e = o.Update(r.Merchant, "phone_number", "LastUpdatedAt", "LastUpdatedBy"); e == nil {
		r.UserMerchant.ForceLogout = 1
		r.UserMerchant.Verification = 1
		if _, e = o.Update(r.UserMerchant, "ForceLogout", "Verification"); e == nil {
			if e = log.AuditLogByUser(r.Session.Staff, r.Merchant.ID, "agent", "update phone number", "previous: "+oldPoneNumber); e != nil {
				return nil, e
			}
		} else {
			o.Rollback()
			return nil, e
		}

	} else {
		o.Rollback()
		return nil, e
	}
	o.Commit()
	r.Merchant.TagCustomer = util.EncIdInStr(r.Merchant.TagCustomer)

	return r.Merchant, nil
}

func UpdateSalesperson(r updateSalespersonRequest) (branch *model.Branch, e error) {
	branchs, _, e := repository.GetBranchsByMerchantId(r.ID)

	for _, v := range branchs {
		branch = &model.Branch{
			ID:            v.ID,
			Salesperson:   r.Salesperson,
			LastUpdatedAt: time.Now(),
			LastUpdatedBy: r.Session.Staff.ID,
		}

		if e = branch.Save("salesperson_id", "LastUpdatedAt", "LastUpdatedBy"); e != nil {
			return nil, e
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.ID, "agent", "update salesperson", "previous:"+r.PrevSalesperson); e != nil {
		return nil, e
	}

	return
}

func UpdateArchetype(r updateArchetypeRequest) (branch *model.Branch, e error) {
	branchs, _, e := repository.GetBranchsByMerchantId(r.ID)

	for _, v := range branchs {
		branch = &model.Branch{
			ID:            v.ID,
			Archetype:     r.Archetype,
			LastUpdatedAt: time.Now(),
			LastUpdatedBy: r.Session.Staff.ID,
		}

		if e = branch.Save("archetype_id", "LastUpdatedAt", "LastUpdatedBy"); e != nil {
			return nil, e
		}
	}

	merchant := &model.Merchant{
		ID:            branchs[0].Merchant.ID,
		BusinessType:  r.BusinessType,
		LastUpdatedAt: time.Now(),
		LastUpdatedBy: r.Session.Staff.ID,
	}

	if e = merchant.Save("BusinessType", "LastUpdatedAt", "LastUpdatedBy"); e != nil {
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.ID, "agent", "update archetype", "previous:"+r.PrevArchetype); e != nil {
		return nil, e
	}

	return
}

func SaveShippingAddress(r createShippingAddressRequest) (branch *model.Branch, e error) {
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	r.Code, e = util.GenerateCustomerCode(r.Code, "branch")

	var mps []model.MerchantPriceSet
	var ps *model.PriceSet
	orSelect.Raw("select * from merchant_price_set mps where mps.merchant_id = ?", r.Merchant.ID).QueryRows(&mps)

	for _, v := range mps {
		if v.Area.ID == r.Area.ID {
			ps = v.PriceSet
		}
	}

	branch = &model.Branch{
		Code:            r.Code,
		Name:            r.Merchant.Name,
		PicName:         r.RecipientName,
		PhoneNumber:     r.RecipientPhoneNumber,
		AltPhoneNumber:  r.RecipientAltPhoneNumber,
		AddressName:     r.RecipientName,
		ShippingAddress: r.ShippingAddress,
		Note:            r.ShippingNote,
		MainBranch:      int8(2),
		Status:          int8(1),
		Merchant:        r.Merchant,
		Archetype:       r.Archetype,
		PriceSet:        ps,
		Salesperson:     r.Salesperson,
		Area:            r.Area,
		Warehouse:       r.WarehouseCoverage.Warehouse,
		SubDistrict:     r.SubDistrict,
		CreatedAt:       time.Now(),
		CreatedBy:       r.Session.Staff.ID,
	}

	if e = branch.Save(); e == nil {
		if e = log.AuditLogByUser(r.Session.Staff, branch.ID, "shipping_address", "create shipping address", ""); e != nil {
			return nil, e
		}
	}

	return
}

func UpdateShippingAddress(r updateShippingAddressRequest) (branch *model.Branch, e error) {
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	m := &model.Branch{
		ID: r.ID,
	}
	m.Read("ID")
	var mps []model.MerchantPriceSet
	var ps *model.PriceSet
	orSelect.Raw("select * from merchant_price_set mps where mps.merchant_id = ?", m.ID).QueryRows(&mps)

	for _, v := range mps {
		if v.Area.ID == r.Area.ID {
			ps = v.PriceSet
		}
	}
	branch = &model.Branch{
		ID:              r.ID,
		PicName:         r.RecipientName,
		PhoneNumber:     r.RecipientPhoneNumber,
		AltPhoneNumber:  r.RecipientAltPhoneNumber,
		AddressName:     r.RecipientName,
		ShippingAddress: r.ShippingAddress,
		Note:            r.ShippingNote,
		Area:            r.Area,
		Warehouse:       r.WarehouseCoverage.Warehouse,
		SubDistrict:     r.SubDistrict,
		PriceSet:        ps,
		LastUpdatedAt:   time.Now(),
		LastUpdatedBy:   r.Session.Staff.ID,
	}

	if e = branch.Save("pic_name", "phone_number", "alt_phone_number", "address_name", "shipping_address", "note", "area_id", "warehouse_id", "sub_district_id", "LastUpdatedAt", "LastUpdatedBy"); e == nil {
		if e = log.AuditLogByUser(r.Session.Staff, branch.ID, "shipping_address", "update shipping address", ""); e != nil {
			return nil, e
		}
	}

	return
}
