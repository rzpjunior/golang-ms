// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package invoice

import (
	"fmt"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// updateRequest : struct to hold price set request data
type updateRequest struct {
	ID                 int64           `json:"-" valid:"required"`
	SalesOrderID       string          `json:"sales_order_id" valid:"required"`
	CodeExt            string          `json:"code_ext"`
	RecognitionDate    string          `json:"recognition_date" valid:"required"`
	DueDate            string          `json:"due_date" valid:"required"`
	BillingAddress     string          `json:"billing_address" valid:"required"`
	Note               string          `json:"note"`
	DeliveryFee        float64         `json:"delivery_fee" valid:"required"`
	AdjustmentAmount   float64         `json:"adj_amount"`
	Adjustment         int8            `json:"adjustment"`
	AdjustmentNote     string          `json:"adj_note"`
	InvoiceItems       []*invoiceItems `json:"invoice_items" valid:"required"`
	TotalSkuDiscAmount float64         `json:"-"`

	IsCreateCreditLimitLog int64
	CreditLimitBefore      float64
	CreditLimitAfter       float64
	OldTotalCharge         float64
	TotalPrice             float64
	TotalCharge            float64
	RecognitionAt          time.Time
	DueDateAt              time.Time
	SalesOrder             *model.SalesOrder
	SalesPaymentTerm       *model.SalesTerm
	InvoiceTerm            *model.InvoiceTerm
	PaymentGroup           *model.PaymentGroup
	Session                *auth.SessionData
	SalesInvoice           *model.SalesInvoice
}

// Validate : function to validate update sales invoice request data
func (u *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	var (
		isExistDocInvoice      bool
		e                      error
		soID                   int64
		totalChargeDifferences float64 = 0
		countInProgessPayment  int8
	)

	//check sales invoice must be active status
	if u.SalesInvoice.Status != 1 {
		o.Failure("id.invalid", util.ErrorActive("sales invoice"))
	}

	u.OldTotalCharge = u.SalesInvoice.TotalCharge

	// validation and read for sales order
	if soID, e = common.Decrypt(u.SalesOrderID); e != nil {
		o.Failure("sales_order_id.invalid", util.ErrorInvalidData("sales order"))
		return o
	}

	// if get sales order success
	if u.SalesOrder, e = repository.ValidSalesOrder(soID); e != nil {
		o.Failure("sales_order_id.invalid", util.ErrorInvalidData("sales order"))
		return o
	}

	// count total document invoice
	if e = orSelect.Raw("SELECT EXISTS(SELECT id FROM sales_invoice WHERE sales_order_id = ? AND id != ? AND status IN ('1','2','6'))", u.SalesOrder.ID, u.ID).QueryRow(&isExistDocInvoice); e == nil && isExistDocInvoice {
		o.Failure("id.invalid", util.ErrorCreateDoc("sales invoice", "sales order"))
		return o
	}

	if len(u.AdjustmentNote) > 0 && u.Adjustment < 0 {
		o.Failure("adj_note.invalid", util.ErrorInputRequired("adjustment note"))
	}

	if e = u.SalesOrder.Branch.Read("ID"); e != nil {
		o.Failure("branch_id.invalid", util.ErrorInvalidData("branch"))
	}
	if e = u.SalesOrder.Branch.Merchant.Read("ID"); e != nil {
		o.Failure("merchant.invalid", util.ErrorInvalidData("merchant"))
	}

	var duplicated = make(map[string]bool) // variable for check duplicate product
	for n, row := range u.InvoiceItems {
		var productID int64

		if row.ProductID == "" {
			o.Failure(fmt.Sprintf("invoice_items.%d.product_id.required", n), util.ErrorSelectRequired("product"))
			continue
		}

		if duplicated[row.ProductID] {
			o.Failure(fmt.Sprintf("invoice_items.%d.product_id.required", n), util.ErrorDuplicate("product"))
			continue
		}

		if row.SalesInvoiceItem, e = repository.ValidSalesInvoiceItem(row.ID); e != nil {
			o.Failure(fmt.Sprintf("invoice_items.%d.product_id.required", n), util.ErrorInvalidData("sales invoice item"))
			continue
		}

		if productID, e = common.Decrypt(row.ProductID); e != nil {
			o.Failure(fmt.Sprintf("invoice_items.%d.product_id.required", n), util.ErrorInvalidData("product"))
			continue
		}

		if row.Product, e = repository.ValidProduct(productID); e != nil {
			o.Failure(fmt.Sprintf("invoice_items.%d.product_id.invalid", n), util.ErrorInvalidData("product"))
			continue
		}

		row.Subtotal = row.InvoiceQty * row.UnitPrice

		// get sales order item to check if soi has sku discount
		if e = row.SalesInvoiceItem.SalesOrderItem.Read("ID"); e != nil {
			o.Failure(fmt.Sprintf("invoice_items.%d.product_id.invalid", n), util.ErrorInvalidData("sales order item"))
			continue
		}

		if row.SalesInvoiceItem.SalesOrderItem.SkuDiscountItem != nil {
			if e = row.SalesInvoiceItem.SalesOrderItem.SkuDiscountItem.Read("ID"); e != nil {
				o.Failure(fmt.Sprintf("invoice_items.%d.product_id.invalid", n), util.ErrorInvalidData("sku discount item"))
				continue
			}

			discQty := row.SalesInvoiceItem.SalesOrderItem.DiscountQty
			if row.InvoiceQty < discQty {
				discQty = row.InvoiceQty
			}
			row.SkuDiscAmount = discQty * row.SalesInvoiceItem.SalesOrderItem.UnitPriceDiscount
			row.Subtotal -= row.SkuDiscAmount
		}

		u.TotalSkuDiscAmount += row.SkuDiscAmount
		u.TotalPrice += row.Subtotal

		duplicated[row.ProductID] = true

		// validate for length note item
		if len(row.Note) > 100 {
			o.Failure("note.invalid", util.ErrorCharLength("invoice_items.%d.note.invalid", 100))
			continue
		}
	}

	if u.Adjustment == 2 {
		u.AdjustmentAmount = u.AdjustmentAmount * -1
	}

	u.TotalCharge = u.TotalPrice + u.DeliveryFee + u.AdjustmentAmount

	if u.SalesOrder.VouRedeemCode != "" {
		u.TotalCharge = u.TotalCharge - u.SalesOrder.VouDiscAmount
	}

	if u.SalesOrder.PointRedeemAmount != 0 {
		u.TotalCharge = u.TotalCharge - u.SalesOrder.PointRedeemAmount
	}

	if u.TotalCharge < 0 {
		o.Failure("id.invalid", util.ErrorEqualGreater("grand total invoice", "0"))
	}

	// validate for parse data recognition date
	if u.RecognitionDate != "" {
		if u.RecognitionAt, e = time.Parse("2006-01-02", u.RecognitionDate); e != nil {
			o.Failure("recognition_date.invalid", util.ErrorInvalidData("invoice date"))
		}
	}

	// validate for parse data due date
	if u.DueDate != "" {
		if u.DueDateAt, e = time.Parse("2006-01-02", u.DueDate); e != nil {
			o.Failure("due_date.invalid", util.ErrorInvalidData("due date"))
		}
	}

	// validate for length billing address
	if len(u.BillingAddress) > 350 {
		o.Failure("billing_address.invalid", util.ErrorCharLength("billing address", 350))
	}

	// validate for length note
	if len(u.Note) > 250 {
		o.Failure("note.invalid", util.ErrorCharLength("note", 250))
	}

	// Validate there is sales payment status in progress in sales invoice
	if countInProgessPayment, e = repository.CheckInProgressPayment(u.SalesInvoice.ID); e != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales invoice"))
	}

	if countInProgessPayment > 0 && u.OldTotalCharge != u.TotalCharge {
		o.Failure("id.invalid", util.ErrorRelated("sales payment", "in progress", "sales invoice"))
	}

	u.CreditLimitBefore = u.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount
	if u.SalesOrder.Branch.Merchant.CreditLimitAmount > 0 || u.CreditLimitBefore < 0 {
		u.IsCreateCreditLimitLog = 1

		totalChargeDifferences = u.TotalCharge - u.SalesInvoice.TotalCharge
		u.CreditLimitAfter = u.CreditLimitBefore - totalChargeDifferences

		if u.CreditLimitAfter < 0 && u.CreditLimitBefore > 0 {
			o.Failure("credit_limit_amount.invalid", util.ErrorCreditLimitExceeded(u.SalesOrder.Branch.Merchant.Name))
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"sales_order_id.required":       util.ErrorInputRequired("sales order"),
		"term_payment_sls_id.required":  util.ErrorSelectRequired("term payment"),
		"term_invoice_sls_id.required":  util.ErrorSelectRequired("term invoice"),
		"payment_group_sls_id.required": util.ErrorSelectRequired("payment group"),
		"recognition_date.required":     util.ErrorInputRequired("invoice date"),
		"due_date.required":             util.ErrorInputRequired("due date"),
		"billing_address.required":      util.ErrorInputRequired("billing address"),
		"delivery_fee.required":         util.ErrorInputRequired("delivery fee"),
	}
}
