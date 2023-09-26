// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package supplier_return

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"time"
)

// Save : function to insert data requested into database
func Save(r createRequest) (sr *model.SupplierReturn, e error) {
	r.Code, e = util.GenerateDocCode("SR", r.Supplier.Code, "supplier_return")
	o := orm.NewOrm()
	o.Begin()

	var arrSupplierReturn []*model.SupplierReturnItem

	sr = &model.SupplierReturn{
		Code:            r.Code,
		RecognitionDate: r.RecognitionDateAt,
		Warehouse:       r.Warehouse,
		GoodsReceipt:    r.GoodsReceipt,
		Supplier:        r.Supplier,
		Note:            r.Note,
		Status:          1,
		CreatedAt:       time.Now(),
		CreatedBy:       r.Session.Staff,
	}

	if _, e = o.Insert(sr); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, v := range r.SupplierReturnItem {
		sri := &model.SupplierReturnItem{
			SupplierReturn: sr,
			Product:        v.Product,
			ReturnGoodsQty: v.ReturnGoodQty,
			ReceivedQty:    v.ReceivedQty,
			Note:           v.Note,
		}

		arrSupplierReturn = append(arrSupplierReturn, sri)
	}

	if _, e = o.InsertMulti(100, &arrSupplierReturn); e != nil {
		o.Rollback()
		return nil, e
	}

	gr := orm.Params{
		"valid_supplier_return": 1,
	}

	if _, e = o.QueryTable(new(model.GoodsReceipt)).Filter("id", r.GoodsReceipt.ID).Update(gr); e != nil {
		o.Rollback()
		return nil, e
	}

	// region insert to debit note
	var (
		debitNoteCode    string
		dn               *model.DebitNote
		arrDebitNoteItem []*model.DebitNoteItem
	)
	debitNoteCode, e = util.GenerateDocCode("DN", r.Supplier.Code, "debit_note")

	dn = &model.DebitNote{
		Code:                  debitNoteCode,
		SupplierReturn:        sr,
		RecognitionDate:       r.RecognitionDateAt,
		Note:                  r.Note,
		Status:                1,
		CreatedAt:             time.Now(),
		CreatedBy:             r.Session.Staff,
		UsedInPurchaseInvoice: 2,
	}

	if _, e = o.Insert(dn); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, v := range r.SupplierReturnItem {
		var unitPrice float64

		if v.GoodsReceiptItem.PurchaseOrderItem.TaxAmount == 0 {
			unitPrice = v.GoodsReceiptItem.PurchaseOrderItem.UnitPrice
		} else {
			unitPrice = v.GoodsReceiptItem.PurchaseOrderItem.UnitPriceTax
		}
		dni := &model.DebitNoteItem{
			DebitNote:   dn,
			Product:     v.Product,
			Note:        v.Note,
			UnitPrice:   unitPrice,
			Subtotal:    unitPrice * v.ReturnGoodQty,
			ReturnQty:   v.ReturnGoodQty,
			ReceivedQty: v.ReceivedQty,
		}

		r.TotalPrice += unitPrice * v.ReturnGoodQty

		arrDebitNoteItem = append(arrDebitNoteItem, dni)
	}

	dn.TotalPrice = r.TotalPrice
	if _, e = o.Update(dn, "TotalPrice"); e != nil {
		o.Rollback()
		return nil, e
	}

	if _, e = o.InsertMulti(100, &arrDebitNoteItem); e != nil {
		o.Rollback()
		return nil, e
	}
	e = log.AuditLogByUser(r.Session.Staff, dn.ID, "debit_note", "create", "")
	// endregion

	e = log.AuditLogByUser(r.Session.Staff, sr.ID, "supplier_return", "create", "")

	o.Commit()

	return
}

