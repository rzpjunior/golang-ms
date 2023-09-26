// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package invoice

import (
	"fmt"
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// createRequest : struct to hold price set request data
type createRequest struct {
	Code                   string    `json:"-"`
	SalesOrderID           string    `json:"sales_order_id" valid:"required"`
	SalesPaymentTermID     string    `json:"term_payment_sls_id" valid:"required"`
	InvoiceTermID          string    `json:"term_invoice_sls_id" valid:"required"`
	PaymentGroupID         string    `json:"payment_group_sls_id" valid:"required"`
	RecognitionDateStr     string    `json:"recognition_date" valid:"required"`
	DueDateStr             string    `json:"due_date" valid:"required"`
	BillingAddress         string    `json:"billing_address" valid:"required"`
	Note                   string    `json:"note"`
	DeliveryFee            float64   `json:"delivery_fee" valid:"required"`
	AdjustmentAmount       float64   `json:"adj_amount"`
	Adjustment             int8      `json:"adjustment"`
	AdjustmentNote         string    `json:"adj_note"`
	DiscountAmount         float64   `json:"disc_amount"`
	TotalPrice             float64   `json:"-"`
	TotalCharge            float64   `json:"-"`
	RecognitionDate        time.Time `json:"-"`
	DueDate                time.Time `json:"-"`
	TotalSkuDiscAmount     float64   `json:"-"`
	CreditLimitBefore      float64   `json:"-"`
	CreditLimitAfter       float64   `json:"-"`
	IsCreateCreditLimitLog int64     `json:"-"`

	SalesOrder       *model.SalesOrder   `json:"-"`
	SalesPaymentTerm *model.SalesTerm    `json:"-"`
	InvoiceTerm      *model.InvoiceTerm  `json:"-"`
	PaymentGroup     *model.PaymentGroup `json:"-"`

	InvoiceItems []*invoiceItems `json:"invoice_items" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

// invoiceItems : struct to hold sales invoice item set request data
type invoiceItems struct {
	ID               string  `json:"id" `
	ProductID        string  `json:"product_id" valid:"required"`
	SalesOrderItemID string  `json:"sales_order_item_id"`
	OrderQty         float64 `json:"order_qty" valid:"required"`
	DeliverQty       float64 `json:"deliver_qty" valid:"required"`
	ReceiveQty       float64 `json:"receive_qty"`
	InvoiceQty       float64 `json:"invoice_qty" valid:"required"`
	UnitPrice        float64 `json:"unit_price" valid:"required"`
	SkuDiscAmount    float64 `json:"-"`
	Note             string  `json:"note"`
	TaxPercentage    float64 `json:"tax_percentage"`

	TaxableItem             int8 `json:"-"`
	Subtotal                float64
	SalesOrderItemIDConvert int64
	Product                 *model.Product
	SalesInvoiceItem        *model.SalesInvoiceItem
	SalesOrderItem          *model.SalesOrderItem
}

// Validate : function to validate uom request data
func (c *createRequest) Validate() *validation.Output {
	var (
		e                      error
		isExistDocInvoice      bool
		soID, productID        int64
		totalChargeDifferences float64 = 0
	)
	o := &validation.Output{Valid: true}
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	statusErrorMsg := map[int8]string{
		1:  "active",
		3:  "cancelled",
		9:  "invoiced not delivered",
		10: "invoiced on delivery",
		11: "invoiced delivered",
		12: "paid not delivered",
	}

	// validation and read for sales order
	if soID, e = common.Decrypt(c.SalesOrderID); e != nil {
		o.Failure("sales_order_id.invalid", util.ErrorInvalidData("sales order"))
		return o
	}

	// if get sales order success
	if c.SalesOrder, e = repository.ValidSalesOrder(soID); e != nil {
		o.Failure("sales_order_id.invalid", util.ErrorInvalidData("sales order"))
		return o
	}

	// check so status
	if c.SalesOrder.Status != 1 && c.SalesOrder.Status != 7 && c.SalesOrder.Status != 8 {
		//var statusField
		o.Failure("id.invalid", util.ErrorCreateDocStatus("sales invoice", "sales order", "paid on delivery"))
		if statusErrorMsg[c.SalesOrder.Status] != "" {
			o.Failure("id.invalid", util.ErrorCreateDocStatus("sales invoice", "sales order", statusErrorMsg[c.SalesOrder.Status]))
		}

		return o
	}

	if e = c.SalesOrder.Branch.Read("ID"); e != nil {
		o.Failure("sales_order_id.invalid", util.ErrorInvalidData("sales order"))
		return o
	}

	// check if so already have si
	if e = orSelect.Raw("SELECT EXISTS(SELECT id FROM sales_invoice WHERE sales_order_id = ? AND status IN ('1','2','6'))", c.SalesOrder.ID).QueryRow(&isExistDocInvoice); e == nil && isExistDocInvoice {
		o.Failure("id.invalid", util.ErrorCreateDoc("sales invoice", "sales order"))
		return o
	}

	if len(c.AdjustmentNote) > 0 && c.Adjustment < 0 {
		o.Failure("adj_note.invalid", util.ErrorInputRequired("adjustment note"))
	}

	var duplicated = make(map[string]bool) // variable for check duplicate product
	for n, row := range c.InvoiceItems {
		if row.ProductID == "" {
			o.Failure(fmt.Sprintf("invoice_items.%d.product_id.required", n), util.ErrorSelectRequired("product"))
			continue
		}

		if duplicated[row.ProductID] {
			o.Failure(fmt.Sprintf("invoice_items.%d.product_id.required", n), util.ErrorDuplicate("product"))
			continue
		}

		if row.SalesOrderItemID == "" {
			o.Failure(fmt.Sprintf("invoice_items.%d.product_id.invalid", n), util.ErrorInputRequired("sales order item"))
			continue
		}

		if row.SalesOrderItemIDConvert, e = common.Decrypt(row.SalesOrderItemID); e != nil {
			o.Failure(fmt.Sprintf("invoice_items.%d.product_id.invalid", n), util.ErrorInputRequired("sales order item"))
			continue
		}

		if row.SalesOrderItem, e = repository.ValidSalesOrderItem(row.SalesOrderItemIDConvert); e != nil {
			o.Failure(fmt.Sprintf("invoice_items.%d.product_id.invalid", n), util.ErrorInputRequired("sales order item"))
			continue
		}

		if productID, e = common.Decrypt(row.ProductID); e != nil {
			o.Failure(fmt.Sprintf("invoice_items.%d.product_id.invalid", n), util.ErrorInvalidData("product"))
			continue
		}

		if row.Product, e = repository.ValidProduct(productID); e != nil {
			o.Failure(fmt.Sprintf("invoice_items.%d.product_id.invalid", n), util.ErrorInvalidData("product"))
			continue
		}

		row.TaxableItem = row.SalesOrderItem.TaxableItem
		row.TaxPercentage = row.SalesOrderItem.TaxPercentage
		row.Subtotal = (row.InvoiceQty * row.UnitPrice)

		// check if so use sku discount
		if row.SalesOrderItem.SkuDiscountItem != nil {
			if e = row.SalesOrderItem.SkuDiscountItem.Read("ID"); e != nil {
				o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("sku discount item"))
				continue
			}

			discQty := row.SalesOrderItem.DiscountQty
			if row.InvoiceQty < discQty {
				discQty = row.InvoiceQty
			}
			row.SkuDiscAmount = discQty * row.SalesOrderItem.UnitPriceDiscount
			row.Subtotal -= row.SkuDiscAmount
		}
		c.TotalSkuDiscAmount += row.SkuDiscAmount
		c.TotalPrice += row.Subtotal

		duplicated[row.ProductID] = true

		// validate for length note item
		if len(row.Note) > 100 {
			o.Failure(fmt.Sprintf("invoice_items.%d.note.invalid", n), util.ErrorCharLength("item note", 100))
		}
	}

	if c.Adjustment == 2 {
		c.AdjustmentAmount = c.AdjustmentAmount * -1
	}

	c.TotalCharge = c.TotalPrice + c.DeliveryFee + c.AdjustmentAmount

	if c.SalesOrder.VouRedeemCode != "" {
		c.SalesOrder.Voucher.Read("ID")
		c.TotalCharge = c.TotalCharge - c.SalesOrder.VouDiscAmount
	}

	if c.SalesOrder.PointRedeemAmount != 0 {
		c.TotalCharge = c.TotalCharge - c.SalesOrder.PointRedeemAmount
	}

	if c.TotalCharge < 0 {
		o.Failure("id.invalid", util.ErrorEqualGreater("grand total invoice", "0"))
	}
	// validation and read for sales payment term
	if slsPaymentTermID, e := common.Decrypt(c.SalesPaymentTermID); e != nil {
		o.Failure("term_payment_sls_id.invalid", util.ErrorInvalidData("term payment sls"))
	} else {
		if c.SalesPaymentTerm, e = repository.ValidSalesTerm(slsPaymentTermID); e != nil {
			o.Failure("term_payment_sls_id.invalid", util.ErrorInvalidData("term payment sls"))
		} else {
			if c.SalesPaymentTerm.Status != int8(1) {
				o.Failure("term_payment_sls_id.active", util.ErrorActive("term payment sls"))
			}
		}
	}

	// validation and read for invoice payment
	if invoiceTermID, e := common.Decrypt(c.InvoiceTermID); e != nil {
		o.Failure("term_invoice_sls_id.invalid", util.ErrorInvalidData("term invoice sls"))
	} else {
		if c.InvoiceTerm, e = repository.ValidInvoiceTerm(invoiceTermID); e != nil {
			o.Failure("term_invoice_sls_id.invalid", util.ErrorInvalidData("term invoice sls"))
		} else {
			if c.InvoiceTerm.Status != int8(1) {
				o.Failure("term_invoice_sls_id.active", util.ErrorActive("term invoice sls"))
			}
		}
	}

	// validation and read for payment group
	if paymentGroupID, e := common.Decrypt(c.PaymentGroupID); e != nil {
		o.Failure("payment_group_sls_id.invalid", util.ErrorInvalidData("payment group sls"))
	} else {
		if c.PaymentGroup, e = repository.ValidPaymentGroup(paymentGroupID); e != nil {
			o.Failure("payment_group_sls_id.invalid", util.ErrorInvalidData("payment group sls"))
		} else {
			if c.PaymentGroup.Status != int8(1) {
				o.Failure("payment_group_sls_id.active", util.ErrorActive("payment group sls"))
			}
		}
	}

	// validate for parse data recognition date
	if c.RecognitionDateStr != "" {
		if c.RecognitionDate, e = time.Parse("2006-01-02", c.RecognitionDateStr); e != nil {
			o.Failure("recognition_date.invalid", util.ErrorInvalidData("invoice date"))
		}
	}

	// validate for parse data due date
	if c.DueDateStr != "" {
		if c.DueDate, e = time.Parse("2006-01-02", c.DueDateStr); e != nil {
			o.Failure("due_date.invalid", util.ErrorInvalidData("due date"))
		}
	}

	// validate for length billing address
	if len(c.BillingAddress) > 350 {
		o.Failure("billing_address.invalid", util.ErrorCharLength("billing address", 350))
	}

	// validate for length note
	if len(c.Note) > 250 {
		o.Failure("note.invalid", util.ErrorCharLength("note", 250))
	}

	if e = c.SalesOrder.Branch.Merchant.Read("ID"); e != nil {
		o.Failure("merchant.invalid", util.ErrorInvalidData("merchant"))
	}

	if e = c.SalesOrder.Branch.Merchant.PaymentTerm.Read("ID"); e != nil {
		o.Failure("payment_term_id.invalid", util.ErrorInvalidData("payment term"))
	}

	if e = c.SalesOrder.Branch.Merchant.BusinessType.Read("ID"); e != nil {
		o.Failure("business_type_id.invalid", util.ErrorInvalidData("business type"))
	}

	c.CreditLimitBefore = c.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount
	if c.SalesOrder.Branch.Merchant.CreditLimitAmount > 0 || c.CreditLimitBefore < 0 {
		c.IsCreateCreditLimitLog = 1

		c.CreditLimitAfter = c.CreditLimitBefore

		if c.TotalCharge > c.SalesOrder.TotalCharge {
			totalChargeDifferences = c.TotalCharge - c.SalesOrder.TotalCharge
			c.CreditLimitAfter = c.CreditLimitBefore - totalChargeDifferences
		}

		if c.TotalCharge < c.SalesOrder.TotalCharge {
			totalChargeDifferences = c.SalesOrder.TotalCharge - c.TotalCharge
			c.CreditLimitAfter = c.CreditLimitBefore + totalChargeDifferences
		}

		if c.CreditLimitAfter < 0 && c.CreditLimitBefore > 0 {
			o.Failure("credit_limit_amount.invalid", util.ErrorCreditLimitExceeded(c.SalesOrder.Branch.Merchant.Name))
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
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
