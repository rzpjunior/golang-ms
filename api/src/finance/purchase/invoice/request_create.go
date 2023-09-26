// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package invoice

import (
	"math"
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"

	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// createRequest : struct to hold Create Purchase Invoice request data
type createRequest struct {
	Code            string                 `json:"-"`
	PurchaseOrderID string                 `json:"purchase_order_id" valid:"required"`
	RecognitionDate string                 `json:"recognition_date" valid:"required"`
	DueDate         string                 `json:"due_date" valid:"required"`
	Note            string                 `json:"note"`
	TaxPct          float64                `json:"tax_pct"`
	DeliveryFee     float64                `json:"delivery_fee"`
	AdjAmount       float64                `json:"adj_amount"`
	Adjustment      int8                   `json:"-"`
	Deduction       int8                   `json:"deduction"`
	AdjNote         string                 `json:"adj_note"`
	TotalPrice      float64                `json:"-"`
	TaxAmount       float64                `json:"-"`
	TotalCharge     float64                `json:"-"`
	InvoiceItems    []*purchaseInvoiceItem `json:"invoice_items" valid:"required"`
	DebitNotes      []*debitNotes          `json:"debit_notes"`
	DebitNoteArr    []*model.DebitNote     `json:"-"`

	RecognitionAt time.Time
	DueDateAt     time.Time

	PurchaseOrder *model.PurchaseOrder

	Session *auth.SessionData `json:"-"`
}

type debitNotes struct {
	DebitNoteID string  `json:"debit_note_id"`
	TotalPrice  float64 `json:"total_price"`
}

type purchaseInvoiceItem struct {
	ID                  string  `json:"id"`
	ProductID           string  `json:"product_id" valid:"required"`
	PurchaseOrderItemID string  `json:"purchase_order_item_id" valid:"required"`
	InvoiceQty          float64 `json:"invoice_qty" valid:"required"`
	UnitPrice           float64 `json:"unit_price" valid:"required"`
	Note                string  `json:"note"`
	TaxPercentage       float64 `json:"tax_percentage"`
	IncludeTax          int8    `json:"include_tax"`

	TaxableItem         int8    `json:"-"`
	TaxAmount           float64 `json:"-"`
	UnitPriceTax        float64 `json:"-"`
	Subtotal            float64 `json:"-"`
	Product             *model.Product
	PurchaseOrderItem   *model.PurchaseOrderItem
	PurchaseInvoiceItem *model.PurchaseInvoiceItem
}

// Validate : function to validate Create Purchase Invoice request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	var adjAmountTemp, ded float64
	var totalDocInvoice int64

	adjAmountTemp = 0
	ded = 1

	poID, e := common.Decrypt(c.PurchaseOrderID)
	if e != nil {
		o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase order"))
	}

	if c.PurchaseOrder, e = repository.ValidPurchaseOrder(poID); e != nil {
		o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase order"))
	}

	if c.PurchaseOrder.Status != 1 {
		o.Failure("id.invalid", util.ErrorDocStatus("purchase order", "active"))
	}

	filter := map[string]interface{}{"status__in": []int{1, 2, 6}, "purchase_order_id": c.PurchaseOrder.ID}
	exclude := map[string]interface{}{}
	_, totalDocInvoice, e = repository.GetDataPurchaseInvoice(filter, exclude)

	if e != nil || totalDocInvoice >= 1 {
		o.Failure("id.invalid", util.ErrorCreateDoc("purchase invoice", "purchase order"))
	}

	if e = c.PurchaseOrder.Supplier.Read("ID"); e != nil {
		o.Failure("supplier_id.invalid", util.ErrorInvalidData("supplier"))
	}

	if c.AdjAmount < 0 {
		o.Failure("adj_amount.invalid", util.ErrorEqualGreater("adjustment amount", "0"))
	}

	if c.AdjAmount > 0 && c.AdjNote == "" {
		o.Failure("adj_note.invalid", util.ErrorInputRequired("adjustment note"))
	}

	if c.TaxPct < 0 {
		o.Failure("tax_pct.invalid", util.ErrorEqualGreater("tax percentage", "0"))
	}

	for _, i := range c.InvoiceItems {

		pID, e := common.Decrypt(i.ProductID)
		if e != nil {
			o.Failure("product_id.invalid", util.ErrorInvalidData("product"))
		}

		i.Product, e = repository.ValidProduct(pID)
		if e != nil {
			o.Failure("product_id.invalid", util.ErrorInvalidData("product"))
		}

		poiID, e := common.Decrypt(i.PurchaseOrderItemID)
		if e != nil {
			o.Failure("purchase_order_item_id.invalid", util.ErrorInvalidData("purchase order item"))
		}

		i.PurchaseOrderItem, e = repository.ValidPurchaseOrderItem(poiID)
		if e != nil {
			o.Failure("purchase_order_item_id.invalid", util.ErrorInvalidData("purchase order item"))
		}

		if i.InvoiceQty < 0 {
			o.Failure("invoice_qty.invalid", util.ErrorEqualGreater("invoice qty", "0"))
		}

		unitPriceInput := i.UnitPrice

		i.TaxableItem = i.PurchaseOrderItem.TaxableItem
		i.TaxAmount = math.Round(unitPriceInput * i.TaxPercentage / 100 * i.InvoiceQty)
		i.UnitPriceTax = math.Round(unitPriceInput * (100 + i.TaxPercentage) / 100)

		isIncludeTax := i.IncludeTax == 1
		isNotTaxableItem := i.TaxableItem != 1

		if isIncludeTax {
			unitPriceNonTax := math.Round(unitPriceInput * 100 / (100 + i.TaxPercentage))
			unitPriceTax := unitPriceInput

			i.TaxAmount = math.Round((unitPriceTax - unitPriceNonTax) * i.InvoiceQty)
			i.UnitPriceTax = unitPriceTax
			i.UnitPrice = unitPriceNonTax
		}

		if isNotTaxableItem {
			i.TaxAmount = 0
			i.UnitPriceTax = 0
		}

		i.Subtotal = i.UnitPrice * i.InvoiceQty

		// Summarize all the item tax amount
		c.TaxAmount += i.TaxAmount
		c.TotalPrice += i.Subtotal

		if len(i.Note) > 100 {
			o.Failure("note", util.ErrorCharLength("note", 100))
		}
	}

	// region debit note validation
	duplicateDebitNote := make(map[int64]bool, 0)
	var totalPriceDebitNote float64
	for i, v := range c.DebitNotes {
		var dn *model.DebitNote
		var debitNoteID int64
		if debitNoteID, e = common.Decrypt(v.DebitNoteID); e != nil {
			o.Failure("debit_note_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("debit note"))
		}
		dn = &model.DebitNote{
			ID: debitNoteID,
		}
		if e = dn.Read("ID"); e != nil {
			o.Failure("debit_note_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("debit note"))
		}
		if dn.Status != 1 {
			o.Failure("debit_note_id"+strconv.Itoa(i)+".invalid", util.ErrorActive("debit note"))
		}
		if dn.UsedInPurchaseInvoice == 1 {
			o.Failure("debit_note_id.invalid", util.ErrorIsBeingUsed("debit note"))
			return o
		}

		if v.TotalPrice != dn.TotalPrice {
			o.Failure("debit_note_id"+strconv.Itoa(i)+".invalid", util.ErrorMustBeSame("total price", "latest total price"))
		}
		if _, ok := duplicateDebitNote[debitNoteID]; ok {
			o.Failure("debit_note_id"+strconv.Itoa(i)+".invalid", util.ErrorDuplicate("debit note"))
		}
		duplicateDebitNote[debitNoteID] = true

		c.DebitNoteArr = append(c.DebitNoteArr, dn)
		totalPriceDebitNote += v.TotalPrice

	}

	// endregion
	// validasi adjamount act as deduction or addition
	adjAmountTemp = adjAmountTemp + c.AdjAmount
	if c.AdjAmount > 0 {
		c.Adjustment = 1
	}
	if c.Deduction == 1 {
		c.Adjustment = 2
		ded = -1
	}
	c.TotalCharge = c.TotalPrice + ((c.TaxPct / 100) * c.TotalPrice) + c.DeliveryFee + (adjAmountTemp * ded) + c.TaxAmount - totalPriceDebitNote

	if c.TotalCharge < 0 {
		o.Failure("total_charge.invalid", util.ErrorEqualGreater("total charge", "0"))
	}

	// validate for parse data recognition date
	if c.RecognitionDate != "" {
		if c.RecognitionAt, e = time.Parse("2006-01-02", c.RecognitionDate); e != nil {
			o.Failure("recognition_date.invalid", util.ErrorInvalidData("invoice date"))
		}
	}

	// validate for parse data due date
	if c.DueDate != "" {
		if c.DueDateAt, e = time.Parse("2006-01-02", c.DueDate); e != nil {
			o.Failure("due_date.invalid", util.ErrorInvalidData("due date"))
		}
	}
	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"recognition_date.required": util.ErrorInputRequired("recognition date"),
		"delivery_fee.required":     util.ErrorInputRequired("delivery fee"),
		"adj_note.required":         util.ErrorInputRequired("adj note"),
	}

	for i, _ := range c.InvoiceItems {
		messages["item."+strconv.Itoa(i)+".product_id.required"] = util.ErrorInputRequired("product")
		messages["item."+strconv.Itoa(i)+".purchase_order_item_id.required"] = util.ErrorInputRequired("purchase order item")
		messages["item."+strconv.Itoa(i)+".invoice_qty.required"] = util.ErrorInputRequired("invoice qty")
		messages["item."+strconv.Itoa(i)+".unit_price.required"] = util.ErrorInputRequired("unit price")
	}

	return messages
}
