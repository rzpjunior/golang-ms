// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package payment

import (
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type bulkConfirmPaymentRequest struct {
	PaymentDateStr string    `json:"payment_date" valid:"required"`
	PaymentDate    time.Time `json:"-"`
	BankReceiveNum string    `json:"bank_receive_num" valid:"required"`

	SalesPaymentItems []*bulkConfirmPaymentItems `json:"items" valid:"required"`
	Session           *auth.SessionData          `json:"-"`
}

type bulkConfirmPaymentItems struct {
	ID                     string  `json:"id"`
	Amount                 float64 `json:"amount" valid:"required"`
	PaidOff                int8    `json:"paid_off" valid:"required"`
	Note                   string  `json:"note"`
	HaveCreditLimit        bool    `json:"-"`
	CreditLimitBefore      float64 `json:"-"`
	CreditLimitAfter       float64 `json:"-"`
	RemainingInvoiceAmount float64 `json:"-"`
	CountInProgressPayment int8    `json:"-"`
	TotalPaidAmount        float64 `json:"-"`

	SalesPayment *model.SalesPayment `json:"-"`
	SalesInvoice *model.SalesInvoice `json:"-"`
}

func (r *bulkConfirmPaymentRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	var creditLimitAmount float64
	TotalInProgress := make(map[int64]int8)
	isAnyPaidOff := make(map[int64]bool)

	if len(r.BankReceiveNum) > 50 {
		o.Failure("bank_receive_num", util.ErrorCharLength("bank receive number", 50))
	}

	for i, item := range r.SalesPaymentItems {
		if len(item.Note) > 250 {
			o.Failure("note", util.ErrorCharLength("note", 250))
		}

		if item.Amount < 0 {
			o.Failure("amount_"+strconv.Itoa(i)+".equalorgreater", util.ErrorEqualGreater("paid amount", "0"))
		}

		salesPaymentID, err := common.Decrypt(item.ID)
		if err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("sales payment"))
		}
		item.SalesPayment = &model.SalesPayment{ID: salesPaymentID}

		layout := "2006-01-02"
		if r.PaymentDate, err = time.Parse(layout, r.PaymentDateStr); err != nil {
			o.Failure("payment_date.invalid", util.ErrorInvalidData("payment date"))
		}

		if err = item.SalesPayment.Read("ID"); err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("sales payment"))
		}

		if item.SalesPayment.Status != 1 && item.SalesPayment.Status != 5 {
			o.Failure("status.inactive", util.ErrorDocStatus("sales payment", "active or in progress"))
			return o
		}

		salesInvoiceId := item.SalesPayment.SalesInvoice.ID
		item.SalesInvoice = &model.SalesInvoice{ID: salesInvoiceId}

		if err = item.SalesInvoice.Read("ID"); err != nil {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("sales invoice"))
			return o
		}

		if item.SalesInvoice.Status != 1 && item.SalesInvoice.Status != 6 {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorStatusDoc("sales payment", "created", "Sales Invoice"))
		}

		if err = item.SalesInvoice.SalesOrder.Read("ID"); err != nil {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("sales order"))
			return o
		}

		if err = item.SalesInvoice.SalesOrder.Branch.Read("ID"); err != nil {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("branch"))
			return o
		}

		if err = item.SalesInvoice.SalesOrder.Branch.Merchant.Read("ID"); err != nil {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("merchant"))
			return o
		}

		if err = item.SalesInvoice.SalesOrder.Branch.Merchant.UserMerchant.Read("ID"); err != nil {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("user merchant"))
			return o
		}

		if item.RemainingInvoiceAmount, err = repository.CheckRemainingSalesInvoiceAmount(salesInvoiceId); err != nil {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("sales invoice"))
			return o
		}

		// validation cannot create payment paid off if any in progress payment
		if item.CountInProgressPayment, err = repository.CheckInProgressPayment(salesInvoiceId); err != nil {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("sales invoice"))
			return o
		}

		// Calculate actual Count "In Progress" Payment in each Sales Invoice
		if _, ok := TotalInProgress[salesInvoiceId]; ok && item.SalesPayment.Status == 5 {
			TotalInProgress[salesInvoiceId] -= 1
			item.CountInProgressPayment = TotalInProgress[salesInvoiceId]
		}

		if _, ok := TotalInProgress[salesInvoiceId]; !ok && item.SalesPayment.Status == 5 {
			TotalInProgress[salesInvoiceId] = item.CountInProgressPayment
			item.CountInProgressPayment = TotalInProgress[salesInvoiceId]
		}

		if item.CountInProgressPayment > 0 && item.PaidOff == 1 {
			o.Failure("id.invalid", util.ErrorCannotCreatePaidOff())
			return o
		}

		// Prevent Any Confirm Active Payment if There is Paid off payment
		if _, ok := isAnyPaidOff[salesInvoiceId]; ok {
			o.Failure("id.invalid", util.ErrorCannotConfirmBulkPaidOff())
			return o
		}

		if item.PaidOff == 1 {
			isAnyPaidOff[salesInvoiceId] = true
		}

		creditLimitAmount = item.SalesInvoice.SalesOrder.Branch.Merchant.CreditLimitAmount
		item.CreditLimitBefore = item.SalesInvoice.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount
		if creditLimitAmount > 0 || item.CreditLimitBefore < 0 {
			item.HaveCreditLimit = true
		}
	}

	return o
}

// Messages : function to return error messages after validation
func (r *bulkConfirmPaymentRequest) Messages() map[string]string {
	return map[string]string{
		"payment_date.required":     util.ErrorInputRequired("payment date"),
		"bank_receive_num.required": util.ErrorInputRequired("bank receive number"),
		"items.required":            util.ErrorInputRequired("sales payment items"),
	}
}
