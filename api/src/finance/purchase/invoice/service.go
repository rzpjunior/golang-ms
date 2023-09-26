package invoice

import (
	"math"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"

	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (pi *model.PurchaseInvoice, e error) {
	o := orm.NewOrm()
	o.Begin()

	var adjNote string
	adjNote = r.AdjNote

	if r.Adjustment == 0 {
		adjNote = ""
	}

	var DebitNoteStr string
	for _, v := range r.DebitNoteArr {
		v.UsedInPurchaseInvoice = 1
		if _, e = o.Update(v, "UsedInPurchaseInvoice"); e != nil {
			o.Rollback()
			return nil, e
		}
		DebitNoteStr += strconv.Itoa(int(v.ID)) + ","
	}
	DebitNoteStr = strings.TrimSuffix(DebitNoteStr, ",")

	docCode, _ := util.GenerateDocCode("PI", r.PurchaseOrder.Supplier.Code, "purchase_invoice")
	pi = &model.PurchaseInvoice{
		Code:            docCode,
		PurchaseOrder:   r.PurchaseOrder,
		RecognitionDate: r.RecognitionAt,
		Status:          1,
		DueDate:         r.DueDateAt,
		DeliveryFee:     r.DeliveryFee,
		Adjustment:      r.Adjustment,
		AdjAmount:       math.Abs(r.AdjAmount),
		AdjNote:         adjNote,
		TotalPrice:      r.TotalPrice,
		TotalCharge:     r.TotalCharge,
		Note:            r.Note,
		TaxAmount:       r.TaxAmount,
		TaxPct:          r.TaxPct,
		PurchaseTerm:    &model.PurchaseTerm{ID: r.PurchaseOrder.TermPaymentPur.ID},
		DebitNoteIDs:    DebitNoteStr,
		CreatedAt:       time.Now(),
		CreatedBy:       r.Session.Staff,
	}
	if _, e := o.Insert(pi); e != nil {
		o.Rollback()
	}

	var arrPii []*model.PurchaseInvoiceItem
	for _, row := range r.InvoiceItems {
		item := &model.PurchaseInvoiceItem{
			PurchaseInvoice:   pi,
			Subtotal:          row.Subtotal,
			InvoiceQty:        row.InvoiceQty,
			PurchaseOrderItem: row.PurchaseOrderItem,
			UnitPrice:         row.UnitPrice,
			Note:              row.Note,
			Product:           row.Product,
			IncludeTax:        row.IncludeTax,
			UnitPriceTax:      row.UnitPriceTax,
			TaxableItem:       row.TaxableItem,
			TaxPercentage:     row.TaxPercentage,
			TaxAmount:         row.TaxAmount,
		}
		arrPii = append(arrPii, item)
	}

	if _, e := o.InsertMulti(100, &arrPii); e != nil {
		o.Rollback()
	}

	if e = log.AuditLogByUser(r.Session.Staff, pi.ID, "purchase_invoice", "create", ""); e != nil {
		o.Rollback()
		return nil, e
	}
	// update total_invoice in purchase order
	po := &model.PurchaseOrder{
		ID:           r.PurchaseOrder.ID,
		TotalInvoice: r.TotalCharge,
	}

	if _, e = o.Update(po, "TotalInvoice"); e != nil {
		o.Rollback()
	}

	o.Commit()
	return pi, e
}

//Update : function to update data requested into database
func Update(u updateRequest) (pi *model.PurchaseInvoice, e error) {
	o := orm.NewOrm()
	o.Begin()

	var adjNote string
	adjNote = u.AdjNote

	if u.Adjustment == 0 {
		adjNote = ""
	}

	var DebitNoteStr string
	for _, v := range u.DebitNoteArr {
		v.UsedInPurchaseInvoice = 1
		if _, e = o.Update(v, "UsedInPurchaseInvoice"); e != nil {
			o.Rollback()
			return nil, e
		}
		DebitNoteStr += strconv.Itoa(int(v.ID)) + ","
	}
	DebitNoteStr = strings.TrimSuffix(DebitNoteStr, ",")

	pi = &model.PurchaseInvoice{
		ID:              u.ID,
		PurchaseOrder:   u.PurchaseOrder,
		RecognitionDate: u.RecognitionAt,
		DueDate:         u.DueDateAt,
		DebitNoteIDs:    DebitNoteStr,
		DeliveryFee:     u.DeliveryFee,
		Adjustment:      u.Adjustment,
		AdjAmount:       math.Abs(u.AdjAmount),
		AdjNote:         adjNote,
		TotalPrice:      u.TotalPrice,
		TotalCharge:     u.TotalCharge,
		TaxAmount:       u.TaxAmount,
		Note:            u.Note,
		TaxPct:          u.TaxPct,
	}
	if _, e = o.Update(pi, "PurchaseOrder", "RecognitionDate", "DueDate", "DebitNoteIDs", "DeliveryFee", "Adjustment", "AdjAmount", "AdjNote", "TotalPrice", "TaxAmount", "TotalCharge", "Note", "TaxPct"); e != nil {
		o.Rollback()
	}

	for _, row := range u.InvoiceItems {
		item := &model.PurchaseInvoiceItem{
			ID:                row.PurchaseInvoiceItem.ID,
			PurchaseInvoice:   pi,
			Subtotal:          row.Subtotal,
			InvoiceQty:        row.InvoiceQty,
			PurchaseOrderItem: row.PurchaseOrderItem,
			UnitPrice:         row.UnitPrice,
			Note:              row.Note,
			Product:           row.Product,
			IncludeTax:        row.IncludeTax,
			UnitPriceTax:      row.UnitPriceTax,
			TaxableItem:       row.TaxableItem,
			TaxPercentage:     row.TaxPercentage,
			TaxAmount:         row.TaxAmount,
		}
		if _, e := o.Update(item, "PurchaseInvoice", "Subtotal", "InvoiceQty", "PurchaseOrderItem", "UnitPrice", "Note", "Product", "IncludeTax", "UnitPriceTax", "TaxableItem", "TaxPercentage", "TaxAmount"); e != nil {
			o.Rollback()
		}

		u.PurchaseOrder.TotalInvoice = u.TotalCharge
		if _, e := o.Update(u.PurchaseOrder, "TotalInvoice"); e != nil {
			o.Rollback()
		}

		if e = log.AuditLogByUser(u.Session.Staff, pi.ID, "purchase_invoice", "update"); e != nil {
			o.Rollback()
			return nil, e
		}

	}
	dnUsed := orm.Params{
		"used_in_purchase_invoice": 2,
	}
	if _, err := o.QueryTable(new(model.DebitNote)).Filter("ID__in", u.ExcludedDebitNoteIDs).Update(dnUsed); err != nil {
		o.Rollback()
		return nil, e
	}

	// update total_invoice in purchase order
	po := &model.PurchaseOrder{
		ID:           u.PurchaseOrder.ID,
		TotalInvoice: u.TotalCharge,
	}
	if _, e = o.Update(po, "TotalInvoice"); e != nil {
		o.Rollback()
	}

	o.Commit()
	return pi, e
}

// Cancel: function to change data status into 3
func Cancel(r cancelRequest) (pi *model.PurchaseInvoice, e error) {
	o := orm.NewOrm()
	o.Begin()

	debitNoteIDs := strings.Split(r.PurchaseInvoice.DebitNoteIDs, ",")
	dnUsed := orm.Params{
		"used_in_purchase_invoice": 2,
	}
	if _, e = o.QueryTable(new(model.DebitNote)).Filter("ID__in", debitNoteIDs).Update(dnUsed); e != nil {
		o.Rollback()
		return nil, e
	}

	r.PurchaseInvoice.Status = int8(3)
	r.PurchaseInvoice.TotalCharge = r.TotalCharge
	r.PurchaseInvoice.DebitNoteIDs = ""
	if _, e = o.Update(r.PurchaseInvoice, "Status", "TotalCharge", "DebitNoteIDs"); e != nil {
		o.Rollback()
	}

	r.PurchaseInvoice.PurchaseOrder.TotalInvoice = float64(0)
	if _, e = o.Update(r.PurchaseInvoice.PurchaseOrder, "total_invoice"); e != nil {
		o.Rollback()
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.PurchaseInvoice.ID, "purchase_invoice", "cancel", r.CancellationNote); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return r.PurchaseInvoice, nil
}

//AddTaxInvoice : function to Add Tax Invoice data requested into database
func AddTaxInvoice(u addTaxInvoiceRequest) (pi *model.PurchaseInvoice, e error) {
	o := orm.NewOrm()
	o.Begin()

	pi = &model.PurchaseInvoice{
		ID:               u.ID,
		TaxInvoiceURL:    u.TaxInvoiceURL,
		TaxInvoiceNumber: u.TaxInvoiceNumber,
	}
	if _, e = o.Update(pi, "TaxInvoiceURL", "TaxInvoiceNumber"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(u.Session.Staff, pi.ID, "purchase_invoice", "add_tax_invoice"); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return pi, e
}
