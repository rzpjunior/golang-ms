// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type cancelRequest struct {
	ID   int64  `json:"-"`
	Note string `json:"note" valid:"required"`

	PurchaseOrder      *model.PurchaseOrder       `json:"-"`
	PurchasePlan       *model.PurchasePlan        `json:"-"`
	PurchaseOrderItems []*model.PurchaseOrderItem `json:"-"`
	Session            *auth.SessionData          `json:"-"`
}

func (r *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var countPI, countGr int
	var filter, exclude map[string]interface{}
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	orSelect.Raw("select count(pi2.id) from purchase_invoice pi2 where pi2.purchase_order_id = ? and pi2.status in(1,2)", r.PurchaseOrder.ID).QueryRow(&countPI)
	orSelect.Raw("select count(gr.id) from goods_receipt gr where gr.purchase_order_id = ? and gr.status in(1,2)", r.PurchaseOrder.ID).QueryRow(&countGr)

	if countGr > 0 || countPI > 0 {
		o.Failure("id.invalid", util.ErrorPORelatedDoc())
	}

	if r.PurchaseOrder.Status != 5 && r.PurchaseOrder.Status != 1 {
		o.Failure("status.inactive", util.ErrorDocStatus("purchase order", "active or draft"))
		return o
	}

	if r.PurchaseOrder.PurchasePlan != nil {
		r.PurchasePlan, err = repository.ValidPurchasePlan(r.PurchaseOrder.PurchasePlan.ID)
		if err != nil {
			o.Failure("purchase_plan_id.invalid", util.ErrorInvalidData("purchase plan"))
		}

		filter = map[string]interface{}{"purchase_order_id": r.PurchaseOrder.ID}
		r.PurchaseOrderItems, _, err = repository.CheckPurchaseOrderItemData(filter, exclude)
		if err != nil {
			o.Failure("purchase_order_item.invalid", util.ErrorInvalidData("purchase order item"))
		}
	}

	if r.PurchaseOrder.ConsolidatedShipment != nil {
		o.Failure("purchase_order.invalid", util.ErrorCannotCancelAfter("purchase order", "consolidated"))
	}

	if r.PurchaseOrder.DeltaPrint > 0 {
		o.Failure("purchase_order.invalid", util.ErrorCannotCancelAfter("purchase order", "printed"))
	}

	return o
}

func (r *cancelRequest) Messages() map[string]string {
	return map[string]string{
		"note.required": util.ErrorInputRequired("cancellation note"),
	}
}
