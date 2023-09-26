// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchase_deliver

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// consolidateRequest : struct to hold input consolidate purchase deliver request data
type consolidateRequest struct {
	Code             string             `json:"-"`
	DriverName       string             `json:"driver_name" valid:"required|alpha_num_space|lte:100"`
	VehicleNumber    string             `json:"vehicle_number" valid:"required|alpha_num|range:5,9"`
	PurchaseDelivers []*purchaseDeliver `json:"purchase_delivers" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

type purchaseDeliver struct {
	PurchaseDeliverID string `json:"purchase_deliver_id" valid:"required"`

	PurchaseDeliver *model.PurchaseDeliver `json:"-"`
	PurchaseOrder   *model.PurchaseOrder   `json:"-"`
	Warehouse       *model.Warehouse       `json:"-"`
}

// Validate : function to validate create field purchase order request data
func (c *consolidateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.Code, err = util.CheckTable("consolidated_purchase_deliver"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	for k, v := range c.PurchaseDelivers {
		purchaseDeliverID, err := common.Decrypt(v.PurchaseDeliverID)
		if err != nil {
			o.Failure("purchase_deliver_id.invalid", util.ErrorInvalidData("purchase deliver id"))
		}

		v.PurchaseDeliver = &model.PurchaseDeliver{ID: purchaseDeliverID}
		if err = v.PurchaseDeliver.Read("ID"); err != nil {
			o.Failure("purchase_deliver_id.invalid", util.ErrorInvalidData("purchase deliver id"))
		}

		v.PurchaseOrder = &model.PurchaseOrder{ID: v.PurchaseDeliver.PurchaseOrder.ID}
		if err = v.PurchaseDeliver.PurchaseOrder.Read("ID"); err != nil {
			o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase order id"))
		}

		if err := v.PurchaseDeliver.PurchaseOrder.Warehouse.Read("ID"); err != nil {
			o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse id"))
		}

		if k > 0 {
			if c.PurchaseDelivers[k].PurchaseDeliver.PurchaseOrder.Supplier.ID != c.PurchaseDelivers[k-1].PurchaseDeliver.PurchaseOrder.Supplier.ID {
				o.Failure("supplier.invalid", util.ErrorMustBeSameInOneDocument("supplier", "consolidated surat jalan"))
			}

			if c.PurchaseDelivers[k].PurchaseDeliver.PurchaseOrder.Warehouse.ID != c.PurchaseDelivers[k-1].PurchaseDeliver.PurchaseOrder.Warehouse.ID {
				o.Failure("warehouse.invalid", util.ErrorMustBeSameInOneDocument("warehouse", "consolidated surat jalan"))
			}
		}
	}

	// get warehouse code
	if err := c.PurchaseDelivers[0].PurchaseDeliver.PurchaseOrder.Warehouse.Read("ID"); err != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse id"))
	}

	return o
}

func (c *consolidateRequest) Messages() map[string]string {
	return map[string]string{
		"driver_name.required":         util.ErrorInputRequired("driver name"),
		"vehicle_number.required":      util.ErrorInputRequired("vehicle number"),
		"purchase_delivers.required":   util.ErrorSelectRequired("surat jalan"),
		"purchase_deliver_id.required": util.ErrorInputRequired("surat jalan id"),
		"driver_name.alpha_num_space":  util.ErrorAlphaNum("driver name"),
		"driver_name.lte":              util.ErrorEqualLess("driver name", "100"),
		"vehicle_number.alpha_num":     util.ErrorAlphaNum("vehicle number"),
		"vehicle_number.range":         util.ErrorRangeChar("vehicle number", "5", "9"),
	}
}
