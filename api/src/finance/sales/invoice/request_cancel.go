// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package invoice

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createRequest : struct to hold price set request data
type cancelRequest struct {
	ID                     int64   `json:"-" valid:"required"`
	CancellationNote       string  `json:"note" valid:"required"`
	CreditLimitBefore      float64 `json:"-"`
	CreditLimitAfter       float64 `json:"-"`
	IsCreateCreditLimitLog int64   `json:"-"`

	SalesInvoice *model.SalesInvoice
	Session      *auth.SessionData
}

// Validate : function to validate uom request data
func (c *cancelRequest) Validate() *validation.Output {
	var totalChargeDifferences float64 = 0
	var e error
	o := &validation.Output{Valid: true}
	c.SalesInvoice.SalesOrder.Read("ID")

	if c.SalesInvoice.Status != 1 {
		o.Failure("status.inactive", util.ErrorActive("sales invoice"))
	}
	if c.SalesInvoice.SalesOrder.Status != 9 && c.SalesInvoice.SalesOrder.Status != 10 && c.SalesInvoice.SalesOrder.Status != 11 {
		if c.SalesInvoice.SalesOrder.Status == 1 {
			o.Failure("id.invalid", util.ErrorCreateDocStatus("sales invoice", "sales order", "active"))
		} else if c.SalesInvoice.SalesOrder.Status == 2 {
			o.Failure("id.invalid", util.ErrorCreateDocStatus("sales invoice", "sales order", "finished"))
		} else if c.SalesInvoice.SalesOrder.Status == 3 {
			o.Failure("id.invalid", util.ErrorCreateDocStatus("sales invoice", "sales order", "cancelled"))
		} else if c.SalesInvoice.SalesOrder.Status == 7 {
			o.Failure("id.invalid", util.ErrorCreateDocStatus("sales invoice", "sales order", "on delivery"))
		} else if c.SalesInvoice.SalesOrder.Status == 8 {
			o.Failure("id.invalid", util.ErrorCreateDocStatus("sales invoice", "sales order", "delivered"))
		} else if c.SalesInvoice.SalesOrder.Status == 12 {
			o.Failure("id.invalid", util.ErrorCreateDocStatus("sales invoice", "sales order", "paid not delivered"))
		} else if c.SalesInvoice.SalesOrder.Status == 13 {
			o.Failure("id.invalid", util.ErrorCreateDocStatus("sales invoice", "sales order", "paid on delivery"))
		}
	}

	if c.SalesInvoice.SalesOrder.Branch, e = repository.ValidBranch(c.SalesInvoice.SalesOrder.Branch.ID); e != nil {
		o.Failure("branch.invalid", util.ErrorInvalidData("branch"))
	}

	if c.SalesInvoice.SalesOrder.Branch.Merchant, e = repository.ValidMerchant(c.SalesInvoice.SalesOrder.Branch.Merchant.ID); e != nil {
		o.Failure("merchant.invalid", util.ErrorInvalidData("merchant"))
	}

	c.CreditLimitBefore = c.SalesInvoice.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount
	if c.SalesInvoice.SalesOrder.Branch.Merchant.CreditLimitAmount > 0 || c.CreditLimitBefore < 0 {
		c.IsCreateCreditLimitLog = 1

		c.CreditLimitAfter = c.CreditLimitBefore

		if c.SalesInvoice.TotalCharge > c.SalesInvoice.SalesOrder.TotalCharge {
			totalChargeDifferences = c.SalesInvoice.TotalCharge - c.SalesInvoice.SalesOrder.TotalCharge
		}

		if c.SalesInvoice.SalesOrder.TotalCharge > c.SalesInvoice.TotalCharge {
			totalChargeDifferences = c.SalesInvoice.TotalCharge - c.SalesInvoice.SalesOrder.TotalCharge
		}

		if totalChargeDifferences != 0 {
			c.CreditLimitAfter = c.CreditLimitBefore + totalChargeDifferences
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *cancelRequest) Messages() map[string]string {
	return map[string]string{
		"note.required": util.ErrorInputRequired("cancellation note"),
	}
}
