// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchase_payment

import (
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// createBulkRequest : struct to hold request data
type createBulkRequest struct {
	Data     []*paymentItems `json:"data" valid:"required"`
	Requests []createRequest `json:"-"`

	Session *auth.SessionData
}

type paymentItems struct {
	PaymentDateStr           string  `json:"payment_date"`
	PaymentMethodID          string  `json:"payment_method_id"`
	PurchaseInvoiceCode      string  `json:"purchase_invoice_code"`
	BankPaymentVoucherNumber string  `json:"bank_payment_voucher_number"`
	Amount                   float64 `json:"amount"`
	Note                     string  `json:"note"`
	PaidOff                  string  `json:"paid_off"`
	RemainingAmount          float64 `json:"-"`

	PaymentDate     time.Time              `json:"-"`
	PurchaseInvoice *model.PurchaseInvoice `json:"-"`
	GoodsReceipt    *model.GoodsReceipt    `json:"-"`
	PaymentMethod   *model.PaymentMethod   `json:"-"`
}

// Validate : function to validate request data
func (c *createBulkRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	layout := "2006-01-02"
	var filter, exclude map[string]interface{}
	var paidAmount float64
	var pmID int64
	var purchaseInvoice []*model.PurchaseInvoice
	var request createRequest
	var paidOff int8
	finishedInvoice := make(map[string]int8)

	for k, v := range c.Data {
		// payment date, payment method, purchase invoice code, bank payment voucher number are required
		if v.PaymentDateStr == "" {
			o.Failure("id.invalid", util.ErrorInputRequired("Payment Date "+strconv.Itoa(k+1)))
			return o
		}

		if v.PaymentMethodID == "" {
			o.Failure("id.invalid", util.ErrorInputRequired("Payment Method "+strconv.Itoa(k+1)))
			return o
		}

		if v.PurchaseInvoiceCode == "" {
			o.Failure("id.invalid", util.ErrorInputRequired("Purchase Invoice Code "+strconv.Itoa(k+1)))
			return o
		}

		if v.BankPaymentVoucherNumber == "" {
			o.Failure("id.invalid", util.ErrorInputRequired("Bank Payment Voucher Number "+strconv.Itoa(k+1)))
			return o
		}

		if _, isExist := finishedInvoice[v.PurchaseInvoiceCode]; isExist {
			continue
		}

		paidOff = 2

		if v.PaymentDate, err = time.Parse(layout, v.PaymentDateStr); err != nil {
			continue
		}

		if len(v.Note) > 100 {
			continue
		}

		v.PaidOff = strings.ToLower(v.PaidOff)

		filter = map[string]interface{}{"code": v.PurchaseInvoiceCode}
		exclude = map[string]interface{}{}
		if purchaseInvoice, _, err = repository.GetDataPurchaseInvoice(filter, exclude); err != nil || (err == nil && len(purchaseInvoice) == 0) {
			continue
		} else {
			v.PurchaseInvoice = purchaseInvoice[0]
			if v.PurchaseInvoice.Status != int8(1) && v.PurchaseInvoice.Status != int8(6) {
				continue
			} else {
				v.PurchaseInvoice.PurchaseOrder.Read("ID")
				v.PurchaseInvoice.PurchaseOrder.Supplier.Read("ID")

				v.GoodsReceipt = &model.GoodsReceipt{PurchaseOrder: v.PurchaseInvoice.PurchaseOrder, Status: int8(2)}
				v.GoodsReceipt.Read("PurchaseOrder", "Status")

				_, paidAmount, err = repository.CheckPurchasePaymentAmount(v.PurchaseInvoice.ID)
				v.RemainingAmount = v.PurchaseInvoice.TotalCharge - paidAmount
			}
		}

		if pmID, err = common.Decrypt(v.PaymentMethodID); err != nil {
			continue
		} else {
			if v.PaymentMethod, err = repository.ValidPaymentMethod(pmID); err != nil {
				continue
			} else {
				if v.PaymentMethod.Status != int8(1) {
					continue
				}
			}
		}

		if v.Amount < 0 {
			continue
		}

		if v.PaidOff == "y" {
			paidOff = int8(1)
		}

		v.Amount = common.Rounder(v.Amount, 0.5, 2)

		if paidOff == 1 || (paidOff == 2 && v.Amount >= v.RemainingAmount) {
			v.PurchaseInvoice.Status = 2
			if v.GoodsReceipt.ID != 0 {
				v.PurchaseInvoice.PurchaseOrder.Status = 2
			}

			finishedInvoice[v.PurchaseInvoiceCode] = 1
		} else {
			v.PurchaseInvoice.Status = 6
		}

		request.PurchaseInvoice = v.PurchaseInvoice
		request.PaymentMethod = v.PaymentMethod
		request.Note = v.Note
		request.RecognitionDate = v.PaymentDate
		request.Amount = v.Amount
		request.PaidOff = paidOff
		request.ImageUrl = ""
		request.BankPaymentVoucherNumber = v.BankPaymentVoucherNumber
		request.Session = c.Session
		request.RemainingAmount = v.RemainingAmount
		request.GoodsReceipt = v.GoodsReceipt

		c.Requests = append(c.Requests, request)
	}

	if len(c.Requests) == 0 {
		o.Failure("id.invalid", "No data has been saved successfully")
	}

	return o
}

// Messages : function to return error validation messages after validation
func (c *createBulkRequest) Messages() map[string]string {
	return map[string]string{
		"data.required": util.ErrorInputRequired("data"),
	}
}
