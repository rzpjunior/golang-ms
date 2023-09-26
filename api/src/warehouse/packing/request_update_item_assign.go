// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packing

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"strconv"
)

// createRequest : struct to hold price set request data
type updateItemAssignRequest struct {
	ID              	 int64  			`json:"-"`
	PackingOrderItemID   string   			`json:"packing_order_item_id" valid:"required"`
	Helper 				 []string   		`json:"helper"`
	HelperDec 			 []int64
	Session            	 *auth.SessionData

	Poi            		 *model.PackingOrderItem
}

// Validate : function to validate uom request data
func (c *updateItemAssignRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	poiID, _ := common.Decrypt(c.PackingOrderItemID)
	c.Poi = &model.PackingOrderItem{ID: poiID}
	c.Poi.Read("ID")

	c.Poi.PackingOrder.Read("ID")

	if c.Poi.PackingOrder.Status != 1 {
		o.Failure("id.invalid", util.ErrorDocStatus("active", "packing order"))
	}

	if len(c.Helper) > 0 {
		for n, row := range c.Helper {
			helperId, _ := common.Decrypt(row)

			if helperPerson, err := repository.ValidStaff(helperId); err != nil {
				o.Failure("helper"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("helper"))
			} else {
				if helperPerson.Status != int8(1) {
					o.Failure("helper"+strconv.Itoa(n)+".active", util.ErrorActive("helper"))
				}

				c.HelperDec = append(c.HelperDec, helperId)
			}
		}
		if _, total, err := repository.CheckPackingOrderItemHelper(poiID, c.HelperDec...); err == nil {
			if total > 0 {
				o.Failure("id.invalid", util.ErrorHelperAssign())
			}
		}
	} else {
		if _, total, err := repository.CheckPackingOrderItemHelper(poiID); err == nil {
			if total > 0 {
				o.Failure("id.invalid", util.ErrorHelperAssign())
			}
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *updateItemAssignRequest) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
