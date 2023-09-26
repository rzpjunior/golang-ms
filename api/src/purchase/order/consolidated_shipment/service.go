// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package consolidated_shipment

import (
	"time"

	"git.edenfarm.id/cuxs/orm"

	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Sign : function to sign purchase deliver
func Sign(r signRequest) (m *model.ConsolidatedShipment, err error) {
	o := orm.NewOrm()
	o.Begin()

	css := &model.ConsolidatedShipmentSignature{
		ConsolidatedShipment: r.ConsolidatedShipment,
		JobFunction:          r.JobFunction,
		Name:                 r.Name,
		SignatureURL:         r.SignatureURL,
		CreatedAt:            time.Now(),
		CreatedBy:            r.Session.Staff.ID,
	}

	if _, err := o.Insert(css); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, r.ConsolidatedShipment.ID, "consolidated shipment", "signed", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return r.ConsolidatedShipment, err
}

// Print : function to count copy of print
func Print(r printRequest) (pd *model.ConsolidatedShipment, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.ConsolidatedShipment.DeltaPrint += 1
	r.ConsolidatedShipment.Status = 2
	if _, err := o.Update(r.ConsolidatedShipment, "DeltaPrint", "Status"); err != nil {
		o.Rollback()
		return nil, err
	}

	if err := log.AuditLogByUser(r.Session.Staff, r.ConsolidatedShipment.ID, "consolidated shipment", "print", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	if err := log.AuditLogByUser(r.Session.Staff, r.ConsolidatedShipment.ID, "consolidated shipment", "confirm", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return r.ConsolidatedShipment, e
}

// Consolidate : function to consolidate purchase order
func Consolidate(r consolidateRequest) (cs *model.ConsolidatedShipment, err error) {
	o := orm.NewOrm()
	o.Begin()

	r.Code, err = util.GenerateDocCode(r.Code, r.PurchaseOrders[0].PurchaseOrder.Warehouse.Code, "consolidated_shipment")
	if err != nil {
		o.Rollback()
		return nil, err
	}

	cs = &model.ConsolidatedShipment{
		Code:              r.Code,
		DriverName:        r.DriverName,
		VehicleNumber:     r.VehicleNumber,
		DriverPhoneNumber: r.DriverPhoneNumber,
		DeltaPrint:        0,
		Status:            1,
		CreatedAt:         time.Now(),
		CreatedBy:         r.Session.Staff,
	}

	if _, err := o.Insert(cs); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, cs.ID, "consolidated shipment", "created", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	for _, v := range r.PurchaseOrders {
		v.PurchaseOrder.ConsolidatedShipment = cs

		if _, err := o.Update(v.PurchaseOrder, "ConsolidatedShipment"); err != nil {
			o.Rollback()
			return nil, err
		}

		if err = log.AuditLogByUser(r.Session.Staff, v.PurchaseOrder.ID, "purchase order", "consolidated", ""); err != nil {
			o.Rollback()
			return nil, err
		}
	}

	o.Commit()

	return cs, err
}

// Update : function to update consolidated shipment
func Update(r updateRequest) (cs *model.ConsolidatedShipment, err error) {
	o := orm.NewOrm()
	o.Begin()

	var keepItemsId []int64

	cs = &model.ConsolidatedShipment{
		ID:                r.ID,
		DriverName:        r.DriverName,
		VehicleNumber:     r.VehicleNumber,
		DriverPhoneNumber: r.DriverPhoneNumber,
	}

	if _, err := o.Update(cs, "DriverName", "VehicleNumber", "DriverPhoneNumber"); err != nil {
		o.Rollback()
		return nil, err
	}

	for _, v := range r.PurchaseOrders {
		v.PurchaseOrder.ConsolidatedShipment = cs

		if _, err := o.Update(v.PurchaseOrder, "ConsolidatedShipment"); err != nil {
			o.Rollback()
			return nil, err
		}

		keepItemsId = append(keepItemsId, v.PurchaseOrder.ID)
	}

	if _, e := o.QueryTable(new(model.PurchaseOrder)).Filter("consolidated_shipment_id", cs.ID).Filter("status", 1).Exclude("ID__in", keepItemsId).Update(orm.Params{"ConsolidatedShipment": nil}); e != nil {
		o.Rollback()
		return nil, e
	}

	if err = log.AuditLogByUser(r.Session.Staff, cs.ID, "consolidated shipment", "updated", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return cs, err
}
