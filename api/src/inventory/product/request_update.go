// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package product

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type updateRequest struct {
	ID                      int64    `json:"-" valid:"required"`
	Code                    string   `json:"code" valid:"required"`
	Name                    string   `json:"name" valid:"required"`
	UomId                   string   `json:"uom_id" valid:"required"`
	UnivProductCode         string   `json:"up_code"`
	Description             string   `json:"description"`
	Note                    string   `json:"note"`
	TagProduct              []string `json:"product_tag"`
	OrderChannelRestriction []string `json:"order_channel_restriction"`
	ImagesUrl               []string `json:"images" valid:"required"`
	CategoryId              string   `json:"category_id" valid:"required"`
	OrderMinQty             float64  `json:"order_min_qty" valid:"required|gt:0"`
	Taxable                 int8     `json:"taxable" valid:"required"`
	TaxPercentage           float64  `json:"tax_percentage"`
	SparePercentage         float64  `json:"spare_percentage"`
	OrderMaxQty             float64  `json:"order_max_qty"`
	ExcludeArchetype        []string `json:"exclude_archetype"`
	ExcludeArchetypeStr     string
	MaxDayDeliveryDate      int64 `json:"max_day_delivery_date"`

	Category         *model.Category       `json:"-"`
	StorabilityStock []*storabilityStock   `json:"storability_stock"`
	Images           []*model.ProductImage `json:"-"`

	Session *auth.SessionData `json:"-"`
}

type storabilityStock struct {
	WarehouseId    string             `json:"warehouse_id"`
	StockQty       float64            `json:"stock_qty" valid:"required"`
	WarehouseStock map[string]float64 `json:"-"`

	Warehouse *model.Warehouse `json:"-"`
}

// Validate : function to validate supplier request data
func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	isValidArchetype := true
	uomId, _ := common.Decrypt(c.UomId)
	filter := map[string]interface{}{"name": c.Name, "uom_id": uomId}
	exclude := map[string]interface{}{"id": c.ID, "status": int8(3)}
	if _, countName, err := repository.CheckProductData(filter, exclude); err != nil {
		o.Failure("name.invalid", util.ErrorInvalidData("name"))
	} else if countName > 0 {
		o.Failure("name.unique", util.ErrorUniqueProduct())
	}

	if len(strings.TrimSpace(c.Code)) != 8 {
		o.Failure("code.invalid", util.ErrorCharLength("code", 8))
	}

	filterCode := map[string]interface{}{"code": c.Code}
	excludeStatus := map[string]interface{}{"id": c.ID, "status": int8(3)}
	if _, countCode, err := repository.CheckProductData(filterCode, excludeStatus); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	} else if countCode > 0 {
		o.Failure("code", util.ErrorDuplicate("code"))
	}

	categoryId, e := common.Decrypt(c.CategoryId)
	if e != nil {
		o.Failure("category.invalid", util.ErrorInvalidData("category"))
	} else {
		if c.Category, e = repository.ValidCategory(categoryId); e != nil {
			o.Failure("category.invalid", util.ErrorInvalidData("category"))
		} else {
			if c.Category.Status != int8(1) {
				o.Failure("category.active", util.ErrorActive("category"))
			}
		}
	}

	for _, v := range c.StorabilityStock {
		v.WarehouseId = common.Encrypt(v.WarehouseId)
		warehouseId, _ := strconv.Atoi(v.WarehouseId)
		v.Warehouse = &model.Warehouse{
			ID: int64(warehouseId),
		}

		if err := v.Warehouse.Read("ID"); err != nil {
			o.Failure("warehouse.invalid", util.ErrorInvalidData("warehouse"))
		}

		v.WarehouseStock = map[string]float64{v.WarehouseId: v.StockQty}
	}

	excludeArchetypeInt := make([]int, len(c.ExcludeArchetype))
	for i, v := range c.ExcludeArchetype {
		vInt, _ := common.Decrypt(v)
		if _, e = repository.ValidArchetype(vInt); e != nil {
			o.Failure("archetype.invalid", util.ErrorInvalidData("archetype"))
			isValidArchetype = false
			break
		} else {
			excludeArchetypeInt[i] = int(vInt)
		}
	}

	if isValidArchetype {
		sort.Ints(excludeArchetypeInt)
		excludeArchetypeJson, _ := json.Marshal(excludeArchetypeInt)
		c.ExcludeArchetypeStr = strings.Trim(string(excludeArchetypeJson), "[]")
	}

	if c.OrderChannelRestriction != nil {
		for k, v := range c.OrderChannelRestriction {
			valInt, err := strconv.Atoi(v)
			if err != nil {
				o.Failure(fmt.Sprintf("order_channel_restriction.invalid[%d]", k), util.ErrorInvalidData("order_channel_restriction"))
			}
			orderChanel := util.IsOrderChannel(valInt)
			if !orderChanel {
				o.Failure(fmt.Sprintf("order_channel_restriction.invalid[%d]", k), util.ErrorInvalidData("order_channel_restriction"))
			}
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":        util.ErrorInputRequired("name"),
		"category_id.required": util.ErrorInputRequired("category"),
		"order_min_qty.gt":     util.ErrorGreater("minimum order quantity", "0"),
	}
}
