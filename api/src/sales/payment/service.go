// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package payment

import (
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
)

// Save : function to insert data requested into database
func Save(r createRequest) (sp *model.SalesPayment, e error) {
	o := orm.NewOrm()
	o.Begin()

	if r.Code, e = util.GenerateDocCode("SP", r.SalesInvoice.SalesOrder.Branch.Code, "sales_payment"); e != nil {
		return nil, e
	}

	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	sp = &model.SalesPayment{
		Code:            r.Code,
		RecognitionDate: r.PaymentDate,
		Amount:          r.Amount,
		PaidOff:         r.PaidOff,
		PaymentMethod:   r.PaymentMethod,
		Note:            r.Note,
		SalesInvoice:    r.SalesInvoice,
		BankReceiveNum:  r.BankReceiveNum,
		Status:          int8(2),
		CreatedAt:       time.Now(),
		CreatedBy:       r.Session.Staff.ID,
	}

	if r.PaymentChannelID != "" {
		sp.PaymentChannel = r.PaymentChannel
	}

	if r.ImageUrl != "" {
		sp.ImageUrl = r.ImageUrl
	}

	if _, e = o.Insert(sp); e != nil {
		o.Rollback()
		return
	}

	if _, e = o.Update(r.SalesInvoice, "Status"); e != nil {
		o.Rollback()
		return
	}

	if _, e = o.Update(r.SalesInvoice.SalesOrder, "Status", "FinishedAt"); e != nil {
		o.Rollback()
		return
	}

	if r.HaveCreditLimit {
		if e = log.CreditLimitLogByStaff(r.SalesInvoice.SalesOrder.Branch.Merchant, sp.ID, "sales_payment", r.CreditLimitBefore, r.CreditLimitAfter, r.Session.Staff.ID, "create sales payment"); e != nil {
			o.Rollback()
			return nil, e
		}
		r.SalesInvoice.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount = r.CreditLimitAfter
		if _, e = o.Update(r.SalesInvoice.SalesOrder.Branch.Merchant, "credit_limit_remaining"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, sp.ID, "sales_payment", "create", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	messageNotif := &util.MessageNotification{}

	r.SalesInvoice.SalesOrder.Branch.Merchant.Read("ID")
	r.SalesInvoice.SalesOrder.Branch.Merchant.UserMerchant.Read("ID")
	if r.SalesInvoice.SalesOrder.Status == 13 || r.SalesInvoice.SalesOrder.Status == 12 {
		orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0004'").QueryRow(&messageNotif)
	} else if r.SalesInvoice.SalesOrder.Status == 2 {
		orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0006'").QueryRow(&messageNotif)
	}

	if r.SalesInvoice.SalesOrder.Status == 13 || r.SalesInvoice.SalesOrder.Status == 12 || r.SalesInvoice.SalesOrder.Status == 2 {
		messageNotif.Message = util.ReplaceNotificationSalesOrder(messageNotif.Message, "#sales_order_code#", r.SalesInvoice.SalesOrder.Code)

		mn := &util.ModelNotification{
			SendTo:     r.SalesInvoice.SalesOrder.Branch.Merchant.UserMerchant.FirebaseToken,
			Title:      messageNotif.Title,
			Message:    messageNotif.Message,
			Type:       "1",
			RefID:      r.SalesInvoice.SalesOrder.ID,
			MerchantID: r.SalesInvoice.SalesOrder.Branch.Merchant.ID,
			ServerKey:  util.ServerKeyFireBase,
		}
		util.PostModelNotification(mn)
	}

	//send webhook finished
	if r.SalesInvoice.SalesOrder.OrderChannel == 6 && r.SalesInvoice.SalesOrder.Status == 2 {
		util.SendIDToOapi(common.Encrypt(r.SalesInvoice.SalesOrder.ID), r.Session.Token)
	}

	return sp, e
}

// Cancel : function to change data status into 3
func Cancel(r cancelRequest) (sp *model.SalesPayment, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.SalesPayment.Status = 3
	if _, e = o.Update(r.SalesPayment, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	if _, e = o.Update(r.SalesPayment.SalesInvoice, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	if _, e = o.Update(r.SalesPayment.SalesInvoice.SalesOrder, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	if r.HaveCreditLimit {
		if e = log.CreditLimitLogByStaff(r.SalesPayment.SalesInvoice.SalesOrder.Branch.Merchant, r.SalesPayment.ID, "sales_payment", r.CreditLimitBefore, r.CreditLimitAfter, r.Session.Staff.ID, "cancel sales payment"); e != nil {
			o.Rollback()
			return nil, e
		}
		r.SalesPayment.SalesInvoice.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount = r.CreditLimitAfter
		if _, e = o.Update(r.SalesPayment.SalesInvoice.SalesOrder.Branch.Merchant, "credit_limit_remaining"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.ID, "sales_payment", "cancel", r.Note); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return r.SalesPayment, nil
}

// Cancel Active : function to change data status from 1 into 3
func CancelActive(r cancelActiveRequest) (sp *model.SalesPayment, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.SalesPayment.Status = 3
	if _, e = o.Update(r.SalesPayment, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.ID, "sales_payment", "cancel active", r.Note); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return r.SalesPayment, nil
}

// BulkCreatePayment : function to insert data requested into database
func BulkCreatePayment(r bulkPaymentRequest) (sp *model.SalesPayment, e error) {
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	// item loops
	for _, k := range r.SalesInvoice {
		o := orm.NewOrm()

		o.Begin()

		code, e := util.GenerateDocCode("SP", k.SalesInvoice.SalesOrder.Branch.Code, "sales_payment")
		if e != nil {
			o.Rollback()
			return nil, e
		}

		if k.HaveCreditLimit {
			if k.CreditLimitBefore, e = repository.GetCreditLimitRemainingMerchant(k.SalesInvoice.SalesOrder.Branch.Merchant.ID); e != nil {
				o.Rollback()
				return nil, e
			}
		}

		if k.RemainingInvoiceAmount < 0 {
			k.RemainingInvoiceAmount = 0
		}

		if k.PaidOff == 1 {
			switch k.SalesInvoice.SalesOrder.Status {
			case 9:
				k.SalesInvoice.SalesOrder.Status = 12
			case 10:
				k.SalesInvoice.SalesOrder.Status = 13
			case 11:
				k.SalesInvoice.SalesOrder.Status = 2
				k.SalesInvoice.SalesOrder.FinishedAt = time.Now()
			}

			if k.SalesInvoice.Status == 6 {
				k.CreditLimitAfter = k.CreditLimitBefore + k.RemainingInvoiceAmount
			}

			if k.SalesInvoice.Status == 1 {
				k.CreditLimitAfter = k.CreditLimitBefore + k.SalesInvoice.TotalCharge
			}

			k.SalesInvoice.Status = 2
			k.SalesInvoice.RemainingAmount = 0
		} else {
			k.CreditLimitAfter = k.CreditLimitBefore + k.Amount

			if k.Amount >= k.RemainingInvoiceAmount {
				k.CreditLimitAfter = k.CreditLimitBefore + k.RemainingInvoiceAmount
			}

			if k.Amount >= k.RemainingInvoiceAmount && k.CountInProgressPayment == 0 {
				switch k.SalesInvoice.SalesOrder.Status {
				case 9:
					k.SalesInvoice.SalesOrder.Status = 12
				case 10:
					k.SalesInvoice.SalesOrder.Status = 13
				case 11:
					k.SalesInvoice.SalesOrder.Status = 2
					k.SalesInvoice.SalesOrder.FinishedAt = time.Now()
				}

				k.SalesInvoice.RemainingAmount = 0
				k.SalesInvoice.Status = 2
			} else {

				k.SalesInvoice.RemainingAmount = k.RemainingInvoiceAmount - k.Amount
				k.SalesInvoice.Status = 6
			}
		}

		sp = &model.SalesPayment{
			Code:            code,
			RecognitionDate: r.PaymentDate,
			Amount:          k.Amount,
			PaidOff:         k.PaidOff,
			PaymentMethod:   r.PaymentMethod,
			Note:            k.Note,
			SalesInvoice:    k.SalesInvoice,
			BankReceiveNum:  r.BankReceiveNum,
			Status:          int8(2),
			CreatedAt:       time.Now(),
			CreatedBy:       r.Session.Staff.ID,
		}

		if r.PaymentChannelID != "" {
			sp.PaymentChannel = r.PaymentChannel
		}

		if k.ImageUrl != "" {
			sp.ImageUrl = k.ImageUrl
		}

		if k.HaveCreditLimit {
			if e = log.CreditLimitLogByStaff(k.SalesInvoice.SalesOrder.Branch.Merchant, sp.ID, "sales_payment", k.CreditLimitBefore, k.CreditLimitAfter, r.Session.Staff.ID, "bulk create payment"); e != nil {
				o.Rollback()
				return nil, e
			}
			k.SalesInvoice.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount = k.CreditLimitAfter
			if _, e = o.Update(k.SalesInvoice.SalesOrder.Branch.Merchant, "credit_limit_remaining"); e != nil {
				o.Rollback()
				return nil, e
			}
		}

		if _, e = o.Insert(sp); e != nil {
			o.Rollback()
			return nil, e
		}

		if _, e = o.Update(k.SalesInvoice, "Status"); e != nil {
			o.Rollback()
			return nil, e
		}

		if _, e = o.Update(k.SalesInvoice.SalesOrder, "Status", "FinishedAt"); e != nil {
			o.Rollback()
			return nil, e
		}

		if e = log.AuditLogByUser(r.Session.Staff, sp.ID, "sales_payment", "create", ""); e != nil {
			o.Rollback()
			return nil, e
		}

		o.Commit()

		messageNotif := &util.MessageNotification{}
		k.SalesInvoice.SalesOrder.Branch.Merchant.Read("ID")
		k.SalesInvoice.SalesOrder.Branch.Merchant.UserMerchant.Read("ID")
		if k.SalesInvoice.SalesOrder.Status == 13 || k.SalesInvoice.SalesOrder.Status == 12 {
			orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0004'").QueryRow(&messageNotif)
		} else if k.SalesInvoice.SalesOrder.Status == 2 {
			orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0006'").QueryRow(&messageNotif)
		}

		if k.SalesInvoice.SalesOrder.Status == 13 || k.SalesInvoice.SalesOrder.Status == 12 || k.SalesInvoice.SalesOrder.Status == 2 {
			messageNotif.Message = util.ReplaceNotificationSalesOrder(messageNotif.Message, "#sales_order_code#", k.SalesInvoice.SalesOrder.Code)

			mn := &util.ModelNotification{
				SendTo:     k.SalesInvoice.SalesOrder.Branch.Merchant.UserMerchant.FirebaseToken,
				Title:      messageNotif.Title,
				Message:    messageNotif.Message,
				Type:       "1",
				RefID:      k.SalesInvoice.SalesOrder.ID,
				MerchantID: k.SalesInvoice.SalesOrder.Branch.Merchant.ID,
				ServerKey:  util.ServerKeyFireBase,
			}
			util.PostModelNotification(mn)
		}

		//send webhook finished
		if k.SalesInvoice.SalesOrder.OrderChannel == 6 && k.SalesInvoice.SalesOrder.Status == 2 {
			util.SendIDToOapi(common.Encrypt(k.SalesInvoice.SalesOrder.ID), r.Session.Token)
		}

	}

	return
}

func BulkCreateActivePayment(r bulkCreateActivePaymentRequest) (sp *model.SalesPayment, e error) {

	for _, k := range r.SalesInvoice {

		o := orm.NewOrm()
		o.Begin()

		code, e := util.GenerateDocCode("SP", k.SalesInvoice.SalesOrder.Branch.Code, "sales_payment")

		if e != nil {
			o.Rollback()
			return nil, e
		}

		if k.RemainingInvoiceAmount, e = repository.CheckRemainingSalesInvoiceAmount(k.SalesInvoice.ID); e != nil {
			o.Rollback()
			return nil, e
		}

		if k.RemainingInvoiceAmount < 0 {
			k.RemainingInvoiceAmount = 0
		}

		sp = &model.SalesPayment{
			Code:          code,
			ReceivedDate:  r.ReceivedDate,
			Amount:        k.Amount,
			PaymentMethod: k.PaymentMethod,
			Note:          k.Note,
			SalesInvoice:  k.SalesInvoice,
			Status:        k.PaymentStatus,
			CreatedAt:     time.Now(),
			CreatedBy:     r.Session.Staff.ID,
		}

		if k.ImageUrl != "" {
			sp.ImageUrl = k.ImageUrl
		}

		if _, e = o.Insert(sp); e != nil {
			o.Rollback()
		}

		if e = log.AuditLogByUser(r.Session.Staff, sp.ID, "sales_payment", "create", "bulk_create_active"); e != nil {
			o.Rollback()
			return nil, e
		}

		if e = k.SalesInvoice.SalesOrder.Branch.Merchant.Read("ID"); e != nil {
			return nil, e
		}

		if k.HaveCreditLimit && k.PaymentStatus == 5 {
			k.CreditLimitBefore = k.SalesInvoice.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount
			k.CreditLimitAfter = k.CreditLimitBefore + k.Amount
			if k.Amount > k.RemainingInvoiceAmount {
				k.CreditLimitAfter = k.CreditLimitBefore + k.RemainingInvoiceAmount
			}
			if e = log.CreditLimitLogByStaff(k.SalesInvoice.SalesOrder.Branch.Merchant, sp.ID, "sales_payment", k.CreditLimitBefore, k.CreditLimitAfter, r.Session.Staff.ID, "bulk create active payment"); e != nil {
				o.Rollback()
				return nil, e
			}
			k.SalesInvoice.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount = k.CreditLimitAfter
			if _, e = o.Update(k.SalesInvoice.SalesOrder.Branch.Merchant, "credit_limit_remaining"); e != nil {
				o.Rollback()
				return nil, e
			}
		}
		o.Commit()

	}

	return
}

func BulkConfirmPayment(r bulkConfirmPaymentRequest) (sp *model.SalesPayment, e error) {
	o := orm.NewOrm()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	for _, k := range r.SalesPaymentItems {

		o.Begin()

		sp = k.SalesPayment

		if k.Note != "" {
			sp.Note = k.Note
		}

		// This query 'select' only use in this unique case, in other case you must use 'orm read only' to get data
		if e = o.Raw("SELECT credit_limit_remaining FROM merchant WHERE id = ?", k.SalesInvoice.SalesOrder.Branch.Merchant.ID).QueryRow(k.SalesInvoice.SalesOrder.Branch.Merchant); e != nil {
			o.Rollback()
			return nil, e
		}

		// This query 'select' only use in this unique case, in other case you must use 'orm read only' to get data
		if k.TotalPaidAmount, e = repository.CheckTotalPaidPaymentAmount(k.SalesInvoice.ID); e != nil {
			o.Rollback()
			return nil, e
		}

		k.RemainingInvoiceAmount = k.SalesInvoice.TotalCharge - k.TotalPaidAmount

		if k.RemainingInvoiceAmount < 0 {
			k.RemainingInvoiceAmount = 0
		}

		k.CreditLimitBefore = k.SalesInvoice.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount

		if k.PaidOff == 1 {
			switch k.SalesInvoice.SalesOrder.Status {
			case 9:
				k.SalesInvoice.SalesOrder.Status = 12
			case 10:
				k.SalesInvoice.SalesOrder.Status = 13
			case 11:
				k.SalesInvoice.SalesOrder.Status = 2
				k.SalesInvoice.SalesOrder.FinishedAt = time.Now()
			}

			k.CreditLimitAfter = k.CreditLimitBefore + k.RemainingInvoiceAmount

			k.SalesInvoice.Status = 2

		} else {

			isAmountFinishInvoice := k.Amount >= k.RemainingInvoiceAmount && k.CountInProgressPayment == 0

			k.CreditLimitAfter = k.CreditLimitBefore + k.Amount

			if k.Amount > k.RemainingInvoiceAmount {
				k.CreditLimitAfter = k.CreditLimitBefore + k.RemainingInvoiceAmount
			}

			if k.SalesPayment.Status == 5 {
				diffencePaymentAmount := k.SalesPayment.Amount - k.Amount
				k.CreditLimitAfter = k.CreditLimitBefore - diffencePaymentAmount
				k.RemainingInvoiceAmount += diffencePaymentAmount
				if k.RemainingInvoiceAmount < 0 {
					k.CreditLimitAfter = k.CreditLimitBefore
				}

				isAmountFinishInvoice = k.RemainingInvoiceAmount <= 0 && k.CountInProgressPayment == 1
			}

			if isAmountFinishInvoice {
				switch k.SalesInvoice.SalesOrder.Status {
				case 9:
					k.SalesInvoice.SalesOrder.Status = 12
				case 10:
					k.SalesInvoice.SalesOrder.Status = 13
				case 11:
					k.SalesInvoice.SalesOrder.Status = 2
					k.SalesInvoice.SalesOrder.FinishedAt = time.Now()
				}

				k.SalesInvoice.Status = 2
			} else {
				k.SalesInvoice.Status = 6
			}
		}

		sp.Status = 2
		sp.Amount = k.Amount
		sp.PaidOff = k.PaidOff
		sp.BankReceiveNum = r.BankReceiveNum
		sp.RecognitionDate = r.PaymentDate

		if _, e = o.Update(sp, "Status", "Amount", "PaidOff", "BankReceiveNum", "Note", "RecognitionDate"); e != nil {
			o.Rollback()
			return nil, e
		}

		if _, e = o.Update(k.SalesInvoice, "Status"); e != nil {
			o.Rollback()
			return nil, e
		}

		if _, e = o.Update(k.SalesInvoice.SalesOrder, "Status", "FinishedAt"); e != nil {
			o.Rollback()
			return nil, e
		}

		if e = log.AuditLogByUser(r.Session.Staff, sp.ID, "sales_payment", "confirm", ""); e != nil {
			o.Rollback()
			return nil, e
		}

		messageNotif := &util.MessageNotification{}
		shouldSendANotification := false

		if e = k.SalesInvoice.SalesOrder.Branch.Merchant.UserMerchant.Read("ID"); e != nil {
			o.Rollback()
			return nil, e
		}

		switch k.SalesInvoice.SalesOrder.Status {
		case 2:
			orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0006'").QueryRow(&messageNotif)
			shouldSendANotification = true
		case 12:
			orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0004'").QueryRow(&messageNotif)
			shouldSendANotification = true
		case 13:
			orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0004'").QueryRow(&messageNotif)
			shouldSendANotification = true
		}

		if shouldSendANotification {
			messageNotif.Message = util.ReplaceNotificationSalesOrder(messageNotif.Message, "#sales_order_code#", k.SalesInvoice.SalesOrder.Code)

			mn := &util.ModelNotification{
				SendTo:     k.SalesInvoice.SalesOrder.Branch.Merchant.UserMerchant.FirebaseToken,
				Title:      messageNotif.Title,
				Message:    messageNotif.Message,
				Type:       "1",
				RefID:      k.SalesInvoice.SalesOrder.ID,
				MerchantID: k.SalesInvoice.SalesOrder.Branch.Merchant.ID,
				ServerKey:  util.ServerKeyFireBase,
			}
			util.PostModelNotification(mn)
		}

		//send webhook finished
		if k.SalesInvoice.SalesOrder.OrderChannel == 5 && k.SalesInvoice.SalesOrder.Status == 2 {
			util.SendIDToOapi(common.Encrypt(k.SalesInvoice.SalesOrder.ID), r.Session.Token)
		}

		if k.HaveCreditLimit {
			if e = log.CreditLimitLogByStaff(k.SalesInvoice.SalesOrder.Branch.Merchant, k.SalesPayment.ID, "sales_payment", k.CreditLimitBefore, k.CreditLimitAfter, r.Session.Staff.ID, "confirm bulk payment"); e != nil {
				o.Rollback()
				return nil, e
			}
			k.SalesInvoice.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount = k.CreditLimitAfter
			if _, e = o.Update(k.SalesInvoice.SalesOrder.Branch.Merchant, "credit_limit_remaining"); e != nil {
				o.Rollback()
				return nil, e
			}
		}
		o.Commit()
	}

	return
}

// AddPaymentProof : function to Add Payment Proof ( image_url ) on the Sales Payment
func AddPaymentProof(r addPaymentProofRequest) (sp *model.SalesPayment, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.SalesPayment.ImageUrl = r.ImageUrl
	if _, e = o.Update(r.SalesPayment, "ImageUrl"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.ID, "sales_payment", "add payment proof", r.ImageUrl); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return r.SalesPayment, nil
}
