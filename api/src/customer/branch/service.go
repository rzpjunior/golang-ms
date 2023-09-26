// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package branch

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/tealeg/xlsx"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/service/xendit"
	"git.edenfarm.id/project-version2/api/util"
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
		if _, e = o.Insert(um); e == nil {
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

				if _, e = o.Update(r.ProspectCustomer, "RegStatus"); e == nil {
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

				} else {
					o.Rollback()
					return nil, e
				}
			}

			if _, e = o.Insert(m); e == nil {
				e = log.AuditLogByUser(r.Session.Staff, m.ID, "main_outlet", "create", "")

				if r.IsCreateCreditLimitLog == true {
					if e = log.CreditLimitLogByStaff(m, m.ID, "merchant", r.MerchantCreditLimitBefore, r.MerchantCreditLimitAfter, r.Session.Staff.ID, "create branch"); e != nil {
						o.Rollback()
						return nil, e
					}
				}
			} else {
				o.Rollback()
				return nil, e
			}

			mainBranch = 1
			merchant = m
		} else {
			o.Rollback()
			return nil, e
		}
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

			if _, e = o.Update(r.Merchant, "BusinessType", "FinanceArea", "InvoiceTerm", "PaymentTerm", "PaymentGroup", "BillingAddress", "CustomerGroup", "ReferrerCode", "ReferenceInfo", "Referrer", "Name", "PicName", "PhoneNumber", "AltPhoneNumber", "Email", "TagCustomer", "Note", "UpgradeStatus", "LastUpdatedAt", "LastUpdatedBy", "CustomCreditLimit", "CreditLimitAmount", "RemainingCreditLimitAmount"); e == nil {
				if _, e = o.QueryTable(new(model.Branch)).Filter("merchant_id", r.Merchant.ID).Update(updateParams); e == nil {
					r.ProspectCustomer.RegStatus = 2
					if _, e = o.Update(r.ProspectCustomer, "RegStatus"); e == nil {
						e = log.AuditLogByUser(r.Session.Staff, r.ProspectCustomer.ID, "prospect_customer", "upgrade", "")

						if r.IsCreateCreditLimitLog == true {
							if e = log.CreditLimitLogByStaff(r.Merchant, r.Merchant.ID, "merchant", r.MerchantCreditLimitBefore, r.MerchantCreditLimitAfter, r.Session.Staff.ID, "upgrade branch"); e != nil {
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

					} else {
						o.Rollback()
						return nil, e
					}

					mainBranch = 1
					auditTypes = "main_outlet"
					auditFunction = "upgrade_to_business"
					auditNote = "previous : " + strconv.Itoa(int(r.ProspectCustomer.ID))
				} else {
					o.Rollback()
					return nil, e
				}
			} else {
				o.Rollback()
				return nil, e
			}
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

	if _, e = o.Insert(b); e == nil {
		if auditTypes == "outlet" {
			auditReff = b.ID
		} else {
			auditReff = r.Merchant.ID
		}

		e = log.AuditLogByUser(r.Session.Staff, auditReff, auditTypes, auditFunction, auditNote)

	} else {
		o.Rollback()
		return nil, e
	}

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

// Update : function to update data requested into database
func Update(r updateRequest) (u *model.Branch, e error) {
	u = &model.Branch{
		ID:              r.ID,
		PriceSet:        r.PriceSet,
		PicName:         r.PicName,
		PhoneNumber:     r.PhoneNumber,
		AltPhoneNumber:  r.AltPhoneNumber,
		ShippingAddress: r.ShippingAddress,
		Note:            r.Note,
		SubDistrict:     r.SubDistrict,
		Area:            r.Area,
		Warehouse:       r.Warehouse,
		LastUpdatedAt:   time.Now(),
		LastUpdatedBy:   r.Session.Staff.ID,
	}
	if e = u.Save("PriceSet", "PicName", "PhoneNumber", "AltPhoneNumber", "ShippingAddress", "Note", "SubDistrict", "Warehouse", "LastUpdatedAt", "LastUpdatedBy"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "outlet", "update", r.NotePriceSetChange)
	}

	return
}

// UpdateSalesPerson : function to update data requested into database
func UpdateSalesPerson(r updatesalespersonRequest) (u *model.Branch, e error) {
	prevSalesperson := r.Branch.Salesperson
	u = &model.Branch{
		ID:            r.ID,
		Salesperson:   r.Salesperson,
		LastUpdatedAt: time.Now(),
		LastUpdatedBy: r.Session.Staff.ID,
	}
	if e = u.Save("Salesperson", "LastUpdatedAt", "LastUpdatedBy"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "outlet", "update_salesperson", "previous:"+prevSalesperson.DisplayName)
	}

	return
}

