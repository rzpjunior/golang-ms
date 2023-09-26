// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package koli_increment

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// printRequest : struct to hold label print request data
type printRequest struct {
	IncrementPrints []*int8 `json:"increment_prints" valid:"required"`

	SalesOrderCode        string                         `json:"sales_order_code" valid:"required"`
	DeliveryKoliIncrement []*model.DeliveryKoliIncrement `json:"-"`

	SalesOrder         *model.SalesOrder         `json:"-"`
	PickingOrderAssign *model.PickingOrderAssign `json:"-"`
	Session            *auth.SessionData         `json:"-"`
}

type IncrementPrints struct {
	Increment int8 `json:"increment" valid:"value" valid:"required"`
}

// Validate : function to validate uom request data
func (c *printRequest) Validate() *validation.Output {
	var err error
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")

	if err = o1.Raw("SELECT id, branch_id, term_payment_sls_id, term_invoice_sls_id, salesperson_id, sales_group_id, sub_district_id, warehouse_id, wrt_id, area_id, voucher_id, price_set_id, payment_group_sls_id, archetype_id, order_type_sls_id, order_channel, code, status, recognition_date, delivery_date, billing_address, shipping_address, shipping_address_note, delivery_fee, vou_redeem_code, vou_disc_amount, point_redeem_amount, point_redeem_id, total_price, total_charge, total_weight, note, reload_packing, payment_reminder, is_locked, has_ext_invoice, has_picking_assigned, cancel_type, created_at, created_by, last_updated_at, last_updated_by, finished_at, locked_by "+
		"FROM eden_v2.sales_order where code = ?", c.SalesOrderCode).QueryRow(&c.SalesOrder); err != nil {
		o.Failure("sales_order_code", util.ErrorInvalidData("sales order"))
	}

	if err = o1.Raw("SELECT id, picking_order_id, sales_order_id, staff_id, courier_id, courier_vendor_id, dispatcher_id, picking_list_id, status, dispatch_status, dispatch_timestamp, assign_timestamp, planning_vendor, been_rejected, note, checkin_timestamp, checkout_timestamp, checker_in_timestamp, checker_out_timestamp, total_koli, total_scan_dispatch, checked_at, checked_by "+
		"FROM eden_v2.picking_order_assign where sales_order_id = ?", c.SalesOrder.ID).QueryRow(&c.PickingOrderAssign); err != nil {
		o.Failure("picking_order_assign", util.ErrorInvalidData("picking order assign"))
	}

	for _, v := range c.IncrementPrints {
		var dki *model.DeliveryKoliIncrement
		o1.Raw("SELECT id, sales_order_id, `increment`, is_read, print_label "+
			"FROM eden_v2.delivery_koli_increment where sales_order_id = ? and `increment` = ?", c.SalesOrder.ID, v).QueryRow(&dki)

		dki.TotalKoli = c.PickingOrderAssign.TotalKoli
		dki.Helper, _ = repository.ValidStaff(c.PickingOrderAssign.Helper.ID)
		dki.HelperCode = dki.Helper.Code
		dki.SalesOrder.Read("ID")
		dki.SalesOrder.Branch.Read("ID")
		dki.SalesOrder.Branch.Merchant.Read("ID")

		c.DeliveryKoliIncrement = append(c.DeliveryKoliIncrement, dki)
	}

	return o
}

// Messages : function to return error validation messages
func (c *printRequest) Messages() map[string]string {
	messages := map[string]string{
		"sales_order_code.required": util.ErrorInputRequired("sales order"),
		"increment_prints.required": util.ErrorInputRequired("increment"),
	}

	return messages
}
