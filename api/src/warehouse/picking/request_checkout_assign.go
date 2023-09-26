// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"strconv"
)

// checkoutRequestAssign : struct to hold picking assign request data
type checkoutRequestAssign struct {
	ID int64 `json:"-"`

	TotalColly        float64                 `json:"total_colly" valid:"required"`
	DeliveryKolies    []*DeliveryKoli         `json:"delivery_kolies" valid:"required"`
	PickingOrderItems []*approvalPickingItems `json:"picking_order_items"`

	PickingOrderAssign *model.PickingOrderAssign `json:"-"`
	SalesOrder         *model.SalesOrder         `json:"-"`

	IsFinished bool `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate picking assign request data
func (r *checkoutRequestAssign) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if r.TotalColly <= 0 {
		o.Failure("id.invalid", util.ErrorEqualGreaterInd("total koli", "0"))
		return o
	}

	if r.PickingOrderAssign, err = repository.ValidPickingOrderAssign(r.ID); err != nil {
		o.Failure("picking_order_assign.invalid", util.ErrorInvalidData("picking order assign"))
		return o
	}

	if r.PickingOrderAssign.Status == int8(5) {
		o.Failure("picking_order_assign.active", util.ErrorPickingStatus("new", "on progress"))
	}

	r.PickingOrderAssign.SalesOrder.Read("ID")
	r.PickingOrderAssign.PickingOrder.Read("ID")
	if r.PickingOrderAssign.SalesOrder.Status != 1 &&
		r.PickingOrderAssign.SalesOrder.Status != 9 &&
		r.PickingOrderAssign.SalesOrder.Status != 12 {
		o.Failure("status.invalid", util.ErrorSalesOrderOnPicking())
	}

	for i2, v2 := range r.DeliveryKolies {
		if koliID, err := common.Decrypt(v2.KoliID); err == nil {
			if v2.Koli, err = repository.ValidKoli(koliID); err != nil {
				o.Failure("koli_id"+strconv.Itoa(i2)+".invalid", util.ErrorInvalidData("koli"))
			}
		} else {
			o.Failure("koli_id"+strconv.Itoa(i2)+".invalid", util.ErrorInvalidData("koli"))
		}
		v2.SalesOrder = r.PickingOrderAssign.SalesOrder
	}

	for i, v := range r.PickingOrderItems {
		if pickingOrderAssign, err := common.Decrypt(v.PickingOrderItemID); err == nil {
			if v.PickingOrderItem, err = repository.ValidPickingOrderItem(pickingOrderAssign); err != nil {
				o.Failure("picking_order_item_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("picking order item"))
			}
		} else {
			o.Failure("picking_order_item_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("picking order item"))
		}
	}

	if r.SalesOrder, err = repository.ValidSalesOrder(r.PickingOrderAssign.SalesOrder.ID); err != nil {
		o.Failure("sales_order_id.invalid", util.ErrorInvalidData("sales order"))
		return o
	}

	return o
}

// Messages : function to return error validation messages
func (r *checkoutRequestAssign) Messages() map[string]string {
	messages := map[string]string{
		"total_koli.required": util.ErrorInputRequired("total koli"),
	}

	return messages
}