// ConvertArchetype : function to update data requested into database
func ConvertArchetype(r convertarchetypeRequest) (u *model.Branch, e error) {
	prevArchetype := r.Branch.Archetype.Name

	u = &model.Branch{
		ID:            r.ID,
		Archetype:     r.Archetype,
		LastUpdatedAt: time.Now(),
		LastUpdatedBy: r.Session.Staff.ID,
	}
	if e = u.Save("Archetype", "LastUpdatedAt", "LastUpdatedBy"); e == nil {
		r.Branch.Merchant.CustomerGroup = r.Archetype.CustomerGroup
		if e = r.Branch.Merchant.Save("CustomerGroup"); e == nil {
			e = log.AuditLogByUser(r.Session.Staff, u.ID, "outlet", "convert_archetype", "previous:"+prevArchetype)
		}
	}

	archetype, e := repository.GetArchetype("id", r.Archetype.ID)
	if e != nil {
		return nil, e
	}

	branch, e := repository.GetBranch("id", r.ID)
	if e != nil {
		return nil, e
	}

	m := &model.Merchant{
		ID:            branch.Merchant.ID,
		BusinessType:  archetype.BusinessType,
		LastUpdatedAt: time.Now(),
		LastUpdatedBy: r.Session.Staff.ID,
	}
	if e = m.Save("BusinessType", "LastUpdatedAt", "LastUpdatedBy"); e != nil {
		return nil, e
	}

	return
}

// Archive : function to update status data into archive
func Archive(r archiveRequest) (branch *model.Branch, e error) {
	branch = &model.Branch{
		ID:     r.ID,
		Status: int8(2),
	}

	if e = branch.Save("Status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, branch.ID, "outlet", "archive", "")

		if r.ArchiveMerchant == 1 {
			r.Merchant.Status = int8(2)

			if e = r.Merchant.Save("Status"); e == nil {
				e = log.AuditLogByUser(r.Session.Staff, r.Merchant.ID, "main_outlet", "archive", "")

				userMerchant := &model.UserMerchant{
					ID:     r.Merchant.UserMerchant.ID,
					Status: int8(2),
				}

				if e = userMerchant.Save("Status"); e == nil {
					e = log.AuditLogByUser(r.Session.Staff, userMerchant.ID, "user_merchant", "archive", "")
				}
			}
		}
	}
	return
}

// Unarchive : function to update status data into active
func Unarchive(r unarchiveRequest) (branch *model.Branch, e error) {
	branch = &model.Branch{
		ID:     r.ID,
		Status: int8(1),
	}

	if e = branch.Save("Status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, branch.ID, "outlet", "unarchive", "")

		if r.Merchant.Status == 2 {
			r.Merchant.Status = int8(1)

			if e = r.Merchant.Save("Status"); e == nil {
				e = log.AuditLogByUser(r.Session.Staff, r.Merchant.ID, "main_outlet", "unarchive", "")

				userMerchant := &model.UserMerchant{
					ID:     r.Merchant.UserMerchant.ID,
					Status: int8(1),
				}

				if userMerchant.Save("Status"); e == nil {
					e = log.AuditLogByUser(r.Session.Staff, userMerchant.ID, "user_merchant", "unarchive", "")
				}
			}
		}
	}

	return branch, e
}

// getBranchFilterBySalespersonXls : function to create excel file of task assignment
func getBranchFilterBySalespersonXls(data []*templateBranchBySalesperson) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := env.GetString("EXPORT_DIRECTORY", "")

	filename := fmt.Sprintf("SalespersonUpdate%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))
	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Sales Group"
		row.AddCell().Value = "Outlet_Code*"
		row.AddCell().Value = "Outlet Name"
		row.AddCell().Value = "Staff Code"
		row.AddCell().Value = "Staff Name"
		row.AddCell().Value = "New_Sales_Person_Code*"
		row.AddCell().Value = "New Sales Person Name"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.AreaName       // Area Name
			row.AddCell().Value = v.SalesGroupName // Sales Group
			row.AddCell().Value = v.BranchCode     // Outlet Code
			row.AddCell().Value = v.BranchName     // Outlet Name
			row.AddCell().Value = v.StaffCode      // Staff Code
			row.AddCell().Value = v.StaffName      // Staff Name
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// UpdateBulkSalesPerson : function to update data requested into database
func UpdateBulkSalesPerson(r updateBulkSalespersonReq) (e error) {
	for _, v := range r.Data {
		prevSalesperson := v.Branch.Salesperson
		newData := &model.Branch{
			ID:            v.Branch.ID,
			Salesperson:   v.NewSalesPerson,
			LastUpdatedAt: time.Now(),
			LastUpdatedBy: r.Session.Staff.ID,
		}
		if e = newData.Save("Salesperson", "LastUpdatedAt", "LastUpdatedBy"); e == nil {
			e = log.AuditLogByUser(r.Session.Staff, v.Branch.ID, "outlet", "update_salesperson", "previous:"+prevSalesperson.DisplayName)
		} else {
			return e
		}
	}

	return e
}
