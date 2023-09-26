// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package field_purchaser

import (
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/orm"

	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// SaveOrder : function to create field purchase order and generate surat jalan
func SaveOrder(r createOrderRequest) (pd *model.PurchaseDeliver, err error) {
	o := orm.NewOrm()
	o.Begin()

	r.Code, err = util.GenerateDocCode(r.Code, r.PurchaseOrder.Warehouse.Code, "field_purchase_order")
	if err != nil {
		o.Rollback()
		return nil, err
	}

	fpo := &model.FieldPurchaseOrder{
		Code:          r.Code,
		PurchaseOrder: r.PurchaseOrder.ID,
		Stall:         r.Stall,
		TotalPrice:    r.TotalPrice,
		TotalItem:     r.TotalItem,
		PaymentMethod: r.PaymentMethod,
		Latitude:      r.Latitude,
		Longitude:     r.Longitude,
		CreatedAt:     time.Now(),
		CreatedBy:     r.Session.Staff.ID,
	}

	if _, err = o.Insert(fpo); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, fpo.ID, "field purchase order", "create", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	for _, v := range r.Items {
		fpoi := &model.FieldPurchaseOrderItem{
			FieldPurchaseOrder: fpo,
			PurchaseOrderItem:  v.PurchaseOrderItem,
			PurchaseQty:        v.PurchaseQty,
			Product:            v.Product,
			UnitPrice:          v.UnitPrice,
		}

		if _, err = o.Insert(fpoi); err != nil {
			o.Rollback()
			return nil, err
		}

		if _, err = o.Update(v.PurchaseOrderItem, "PurchaseQty", "UnitPrice", "Subtotal", "UnitPriceTax", "TaxAmount"); err != nil {
			o.Rollback()
			return nil, err
		}

		if err = log.AuditLogByUser(r.Session.Staff, v.PurchaseOrderItem.ID, "purchase_order_item", "update", strconv.FormatFloat(v.PurchaseOrderItem.UnitPrice, 'f', 2, 64)); err != nil {
			o.Rollback()
			return nil, err
		}
	}

	if _, err = o.Update(r.PurchaseOrder, "TotalPrice", "TotalCharge", "TaxAmount"); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, r.PurchaseOrder.ID, "purchase order", "update", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	r.CodePurchaseDeliver, err = util.GenerateDocCode(r.CodePurchaseDeliver, r.PurchaseOrder.Warehouse.Code, "purchase_deliver")
	if err != nil {
		o.Rollback()
		return nil, err
	}

	pd = &model.PurchaseDeliver{
		Code:               r.CodePurchaseDeliver,
		PurchaseOrder:      r.PurchaseOrder,
		FieldPurchaseOrder: fpo,
		Stall:              r.Stall,
		DeltaPrint:         1,
		CreatedAt:          time.Now(),
		CreatedBy:          r.Session.Staff.ID,
	}

	if _, err = o.Insert(pd); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, pd.ID, "purchase deliver", "create", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	if r.Signature != "" {
		pds := &model.PurchaseDeliverSignature{
			PurchaseDeliver: pd,
			Role:            r.Stall.Name,
			Name:            r.Name,
			Signature:       r.Signature,
			CreatedAt:       time.Now(),
			CreatedBy:       r.Session.Staff.ID,
		}

		if _, err = o.Insert(pds); err != nil {
			o.Rollback()
			return nil, err
		}
	}

	o.Commit()

	return pd, err
}

// Update : function to update data in database
func Update(r updateRequest) (fpo *model.FieldPurchaseOrder, err error) {
	o := orm.NewOrm()
	o.Begin()

	fpo = &model.FieldPurchaseOrder{
		ID:            r.ID,
		TotalPrice:    r.TotalPrice,
		PaymentMethod: r.PaymentMethod,
	}

	if _, err = o.Update(fpo, "TotalPrice", "PaymentMethod"); err != nil {
		o.Rollback()
		return nil, err
	}

	for _, v := range r.FieldPurchaseOrderItems {

		fpoi := &model.FieldPurchaseOrderItem{
			ID:          v.FieldPurchaseOrderItem.ID,
			PurchaseQty: v.PurchaseQty,
			UnitPrice:   v.UnitPrice,
		}

		if _, err = o.Update(fpoi, "PurchaseQty", "UnitPrice"); err != nil {
			o.Rollback()
			return nil, err
		}

		if _, err = o.Update(v.PurchaseOrderItem, "PurchaseQty", "UnitPrice", "Subtotal", "UnitPriceTax", "TaxAmount"); err != nil {
			o.Rollback()
			return nil, err
		}

		if err = log.AuditLogByUser(r.Session.Staff, v.PurchaseOrderItem.ID, "purchase_order_item", "update", strconv.FormatFloat(v.PurchaseOrderItem.UnitPrice, 'f', 2, 64)); err != nil {
			o.Rollback()
			return nil, err
		}
	}

	if _, err = o.Update(r.PurchaseOrder, "TotalPrice", "TotalCharge", "TaxAmount"); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, fpo.ID, "field_purchase_order", "update", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()
	return fpo, err
}
