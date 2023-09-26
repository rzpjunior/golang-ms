// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package paymentedn

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// cancelRequest : struct to hold payment set request data
type cancelRequest struct {
	ID                     int64   `json:"-"`
	Note                   string  `json:"note" valid:"required"`
	HaveCreditLimit        bool    `json:"-"`
	CreditLimitBefore      float64 `json:"-"`
	CreditLimitAfter       float64 `json:"-"`
	RemainingInvoiceAmount float64 `json:"-"`

	SalesPayment *model.SalesPayment `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate uom request data
func (r *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var (
		err                                        error
		totalAmountPayment, remainingAmountPayment float64
		countInProgressPayment                     int8
		countFinishedPayment                       int64
	)

	r.SalesPayment = &model.SalesPayment{ID: r.ID}
	if err = r.SalesPayment.Read("ID"); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales payment"))
		return o
	}

	if r.SalesPayment.Status != 5 {
		o.Failure("status.inactive", util.ErrorDocStatus("sales payment", "in progress"))
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

	// Only For Payment For Order Type EDN Sales
	if r.SalesPayment.SalesInvoice.SalesOrder.OrderType.ID != 13 {
		o.Failure("id.invalid", util.ErrorDocStatus("order type", "EDN Sales"))
		return o
	}

	if totalAmountPayment, err = repository.CheckAmountFinishAndInprogressPayment(r.SalesPayment.SalesInvoice.ID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales invoice"))
		return o
	}

	// validation cannot create payment paid off if any in progress payment
	if countInProgressPayment, err = repository.CheckInProgressPayment(r.SalesPayment.SalesInvoice.ID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales invoice"))
		return o
	}

	filter := map[string]interface{}{"sales_invoice_id": r.SalesPayment.SalesInvoice.ID, "status": int8(2)}
	exclude := map[string]interface{}{"ID": r.SalesPayment.ID}
	if _, countFinishedPayment, err = repository.CheckSalesPaymentData(filter, exclude); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales invoice"))
		return o
	}

	if r.SalesPayment.Status == 5 {
		countInProgressPayment -= 1
	}

	r.CreditLimitBefore = r.SalesPayment.SalesInvoice.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount

	remainingAmountPayment = totalAmountPayment - r.SalesPayment.Amount

	r.SalesPayment.SalesInvoice.Status = 1

	r.CreditLimitAfter = r.CreditLimitBefore - r.SalesPayment.Amount

	if r.SalesPayment.SalesInvoice.SalesOrder.Branch.Merchant.CreditLimitAmount > 0 || r.CreditLimitBefore < 0 {
		r.HaveCreditLimit = true
	}

	if r.SalesPayment.Amount > r.SalesPayment.SalesInvoice.TotalCharge && remainingAmountPayment == 0 {
		r.CreditLimitAfter = r.CreditLimitBefore - r.SalesPayment.SalesInvoice.TotalCharge
		return o
	}

	if remainingAmountPayment >= r.SalesPayment.SalesInvoice.TotalCharge {
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
		return o
	}

	if remainingAmountPayment < r.SalesPayment.SalesInvoice.TotalCharge && remainingAmountPayment != 0 {
		if totalAmountPayment > r.SalesPayment.SalesInvoice.TotalCharge {
			r.CreditLimitAfter = r.CreditLimitBefore - (r.SalesPayment.SalesInvoice.TotalCharge - remainingAmountPayment)
		}
		if countFinishedPayment > 0 {
			r.SalesPayment.SalesInvoice.Status = 6
		}
		return o
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
