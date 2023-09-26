package distribution_network

import (
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/service/xendit"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (b *model.Branch, e error) {
	o := orm.NewOrm()
	o.Begin()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	mainBranch := int8(2)
	merchant := r.Merchant
	auditTypes := "outlet"
	auditFunction := "create"
	auditNote := ""
	auditReff := int64(0)

	if r.NewMerchantCheck == "true" {
		r.CodeUserMerchant, e = util.GenerateCode(r.CodeUserMerchant, "user_merchant")
		um := &model.UserMerchant{
			Code:          r.CodeUserMerchant,
			FirebaseToken: "-",
			LoginToken:    "-",
			Verification:  int8(1),
			Status:        int8(1),
		}
		if _, e = o.Insert(um); e != nil {
			o.Rollback()
			return nil, e
		}

		r.CodeMerchant, e = util.GenerateCustomerCode(r.CodeMerchant, "merchant")
		codeReferral, e := util.GenerateCodeReferral()
		m := &model.Merchant{
			Code:                       r.CodeMerchant,
			UserMerchant:               um,
			BusinessType:               r.MerchantBusinessType,
			BusinessTypeCreditLimit:    r.BusinessTypeCreditLimit,
			Name:                       r.MerchantName,
			PicName:                    r.MerchantPicName,
			PhoneNumber:                strings.TrimPrefix(r.MerchantPhoneNumber, "0"),
			AltPhoneNumber:             r.MerchantAltPhoneNumber,
			Email:                      r.MerchantEmail,
			Password:                   r.MerchantPassword,
			BillingAddress:             r.BillingAddress,
			Note:                       r.MerchantNote,
			FinanceArea:                r.FinanceArea,
			InvoiceTerm:                r.InvoiceTerm,
			PaymentTerm:                r.PaymentTerm,
			PaymentGroup:               r.PaymentGroup,
			CustomerGroup:              r.BranchArchetype.CustomerGroup,
			Status:                     int8(1),
			ReferralCode:               codeReferral,
			CreatedAt:                  time.Now(),
			CreatedBy:                  r.Session.Staff.ID,
			Suspended:                  2,
			CustomCreditLimit:          2,
			CreditLimitAmount:          r.MerchantCreditLimitAmount,
			RemainingCreditLimitAmount: r.MerchantCreditLimitAmount,
			KTPPhotosUrl:               r.KTPPhotosStr,
			MerchantPhotosUrl:          r.MerchantPhotosStr,
		}

		if r.ReferrerCode != "" {
			m.Referrer = r.Referrer
			m.ReferrerCode = r.ReferrerCode
		}

		if r.ReferenceInfo != "" {
			m.ReferenceInfo = r.ReferenceInfo
		}

		if len(r.CustomerTag) > 0 {
			m.TagCustomer = r.CustomerTagStr
		}

		if r.ProspectCustomer != nil {
			m.ProspectCustomer = r.ProspectCustomer
			r.ProspectCustomer.RegStatus = 2

			if _, e = o.Update(r.ProspectCustomer, "RegStatus"); e != nil {
				o.Rollback()
				return nil, e
			}
			e = log.AuditLogByUser(r.Session.Staff, r.ProspectCustomer.ID, "prospect_customer", "register", "")

			pr := &model.ProspectCustomer{
				ID:          r.ProspectCustomer.ID,
				ProcessedAt: time.Now(),
				ProcessedBy: r.Session.Staff.ID,
			}

			if _, e = o.Update(pr, "ProcessedAt", "ProcessedBy"); e != nil {
				o.Rollback()
				return nil, e
			}

		}

		if _, e = o.Insert(m); e != nil {
			o.Rollback()
			return nil, e
		}
		e = log.AuditLogByUser(r.Session.Staff, m.ID, "distribution_network", "create", "")

		if r.IsCreateCreditLimitLog {
			if e = log.CreditLimitLogByMerchant(m, m.ID, "merchant", r.MerchantCreditLimitBefore, r.MerchantCreditLimitAfter, "create branch"); e != nil {
				o.Rollback()
				return nil, e
			}
		}

		mainBranch = 1
		merchant = m

	} else {
		if r.ProspectCustomerID != "" && r.Merchant.UpgradeStatus == 1 {
			updateParams := orm.Params{"status": 3, "main_branch": 0}

			r.Merchant.UpgradeStatus = 2
			r.Merchant.BusinessType = r.MerchantBusinessType
			r.Merchant.LastUpdatedAt = time.Now()
			r.Merchant.LastUpdatedBy = r.Session.Staff.ID
			r.Merchant.CustomCreditLimit = 2
			r.Merchant.CreditLimitAmount = r.MerchantCreditLimitAmount
			r.Merchant.RemainingCreditLimitAmount = r.MerchantCreditLimitAfter

			if _, e = o.Update(r.Merchant, "BusinessType", "FinanceArea", "InvoiceTerm", "PaymentTerm", "PaymentGroup", "BillingAddress", "CustomerGroup", "ReferrerCode", "ReferenceInfo", "Referrer", "Name", "PicName", "PhoneNumber", "AltPhoneNumber", "Email", "TagCustomer", "Note", "UpgradeStatus", "LastUpdatedAt", "LastUpdatedBy", "CustomCreditLimit", "CreditLimitAmount", "RemainingCreditLimitAmount"); e != nil {
				o.Rollback()
				return nil, e
			}

			if _, e = o.QueryTable(new(model.Branch)).Filter("merchant_id", r.Merchant.ID).Update(updateParams); e != nil {
				o.Rollback()
				return nil, e
			}
			r.ProspectCustomer.RegStatus = 2
			if _, e = o.Update(r.ProspectCustomer, "RegStatus"); e != nil {
				o.Rollback()
				return nil, e
			}
			e = log.AuditLogByUser(r.Session.Staff, r.ProspectCustomer.ID, "prospect_customer", "upgrade", "")

			if r.IsCreateCreditLimitLog {
				if e = log.CreditLimitLogByMerchant(r.Merchant, r.Merchant.ID, "merchant", r.MerchantCreditLimitBefore, r.MerchantCreditLimitAfter, "upgrade branch"); e != nil {
					o.Rollback()
					return nil, e
				}
			}

			pr := &model.ProspectCustomer{
				ID:          r.ProspectCustomer.ID,
				ProcessedAt: time.Now(),
				ProcessedBy: r.Session.Staff.ID,
			}

			if _, e = o.Update(pr, "ProcessedAt", "ProcessedBy"); e != nil {
				o.Rollback()
				return nil, e
			}

			mainBranch = 1
			auditTypes = "main_outlet"
			auditFunction = "upgrade_to_business"
			auditNote = "previous : " + strconv.Itoa(int(r.ProspectCustomer.ID))
		}
	}

	r.CodeBranch, e = util.GenerateCustomerCode(r.CodeBranch, "branch")
	b = &model.Branch{
		Merchant:        merchant,
		Area:            r.BranchArea,
		Archetype:       r.BranchArchetype,
		PriceSet:        r.BranchPriceSet,
		Warehouse:       r.Warehouse,
		Salesperson:     r.BranchSalesPerson,
		SubDistrict:     r.SubDistrict,
		Code:            r.CodeBranch,
		Name:            r.BranchName,
		PicName:         r.BranchPicName,
		PhoneNumber:     r.BranchPhoneNumber,
		AltPhoneNumber:  r.BranchAltPhoneNumber,
		AddressName:     r.BranchName,
		ShippingAddress: r.BranchShippingAddress,
		Note:            r.BranchNote,
		MainBranch:      mainBranch,
		Status:          int8(1),
		CreatedAt:       time.Now(),
		CreatedBy:       r.Session.Staff.ID,
	}

	if _, e = o.Insert(b); e != nil {
		o.Rollback()
		return nil, e
	}

	if auditTypes == "outlet" {
		auditReff = b.ID
	} else {
		auditReff = r.Merchant.ID
	}

	e = log.AuditLogByUser(r.Session.Staff, auditReff, auditTypes, auditFunction, auditNote)

	o.Commit()

	if r.ProspectCustomerID != "" && r.ProspectCustomer.SalespersonID != 0 {
		// notification FS Apps when Register Prospective Customer
		messageNotif := &util.MessageNotification{}
		orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0017'").QueryRow(&messageNotif)
		messageNotif.Message = util.ReplaceNotificationSalesOrder(messageNotif.Message, "#Name#", r.MerchantName)

		mn := &util.ModelNotification{
			SendTo:    r.Staff.User.SalesAppNotifToken,
			Title:     messageNotif.Title,
			Message:   messageNotif.Message,
			Type:      "2",
			RefID:     r.ProspectCustomer.ID,
			ServerKey: util.FieldSalesServerKeyFireBase,
			StaffID:   r.Staff.ID,
		}
		util.PostModelNotificationFieldSales(mn)
	}

	if e == nil && r.NewMerchantCheck == "true" {
		xendit.BCAXenditFixedVA(merchant) // ini ke package xendit internal
		xendit.PermataXenditFixedVA(merchant)
	}

	return b, e
}

