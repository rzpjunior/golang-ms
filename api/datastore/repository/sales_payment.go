// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"

	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetSalesPayment find a single data sales term using field and value condition.
func GetSalesPayment(field string, values ...interface{}) (*model.SalesPayment, error) {
	m := new(model.SalesPayment)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetSalesPayments : function to get data from database based on parameters
func GetSalesPayments(rq *orm.RequestQuery, createdByFilter ...int64) (m []*model.SalesPayment, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	q, _ := rq.QueryReadOnly(new(model.SalesPayment))
	var mx []*model.SalesPayment

	// if QueryParam createdByFilter is exist
	condFilterByCreatedBy := q.GetCond()
	if len(createdByFilter) > 0 {
		cond := orm.NewCondition()
		cond = cond.And("created_by__in", createdByFilter)
		condFilterByCreatedBy = condFilterByCreatedBy.AndCond(cond)
	}
	q = q.SetCond(condFilterByCreatedBy)

	if total, err = q.Filter("status__in", 1, 2, 3, 5).OrderBy("-id").All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			o.Raw("SELECT note from audit_log where ref_id = ? and type = 'sales_payment' and function = 'cancel'", v.ID).QueryRow(&v.CancellationNote)
		}
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterSalesPayments : function to get data from database based on parameters with filtered permission
func GetFilterSalesPayments(rq *orm.RequestQuery) (m []*model.SalesPayment, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SalesPayment))

	var mx []*model.SalesPayment
	if total, err = q.Filter("status__in", 1, 2, 3, 5).OrderBy("-id").All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidSalesPayment : function to check if id is valid in database
func ValidSalesPayment(id int64) (SalesPayment *model.SalesPayment, e error) {
	SalesPayment = &model.SalesPayment{ID: id}
	e = SalesPayment.Read("ID")

	return
}

// CheckSalesPaymentData : function to check data based on filter and exclude parameters
func CheckSalesPaymentData(filter map[string]interface{}, exclude map[string]interface{}) (sp []*model.SalesPayment, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.SalesPayment))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err := o.All(&sp); err == nil {
		return sp, total, nil
	}

	return nil, 0, err
}

// Get Grand Total Amount Sales Payment Where SP Status (2, 5) and SI status (1, 6)
func GetGrandTotalAmountSP(merchantID int64) (totalAmount float64, e error) {
	o := orm.NewOrm()
	o.Using("read_only")

	// filter payment status finished and in progress
	paymentStatus := []int{2, 5}
	// filter invoice active and partial
	invoiceStatus := []int{1, 6}
	totalPaymentSI := make(map[int64]float64)
	totalChargeSI := make(map[int64]float64)

	rq := o.QueryTable(new(model.SalesPayment))

	cond1 := orm.NewCondition()
	cond1 = cond1.And("status__in", paymentStatus).And("salesinvoice__status__in", invoiceStatus).And("salesinvoice__salesorder__branch__merchant__id", merchantID)

	rq = rq.SetCond(cond1)

	var mx []*model.SalesPayment

	if _, e = rq.All(&mx); e != nil {
		return totalAmount, e
	}

	for _, v := range mx {
		// Calculate total amount payment of sales invoice
		if _, ok := totalPaymentSI[v.SalesInvoice.ID]; ok {
			totalPaymentSI[v.SalesInvoice.ID] += v.Amount
		} else {
			totalPaymentSI[v.SalesInvoice.ID] = v.Amount
		}

		// Intitate Total Charge SI to map
		if _, ok := totalChargeSI[v.SalesInvoice.ID]; !ok {
			if e = v.SalesInvoice.Read("ID"); e != nil {
				return totalAmount, e
			}
			totalChargeSI[v.SalesInvoice.ID] = v.SalesInvoice.TotalCharge
		}

		// change total payment to total charge if total payment more than total charge invoice
		if totalPaymentSI[v.SalesInvoice.ID] > totalChargeSI[v.SalesInvoice.ID] {
			totalPaymentSI[v.SalesInvoice.ID] = totalChargeSI[v.SalesInvoice.ID]
		}
	}

	// Calculate total amount payment
	for _, amount := range totalPaymentSI {
		totalAmount += amount
	}

	return totalAmount, nil
}

// CheckPaidOff : Check There is Payment Paid Off Or not in a SI
func CheckPaidOff(SalesInvoiceID int64) (isAnyPaidOff bool) {

	o := orm.NewOrm()
	o.Using("read_only")

	isAnyPaidOff = o.QueryTable("sales_payment").Filter("status", 2).Filter("paid_off", 1).Filter("sales_invoice_id", SalesInvoiceID).Exist()

	return

}

// CheckInProgressPayment: Check There is Payment in Progress Or not in a SI
func CheckInProgressPayment(SalesInvoiceID int64) (CountInProgressPayment int8, e error) {

	o := orm.NewOrm()
	o.Using("read_only")

	var count int64

	if count, e = o.QueryTable("sales_payment").Filter("status", 5).Filter("sales_invoice_id", SalesInvoiceID).Count(); e != nil {
		return int8(count), e
	}

	return int8(count), nil
}

// CheckAmountFinishAndInprogressPayment: Check Total Amount Finish and In Progress Payment of Sales Invoice
func CheckAmountFinishAndInprogressPayment(SalesInvoiceID int64) (totalAmount float64, e error) {

	o := orm.NewOrm()
	o.Using("read_only")

	q := "SELECT SUM(sp.amount) total_amount FROM sales_invoice si " +
		"LEFT JOIN sales_payment sp on sp.sales_invoice_id = si.id " +
		"WHERE sp.status IN (2,5) AND si.id = ?"

	if e = o.Raw(q, SalesInvoiceID).QueryRow(&totalAmount); e != nil {
		return totalAmount, e
	}

	return totalAmount, nil
}
