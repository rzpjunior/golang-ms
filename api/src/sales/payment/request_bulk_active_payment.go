// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
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

type bulkCreateActivePaymentRequest struct {
	ReceivedDateStr string    `json:"received_date" valid:"required"`
	ReceivedDate    time.Time `json:"-"`
	SalesInvoice    []*items  `json:"sales_invoice_items" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

type items struct {
	PaymentMethodID        string  `json:"payment_method_id" valid:"required"`
	SalesInvoiceID         string  `json:"sales_invoice_id"`
	Amount                 float64 `json:"amount" valid:"required"`
	Note                   string  `json:"note"`
	ImageUrl               string  `json:"image_url"`
	FixVA                  bool    `json:"fix_va"`
	RemainingInvoiceAmount float64 `json:"-"`
	PaymentStatus          int8    `json:"-"`
	HaveCreditLimit        bool    `json:"-"`
	CreditLimitBefore      float64 `json:"-"`
	CreditLimitAfter       float64 `json:"-"`

	SalesInvoice  *model.SalesInvoice  `json:"-"`
	PaymentMethod *model.PaymentMethod `json:"-"`
}

func (r *bulkCreateActivePaymentRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	var err error
	var creditLimitAmount float64

	layout := "2006-01-02"
	if r.ReceivedDate, err = time.Parse(layout, r.ReceivedDateStr); err != nil {
		o.Failure("received_date.invalid", util.ErrorInvalidData("received date"))
		return o
	}

	for i, k := range r.SalesInvoice {
		salesInvoiceID, _ := common.Decrypt(k.SalesInvoiceID)
		k.SalesInvoice = &model.SalesInvoice{ID: salesInvoiceID}

		if err = k.SalesInvoice.Read("ID"); err != nil {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("sales invoice"))
			return o
		}
		if k.SalesInvoice.Status != 1 && k.SalesInvoice.Status != 6 {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorStatusDoc("sales payment", "created", "Sales Invoice"))
			return o
		}
		if err = k.SalesInvoice.SalesOrder.Read("ID"); err != nil {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("sales order"))
			return o
		}
		if err = k.SalesInvoice.SalesOrder.Branch.Read("ID"); err != nil {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("branch"))
			return o
		}
		if err = k.SalesInvoice.SalesOrder.Branch.Merchant.Read("ID"); err != nil {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("merchant"))
			return o
		}

		if k.Amount < 0 {
			o.Failure("amount_"+strconv.Itoa(i)+".equalorgreater", util.ErrorEqualGreater("paid amount", "0"))
			return o
		}

		if paymentMethodID, err := common.Decrypt(k.PaymentMethodID); err == nil {
			if k.PaymentMethod, err = repository.ValidPaymentMethod(paymentMethodID); err != nil {
				o.Failure("payment_method_id.invalid", util.ErrorInvalidData("payment method"))
				return o
			}
			k.PaymentMethod = &model.PaymentMethod{ID: paymentMethodID}
			k.PaymentMethod.Read("ID")
		} else {
			o.Failure("payment_method_id.invalid", util.ErrorInvalidData("payment method"))
			return o
		}

		if k.RemainingInvoiceAmount, err = repository.CheckRemainingSalesInvoiceAmount(salesInvoiceID); err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("sales invoice"))
		}

		k.PaymentStatus = 1

		// Validation Fix VA
		if (k.FixVA && k.PaymentMethod.ID == 2) || k.SalesInvoice.SalesOrder.OrderType.ID == 13 {

			k.PaymentStatus = 5

			creditLimitAmount = k.SalesInvoice.SalesOrder.Branch.Merchant.CreditLimitAmount
			k.CreditLimitBefore = k.SalesInvoice.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount
			if creditLimitAmount > 0 || k.CreditLimitBefore < 0 {
				k.HaveCreditLimit = true
			}
		}
	}

	return o
}

// Messages : function to return error messages after validation
func (r *bulkCreateActivePaymentRequest) Messages() map[string]string {
	return map[string]string{
		"received_date.required":       util.ErrorInputRequired("received date"),
		"sales_invoice_items.required": util.ErrorInputRequired("sales invoice item"),
	}
}
