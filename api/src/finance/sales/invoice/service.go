// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package invoice

import (
	"math"
	"time"

	"git.edenfarm.id/cuxs/orm"

	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (si *model.SalesInvoice, e error) {
	docCode, e := util.GenerateDocCode("SI", r.SalesOrder.Branch.Code, "sales_invoice")
	if e != nil {
		return nil, e
	}

	o := orm.NewOrm()
	o.Begin()

	if e != nil {
		return nil, e
	}

	si = &model.SalesInvoice{
		Code:               docCode,
		SalesOrder:         r.SalesOrder,
		PaymentGroup:       r.PaymentGroup,
		SalesTerm:          r.SalesPaymentTerm,
		InvoiceTerm:        r.InvoiceTerm,
		RecognitionDate:    r.RecognitionDate,
		Status:             1,
		DueDate:            r.DueDate,
		BillingAddress:     r.BillingAddress,
		DeliveryFee:        r.DeliveryFee,
		Adjustment:         r.Adjustment,
		AdjAmount:          math.Abs(r.AdjustmentAmount),
		AdjNote:            r.AdjustmentNote,
		TotalPrice:         r.TotalPrice,
		TotalCharge:        r.TotalCharge,
		Note:               r.Note,
		CodeExt:            docCode,
		DeltaPrint:         0,
		CreatedAt:          time.Now(),
		CreatedBy:          r.Session.Staff.ID,
		PointRedeemAmount:  r.SalesOrder.PointRedeemAmount,
		TotalSkuDiscAmount: r.TotalSkuDiscAmount,
	}

	if r.SalesOrder.VouRedeemCode != "" {
		si.VoucherID = r.SalesOrder.Voucher.ID
		si.VouRedeemCode = r.SalesOrder.Voucher.RedeemCode
		si.VouDiscAmount = r.SalesOrder.Voucher.DiscAmount
	}

	if _, e = o.Insert(si); e != nil {
		o.Rollback()
		return nil, e
	}

	var arrSii []*model.SalesInvoiceItem
	for _, row := range r.InvoiceItems {
		item := &model.SalesInvoiceItem{
			SalesInvoice:   &model.SalesInvoice{ID: si.ID},
			Product:        row.Product,
			InvoiceQty:     row.InvoiceQty,
			UnitPrice:      row.UnitPrice,
			Subtotal:       row.Subtotal,
			SalesOrderItem: row.SalesOrderItem,
			Note:           row.Note,
			TaxableItem:    row.TaxableItem,
			TaxPercentage:  row.TaxPercentage,
			SkuDiscAmount:  row.SkuDiscAmount,
		}
		arrSii = append(arrSii, item)
		si.SalesInvoiceItems = append(si.SalesInvoiceItems, item)
	}
	if _, e = o.InsertMulti(100, &arrSii); e != nil {
		o.Rollback()
		return nil, e
	}

	so := &model.SalesOrder{
		ID: r.SalesOrder.ID,
	}

	switch r.SalesOrder.Status {
	case 1:
		so.Status = 9
	case 7:
		so.Status = 10
	case 8:
		so.Status = 11
	}

	if _, e = o.Update(so, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	if r.IsCreateCreditLimitLog == 1 {

		if e = log.CreditLimitLogByStaff(r.SalesOrder.Branch.Merchant, si.ID, "sales_invoice", r.CreditLimitBefore, r.CreditLimitAfter, r.Session.Staff.ID, "create sales invoice"); e != nil {
			o.Rollback()
			return nil, e
		}
		r.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount = r.CreditLimitAfter
		if _, e = o.Update(r.SalesOrder.Branch.Merchant, "credit_limit_remaining"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, si.ID, "sales_invoice", "create", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return si, e
}

// Update : function to update data requested into database
func Update(u updateRequest) (si *model.SalesInvoice, e error) {
	o := orm.NewOrm()
	o.Begin()

	si = &model.SalesInvoice{
		ID:                 u.ID,
		RecognitionDate:    u.RecognitionAt,
		DueDate:            u.DueDateAt,
		BillingAddress:     u.BillingAddress,
		DeliveryFee:        u.DeliveryFee,
		Adjustment:         u.Adjustment,
		AdjAmount:          math.Abs(u.AdjustmentAmount),
		AdjNote:            u.AdjustmentNote,
		TotalPrice:         u.TotalPrice,
		TotalCharge:        u.TotalCharge,
		Note:               u.Note,
		CodeExt:            u.CodeExt,
		LastUpdatedAt:      time.Now(),
		LastUpdatedBy:      u.Session.Staff.ID,
		TotalSkuDiscAmount: u.TotalSkuDiscAmount,
	}
	if _, e = o.Update(si, "RecognitionDate", "DueDate", "BillingAddress", "DeliveryFee", "Adjustment", "AdjAmount", "AdjNote", "TotalPrice", "TotalCharge", "Note", "CodeExt", "LastUpdatedAt", "LastUpdatedBy", "TotalSkuDiscAmount"); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, row := range u.InvoiceItems {
		item := &model.SalesInvoiceItem{
			ID:            row.SalesInvoiceItem.ID,
			SalesInvoice:  si,
			Product:       row.Product,
			InvoiceQty:    row.InvoiceQty,
			UnitPrice:     row.UnitPrice,
			Subtotal:      row.Subtotal,
			Note:          row.Note,
			SkuDiscAmount: row.SkuDiscAmount,
		}
		if _, e = o.Update(item, "Product", "InvoiceQty", "UnitPrice", "Subtotal", "Note", "SkuDiscAmount"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if u.IsCreateCreditLimitLog == 1 {

		if e = log.CreditLimitLogByStaff(u.SalesOrder.Branch.Merchant, si.ID, "sales_invoice", u.CreditLimitBefore, u.CreditLimitAfter, u.Session.Staff.ID, "update sales invoice"); e != nil {
			o.Rollback()
			return nil, e
		}
		u.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount = u.CreditLimitAfter
		if _, e = o.Update(u.SalesOrder.Branch.Merchant, "credit_limit_remaining"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if e = log.AuditLogByUser(u.Session.Staff, si.ID, "sales_invoice", "update"); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return si, e
}

// Cancel : function to update status data into cancelled
func Cancel(r cancelRequest) (u *model.SalesInvoice, e error) {
	o := orm.NewOrm()
	o.Begin()

	u = &model.SalesInvoice{
		ID:     r.ID,
		Status: 3,
	}

	if _, e = o.Update(u, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	so := &model.SalesOrder{
		ID: r.SalesInvoice.SalesOrder.ID,
	}

	switch r.SalesInvoice.SalesOrder.Status {
	case 9:
		so.Status = 1
	case 10:
		so.Status = 7
	case 11:
		so.Status = 8
	}

	if _, e = o.Update(so, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	if r.IsCreateCreditLimitLog == 1 {

		if e = log.CreditLimitLogByStaff(r.SalesInvoice.SalesOrder.Branch.Merchant, u.ID, "sales_invoice", r.CreditLimitBefore, r.CreditLimitAfter, r.Session.Staff.ID, "cancel sales invoice"); e != nil {
			o.Rollback()
			return nil, e
		}
		r.SalesInvoice.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount = r.CreditLimitAfter
		if _, e = o.Update(r.SalesInvoice.SalesOrder.Branch.Merchant, "credit_limit_remaining"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, u.ID, "sales_invoice", "cancel", r.CancellationNote); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return u, e
}

func Print(si *model.SalesInvoice) (returnSI *model.SalesInvoice, e error) {
	si.DeltaPrint = si.DeltaPrint + 1
	si.Save("DeltaPrint")
	return si, e
}
