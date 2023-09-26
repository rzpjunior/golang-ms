// Copyright 2020 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package payment

import (
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type bulkPaymentRequest struct {
	PaymentMethodID  string    `json:"payment_method_id" valid:"required"`
	PaymentChannelID string    `json:"payment_channel_id"`
	PaymentDateStr   string    `json:"payment_date" valid:"required"`
	PaymentDate      time.Time `json:"-"`
	BankReceiveNum   string    `json:"bank_receive_num" valid:"required"`

	SalesInvoice   []*salesInvoiceItems  `json:"sales_invoice_items" valid:"required"`
	PaymentMethod  *model.PaymentMethod  `json:"-"`
	PaymentChannel *model.PaymentChannel `json:"-"`

	Session *auth.SessionData `json:"-"`
}

type salesInvoiceItems struct {
	SalesInvoiceID         string  `json:"sales_invoice_id"`
	Amount                 float64 `json:"amount" valid:"required"`
	PaidOff                int8    `json:"paid_off"`
	Note                   string  `json:"note"`
	ImageUrl               string  `json:"image_url"`
	RemainingInvoiceAmount float64 `json:"-"`
	CountInProgressPayment int8    `json:"-"`
	HaveCreditLimit        bool    `json:"-"`
	CreditLimitBefore      float64 `json:"-"`
	CreditLimitAfter       float64 `json:"-"`

	SalesInvoice *model.SalesInvoice `json:"-"`
}

func (r *bulkPaymentRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	var err error
	var creditLimitAmount float64

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

	if len(r.BankReceiveNum) > 50 {
		o.Failure("bank_receive_num", util.ErrorCharLength("bank receive number", 50))
	}

	for i, k := range r.SalesInvoice {
		salesInvoiceID, _ := common.Decrypt(k.SalesInvoiceID)
		k.SalesInvoice = &model.SalesInvoice{ID: salesInvoiceID}
		if err = k.SalesInvoice.Read("ID"); err == nil {
			if k.SalesInvoice.Status != 1 && k.SalesInvoice.Status != 6 {
				o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorStatusDoc("sales payment", "created", "Sales Invoice"))
			}
			if err = k.SalesInvoice.SalesOrder.Read("ID"); err == nil {
				if err = k.SalesInvoice.SalesOrder.Branch.Read("ID"); err != nil {
					o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("branch"))
					return o
				}
			} else {
				o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("sales order"))
				return o
			}
		} else {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("sales invoice"))
			return o
		}

		if k.Amount < 0 {
			o.Failure("amount_"+strconv.Itoa(i)+".equalorgreater", util.ErrorEqualGreater("paid amount", "0"))
		}

		if k.RemainingInvoiceAmount, err = repository.CheckRemainingSalesInvoiceAmount(k.SalesInvoice.ID); err != nil {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("sales invoice"))
		}

		if err = k.SalesInvoice.SalesOrder.Branch.Merchant.Read("ID"); err != nil {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("merchant"))
		}

		// validation cannot create payment paid off if any in progress payment
		if k.CountInProgressPayment, err = repository.CheckInProgressPayment(salesInvoiceID); err != nil {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("sales invoice"))
		}

		if k.CountInProgressPayment > 0 && k.PaidOff == 1 {
			o.Failure("id.invalid", util.ErrorCannotCreatePaidOff())
		}

		creditLimitAmount = k.SalesInvoice.SalesOrder.Branch.Merchant.CreditLimitAmount
		k.CreditLimitAfter = k.SalesInvoice.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount
		if creditLimitAmount > 0 || k.CreditLimitBefore < 0 {
			k.HaveCreditLimit = true
		}
	}

	return o
}

// Messages : function to return error messages after validation
func (r *bulkPaymentRequest) Messages() map[string]string {
	return map[string]string{
		"payment_method_id.required":   util.ErrorInputRequired("payment method"),
		"payment_date.required":        util.ErrorInputRequired("payment date"),
		"bank_receive_num.required":    util.ErrorInputRequired("bank receive number"),
		"sales_invoice_items.required": util.ErrorInputRequired("sales invoice item"),
	}
}
