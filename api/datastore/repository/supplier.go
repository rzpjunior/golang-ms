// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetSupplier find a single data supplier using field and value condition.
func GetSupplier(field string, values ...interface{}) (*model.Supplier, error) {
	m := new(model.Supplier)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel(4).Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetSuppliers get all data supplier that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetSuppliers(rq *orm.RequestQuery) (m []*model.Supplier, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Supplier))

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Supplier
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// GetFilterSuppliers : function to get data from database based on parameters with filtered permission
func GetFilterSuppliers(rq *orm.RequestQuery) (m []*model.Supplier, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Supplier))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Supplier
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// CountNonDeletedSupplierByPurchaseTermId : function to check whether purchase term id is still used by any active or archive supplier
func CountNonDeletedSupplierByPurchaseTermId(id int64) (countDeletedSupplier int64, e error) {
	m := new(model.Supplier)
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(m)

	countDeletedSupplier, err := o.Filter("term_payment_pur_id", id).Exclude("status", 3).Count()
	if err != nil {
		return 0, err
	}

	return countDeletedSupplier, nil
}

func ValidSupplier(id int64) (supplier *model.Supplier, e error) {
	supplier = &model.Supplier{ID: id}
	e = supplier.Read("ID")

	return
}

func GetSupplierDetail(field string, values ...interface{}) (*model.Supplier, error) {
	m := new(model.Supplier)
	o := orm.NewOrm()
	var err error

	o.Using("read_only")

	if err = o.QueryTable(m).Filter(field, values...).Limit(1).One(m); err != nil {
		return nil, err
	}

	if m.SupplierCommodity != nil && m.SupplierCommodity.ID != 0 {

		if err = m.SupplierCommodity.Read("ID"); err != nil {
			return nil, err
		}

	}

	if m.SupplierType != nil && m.SupplierType.ID != 0 {

		if err = m.SupplierType.Read("ID"); err != nil {
			return nil, err
		}

	}

	if m.SupplierBadge != nil && m.SupplierBadge.ID != 0 {

		if err = m.SupplierBadge.Read("ID"); err != nil {
			return nil, err
		}

	}

	if m.PaymentMethod != nil && m.PaymentMethod.ID != 0 {

		if err = m.PaymentMethod.Read("ID"); err != nil {
			return nil, err
		}

	}

	if m.PaymentTerm != nil && m.PaymentTerm.ID != 0 {

		if err = m.PaymentTerm.Read("ID"); err != nil {
			return nil, err
		}

	}

	if m.SubDistrict != nil && m.SubDistrict.ID != 0 {

		if err = m.SubDistrict.Read("ID"); err != nil {
			return nil, err
		}

		if err = m.SubDistrict.District.Read("ID"); err != nil {
			return nil, err
		}

		if err = m.SubDistrict.District.City.Read("ID"); err != nil {
			return nil, err
		}

		if err = m.SubDistrict.District.City.Province.Read("ID"); err != nil {
			return nil, err
		}

	}

	if m.SupplierOrganization != nil && m.SupplierOrganization.ID != 0 {

		if err = m.SupplierOrganization.Read("ID"); err != nil {
			return nil, err
		}

	}

	return m, nil
}

// GetSuppliersInFieldPurchaserApps get all data supplier for field purchaser apps that matched with query request parameters.
func GetSuppliersInFieldPurchaserApps(rq *orm.RequestQuery) (m []*model.Supplier, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")
	q, _ := rq.QueryReadOnly(new(model.Supplier))

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Supplier
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err != nil {
		return nil, total, err
	}

	// return error some thing went wrong
	return mx, total, nil
}
