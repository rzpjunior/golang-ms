// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchase_payment

import (
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createRequest : struct to hold request data
type createRequest struct {
	Code                     string  `json:"-"`
	PurchaseInvoiceID        string  `json:"purchase_invoice_id" valid:"required"`
	PaymentMethodID          string  `json:"payment_method_id" valid:"required"`
	RecognitionDateStr       string  `json:"recognition_date" valid:"required"`
	Amount                   float64 `json:"amount"`
	PaidOff                  int8    `json:"paid_off"`
	ImageUrl                 string  `json:"image_url"`
	Note                     string  `json:"note"`
	RemainingAmount          float64 `json:"-"`
	BankPaymentVoucherNumber string  `json:"bank_payment_voucher_number" valid:"required"`

	RecognitionDate time.Time
	PurchaseInvoice *model.PurchaseInvoice
	PaymentMethod   *model.PaymentMethod
	GoodsReceipt    *model.GoodsReceipt
	DebitNote       []*model.DebitNote

	Session *auth.SessionData
}

// Validate : function to validate request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var e error
	debitNoteArr := make([]string, 0)
	layout := "2006-01-02"

	if c.RecognitionDate, e = time.Parse(layout, c.RecognitionDateStr); e != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("purchase payment date"))
	}

	if len(c.Note) > 100 {
		o.Failure("note.invalid", util.ErrorCharLength("note", 100))
	}

	if piID, err := common.Decrypt(c.PurchaseInvoiceID); err != nil {
		o.Failure("purchase_invoice.invalid", util.ErrorInvalidData("purchase invoice"))
	} else {
		if c.PurchaseInvoice, e = repository.ValidPurchaseInvoice(piID); e != nil {
			o.Failure("purchase_invoice.invalid", util.ErrorInvalidData("purchase invoice"))
		} else {
			if c.PurchaseInvoice.Status != int8(1) && c.PurchaseInvoice.Status != int8(6) {
				o.Failure("purchase_invoice.invalid", util.ErrorDocStatus("purchase invoice", "active or partial"))
			} else {
				c.PurchaseInvoice.PurchaseOrder.Read("ID")
				c.PurchaseInvoice.PurchaseOrder.Supplier.Read("ID")

				c.GoodsReceipt = &model.GoodsReceipt{PurchaseOrder: c.PurchaseInvoice.PurchaseOrder, Status: int8(2)}
				c.GoodsReceipt.Read("PurchaseOrder", "Status")

				_, paymentAmount, _ := repository.CheckPurchasePaymentAmount(c.PurchaseInvoice.ID)
				c.RemainingAmount = c.PurchaseInvoice.TotalCharge - paymentAmount
			}
		}
	}

	if strings.TrimSpace(c.PurchaseInvoice.DebitNoteIDs) != "" {
		debitNoteArr = strings.Split(c.PurchaseInvoice.DebitNoteIDs, ",")

		for _, v := range debitNoteArr {
			var dnID int
			if dnID, e = strconv.Atoi(v); e != nil {
				o.Failure("debit_note_id.invalid", util.ErrorInvalidData("debit note"))
			}

			debitNote := &model.DebitNote{
				ID: int64(dnID),
			}
			debitNote.Read("ID")
			c.DebitNote = append(c.DebitNote, debitNote)
		}
	}

	if pmID, err := common.Decrypt(c.PaymentMethodID); err != nil {
		o.Failure("payment_method.invalid", util.ErrorInvalidData("purchase invoice"))
	} else {
		if c.PaymentMethod, e = repository.ValidPaymentMethod(pmID); e != nil {
			o.Failure("payment_method.invalid", util.ErrorInvalidData("purchase invoice"))
		} else {
			if c.PaymentMethod.Status != int8(1) {
				o.Failure("payment_method.invalid", util.ErrorActive("payment method"))
			}
		}
	}

	if c.Amount < 0 {
		o.Failure("amount.invalid", util.ErrorEqualGreater("amount", "0"))
	}

	c.Amount = common.Rounder(c.Amount, 0.5, 2)

	if c.PaidOff == 1 || (c.PaidOff == 2 && c.Amount >= c.RemainingAmount) {
		c.PurchaseInvoice.Status = 2
		if c.GoodsReceipt.ID != 0 {
			c.PurchaseInvoice.PurchaseOrder.Status = 2
		}
	} else {
		c.PurchaseInvoice.Status = 6
	}

	return o
}

// Messages : function to return error validation messages after validation
func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"purchase_invoice_id.required":         util.ErrorInputRequired("purchase invoice"),
		"payment_method_id.required":           util.ErrorInputRequired("payment method"),
		"recognition_date.required":            util.ErrorInputRequired("purchase payment date"),
		"bank_payment_voucher_number.required": util.ErrorInputRequired("bank payment voucher number"),
	}
}
