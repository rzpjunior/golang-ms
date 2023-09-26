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

// createRequest : struct to hold Create Purchase Order request data

type createRequest struct {
	Code               string    `json:"-"`
	SupplierID         string    `json:"supplier_id" valid:"required"`
	PurchaseTermID     string    `json:"term_payment_pur_id" valid:"required"`
	WarehouseID        string    `json:"warehouse_id" valid:"required"`
	WarehouseAddress   string    `json:"warehouse_address" valid:"required"`
	StrRecognitionDate string    `json:"order_date" valid:"required"`
	StrEtaDate         string    `json:"eta_date" valid:"required"`
	EtaTime            string    `json:"eta_time" valid:"required"`
	DeliveryFee        float64   `json:"delivery_fee"`
	Note               string    `json:"note"`
	TaxPct             float64   `json:"tax_pct"`
	RecognitionDate    time.Time `json:"-"`
	EtaDate            time.Time `json:"-"`
	CreatedFrom        int8      `json:"created_from"`
	PurchasePlanID     string    `json:"purchase_plan_id"`
	Latitude           float64   `json:"latitude"`
	Longitude          float64   `json:"longitude"`

	EtaTimeFormat time.Time `json:"-"`

	PurchaseOrderItems []*requestItem `json:"purchase_order_items" valid:"required"`
	Images             []string       `json:"images"`

	TotalPrice    float64             `json:"-"`
	TaxAmount     float64             `json:"-"`
	TotalCharge   float64             `json:"-"`
	TotalWeight   float64             `json:"-"`
	RecognitionAt time.Time           `json:"-"`
	EtaDateAt     time.Time           `json:"-"`
	Supplier      *model.Supplier     `json:"-"`
	Warehouse     *model.Warehouse    `json:"-"`
	PurchaseTerm  *model.PurchaseTerm `json:"-"`
	Area          *model.Area         `json:"-"`
	PurchasePlan  *model.PurchasePlan `json:"-"`
	Session       *auth.SessionData   `json:"-"`
}

type requestItem struct {
	ID                 string            `json:"id"`
	ProductID          string            `json:"product_id" valid:"required`
	OrderQty           float64           `json:"qty" valid:"required`
	UnitPrice          float64           `json:"unit_price" valid:"required`
	Note               string            `json:"note" valid:"lte:500"`
	MarketPurchase     []*marketPurchase `json:"market_purchase"`
	PurchaseQty        float64           `json:"purchase_qty"`
	IncludeTax         int8              `json:"include_tax"`
	TaxPercentage      float64           `json:"tax_percentage"`
	PurchasePlanItemID string            `json:"purchase_plan_item_id"`

	TaxableItem       int8                     `json:"-"`
	TaxAmount         float64                  `json:"-"`
	UnitPriceTax      float64                  `json:"-"`
	Subtotal          float64                  `json:"-"`
	PurchaseOrderItem *model.PurchaseOrderItem `json:"-"`
	Product           *model.Product           `json:"-"`
	Price             *model.Price             `json:"-"`
	PurchasePlanItem  *model.PurchasePlanItem  `json:"-"`
	Uom               *model.Uom               `json:"-"`
}

type marketPurchase struct {
	Stall string  `json:"stall"`
	Qty   float64 `json:"qty"`
	Price float64 `json:"price"`
}

func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	var filter, exclude map[string]interface{}
	productList := make(map[string]string)

	SupplierID, e := common.Decrypt(c.SupplierID)
	if e != nil {
		o.Failure("supplier_id.invalid", util.ErrorInvalidData("supplier"))
	}

	c.Supplier = &model.Supplier{ID: SupplierID}
	if e = c.Supplier.Read("ID"); e != nil {
		o.Failure("supplier_id.invalid", util.ErrorInvalidData("supplier"))
	}

	PurchaseTermID, e := common.Decrypt(c.PurchaseTermID)
	if e != nil {
		o.Failure("term_payment_pur_id.invalid", util.ErrorInvalidData("purchase term"))
	}

	c.PurchaseTerm = &model.PurchaseTerm{ID: PurchaseTermID}
	if e = c.PurchaseTerm.Read("ID"); e != nil {
		o.Failure("term_payment_pur_id.invalid", util.ErrorInvalidData("purchase term"))
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

	if len(c.Note) > 250 {
		o.Failure("note", util.ErrorCharLength("note", 250))
	}

	if c.PurchasePlanID != "" {
		purchasePlanID, e := common.Decrypt(c.PurchasePlanID)
		if e != nil {
			o.Failure("purchase_plan_id.invalid", util.ErrorInvalidData("purchase plan"))
		}

		c.PurchasePlan = &model.PurchasePlan{ID: purchasePlanID}
		if e = c.PurchasePlan.Read("ID"); e != nil {
			o.Failure("purchase_plan_id.invalid", util.ErrorInvalidData("purchase plan"))
		}

		if c.PurchasePlan.Status != 1 {
			o.Failure("purchase_plan_id.invalid", util.ErrorActive("purchase plan"))
		}

		if len(c.Images) == 0 {
			o.Failure("purchase_order_image.required", util.ErrorInputRequired("image"))
		}
	}

	for i, v := range c.PurchaseOrderItems {
		var productID int64
		var strProductID string

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
		if v.Product, e = repository.ValidProduct(productID); e != nil {
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

		unitPriceInput := v.UnitPrice

		v.TaxableItem = v.Product.Taxable
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

		if v.PurchasePlanItemID != "" {
			purchasePlanItemID, e := common.Decrypt(v.PurchasePlanItemID)
			if e != nil {
				o.Failure("purchase_plan_item_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("purchase plan item"))
			}

			v.PurchasePlanItem = &model.PurchasePlanItem{ID: purchasePlanItemID}
			if e = v.PurchasePlanItem.Read("ID"); e != nil {
				o.Failure("purchase_plan_item_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("purchase plan item"))
			}

			totalPurchaseQty := v.OrderQty + v.PurchasePlanItem.PurchaseQty
			if totalPurchaseQty > v.PurchasePlanItem.PurchasePlanQty {
				o.Failure("qty"+strconv.Itoa(i)+".invalid", util.ErrorEqualLess("total purchase qty", "purchase plan qty"))
			}

			v.UnitPrice = unitPriceInput
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

	if len(c.Images) > 4 {
		o.Failure("purchase_order_image.invalid", util.ErrorEqualLess("photo", "4 photos"))
	}

	return o
}

func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"supplier_id.required":         util.ErrorInputRequired("supplier"),
		"term_payment_pur_id.required": util.ErrorInputRequired("payment term"),
		"warehouse_id.required":        util.ErrorInputRequired("warehouse"),
		"warehouse_address.required":   util.ErrorInputRequired("warehouse address"),
		"order_date.required":          util.ErrorInputRequired("order date"),
		"eta_date.required":            util.ErrorInputRequired("eta date"),
		"eta_time.required":            util.ErrorInputRequired("eta time"),
	}

	for i, _ := range c.PurchaseOrderItems {
		messages["item."+strconv.Itoa(i)+".product_id.required"] = util.ErrorInputRequired("product")
		messages["item."+strconv.Itoa(i)+".qty.required"] = util.ErrorInputRequired("qty")
		messages["item."+strconv.Itoa(i)+".unit_price.required"] = util.ErrorInputRequired("unit price")
	}

	return messages
}
