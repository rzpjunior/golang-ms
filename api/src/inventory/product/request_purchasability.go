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

type purchasabilityRequest struct {
	ID             int64  `json:"-" valid:"required"`
	Purchasability int8   `json:"purchasability" valid:"required"`
	WarehouseStr   string `json:"-"`

	WarehouseChecked   []string       `json:"warehouse_checked"`
	WarehouseUnchecked []string       `json:"warehouse_unchecked"`
	Stocks             []*model.Stock `json:"stock"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate supplier request data
func (c *purchasabilityRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var total int64
	var filter, exclude map[string]interface{}

	if dbredis.Redis.CheckExistKey("warehouse_creation") {
		o.Failure("warehouse_creation.invalid", util.ErrorCreationInProgress("warehouse"))
	}

	for i, v := range c.WarehouseUnchecked {
		c.WarehouseUnchecked[i] = util.DecryptIdInStr(v)
	}

	if c.Purchasability == 1 {
		if len(c.WarehouseChecked) == 0 {
			o.Failure("message.invalid", util.ErrorSelectOne("warehouse"))
		} else {
			for i, v := range c.WarehouseChecked {
				v = common.Encrypt(v)
				c.WarehouseChecked[i] = v
				c.WarehouseStr = c.WarehouseStr + v + ","
			}

			filter = map[string]interface{}{"product_id": c.ID, "status": 2, "warehouse_id__in": c.WarehouseChecked}
			if _, total, err = repository.CheckStockData(filter, exclude); err == nil && total > 0 {
				o.Failure("message.invalid", util.ErrorMustExistWarehouse("purchasability assigned warehouse", "storability assigned warehouse"))
			}

			if _, total, err = repository.CheckPurchaseOrderProductStatus(c.ID, 5, c.WarehouseChecked...); err == nil && total > 0 {
				o.Failure("message.invalid", util.ErrorRelated("draft ", "purchase order", "product"))
			}

			c.WarehouseStr = strings.TrimSuffix(c.WarehouseStr, ",")
		}
	} else {
		if len(c.WarehouseChecked) > 0 {
			o.Failure("message.invalid", util.ErrorNoAssignedWh("salability assigned warehouse"))
		} else {
			if _, total, err = repository.CheckPurchaseOrderProductStatus(c.ID, 5); err == nil && total > 0 {
				o.Failure("message.invalid", util.ErrorRelated("draft ", "purchase order", "product"))
			}
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *purchasabilityRequest) Messages() map[string]string {
	return map[string]string{
		"purchasability.required": util.ErrorInputRequired("purchasability"),
	}
}
