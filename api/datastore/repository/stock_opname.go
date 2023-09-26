// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetStockOpname find a single data price set using field and value condition.
func GetStockOpname(field string, values ...interface{}) (*model.StockOpname, error) {
	m := new(model.StockOpname)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "StockOpnameItems", 2)

	o.Raw("select note from audit_log where `type` = 'stock_opname' and `function` = 'cancel' and ref_id = ?", m.ID).QueryRow(&m.CancellationNote)

	for _, row := range m.StockOpnameItems {
		o.Raw("select value_name from glossary where `table` = ? and `attribute` = ? and `value_int` = ?", "stock_opname", "opname_reason", row.OpnameReason).QueryRow(&row.OpnameReasonValue)
	}
	return m, nil
}

// GetStockOpnames : function to get data from database based on parameters
func GetStockOpnames(rq *orm.RequestQuery) (m []*model.StockOpname, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.StockOpname))

	if total, err = q.Filter("status__in", 1, 2, 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.StockOpname
	if _, err = q.Filter("status__in", 1, 2, 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterStockOpnames : function to get data from database based on parameters with filtered permission
func GetFilterStockOpnames(rq *orm.RequestQuery) (m []*model.StockOpname, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.StockOpname))

	if total, err = q.Filter("status__in", 1, 2, 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.StockOpname
	if _, err = q.Filter("status__in", 1, 2, 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidStockOpname : function to check if id is valid in database
func ValidStockOpname(id int64) (so *model.StockOpname, e error) {
	so = &model.StockOpname{ID: id}
	e = so.Read("ID")

	return
}

// CheckStockOpnameProductStatus : function to check if product is valid in table
func CheckStockOpnameProductStatus(productID int64, status int8, warehouseArr ...string) (*model.StockOpnameItem, int64, error) {
	var err error
	o := orm.NewOrm()
	o.Using("read_only")
	soi := new(model.StockOpnameItem)
	q := o.QueryTable(soi).RelatedSel("StockOpname").Filter("StockOpname__Status", status).Filter("product_id", productID)

	if len(warehouseArr) > 0 {
		q = q.Exclude("StockOpname__Warehouse__id__in", warehouseArr)
	}

	if total, err := q.All(soi); err == nil {
		return soi, total, nil
	}

	return nil, 0, err
}

// CheckStockOpnameData : function to check stock opname data base on parameters
func CheckStockOpnameData(filter, exclude map[string]interface{}) (so []*model.StockOpname, total int64, e error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.StockOpname))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err := o.All(&so); err == nil {
		return so, total, nil
	}

	return nil, 0, e
}

// CheckStockOpnameActive : function to check stock opname active base on warehouse id
func CheckStockOpnameActive(warehouseArr ...string) (active bool, err error) {
	o := orm.NewOrm()
	o.Using("read_only")
	stockOpname := new(model.StockOpname)

	total, err := o.QueryTable(stockOpname).Filter("Status", 1).Filter("Warehouse__id__in", warehouseArr).Count()

	if err != nil {
		return false, err
	}

	if total > 0 {
		return true, nil
	}

	return false, nil
}
