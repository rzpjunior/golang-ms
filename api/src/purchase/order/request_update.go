// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"math"
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// updateRequest : struct to hold Update Purchase Invoice request data
type updateRequest struct {
	ID                 int64     `json:"-"`
	Code               string    `json:"-"`
	WarehouseID        string    `json:"warehouse_id" valid:"required"`
	StrRecognitionDate string    `json:"order_date" valid:"required"`
	StrEtaDate         string    `json:"eta_date" valid:"required"`
	EtaTime            string    `json:"eta_time" valid:"required"`
	DeliveryFee        float64   `json:"delivery_fee"`
	Note               string    `json:"note"`
	TaxPct             float64   `json:"tax_pct"`
	RecognitionDate    time.Time `json:"-"`
	EtaDate            time.Time `json:"-"`

	EtaTimeFormat time.Time `json:"-"`

	PurchaseOrderItems []*requestItem `json:"purchase_order_items" valid:"required"`

	TotalPrice    float64              `json:"-"`
	TaxAmount     float64              `json:"-"`
	TotalCharge   float64              `json:"-"`
	TotalWeight   float64              `json:"-"`
	RecognitionAt time.Time            `json:"-"`
	EtaDateAt     time.Time            `json:"-"`
	Supplier      *model.Supplier      `json:"-"`
	Warehouse     *model.Warehouse     `json:"-"`
	PurchaseTerm  *model.PurchaseTerm  `json:"-"`
	PurchaseOrder *model.PurchaseOrder `json:"-"`
	Uom           *model.Uom           `json:"-"`
	Session       *auth.SessionData    `json:"-"`
}

func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	var filter, exclude map[string]interface{}
	productList := make(map[string]string)

	c.PurchaseOrder = &model.PurchaseOrder{ID: c.ID}
	if e = c.PurchaseOrder.Read("ID"); e != nil {
		o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase order"))
	}

	if c.PurchaseOrder.Status != 5 {
		o.Failure("id.invalid", util.ErrorDraft("purchase order"))
	}

	warehouseID, e := common.Decrypt(c.WarehouseID)
	if e != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}

	c.Warehouse = &model.Warehouse{ID: warehouseID}
	if e = c.Warehouse.Read("ID"); e != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}

	if c.RecognitionDate, e = time.Parse("2006-01-02", c.StrRecognitionDate); e != nil {
		o.Failure("order_date.invalid", util.ErrorInvalidData("order date"))
	}

	if c.EtaDate, e = time.Parse("2006-01-02", c.StrEtaDate); e != nil {
		o.Failure("eta_date.invalid", util.ErrorInvalidData("estimated arrival date"))
	}

	// only for checking format time from apps
	if c.EtaTimeFormat, e = time.Parse("15:04", c.EtaTime); e != nil {
		o.Failure("eta_time.invalid", util.ErrorInvalidData("eta time"))
	}

	for i, v := range c.PurchaseOrderItems {
		var productID int64
		var strProductID string
		var purchaseOrderItemID int64

		if v.OrderQty <= 0 {
			o.Failure("qty"+strconv.Itoa(i)+".greater", util.ErrorGreater("product quantity", "0"))
		}

		if v.UnitPrice < 0 {
			o.Failure("unit_price"+strconv.Itoa(i)+".equalorgreater", util.ErrorEqualGreater("product unit price", "0"))
		}

		if v.TaxPercentage < 0 {
			o.Failure("tax_percentage"+strconv.Itoa(i)+".equalorgreater", util.ErrorEqualGreater("product tax percentage", "0"))
		}

		if len(v.Note) > 100 {
			o.Failure("note"+strconv.Itoa(i), util.ErrorCharLength("note", 100))
		}

		if v.ProductID == "" {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInputRequired("product"))
		}

		if productID, e = common.Decrypt(v.ProductID); e != nil {
			o.Failure("product_id.invalid", util.ErrorInvalidData("product"))
		}

		strProductID = strconv.Itoa(int(productID))
		v.Product, e = repository.ValidProduct(productID)
		if e != nil {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
		}

		v.Uom, e = repository.ValidUom(v.Product.Uom.ID)
		if e != nil {
			o.Failure("uom_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("uom"))
		}

		if v.Uom.DecimalEnabled == 2 {
			if math.Mod(v.OrderQty, 1) != 0 {
				o.Failure("order_qty"+strconv.Itoa(i)+".invalid", util.ErrorNotAllowedFor("decimal", "product qty"))
			}
		}

		v.TaxableItem = v.Product.Taxable

		if v.ID != "" {
			if purchaseOrderItemID, e = common.Decrypt(v.ID); e != nil {
				o.Failure("id.invalid", util.ErrorInvalidData("purchase order item"))
			}

			if v.PurchaseOrderItem, e = repository.ValidPurchaseOrderItem(purchaseOrderItemID); e != nil {
				o.Failure("id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
			}

			v.TaxableItem = v.PurchaseOrderItem.TaxableItem
		}

		unitPriceInput := v.UnitPrice

		v.TaxAmount = math.Round((unitPriceInput * v.TaxPercentage / 100) * v.OrderQty)
		v.UnitPriceTax = math.Round(unitPriceInput * (100 + v.TaxPercentage) / 100)

		isIncludeTax := v.IncludeTax == 1
		isNotTaxableItem := v.TaxableItem != 1

		if isIncludeTax {
			unitPriceNonTax := math.Round(unitPriceInput * 100 / (100 + v.TaxPercentage))
			unitPriceTax := unitPriceInput

			v.TaxAmount = math.Round((unitPriceTax - unitPriceNonTax) * v.OrderQty)
			v.UnitPriceTax = unitPriceTax
			v.UnitPrice = unitPriceNonTax
		}

		if isNotTaxableItem {
			v.TaxAmount = 0
			v.UnitPriceTax = 0
		}

		v.Subtotal = v.OrderQty * v.UnitPrice

		// Summarize all the item tax amount
		c.TaxAmount += v.TaxAmount
		c.TotalPrice = c.TotalPrice + v.Subtotal
		c.TotalWeight = c.TotalWeight + (v.OrderQty * v.Product.UnitWeight)

		if _, exist := productList[strProductID]; exist {
			o.Failure("product_id"+strconv.Itoa(i)+".duplicate", util.ErrorDuplicate("product"))
		}

		productList[strProductID] = "t"

		filter = map[string]interface{}{"product_id": productID, "warehouse_id": warehouseID, "purchasable": 1}
		if _, countStock, err := repository.CheckStockData(filter, exclude); err == nil && countStock == 0 {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorProductMustAvailable())
		}
	}

	c.TotalCharge = c.TotalPrice + c.DeliveryFee + (c.TaxPct * c.TotalPrice / 100) + c.TaxAmount

	return o
}

func (c *updateRequest) Messages() map[string]string {
	messages := map[string]string{
		"warehouse_id.required": util.ErrorInputRequired("warehouse"),
		"order_date.required":   util.ErrorInputRequired("order date"),
		"eta_date.required":     util.ErrorInputRequired("eta date"),
		"eta_time.required":     util.ErrorInputRequired("eta time"),
	}

	for i, _ := range c.PurchaseOrderItems {
		messages["item."+strconv.Itoa(i)+".product_id.required"] = util.ErrorInputRequired("product")
		messages["item."+strconv.Itoa(i)+".qty.required"] = util.ErrorInputRequired("qty")
		messages["item."+strconv.Itoa(i)+".unit_price.required"] = util.ErrorInputRequired("unit price")
	}

	return messages
}
