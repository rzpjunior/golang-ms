package purchase_payment

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
)

// Save : function to save data requested into database
func Save(r createRequest) (pp *model.PurchasePayment, e error) {
	r.Code, e = util.GenerateDocCode("PP", r.PurchaseInvoice.PurchaseOrder.Supplier.Code, "purchase_payment")

	o := orm.NewOrm()
	o.Begin()

	if e == nil {
		pp = &model.PurchasePayment{
			Code:                     r.Code,
			PurchaseInvoice:          r.PurchaseInvoice,
			PaymentMethod:            r.PaymentMethod,
			Note:                     r.Note,
			Status:                   int8(2),
			RecognitionDate:          r.RecognitionDate,
			Amount:                   r.Amount,
			PaidOff:                  r.PaidOff,
			ImageUrl:                 r.ImageUrl,
			BankPaymentVoucherNumber: r.BankPaymentVoucherNumber,
			CreatedAt:                time.Now(),
			CreatedBy:                r.Session.Staff,
		}
		if _, e = o.Insert(pp); e != nil {
			o.Rollback()
			return nil, e
		}

		if _, e = o.Update(r.PurchaseInvoice, "Status"); e != nil {
			o.Rollback()
			return nil, e
		}

		// region update status debit note & supplier return
		for _, v := range r.DebitNote {
			v.Status = 2
			v.SupplierReturn.Status = 2

			if _, e = o.Update(v, "Status"); e != nil {
				o.Rollback()
				return nil, e
			}

			if _, e = o.Update(v.SupplierReturn, "Status"); e != nil {
				o.Rollback()
				return nil, e
			}

			if e = log.AuditLogByUser(r.Session.Staff, v.ID, "debit_note", "finish", ""); e != nil {
				o.Rollback()
				return nil, e
			}
			if e = log.AuditLogByUser(r.Session.Staff, v.SupplierReturn.ID, "supplier_return", "finish", ""); e != nil {
				o.Rollback()
				return nil, e
			}
		}
		// endregion
		if _, e = o.Update(r.PurchaseInvoice.PurchaseOrder, "Status"); e != nil {
			o.Rollback()
			return nil, e
		}

		o.Commit()
		if e = log.AuditLogByUser(r.Session.Staff, pp.ID, "purchase_payment", "create", ""); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	return pp, e
}

// Cancel: function to change data status into 3
func Cancel(r cancelRequest) (pp *model.PurchasePayment, e error) {
	var totalPay, totalPayFin, totalPayCan int
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	orSelect.Raw("SELECT COUNT(id) FROM purchase_payment WHERE purchase_invoice_id = ? AND id != ?", r.PurchasePayment.PurchaseInvoice.ID, r.ID).QueryRow(&totalPay)
	orSelect.Raw("SELECT COUNT(id) FROM purchase_payment WHERE purchase_invoice_id = ? AND id != ? AND status = 2", r.PurchasePayment.PurchaseInvoice.ID, r.ID).QueryRow(&totalPayFin)
	orSelect.Raw("SELECT COUNT(id) FROM purchase_payment WHERE purchase_invoice_id = ? AND id != ?  AND status = 3", r.PurchasePayment.PurchaseInvoice.ID, r.ID).QueryRow(&totalPayCan)
	o := orm.NewOrm()
	o.Begin()
	r.PurchasePayment.Status = 3
	if _, e = o.Update(r.PurchasePayment, "Status"); e == nil {
		r.PurchasePayment.PurchaseInvoice.Read("ID")
		r.PurchasePayment.PurchaseInvoice.PurchaseOrder.Read("ID")
		if r.PurchasePayment.PurchaseInvoice.Status == 2 { // if pi status finished
			if totalPayFin > 0 {
				r.PurchasePayment.PurchaseInvoice.Status = 6
				if _, e = o.Update(r.PurchasePayment.PurchaseInvoice, "Status"); e != nil {
					o.Rollback()
				}
			}
			if (totalPayCan > 0 && totalPayFin == 0) || totalPay == 0 { // jika semua cancel atau tidak ada payment selain diri nya
				r.PurchasePayment.PurchaseInvoice.Status = 1
				if _, e = o.Update(r.PurchasePayment.PurchaseInvoice, "Status"); e != nil {
					o.Rollback()
				}
			}
			r.PurchasePayment.PurchaseInvoice.PurchaseOrder.Status = 1
			if _, e = o.Update(r.PurchasePayment.PurchaseInvoice.PurchaseOrder, "Status"); e != nil {
				o.Rollback()
			}
		} else if r.PurchasePayment.PurchaseInvoice.Status == 6 { // if pi status partial
			if (totalPayCan > 0 && totalPayFin == 0) || totalPay == 0 {
				r.PurchasePayment.PurchaseInvoice.Status = 1
				if _, e = o.Update(r.PurchasePayment.PurchaseInvoice, "Status"); e != nil {
					o.Rollback()
				}
			}
		}
	} else {
		o.Rollback()
	}
	if e == nil {
		o.Commit()
		e = log.AuditLogByUser(r.Session.Staff, r.ID, "purchase_payment", "cancel", r.CancellationNote)
	}
	return r.PurchasePayment, nil
}

// SaveBulk : function to save bulk purchase payment
func SaveBulk(r createBulkRequest) (totalSuccess int64, e error) {
	for _, v := range r.Requests {
		purchaseInvoice := &model.PurchaseInvoice{ID: v.PurchaseInvoice.ID}
		purchaseInvoice.Read("ID")

		if purchaseInvoice.Status == 1 || purchaseInvoice.Status == 6 {
			if _, e = Save(v); e == nil {
				totalSuccess++
			}
		}
	}

	return
}
