// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchase_deliver

import (
	"time"

	"git.edenfarm.id/cuxs/orm"

	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Sign : function to sign purchase deliver
func Sign(r signRequest) (pd *model.PurchaseDeliver, err error) {
	o := orm.NewOrm()
	o.Begin()

	pds := &model.PurchaseDeliverSignature{
		PurchaseDeliver: r.PurchaseDeliver,
		Role:            r.Role,
		Name:            r.Name,
		Signature:       r.Signature,
		CreatedAt:       time.Now(),
		CreatedBy:       r.Session.Staff.ID,
	}

	if _, err := o.Insert(pds); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, r.PurchaseDeliver.ID, "purchase deliver", "signed", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return r.PurchaseDeliver, err
}

// Print : function to count copy of print
func Print(r printRequest) (pd *model.PurchaseDeliver, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.PurchaseDeliver.DeltaPrint += 1
	if _, err := o.Update(r.PurchaseDeliver, "DeltaPrint"); err != nil {
		o.Rollback()
		return nil, err
	}

	if err := log.AuditLogByUser(r.Session.Staff, r.PurchaseDeliver.ID, "purchase deliver", "print", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return r.PurchaseDeliver, e
}

// Consolidate : function to consolidate purchase deliver
func Consolidate(r consolidateRequest) (cpd *model.ConsolidatedPurchaseDeliver, err error) {
	o := orm.NewOrm()
	o.Begin()

	r.Code, err = util.GenerateDocCode(r.Code, r.PurchaseDelivers[0].PurchaseDeliver.PurchaseOrder.Warehouse.Code, "purchase_deliver")
	if err != nil {
		o.Rollback()
		return nil, err
	}

	cpd = &model.ConsolidatedPurchaseDeliver{
		Code:          r.Code,
		DriverName:    r.DriverName,
		VehicleNumber: r.VehicleNumber,
		DeltaPrint:    1,
		CreatedAt:     time.Now(),
		CreatedBy:     r.Session.Staff,
	}

	if _, err := o.Insert(cpd); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, cpd.ID, "consolidated purchase deliver", "created", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	for _, v := range r.PurchaseDelivers {
		v.PurchaseDeliver.ConsolidatedPurchaseDeliver = cpd

		if _, err := o.Update(v.PurchaseDeliver, "ConsolidatedPurchaseDeliver"); err != nil {
			o.Rollback()
			return nil, err
		}

		if err = log.AuditLogByUser(r.Session.Staff, v.PurchaseDeliver.ID, "purchase deliver", "consolidated", ""); err != nil {
			o.Rollback()
			return nil, err
		}
	}

	o.Commit()

	return cpd, err
}

// PrintConsolidate : function to count copy of print
func PrintConsolidate(r printConsolidateRequest) (pd *model.ConsolidatedPurchaseDeliver, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.ConsolidatedPurchaseDeliver.DeltaPrint += 1
	if _, err := o.Update(r.ConsolidatedPurchaseDeliver, "DeltaPrint"); err != nil {
		o.Rollback()
		return nil, err
	}

	if err := log.AuditLogByUser(r.Session.Staff, r.ConsolidatedPurchaseDeliver.ID, "consolidated purchase deliver", "print", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return r.ConsolidatedPurchaseDeliver, e
}

// SignConsolidate : function to sign consolidated purchase deliver
func SignConsolidate(r signConsolidateRequest) (cpd *model.ConsolidatedPurchaseDeliver, err error) {
	o := orm.NewOrm()
	o.Begin()

	pds := &model.ConsolidatedPurchaseDeliverSignature{
		ConsolidatedPurchaseDeliver: r.ConsolidatedPurchaseDeliver,
		Role:                        r.Role,
		Name:                        r.Name,
		Signature:                   r.Signature,
		CreatedAt:                   time.Now(),
		CreatedBy:                   r.Session.Staff.ID,
	}

	if _, err := o.Insert(pds); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, r.ConsolidatedPurchaseDeliver.ID, "consolidated purchase deliver", "signed", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return r.ConsolidatedPurchaseDeliver, err
}