func Update(r updateRequest) (sr *model.SupplierReturn, e error) {
	o := orm.NewOrm()
	o.Begin()
	var isItemCreated bool
	var keepItemsId []int64
	var keepItemsDebitNoteId []int64

	r.SupplierReturn.Note = r.Note
	if _, e = o.Update(r.SupplierReturn, "Note"); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, v := range r.SupplierReturnItem {
		var unitPrice float64

		if v.GoodsReceiptItem.PurchaseOrderItem.TaxAmount == 0 {
			unitPrice = v.GoodsReceiptItem.PurchaseOrderItem.UnitPrice
		} else {
			unitPrice = v.GoodsReceiptItem.PurchaseOrderItem.UnitPriceTax
		}
		sri := &model.SupplierReturnItem{
			SupplierReturn: r.SupplierReturn,
			Product:        v.Product,
			ReceivedQty:    v.ReceivedQty,
			ReturnGoodsQty: v.ReturnGoodQty,
			Note:           v.Note,
		}

		if isItemCreated, sri.ID, e = o.ReadOrCreate(sri, "SupplierReturn", "Product"); e != nil {
			o.Rollback()
			return nil, e
		}

		if !isItemCreated {
			sri := &model.SupplierReturnItem{
				ID:             sri.ID,
				ReturnGoodsQty: v.ReturnGoodQty,
				ReceivedQty:    v.ReceivedQty,
				Note:           v.Note,
			}
			if _, e = o.Update(sri, "ReturnGoodsQty", "ReceivedQty", "Note"); e != nil {
				o.Rollback()
				return nil, e
			}
		}

		keepItemsId = append(keepItemsId, sri.ID)

		// region debit note item
		debitNoteItem := &model.DebitNoteItem{
			DebitNote:   r.DebitNote,
			Product:     v.Product,
			Note:        v.Note,
			UnitPrice:   unitPrice,
			Subtotal:    unitPrice * v.ReturnGoodQty,
			ReturnQty:   v.ReturnGoodQty,
			ReceivedQty: v.ReceivedQty,
		}
		r.TotalPrice += unitPrice * v.ReturnGoodQty
		if debitNoteItem.ID, e = o.InsertOrUpdate(debitNoteItem); e != nil {
			o.Rollback()
			return nil, e
		}
		keepItemsDebitNoteId = append(keepItemsDebitNoteId, debitNoteItem.ID)
		// endregion

	}

	// region update debit note
	r.DebitNote.TotalPrice = r.TotalPrice
	r.DebitNote.Note = r.Note
	if _, e = o.Update(r.DebitNote, "TotalPrice", "Note"); e != nil {
		o.Rollback()
		return nil, e
	}
	// endregion

	if _, e = o.QueryTable(new(model.SupplierReturnItem)).Filter("supplier_return_id", r.SupplierReturn.ID).Exclude("ID__in", keepItemsId).Delete(); e != nil {
		o.Rollback()
		return nil, e
	}

	if _, e = o.QueryTable(new(model.DebitNoteItem)).Filter("debit_note_id", r.DebitNote.ID).Exclude("ID__in", keepItemsDebitNoteId).Delete(); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.SupplierReturn.ID, "supplier_return", "update", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.DebitNote.ID, "debit_note", "update", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return r.SupplierReturn, nil
}

func Confirm(r confirmRequest) (sr *model.SupplierReturn, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.SupplierReturn.Status = 2
	r.SupplierReturn.ConfirmedAt = time.Now()
	r.SupplierReturn.ConfirmedBy = r.Session.Staff
	if _, e = o.Update(r.SupplierReturn, "Status", "ConfirmedAt", "ConfirmedBy"); e != nil {
		o.Rollback()
		return nil, e
	}
	for _, v := range r.SupplierReturnItem {
		s := &model.Stock{
			Product:   v.Product,
			Warehouse: r.Warehouse,
		}
		s.Read("Product", "Warehouse")

		sl := &model.StockLog{
			Warehouse:    r.Warehouse,
			Product:      v.Product,
			Ref:          r.SupplierReturn.ID,
			RefType:      8,
			Type:         2,
			InitialStock: s.AvailableStock,
			Quantity:     v.ReturnGoodQty,
			FinalStock:   s.AvailableStock - v.ReturnGoodQty,
			UnitCost:     0,
			Status:       1,
			DocNote:      r.SupplierReturn.Note,
			ItemNote:     v.Note,
			CreatedAt:    time.Now(),
		}

		if _, e = o.Insert(sl); e != nil {
			o.Rollback()
			return nil, e
		}

		s.AvailableStock = sl.FinalStock
		if _, e = o.Update(s, "AvailableStock"); e != nil {
			o.Rollback()
			return nil, e
		}

	}

	e = log.AuditLogByUser(r.Session.Staff, r.SupplierReturn.ID, "supplier_return", "confirm", "")

	o.Commit()
	return r.SupplierReturn, nil
}

func Cancel(r cancelRequest) (sr *model.SupplierReturn, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.SupplierReturn.Status = 3
	if _, e = o.Update(r.SupplierReturn, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	r.DebitNote.Status = 3
	if _, e = o.Update(r.DebitNote, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	gr := orm.Params{
		"valid_supplier_return": 2,
	}

	if _, e = o.QueryTable(new(model.GoodsReceipt)).Filter("id", r.SupplierReturn.GoodsReceipt.ID).Update(gr); e != nil {
		o.Rollback()
		return nil, e
	}

	// region audit log supplier return
	e = log.AuditLogByUser(r.Session.Staff, r.SupplierReturn.ID, "supplier_return", "cancel", r.CancellationNote)
	// endregion

	// region audit log debit note
	e = log.AuditLogByUser(r.Session.Staff, r.DebitNote.ID, "debit_note", "cancel", r.CancellationNote)
	// endregion

	o.Commit()
	return r.SupplierReturn, nil
}
