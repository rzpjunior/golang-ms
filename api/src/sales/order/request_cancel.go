// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"strconv"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// cancelRequest : struct to hold sales order request data
type cancelRequest struct {
	ID                     int64   `json:"-"`
	Note                   string  `json:"note" valid:"required"`
	RecentPoint            float64 `json:"recent_point" valid:"required"`
	CancelType             int8    `json:"cancel_type" valid:"required"`
	IsUseSkuDiscount       int8    `json:"-"`
	CreditLimitBefore      float64 `json:"-"`
	CreditLimitAfter       float64 `json:"-"`
	IsCreateCreditLimitLog int64   `json:"-"`

	PackingOrder            *model.PackingOrder            `json:"-"`
	PickingOrderAssign      *model.PickingOrderAssign      `json:"-"`
	SalesOrder              *model.SalesOrder              `json:"-"`
	SkuDiscountItems        []*model.SkuDiscountItem       `json:"-"`
	SkuDiscountLogs         []*model.SkuDiscountLog        `json:"-"`
	VoucherLog              []*model.VoucherLog            `json:"-"`
	SalesOrderItem          []*model.SalesOrderItem        `json:"-"`
	MerchantPointLog        []*model.MerchantPointLog      `json:"-"`
	MerchantPointExpiration *model.MerchantPointExpiration `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate uom request data
func (r *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var err error
	var filter, exclude map[string]interface{}

	r.SalesOrder = &model.SalesOrder{ID: r.ID}
	if err = r.SalesOrder.Read("ID"); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales order"))
		return o
	}

	if r.SalesOrder.Status != 1 {
		o.Failure("status.inactive", util.ErrorActive("sales order"))
		return o
	}

	if err = r.Session.Staff.Role.Read("ID"); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("role"))
		return o
	}

	if r.Session.Staff.Role.Code == "ROL0008" {
		if err = r.SalesOrder.OrderType.Read("ID"); err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("order type"))
			return o
		}

		if r.SalesOrder.OrderType.Code != "SOT0010" {
			o.Failure("id.invalid", util.ErrorOrderTypeDraft())
			return o
		}

		if r.SalesOrder.IsLocked == 1 {
			o.Failure("id.invalid", util.ErrorSOLocked())
			return o
		}
	}

	if err = r.SalesOrder.Branch.Read("ID"); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("branch"))
		return o
	}

	if err = r.SalesOrder.Branch.Merchant.Read("ID"); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("merchant"))
		return o
	}

	if r.SalesOrder.HasExtInvoice == 2 {
		if err = r.SalesOrder.Branch.Merchant.UserMerchant.Read("ID"); err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("user merchant"))
			return o
		}
	}

	if r.SalesOrder.Voucher != nil {
		r.SalesOrder.Voucher.Read("ID")

		filter := map[string]interface{}{
			"voucher_id":     r.SalesOrder.Voucher.ID,
			"merchant_id":    r.SalesOrder.Branch.Merchant.ID,
			"branch_id":      r.SalesOrder.Branch.ID,
			"sales_order_id": r.SalesOrder.ID,
			"status":         int8(1),
		}
		r.VoucherLog, _, err = repository.CheckVoucherLogData(filter, exclude)
	}

	if r.SalesOrder.PointRedeemID != 0 && r.SalesOrder.PointRedeemAmount != 0 {
		filter := map[string]interface{}{
			"id":             r.SalesOrder.PointRedeemID,
			"merchant_id":    r.SalesOrder.Branch.Merchant.ID,
			"sales_order_id": r.SalesOrder.ID,
			"status":         int8(2),
			"point_value":    r.SalesOrder.PointRedeemAmount,
		}

		r.MerchantPointLog, _, err = repository.CheckMerchantPointLogData(filter, exclude)

		o1.Raw("SELECT recent_point from merchant_point_log where merchant_id = ? order by id desc limit 1 ", r.SalesOrder.Branch.Merchant.ID).QueryRow(&r.RecentPoint)

		r.MerchantPointExpiration = &model.MerchantPointExpiration{ID: r.SalesOrder.Branch.Merchant.ID}
		if err = r.MerchantPointExpiration.Read("ID"); err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("Merchant Point Expiration"))
			return o
		}
	}

	r.PackingOrder = &model.PackingOrder{
		Status:       1,
		DeliveryDate: r.SalesOrder.DeliveryDate,
		Warehouse:    r.SalesOrder.Warehouse,
	}

	if err = r.PackingOrder.Read("Status", "DeliveryDate", "Warehouse"); err == nil {
		filter = map[string]interface{}{"SalesOrder__id": r.ID}
		r.SalesOrderItem, _, _ = repository.CheckSalesOrderItemData(filter, exclude)
	} else {
		r.PackingOrder = nil
	}

	o1.Raw("SELECT id, picking_order_id, sales_order_id, staff_id, courier_id, courier_vendor_id, dispatcher_id, picking_list_id, status, dispatch_status, dispatch_timestamp, assign_timestamp, planning_vendor, been_rejected, note, checkin_timestamp, checkout_timestamp, checker_in_timestamp, checker_out_timestamp, total_koli, total_scan_dispatch, checked_at, checked_by "+
		"FROM eden_v2.picking_order_assign where sales_order_id = ?", r.SalesOrder.ID).QueryRow(&r.PickingOrderAssign)

	// checking if sales order items uses discount and sales order type not zero waste
	if r.SalesOrder.OrderType.Name != "Zero Waste" {
		o1.LoadRelated(r.SalesOrder, "SalesOrderItems", 1)
		for i, v := range r.SalesOrder.SalesOrderItems {
			if v.SkuDiscountItem != nil {
				skuDiscountItem := &model.SkuDiscountItem{ID: v.SkuDiscountItem.ID}
				if err = skuDiscountItem.Read("ID"); err != nil {
					o.Failure("id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("sku discount"))
					continue
				}
				skuDiscountItem.RemOverallQuota += v.DiscountQty
				if skuDiscountItem.IsUseBudget == 1 {
					skuDiscountItem.RemBudget += (v.DiscountQty * v.UnitPriceDiscount)
				}
				r.SkuDiscountItems = append(r.SkuDiscountItems, skuDiscountItem)

				skuDiscountLog := &model.SkuDiscountLog{Branch: r.SalesOrder.Branch, SalesOrderItem: v, SkuDiscountItem: skuDiscountItem, Status: 1}
				if err = skuDiscountLog.Read("Branch", "SalesOrderItem", "SkuDiscountItem", "Status"); err != nil {
					o.Failure("id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("sku discount"))
					continue
				}
				skuDiscountLog.Status = 2
				r.SkuDiscountLogs = append(r.SkuDiscountLogs, skuDiscountLog)
			}
		}
	}

	r.CreditLimitBefore = r.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount
	if r.SalesOrder.Branch.Merchant.CreditLimitAmount > 0 || r.CreditLimitBefore < 0 {
		r.IsCreateCreditLimitLog = 1
		r.CreditLimitAfter = r.CreditLimitBefore + r.SalesOrder.TotalCharge
	}
	return o
}

// Messages : function to return error validation messages
func (r *cancelRequest) Messages() map[string]string {
	messages := map[string]string{
		"note.required":        util.ErrorInputRequired("note"),
		"cancel_type.required": util.ErrorInputRequired("cancel type"),
	}

	return messages
}
