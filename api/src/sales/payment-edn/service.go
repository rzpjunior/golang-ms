package paymentedn

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/datamodel/model"
)

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
		if e = log.CreditLimitLogByStaff(r.SalesPayment.SalesInvoice.SalesOrder.Branch.Merchant, r.SalesPayment.ID, "sales_payment_edn", r.CreditLimitBefore, r.CreditLimitAfter, r.Session.Staff.ID, "cancel payment"); e != nil {
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
