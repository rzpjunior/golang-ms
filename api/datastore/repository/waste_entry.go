// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetWasteEntry find a single data price set using field and value condition.
func GetWasteEntry(field string, values ...interface{}) (*model.WasteEntry, error) {
	m := new(model.WasteEntry)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel(2).Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "WasteEntryItems", 2)

	o.Raw("select note from audit_log where `type` = 'waste_entry' and `function` = 'cancel' and ref_id = ?", m.ID).QueryRow(&m.CancellationNote)

	for _, row := range m.WasteEntryItems {
		o.Raw("select available_stock, waste_stock from stock where `warehouse_id` = ? and `product_id` = ?", m.Warehouse.ID, row.Product.ID).QueryRow(&row.AvailableStock, &row.WasteStock)
		o.Raw("select value_name from glossary where `table` = ? and `attribute` = ? and `value_int` = ?", "all", "waste_reason", row.WasteReason).QueryRow(&row.WasteReasonValue)
	}

	return m, nil
}

// GetWasteEntrys : function to get data from database based on parameters
func GetWasteEntrys(rq *orm.RequestQuery) (m []*model.WasteEntry, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.WasteEntry))

	if total, err = q.Exclude("status", 4).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.WasteEntry
	if _, err = q.Exclude("status", 4).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterWasteEntrys : function to get data from database based on parameters with filtered permission
func GetFilterWasteEntrys(rq *orm.RequestQuery) (m []*model.WasteEntry, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.WasteEntry))

	if total, err = q.Exclude("status", 4).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.WasteEntry
	if _, err = q.Exclude("status", 4).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidWasteEntry : function to check if id is valid in database
func ValidWasteEntry(id int64) (wasteEntry *model.WasteEntry, e error) {
	wasteEntry = &model.WasteEntry{ID: id}
	e = wasteEntry.Read("ID")

	return
}

// CheckWasteEntryProductStatus : function to check if product is valid in table
func CheckWasteEntryProductStatus(productID int64, status int8, warehouseArr ...string) (*model.WasteEntryItem, int64, error) {
	var err error
	o := orm.NewOrm()
	o.Using("read_only")

	wei := new(model.WasteEntryItem)
	q := o.QueryTable(wei).RelatedSel("WasteEntry").Filter("WasteEntry__status", status).Filter("product_id", productID)

	if len(warehouseArr) > 0 {
		q = q.Exclude("WasteEntry__Warehouse__id__in", warehouseArr)
	}

	if total, err := q.All(wei); err == nil {
		return wei, total, nil
	}

	return nil, 0, err
}
