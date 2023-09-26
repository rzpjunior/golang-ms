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

type salabilityRequest struct {
	ID           int64  `json:"-" valid:"required"`
	Salability   int8   `json:"salability" valid:"required"`
	WarehouseStr string `json:"-"`

	Warehouse []string       `json:"warehouse"`
	Stocks    []*model.Stock `json:"stock"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate supplier request data
func (c *salabilityRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var filter, exclude map[string]interface{}

	if dbredis.Redis.CheckExistKey("warehouse_creation") {
		o.Failure("warehouse_creation.invalid", util.ErrorCreationInProgress("warehouse"))
	}

	if c.Salability == 1 {
		if len(c.Warehouse) == 0 {
			o.Failure("message.invalid", util.ErrorSelectOne("warehouse"))
		} else {
			for i, v := range c.Warehouse {
				v = common.Encrypt(v)
				c.Warehouse[i] = v
				c.WarehouseStr = c.WarehouseStr + v + ","
			}

			filter = map[string]interface{}{"product_id": c.ID, "status": 2, "warehouse_id__in": c.Warehouse}
			if _, total, err := repository.CheckStockData(filter, exclude); err == nil && total > 0 {
				o.Failure("message.invalid", util.ErrorMustExistWarehouse("salability assigned warehouse", "storability assigned warehouse"))
			}

			c.WarehouseStr = strings.TrimSuffix(c.WarehouseStr, ",")
		}
	} else {
		if len(c.Warehouse) > 0 {
			o.Failure("message.invalid", util.ErrorNoAssignedWh("salability assigned warehouse"))
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *salabilityRequest) Messages() map[string]string {
	return map[string]string{
		"salability.required": util.ErrorInputRequired("salability"),
	}
}
