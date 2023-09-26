// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// scanRequestChecker : struct to hold picking assign request data
type scanRequestChecker struct {
	ID            int64   `json:"-" valid:"required"`
	OrderQty      float64 `json:"order_qty"`
	CheckOrderQty float64 `json:"check_qty"`

	PickingOrderItem *model.PickingOrderItem `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate picking assign request data
func (r *scanRequestChecker) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var err error

	if r.PickingOrderItem, err = repository.ValidPickingOrderItem(r.ID); err != nil {
		o.Failure("picking_order_item.invalid", util.ErrorInvalidData("picking order item"))
		return o
	}
	if err = r.PickingOrderItem.Product.Read("ID"); err != nil {
		o.Failure("product.invalid", util.ErrorInvalidData("product"))
		return o
	}
	if r.PickingOrderItem.Product.Packability != 1 {
		o.Failure("product.invalid", util.ErrorInvalidData("product"))
		return o
	}

	return o
}

// Messages : function to return error validation messages
func (r *scanRequestChecker) Messages() map[string]string {
	messages := map[string]string{
		"id.required": util.ErrorInputRequired("id"),
	}

	return messages
}
