// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package payment

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// cancelRequest : struct to hold payment set request data
type cancelRequest struct {
	ID                int64   `json:"-"`
	Note              string  `json:"note" valid:"required"`
	HaveCreditLimit   bool    `json:"-"`
	CreditLimitBefore float64 `json:"-"`
	CreditLimitAfter  float64 `json:"-"`

	SalesPayment *model.SalesPayment `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate uom request data
func (r *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var (
		err                                        error
		countFinishedPayment                       int64
		filter, exclude                            map[string]interface{}
		isAnyPaidOff                               bool
		totalAmountPayment, remainingAmountPayment float64
		countInProgressPayment                     int8
	)

	r.SalesPayment = &model.SalesPayment{ID: r.ID}
	if err = r.SalesPayment.Read("ID"); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales payment"))
		return o
	}

	if r.SalesPayment.Status != 1 && r.SalesPayment.Status != 2 && r.SalesPayment.Status != 5 {
		o.Failure("status.inactive", util.ErrorDocStatus("sales payment", "active"))
		return o
	}

	if err = r.SalesPayment.SalesInvoice.Read("ID"); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales invoice"))
		return o
	}

	if err = r.SalesPayment.SalesInvoice.SalesOrder.Read("ID"); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales order"))
		return o
	}

	if err = r.SalesPayment.SalesInvoice.SalesOrder.Branch.Read("ID"); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("branch"))
		return o
	}

	if err = r.SalesPayment.SalesInvoice.SalesOrder.Branch.Merchant.Read("ID"); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("merchant"))
		return o
	}

	filter = map[string]interface{}{"sales_invoice_id": r.SalesPayment.SalesInvoice.ID, "status": int8(2)}
	exclude = map[string]interface{}{"ID": r.SalesPayment.ID}
	_, countFinishedPayment, err = repository.CheckSalesPaymentData(filter, exclude)

	r.CreditLimitBefore = r.SalesPayment.SalesInvoice.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount

	if r.SalesPayment.SalesInvoice.SalesOrder.Branch.Merchant.CreditLimitAmount > 0 || r.CreditLimitBefore < 0 {
		r.HaveCreditLimit = true
	}

	if r.SalesPayment.PaidOff == 1 && countFinishedPayment == 0 {
		r.CreditLimitAfter = r.CreditLimitBefore - r.SalesPayment.SalesInvoice.TotalCharge
	} else {
		r.CreditLimitAfter = r.CreditLimitBefore - r.SalesPayment.Amount
	}

	if totalAmountPayment, err = repository.CheckAmountFinishAndInprogressPayment(r.SalesPayment.SalesInvoice.ID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales_invoice"))
		return o
	}

	// validation cannot create payment paid off if any in progress payment
	if countInProgressPayment, err = repository.CheckInProgressPayment(r.SalesPayment.SalesInvoice.ID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales invoice"))
		return o
	}

	if r.SalesPayment.Status == 5 {
		countInProgressPayment -= 1
	}

	remainingAmountPayment = totalAmountPayment - r.SalesPayment.Amount

	if r.SalesPayment.SalesInvoice.Status == 2 {

		switch r.SalesPayment.SalesInvoice.SalesOrder.Status {
		case 12:
			r.SalesPayment.SalesInvoice.SalesOrder.Status = 9
		case 13:
			r.SalesPayment.SalesInvoice.SalesOrder.Status = 10
		case 2:
			r.SalesPayment.SalesInvoice.SalesOrder.Status = 11
		}
	}

	r.SalesPayment.SalesInvoice.Status = 1

	if r.SalesPayment.Amount > r.SalesPayment.SalesInvoice.TotalCharge && remainingAmountPayment == 0 {
		r.CreditLimitAfter = r.CreditLimitBefore - r.SalesPayment.SalesInvoice.TotalCharge
	} else if remainingAmountPayment >= r.SalesPayment.SalesInvoice.TotalCharge {
		r.CreditLimitAfter = r.CreditLimitBefore
		if countFinishedPayment > 0 {
			if countInProgressPayment == 0 {
				r.SalesPayment.SalesInvoice.Status = 2
				switch r.SalesPayment.SalesInvoice.SalesOrder.Status {
				case 9:
					r.SalesPayment.SalesInvoice.SalesOrder.Status = 12
				case 10:
					r.SalesPayment.SalesInvoice.SalesOrder.Status = 13
				case 11:
					r.SalesPayment.SalesInvoice.SalesOrder.Status = 2
				}
			} else {
				r.SalesPayment.SalesInvoice.Status = 6
			}
		}
	} else if remainingAmountPayment < r.SalesPayment.SalesInvoice.TotalCharge && remainingAmountPayment != 0 {
		if totalAmountPayment > r.SalesPayment.SalesInvoice.TotalCharge {
			r.CreditLimitAfter = r.CreditLimitBefore - (r.SalesPayment.SalesInvoice.TotalCharge - remainingAmountPayment)
		}
		if countFinishedPayment > 0 {
			r.SalesPayment.SalesInvoice.Status = 6
		}
		if r.SalesPayment.PaidOff == 1 {
			differenceActual := r.SalesPayment.SalesInvoice.TotalCharge - (totalAmountPayment - r.SalesPayment.Amount)
			r.CreditLimitAfter = r.CreditLimitBefore - differenceActual
		}
	}

	// Validation Cancel Payment Regular
	if r.SalesPayment.Status == 2 && r.SalesPayment.PaidOff != 1 {
		isAnyPaidOff = repository.CheckPaidOff(r.SalesPayment.SalesInvoice.ID)
		if isAnyPaidOff {
			o.Failure("id.invalid", util.ErrorInvalidData("to cancel the payment."))
			o.Failure("sales_payment.invalid", util.ErrorCannotCancelSalesPayment())
			return o
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *cancelRequest) Messages() map[string]string {
	messages := map[string]string{
		"note.required": util.ErrorInputRequired("cancellation note"),
	}

	return messages
}
