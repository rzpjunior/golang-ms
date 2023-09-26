// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package invoice

import (
	"math"
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

// updateRequest : struct to hold Update Purchase Invoice request data
type updateRequest struct {
	ID                   int64                  `json:"-" valid:"required"`
	PurchaseOrderID      string                 `json:"purchase_order_id" valid:"required"`
	RecognitionDate      string                 `json:"recognition_date" valid:"required"`
	DueDate              string                 `json:"due_date" valid:"required"`
	Note                 string                 `json:"note"`
	TaxPct               float64                `json:"tax_pct"`
	DeliveryFee          float64                `json:"delivery_fee"`
	AdjAmount            float64                `json:"adj_amount"`
	Adjustment           int8                   `json:"-"`
	Deduction            int8                   `json:"deduction"`
	AdjNote              string                 `json:"adj_note"`
	TotalPrice           float64                `json:"-"`
	TaxAmount            float64                `json:"-"`
	TotalCharge          float64                `json:"-"`
	InvoiceItems         []*purchaseInvoiceItem `json:"invoice_items" valid:"required"`
	DebitNotes           []*debitNotes          `json:"debit_notes"`
	DebitNoteArr         []*model.DebitNote     `json:"-"`
	ExcludedDebitNoteIDs []int64                `json:"-"`

	RecognitionAt time.Time
	DueDateAt     time.Time

	PurchaseOrder   *model.PurchaseOrder
	PurchaseInvoice *model.PurchaseInvoice

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate Update Purchase Invoice request data
func (u *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	var adjAmountTemp, ded float64

	adjAmountTemp = 0
	ded = 1

	poID, e := common.Decrypt(u.PurchaseOrderID)
	if e != nil {
		o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase order"))
	}

	u.PurchaseOrder, e = repository.ValidPurchaseOrder(poID)
	if e != nil {
		o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase order"))
	}

	if u.PurchaseOrder.Status != 1 {
		o.Failure("id.invalid", util.ErrorDocStatus("purchase order", "active"))
	}

	e = u.PurchaseOrder.Supplier.Read("ID")
	if e != nil {
		o.Failure("po.supplier_id.invalid", util.ErrorInvalidData("supplier"))
	}

	if u.PurchaseInvoice.Status != 1 {
		o.Failure("id.invalid", util.ErrorCreateDoc("purchase invoice", "purchase order"))
	}

	if u.AdjAmount < 0 {
		o.Failure("adj_amount.invalid", util.ErrorEqualGreater("adjustment amount", "0"))
	}

	if u.AdjAmount > 0 {
		if u.AdjNote == "" {
			o.Failure("adj_note", util.ErrorInputRequired("adjustment note"))
		}
	}

	if u.TaxPct < 0 {
		o.Failure("tax_pct.invalid", util.ErrorEqualGreater("tax percentage", "0"))
	}

	for _, i := range u.InvoiceItems {
		if i.ID == "" {
			o.Failure("id.invalid", util.ErrorInputRequired("id"))
			return o
		}

		piiID, e := common.Decrypt(i.ID)
		if e != nil {
			o.Failure("purchase_invoice_item_id.invalid", util.ErrorInvalidData("purchase invoice item"))
		}

		i.PurchaseInvoiceItem, e = repository.ValidPurchaseInvoiceItem(piiID)
		if e != nil {
			o.Failure("purchase_invoice_item_id.invalid", util.ErrorInvalidData("purchase invoice item"))
		}

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

		i.TaxableItem = i.PurchaseInvoiceItem.TaxableItem
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
		u.TaxAmount += i.TaxAmount
		u.TotalPrice += i.Subtotal

		if len(i.Note) > 100 {
			o.Failure("note", util.ErrorCharLength("note", 100))
		}
	}

	// region debit note validation
	duplicateDebitNote := make(map[int64]bool, 0)
	includedDebitNote := make(map[int64]bool, 0)
	var totalPriceDebitNote float64
	debitNoteIDsOriginal := strings.Split(u.PurchaseInvoice.DebitNoteIDs, ",")
	debitNoteIDsOriginalMap := make(map[int64]bool, 0)
	for _, id := range debitNoteIDsOriginal {
		dnID, _ := strconv.Atoi(id)
		debitNoteIDsOriginalMap[int64(dnID)] = true
	}
	for i, v := range u.DebitNotes {
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
		if _, ok := debitNoteIDsOriginalMap[debitNoteID]; !ok {
			if dn.UsedInPurchaseInvoice == 1 {
				o.Failure("debit_note_id.invalid", util.ErrorIsBeingUsed("debit note"))
				return o
			}
		}

		if v.TotalPrice != dn.TotalPrice {
			o.Failure("debit_note_id"+strconv.Itoa(i)+".invalid", util.ErrorMustBeSame("total price", "latest total price"))
		}
		if _, ok := duplicateDebitNote[debitNoteID]; ok {
			o.Failure("debit_note_id"+strconv.Itoa(i)+".invalid", util.ErrorDuplicate("debit note"))
		}
		duplicateDebitNote[debitNoteID] = true
		includedDebitNote[debitNoteID] = true

		u.DebitNoteArr = append(u.DebitNoteArr, dn)
		totalPriceDebitNote += v.TotalPrice

	}

	for _, v := range debitNoteIDsOriginal {
		dnID, _ := strconv.Atoi(v)
		if _, ok := includedDebitNote[int64(dnID)]; !ok {
			u.ExcludedDebitNoteIDs = append(u.ExcludedDebitNoteIDs, int64(dnID))
		}
	}
	if len(u.ExcludedDebitNoteIDs) == 0 {
		u.ExcludedDebitNoteIDs = append(u.ExcludedDebitNoteIDs, int64(0))
	}

	// endregion

	// validasi adjamount act as deduction or addition
	adjAmountTemp = adjAmountTemp + u.AdjAmount
	if u.AdjAmount > 0 {
		u.Adjustment = 1
	}
	if u.Deduction == 1 {
		u.Adjustment = 2
		ded = -1
	}
	u.TotalCharge = u.TotalPrice + ((u.TaxPct / 100) * u.TotalPrice) + u.DeliveryFee + (adjAmountTemp * ded) + u.TaxAmount - totalPriceDebitNote

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
	return o
}

// Messages : function to return error validation messages
func (c *updateRequest) Messages() map[string]string {
	messages := map[string]string{
		"recognition_date.required": util.ErrorInputRequired("recognition date"),
		"delivery_fee.required":     util.ErrorInputRequired("delivery fee"),
		"adj_note.required":         util.ErrorInputRequired("adj note"),
	}

	for i, _ := range c.InvoiceItems {
		messages["item."+strconv.Itoa(i)+".id.required"] = util.ErrorInputRequired("purchase invoice item")
		messages["item."+strconv.Itoa(i)+".product_id.required"] = util.ErrorInputRequired("product")
		messages["item."+strconv.Itoa(i)+".purchase_order_item_id.required"] = util.ErrorInputRequired("purchase order item")
		messages["item."+strconv.Itoa(i)+".invoice_qty.required"] = util.ErrorInputRequired("invoice qty")
		messages["item."+strconv.Itoa(i)+".unit_price.required"] = util.ErrorInputRequired("unit price")
	}

	return messages
}
