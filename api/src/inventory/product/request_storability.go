// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package product

import (
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/dbredis"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type storabilityRequest struct {
	ID           int64  `json:"-" valid:"required"`
	Storability  int8   `json:"storability" valid:"required"`
	WarehouseStr string `json:"-"`
	Status       int8   `json:"-"`

	WarehouseChecked []string       `json:"warehouse_checked"`
	Stock            []*model.Stock `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate supplier request data
func (c *storabilityRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var total int64
	var filter, exclude map[string]interface{}

	c.Status = 1

	if dbredis.Redis.CheckExistKey("warehouse_creation") {
		o.Failure("warehouse_creation.invalid", util.ErrorCreationInProgress("warehouse"))
	}

	for i, v := range c.WarehouseChecked {
		v = common.Encrypt(v)
		c.WarehouseChecked[i] = v
		c.WarehouseStr = c.WarehouseStr + v + ","
	}
	c.WarehouseStr = strings.TrimSuffix(c.WarehouseStr, ",")

	if _, total, err = repository.CheckDeliveryOrderProductStatus(c.ID, []int8{1, 5, 6, 7}, c.WarehouseChecked...); err == nil && total > 0 {
		o.Failure("message.invalid", util.ErrorRelated("active ", "delivery order", "product"))
	}

	if _, total, err = repository.CheckPurchaseOrderProductStatus(c.ID, 1, c.WarehouseChecked...); err == nil && total > 0 {
		o.Failure("message.invalid", util.ErrorRelated("active ", "purchase order", "product"))
	}

	if _, total, err = repository.CheckGoodsReceiptProductStatus(c.ID, 1, c.WarehouseChecked...); err == nil && total > 0 {
		o.Failure("message.invalid", util.ErrorRelated("active ", "goods receipt", "product"))
	}

	if _, total, err = repository.CheckDeliveryReturnProductStatus(c.ID, 1, c.WarehouseChecked...); err == nil && total > 0 {
		o.Failure("message.invalid", util.ErrorRelated("active ", "delivery return", "product"))
	}

	if _, total, err = repository.CheckWasteEntryProductStatus(c.ID, 1, c.WarehouseChecked...); err == nil && total > 0 {
		o.Failure("message.invalid", util.ErrorRelated("active ", "waste entry", "product"))
	}

	if _, total, err = repository.CheckGoodsTransferProductStatus(c.ID, 1, c.WarehouseChecked...); err == nil && total > 0 {
		o.Failure("message.invalid", util.ErrorRelated("active ", "goods transfer", "product"))
	}

	if _, total, err = repository.CheckStockOpnameProductStatus(c.ID, 1, c.WarehouseChecked...); err == nil && total > 0 {
		o.Failure("message.invalid", util.ErrorRelated("active ", "stock opname", "product"))
	}

	if stock, total, err := repository.CheckStockPerWarehouse(c.ID, c.WarehouseChecked...); err == nil && total > 0 {
		if stock.AvailableStock > 0 {
			o.Failure("message.invalid", util.ErrorMustZero("available stock"))
		}

		if stock.WasteStock > 0 {
			o.Failure("message.invalid", util.ErrorMustZero("waste stock"))
		}
	}

	filter = map[string]interface{}{"product_id": c.ID, "purchasable": 1}
	exclude = map[string]interface{}{}
	if len(c.WarehouseChecked) > 0 {
		exclude = map[string]interface{}{"warehouse_id__in": c.WarehouseChecked}
	}
	if _, total, err := repository.CheckStockData(filter, exclude); err == nil && total > 0 {
		o.Failure("message.invalid", util.ErrorRelated("", "warehouse purchasability", "warehouse storability"))
	}

	filter = map[string]interface{}{"product_id": c.ID, "salable": 1}
	exclude = map[string]interface{}{}
	if len(c.WarehouseChecked) > 0 {
		exclude = map[string]interface{}{"warehouse_id__in": c.WarehouseChecked}
	}
	if _, total, err := repository.CheckStockData(filter, exclude); err == nil && total > 0 {
		o.Failure("message.invalid", util.ErrorRelated("", "warehouse salability", "warehouse storability"))
	}

	if c.Storability == 2 {
		c.Status = 2
	}

	return o
}

// Messages : function to return error validation messages
func (c *storabilityRequest) Messages() map[string]string {
	return map[string]string{
		"storability.required": util.ErrorInputRequired("storability"),
	}
}
