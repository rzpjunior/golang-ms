// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package invoice

import (
	"git.edenfarm.id/cuxs/validation"

	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// addTaxInvoiceRequest : struct to Add Tax Invoice request Data
type addTaxInvoiceRequest struct {
	ID               int64  `json:"-" valid:"required"`
	TaxInvoiceURL    string `json:"tax_invoice_url" valid:"required"`
	TaxInvoiceNumber string `json:"tax_invoice_number" valid:"required"`
	PurchaseInvoice  *model.PurchaseInvoice

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate add tax invoice request data
func (u *addTaxInvoiceRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if len(u.TaxInvoiceNumber) > 50 {
		o.Failure("tax_invoice_number", util.ErrorCharLength("tax invoice number", 50))
	}

	if len(u.TaxInvoiceURL) > 500 {
		o.Failure("tax_invoice_url", util.ErrorCharLength("tax invoice url", 500))
	}

	return o
}

// Messages : function to return error validation messages
func (c *addTaxInvoiceRequest) Messages() map[string]string {
	messages := map[string]string{
		"tax_invoice_url.required":    util.ErrorInputRequired("tax invoice url"),
		"tax_invoice_number.required": util.ErrorInputRequired("tax invoice number"),
	}

	return messages
}
