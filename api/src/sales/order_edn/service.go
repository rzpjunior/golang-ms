// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order_edn

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	xendit2 "git.edenfarm.id/project-version2/api/service/xendit"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to insert data requested into database
func Save(r createRequest) (so *model.SalesOrder, e error) {
	r.Code, e = util.GenerateDocCode("SO", r.Branch.Code, "sales_order")
	if e != nil {
		return nil, e
	}

	o := orm.NewOrm()
	o.Begin()

	// insert sales order
	so = &model.SalesOrder{
		Code:            r.Code,
		Branch:          r.Branch,
		SalesTerm:       r.SalesTerm,
		InvoiceTerm:     r.InvoiceTerm,
		Salesperson:     r.Salesperson,
		SalesGroupID:    r.Salesperson.SalesGroupID,
		Warehouse:       r.Warehouse,
		Wrt:             r.Wrt,
		Area:            r.Branch.Area,
		Voucher:         r.Voucher,
		SubDistrict:     r.Branch.SubDistrict,
		PriceSet:        r.Branch.PriceSet,
		PaymentGroup:    r.PaymentGroup,
		Archetype:       r.Branch.Archetype,
		OrderType:       r.OrderType,
		DeliveryDate:    r.DeliveryDate,
		RecognitionDate: r.RecognitionDate,
		BillingAddress:  r.BillingAddress,
		ShippingAddress: r.ShippingAddress,
		OrderChannel:    int8(1),
		HasExtInvoice:   int8(2),
		Note:            r.Note,
		Status:          int8(11),
		TotalPrice:      r.TotalPrice,
		TotalCharge:     r.TotalCharge,
		TotalWeight:     r.TotalWeight,
		CreatedAt:       time.Now(),
		CreatedBy:       r.Session.Staff.ID,
	}
	if r.Voucher != nil {
		so.Voucher = r.Voucher
		so.VouRedeemCode = r.Voucher.RedeemCode
		so.VouDiscAmount = r.Voucher.DiscAmount
	}
	if _, e = o.Insert(so); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = so.Branch.Read("ID"); e != nil {
		o.Rollback()
		return nil, e
	}
	codeDO, _ := util.GenerateDocCode("DO", so.Branch.Code, "delivery_order")
	if e != nil {
		o.Rollback()
		return nil, e
	}

	do := &model.DeliveryOrder{
		Code:            codeDO,
		SalesOrder:      so,
		Warehouse:       r.Warehouse,
		Wrt:             r.Wrt,
		Status:          2,
		RecognitionDate: r.RecognitionDate,
		ShippingAddress: r.ShippingAddress,
		ReceiptNote:     "",
		TotalWeight:     r.TotalWeight,
		DeltaPrint:      0,
		Note:            r.Note,
		CreatedAt:       time.Now(),
		CreatedBy:       r.Session.Staff.ID,
		ConfirmedAt:     time.Now(),
		ConfirmedBy:     r.Session.Staff.ID,
	}

	if _, e = o.Insert(do); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = so.SalesTerm.Read("ID"); e != nil {
		o.Rollback()
		return nil, e
	}

	dueDate, _ := time.Parse("2006-01-02", r.RecognitionDateStr)
	t2 := dueDate.AddDate(0, 0, int(so.SalesTerm.DaysValue))

	//generate codes for document
	codeSI, _ := util.GenerateDocCode("SI", so.Branch.Code, "sales_invoice")

	si := &model.SalesInvoice{
		SalesOrder:        so,
		Code:              codeSI,
		SalesTerm:         so.SalesTerm,
		InvoiceTerm:       so.InvoiceTerm,
		PaymentGroup:      so.PaymentGroup,
		CodeExt:           codeSI, //customer code
		RecognitionDate:   r.RecognitionDate,
		DueDate:           t2,
		BillingAddress:    so.BillingAddress,
		DeliveryFee:       so.DeliveryFee,
		TotalPrice:        r.TotalPrice,
		TotalCharge:       r.TotalCharge,
		Note:              so.Note,
		AdjNote:           "-",
		Status:            1,
		CreatedAt:         time.Now(),
		CreatedBy:         r.Session.Staff.ID,
		PointRedeemAmount: so.PointRedeemAmount,
	}

	if so.Voucher != nil {
		if e = so.Voucher.Read("ID"); e != nil {
			o.Rollback()
			return nil, e
		}
		si.VouRedeemCode = so.Voucher.RedeemCode
		si.VouDiscAmount = so.Voucher.DiscAmount
		si.VoucherID = so.Voucher.ID
	}

	if _, e = o.Insert(si); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, v := range r.Products {
		// insert sales order item
		soi := &model.SalesOrderItem{
			SalesOrder:    so,
			Product:       v.Product,
			OrderQty:      v.Quantity,
			UnitPrice:     float64(v.UnitPrice),
			ShadowPrice:   float64(v.Price.ShadowPrice),
			Subtotal:      v.Subtotal,
			Weight:        v.Weight,
			Note:          v.Note,
			ProductPush:   v.ProductPush,
			TaxableItem:   v.Product.Taxable,
			TaxPercentage: v.Product.TaxPercentage,
			DefaultPrice:  v.DefaultPrice,
		}

		if _, e = o.Insert(soi); e != nil {
			o.Rollback()
			return nil, e
		}

		sii := &model.SalesInvoiceItem{
			SalesInvoice:   &model.SalesInvoice{ID: si.ID},
			Product:        v.Product,
			InvoiceQty:     v.Quantity,
			UnitPrice:      float64(v.UnitPrice),
			Subtotal:       v.Subtotal,
			SalesOrderItem: soi,
			Note:           v.Note,
			TaxableItem:    v.TaxableItem,
			TaxPercentage:  v.TaxPercentage,
		}

		if _, e = o.Insert(sii); e != nil {
			o.Rollback()
			return nil, e
		}

		item := &model.DeliveryOrderItem{
			DeliveryOrder:  &model.DeliveryOrder{ID: do.ID},
			SalesOrderItem: soi,
			Product:        v.Product,
			DeliverQty:     v.Quantity,
			ReceiveQty:     v.Quantity,
			Weight:         v.Weight,
			OrderItemNote:  v.Note,
		}

		if _, e = o.Insert(item); e != nil {
			o.Rollback()
			return nil, e
		}

		stock := &model.Stock{
			Product:   v.Product,
			Warehouse: do.Warehouse,
		}

		if e = stock.Read("Product", "Warehouse"); e != nil {
			o.Rollback()
			return nil, e
		}

		slOut := &model.StockLog{
			Warehouse:    do.Warehouse,
			Product:      v.Product,
			Ref:          do.ID,
			RefType:      1,
			Type:         2,
			InitialStock: stock.AvailableStock,
			Quantity:     item.DeliverQty,
			FinalStock:   stock.AvailableStock - item.DeliverQty,
			UnitCost:     soi.UnitPrice,
			Status:       1,
			DocNote:      do.Note,
			ItemNote:     "SO EDN Sales",
			CreatedAt:    time.Now(),
		}

		if _, e = o.Insert(slOut); e != nil {
			o.Rollback()
			return nil, e
		}

		stock.AvailableStock = slOut.FinalStock

		if _, e = o.Update(stock, "AvailableStock"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	// insert voucher
	if r.Voucher != nil && r.Voucher.ID != 0 {
		r.Voucher.RemOverallQuota = r.Voucher.RemOverallQuota - 1
		if _, e = o.Update(r.Voucher, "rem_overall_quota"); e != nil {
			o.Rollback()
			return nil, e
		}

		vl := &model.VoucherLog{
			Voucher:           r.Voucher,
			Merchant:          r.Branch.Merchant,
			Branch:            r.Branch,
			SalesOrder:        so,
			TagCustomer:       r.SameTagCustomer,
			VoucherDiscAmount: r.Voucher.DiscAmount,
			Timestamp:         time.Now(),
			Status:            int8(1),
		}

		if _, e = o.Insert(vl); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if r.IsCreateCreditLimitLog {
		if e = log.CreditLimitLogByStaff(r.Branch.Merchant, so.ID, "sales_order_edn", r.CreditLimitBefore, r.CreditLimitAfter, r.Session.Staff.ID, "create sales order edn"); e != nil {
			o.Rollback()
			return nil, e
		}
		r.Branch.Merchant.RemainingCreditLimitAmount = r.CreditLimitAfter
		if _, e = o.Update(r.Branch.Merchant, "credit_limit_remaining"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, so.ID, "sales_order", "create", r.NotePriceChange); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, si.ID, "sales_invoice", "create", "order EDN"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, do.ID, "delivery_order", "create & confirm", "order EDN"); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	if r.IsCreateMerchantVa["bca"] == 1 {
		xendit2.BCAXenditFixedVA(r.Branch.Merchant)
	}

	if r.IsCreateMerchantVa["permata"] == 1 {
		xendit2.PermataXenditFixedVA(r.Branch.Merchant)
	}

	return so, e
}