// Update : function to update data of merchant
func Update(r updateRequest) (u *model.Merchant, e error) {
	u = &model.Merchant{
		ID:                      r.ID,
		PicName:                 r.PicName,
		AltPhoneNumber:          r.AltPhoneNumber,
		Email:                   r.Email,
		Note:                    r.Note,
		InvoiceTerm:             r.InvoiceTerm,
		PaymentTerm:             r.PaymentTerm,
		PaymentGroup:            r.PaymentGroup,
		BillingAddress:          r.BillingAddress,
		LastUpdatedAt:           time.Now(),
		LastUpdatedBy:           r.Session.Staff.ID,
		BusinessTypeCreditLimit: r.BusinessTypeCreditLimit,
		CustomCreditLimit:       r.CustomCreditLimit,
		CreditLimitAmount:       r.CreditLimitAmount,
		KTPPhotosUrl:            r.KTPPhotosStr,
		MerchantPhotosUrl:       r.MerchantPhotosStr,
	}

	if e = u.Save("PicName", "AltPhoneNumber", "Email", "Note", "InvoiceTerm", "PaymentTerm", "PaymentGroup", "BillingAddress", "LastUpdatedAt", "LastUpdatedBy", "BusinessTypeCreditLimit", "CustomCreditLimit", "CreditLimitAmount", "KTPPhotosUrl", "MerchantPhotosUrl"); e != nil {
		return nil, e
	}

	if r.IsCreateCreditLimitLog == 1 {
		if e = log.CreditLimitLogByMerchant(u, u.ID, "merchant", r.CreditLimitBefore, r.CreditLimitAfter, "update merchant"); e != nil {
			return nil, e
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, u.ID, "merchant", "update", ""); e != nil {
		return nil, e
	}

	return
}

func resetPassword(r resetPasswordRequest) (err error) {
	o := orm.NewOrm()
	o.Begin()

	// update last login dari user tersebut
	r.Merchant.Password = r.PasswordHash
	if _, err = o.Update(r.Merchant, "Password"); err != nil {
		o.Rollback()
		return err
	}

	if err = log.AuditLogByUser(r.Session.Staff, r.Merchant.ID, "merchant", "reset password EDN", ""); err != nil {
		return err
	}

	o.Commit()
	return nil
}
