// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package payment

import (
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createRequest : struct to hold sales payment request data
type createRequest struct {
	Code                   string    `json:"-"`
	PaymentDateStr         string    `json:"payment_date" valid:"required"`
	PaymentMethodID        string    `json:"payment_method_id" valid:"required"`
	PaymentChannelID       string    `json:"payment_channel_id"`
	Amount                 float64   `json:"amount" valid:"required"`
	PaidOff                int8      `json:"paid_off"`
	Note                   string    `json:"note"`
	ImageUrl               string    `json:"image_url"`
	SalesInvoiceID         string    `json:"sales_invoice_id"`
	BankReceiveNum         string    `json:"bank_receive_num" valid:"required"`
	PaymentDate            time.Time `json:"-"`
	HaveCreditLimit        bool      `json:"-"`
	CreditLimitBefore      float64   `json:"-"`
	CreditLimitAfter       float64   `json:"-"`
	RemainingInvoiceAmount float64   `json:"-"`

	PaymentMethod  *model.PaymentMethod  `json:"-"`
	PaymentChannel *model.PaymentChannel `json:"-"`
	SalesInvoice   *model.SalesInvoice   `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate sales payment request data
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var (
		err                    error
		CountInProgressPayment int8
	)

	salesInvoiceID, _ := common.Decrypt(r.SalesInvoiceID)

	r.SalesInvoice = &model.SalesInvoice{ID: salesInvoiceID}

	if err = r.SalesInvoice.Read("ID"); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales invoice"))
		return o
	}

	if r.SalesInvoice.Status != 1 && r.SalesInvoice.Status != 6 {
		o.Failure("id.invalid", util.ErrorStatusDoc("sales payment", "created", "Sales Invoice"))
	}

	if err = r.SalesInvoice.SalesOrder.Read("ID"); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales order"))
		return o
	}

	if err = r.SalesInvoice.SalesOrder.Branch.Read("ID"); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("branch"))
		return o
	}

	if err = r.SalesInvoice.SalesOrder.Branch.Merchant.Read("ID"); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("merchant"))
		return o
	}

	if r.Amount < 0 {
		o.Failure("amount.equalorgreater", util.ErrorEqualGreater("paid amount", "0"))
	}

	r.CreditLimitBefore = r.SalesInvoice.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount

	if r.SalesInvoice.SalesOrder.Branch.Merchant.CreditLimitAmount > 0 || r.CreditLimitBefore < 0 {
		r.HaveCreditLimit = true
	}

	paymentMethodID, _ := common.Decrypt(r.PaymentMethodID)
	r.PaymentMethod = &model.PaymentMethod{ID: paymentMethodID}
	r.PaymentMethod.Read("ID")

	if r.PaymentChannelID != "" {
		paymentChannelID, _ := common.Decrypt(r.PaymentChannelID)
		r.PaymentChannel = &model.PaymentChannel{ID: paymentChannelID}
	}

	layout := "2006-01-02"
	if r.PaymentDate, err = time.Parse(layout, r.PaymentDateStr); err != nil {
		o.Failure("payment_date.invalid", util.ErrorInvalidData("payment date"))
	}

	if r.RemainingInvoiceAmount, err = repository.CheckRemainingSalesInvoiceAmount(salesInvoiceID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales invoice"))
	}

	if CountInProgressPayment, err = repository.CheckInProgressPayment(salesInvoiceID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales invoice"))
	}

	if r.RemainingInvoiceAmount < 0 {
		r.RemainingInvoiceAmount = 0
	}

	if r.PaidOff == 1 {
		if CountInProgressPayment > 0 {
			o.Failure("id.invalid", util.ErrorCannotCreatePaidOff())
		}

		if r.SalesInvoice.SalesOrder.Status == 9 {
			r.SalesInvoice.SalesOrder.Status = 12
		} else if r.SalesInvoice.SalesOrder.Status == 10 {
			r.SalesInvoice.SalesOrder.Status = 13
		} else if r.SalesInvoice.SalesOrder.Status == 11 {
			r.SalesInvoice.SalesOrder.Status = 2
			r.SalesInvoice.SalesOrder.FinishedAt = time.Now()
		}

		if r.SalesInvoice.Status == 6 {
			r.CreditLimitAfter = r.CreditLimitBefore + r.RemainingInvoiceAmount
		}

		if r.SalesInvoice.Status == 1 {
			r.CreditLimitAfter = r.CreditLimitBefore + r.SalesInvoice.TotalCharge
		}

		r.SalesInvoice.Status = 2

	} else {
		r.CreditLimitAfter = r.CreditLimitBefore + r.Amount
		r.SalesInvoice.Status = 6
		if r.Amount >= r.RemainingInvoiceAmount {
			r.CreditLimitAfter = r.CreditLimitBefore + r.RemainingInvoiceAmount
			if CountInProgressPayment == 0 {
				if r.SalesInvoice.SalesOrder.Status == 9 {
					r.SalesInvoice.SalesOrder.Status = 12
				} else if r.SalesInvoice.SalesOrder.Status == 10 {
					r.SalesInvoice.SalesOrder.Status = 13
				} else if r.SalesInvoice.SalesOrder.Status == 11 {
					r.SalesInvoice.SalesOrder.Status = 2
					r.SalesInvoice.SalesOrder.FinishedAt = time.Now()
				}

				r.SalesInvoice.Status = 2

			}
		}
	}

	if len(r.BankReceiveNum) > 50 {
		o.Failure("bank_receive_num", util.ErrorCharLength("bank receive number", 50))
	}

	return o
}

// Messages : function to return error validation messages
func (r *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"payment_date.required":      util.ErrorInputRequired("payment date"),
		"payment_method_id.required": util.ErrorInputRequired("payment method"),
		"amount.required":            util.ErrorInputRequired("amount"),
		"bank_receive_num.required":  util.ErrorInputRequired("bank receive voucher number"),
	}

	return messages
}
