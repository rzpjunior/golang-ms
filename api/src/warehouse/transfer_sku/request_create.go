// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package transfer_sku

import (
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/orm"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createRequest : struct to hold Create Transfer SKU request data
type createRequest struct {
	Code        string            `json:"-"`
	WarehouseID string            `json:"warehouse_id" valid:"required"`
	Products    []*productRequest `json:"products" valid:"required"`
	Note        string            `json:"note"`

	Warehouse         *model.Warehouse `json:"-"`
	TotalTransferQty  float64          `json:"-"`
	TotalWasteQty     float64          `json:"-"`
	RecognitionDateAt time.Time        `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// productRequest : struct to hold Product Origin data
type productRequest struct {
	ID                  string               `json:"id" valid:"required"`
	TransferTo          []*transferToRequest `json:"transfer_to" valid:"required"`
	AvailableQty        float64              `json:"available_qty" valid:"required"`
	ConvertWeightParent float64              `json:"-"`
	ConvertWeightChild  float64              `json:"-"`
	Discrepancy         float64              `json:"-"`

	Product *model.Product `json:"-"`
	Stock   *model.Stock   `json:"-"`
}

// transferToRequest : struct to hold Transferred Product data
type transferToRequest struct {
	ProductID   string  `json:"product_id" valid:"required"`
	TransferQty float64 `json:"transfer_qty" valid:"required, numeric"`
	WasteQty    float64 `json:"waste_qty" valid:"required, numeric"`
	WasteReason int8    `json:"waste_reason"`

	Product *model.Product `json:"-"`
	Stock   *model.Stock   `json:"-"`
}

// Validate : function to validate request data
func (c *createRequest) Validate() *validation.Output {

	o := &validation.Output{Valid: true}
	var isProductExist = make(map[string]bool)
	var isTransferProductExist = make(map[string]bool)
	var totalTransferQty float64 // for document level
	var totalWasteQty float64    // for document level
	layout := "2006-01-02"
	var err error
	var o1 = orm.NewOrm()
	o1.Using("read_only")

	if c.Code, err = util.CheckTable("transfer_sku"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	if c.RecognitionDateAt, err = time.Parse(layout, time.Now().Format(layout)); err != nil {
		o.Failure("recognition_date.invalid", util.ErrorInvalidData("transfer sku recognition date"))
	}

	if c.WarehouseID != "" {
		warehouseID, err := common.Decrypt(c.WarehouseID)

		if err != nil {
			o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
		}

		if c.Warehouse, err = repository.ValidWarehouse(warehouseID); err != nil {
			o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
		}

		if c.Warehouse.Status != 1 {
			o.Failure("warehouse_id.invalid", util.ErrorActive("warehouse"))
		}

	}

	for i, v := range c.Products {
		if _, productExist := isProductExist[v.ID]; productExist {
			o.Failure("products_"+strconv.Itoa(i)+"_id.duplicate", util.ErrorDuplicate("product"))
			return o
		}

		var productID int64
		if productID, err = common.Decrypt(v.ID); err != nil {
			o.Failure("products_"+strconv.Itoa(i)+"_id.invalid", util.ErrorInvalidData("product"))
			return o
		}

		if v.Product, err = repository.ValidProduct(productID); err != nil {
			o.Failure("products_"+strconv.Itoa(i)+"_id.invalid", util.ErrorInvalidData("product"))
			return o
		}

		isProductExist[v.ID] = true
		var initTransferProductExist = make(map[string]bool)
		isTransferProductExist = initTransferProductExist

		// region convert qty to kg weight
		v.ConvertWeightParent = v.AvailableQty * v.Product.UnitWeight
		// endregion

		for j, k := range v.TransferTo {
			if _, transferProductExist := isTransferProductExist[k.ProductID]; transferProductExist {
				o.Failure("products_"+strconv.Itoa(i)+"_transfer_"+strconv.Itoa(j)+".duplicate", util.ErrorDuplicate("transfer product"))
				return o
			}

			var transferProductID int64
			if transferProductID, err = common.Decrypt(k.ProductID); err != nil {
				o.Failure("products_"+strconv.Itoa(i)+"_transfer_"+strconv.Itoa(j)+".product_id.invalid", util.ErrorInvalidData("transfer product"))
				return o
			}

			if k.Product, err = repository.ValidProduct(transferProductID); err != nil {
				o.Failure("products_"+strconv.Itoa(i)+"_transfer_"+strconv.Itoa(j)+"product_id.invalid", util.ErrorInvalidData("transfer product"))
				return o
			}

			isTransferProductExist[k.ProductID] = true

			if k.WasteQty > 0 {
				v.ConvertWeightChild += k.WasteQty * k.Product.UnitWeight
				if k.WasteReason == 0 {
					o.Failure("waste_reason_id"+strconv.Itoa(i)+".invalid", util.ErrorInputRequired("waste reason"))
				} else {
					// check waste reason in glossary
					_, e := repository.GetGlossaryMultipleValue("table", "all", "attribute", "waste_reason", "value_int", k.WasteReason)
					if e != nil {
						o.Failure("waste_reason.invalid", util.ErrorInvalidData("waste_reason"))
					}
				}
			} else {
				k.WasteReason = 0
			}

			v.ConvertWeightChild += k.TransferQty * k.Product.UnitWeight

			// region to add total transfer qty for document level
			if v.Product.ID != k.Product.ID {
				totalTransferQty += k.TransferQty
			}
			// region

			// region to add total waste qty for document level
			if v.Product.ID == k.Product.ID {
				totalWasteQty += k.WasteQty
			}
			// region
		}
		v.Discrepancy = v.ConvertWeightParent - v.ConvertWeightChild
		if v.Discrepancy < -0.001 {
			o.Failure("discrepancy.invalid", util.ErrorGreater("discrepancy", "0"))
		}

	}

	c.TotalTransferQty = totalTransferQty
	c.TotalWasteQty = totalWasteQty

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"warehouse_id.required": util.ErrorInputRequired("warehouse"),
		"products_required":     util.ErrorInputRequired("products"),
	}

	for i, v := range c.Products {
		messages["products_"+strconv.Itoa(i)+"_id.required"] = util.ErrorInputRequired("product")
		messages["products_"+strconv.Itoa(i)+"_total_transfer_qty.required"] = util.ErrorInputRequired("total transfer qty")

		for j := range v.TransferTo {
			messages["products_"+strconv.Itoa(i)+"_transfer_to_"+strconv.Itoa(j)+"_product_id.required"] = util.ErrorInputRequired("transfer product")
			messages["products_"+strconv.Itoa(i)+"_transfer_to_"+strconv.Itoa(j)+"_transfer_qty.required"] = util.ErrorInputRequired("transfer qty")
			messages["products_"+strconv.Itoa(i)+"_transfer_to_"+strconv.Itoa(j)+"_waste_qty.required"] = util.ErrorInputRequired("waste qty")
		}
	}

	return messages
}
