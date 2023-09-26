// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetWasteDisposal find a single data price set using field and value condition.
func GetWasteDisposal(field string, values ...interface{}) (*model.WasteDisposal, error) {
	m := new(model.WasteDisposal)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	o.LoadRelated(m, "WasteDisposalItems", 2)

	orm.NewOrm().Raw("SELECT city_name, province_name from adm_division where sub_district_id = ?", m.Warehouse.SubDistrict.ID).QueryRow(&m.City, &m.Province)

	for _, v := range m.WasteDisposalItems {
		stock := new(model.Stock)
		if err := o.QueryTable(new(model.Stock)).Filter("warehouse_id", v.WasteDisposal.Warehouse.ID).Filter("product_id", v.Product.ID).RelatedSel(1).One(stock); err != nil {
			stock = nil
		}

		v.Stock = stock
	}
	return m, nil
}

// GetWasteDisposals : function to get data from database based on parameters
func GetWasteDisposals(rq *orm.RequestQuery) (m []*model.WasteDisposal, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.WasteDisposal))

	if total, err = q.Filter("status__in", 1, 2, 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.WasteDisposal
	if _, err = q.Filter("status__in", 1, 2, 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterWasteDisposals : function to get data from database based on parameters with filtered permission
func GetFilterWasteDisposals(rq *orm.RequestQuery) (m []*model.WasteDisposal, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.WasteDisposal))

	if total, err = q.Filter("status__in", 1, 2, 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.WasteDisposal
	if _, err = q.Filter("status__in", 1, 2, 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidWasteDisposal : function to check if id is valid in database
func ValidWasteDisposal(id int64) (wasteDisposal *model.WasteDisposal, e error) {
	wasteDisposal = &model.WasteDisposal{ID: id}
	e = wasteDisposal.Read("ID")

	return
}
