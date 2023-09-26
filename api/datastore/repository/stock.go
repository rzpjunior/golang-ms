// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"strconv"
	"strings"

	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetStock find a single data sales term using field and value condition.
func GetStock(field string, values ...interface{}) (*model.Stock, error) {
	m := new(model.Stock)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetStocks : function to get data from database based on parameters
func GetStocks(rq *orm.RequestQuery) (m []*model.Stock, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Stock))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Stock
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterStocks : function to get data from database based on parameters with filtered permission
func GetFilterStocks(rq *orm.RequestQuery) (m []*model.Stock, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Stock))

	if total, err = q.RelatedSel(2).Exclude("status", 3).Filter("warehouse__status", 1).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Stock
	if _, err = q.RelatedSel(2).Exclude("status", 3).Filter("warehouse__status", 1).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidStock : function to check if id is valid in database
func ValidStock(id int64) (stock *model.Stock, e error) {
	stock = &model.Stock{ID: id}
	e = stock.Read("ID")

	return
}

// CheckStockPerWarehouse : function to check stock of each warehouse
func CheckStockPerWarehouse(productID int64, warehouseID ...string) (*model.Stock, int64, error) {
	var err error
	o := orm.NewOrm()
	o.Using("read_only")

	s := new(model.Stock)

	cond1 := orm.NewCondition()
	cond1 = cond1.And("product_id", productID)

	if len(warehouseID) > 0 {
		cond1 = cond1.AndNot("Warehouse__in", warehouseID)
	}

	cond2 := orm.NewCondition()
	cond2 = cond2.And("available_stock__gt", 0).Or("waste_stock__gt", 0)

	cond3 := cond1.AndCond(cond2)

	q := o.QueryTable(s).SetCond(cond3)

	if total, err := q.All(s); err == nil {
		return s, total, nil
	}

	return nil, 0, err
}

// GetStocksByProduct : function to check stock by its's product
func GetStocksByProduct(productID ...int64) (stocks []*model.Stock, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	if _, err = o.QueryTable(new(model.Stock)).Filter("product_id__in", productID).All(&stocks); err == nil {
		return stocks, nil
	}

	return nil, err
}

// CheckStockData : function to check stock data based on filter and exclude parameters
func CheckStockData(filter, exclude map[string]interface{}) (stock []*model.Stock, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.Stock))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err = o.All(&stock); err == nil {
		return stock, total, nil
	}

	return nil, 0, err
}

func GetExportForm(rq *orm.RequestQuery, warehouse int64, category int64, classification string) (m []*model.Stock, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Stock))
	o := orm.NewOrm()
	o.Using("read_only")

	var queryString string
	var categoryId []int64
	switch classification {
	case "1":
		queryString = "c.grandparent_id = ?"
	case "2":
		queryString = "c.parent_id = ?"
	case "3":
		queryString = "c.id = ?"
	default:
		queryString = "c.id = ?"
	}
	o.Raw("select id from category c where "+queryString, category).QueryRows(&categoryId)
	if len(categoryId) == 0 {
		categoryId = append(categoryId, 0)
	}

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}
	// get data requested
	var mx []*model.Stock
	if _, err = q.Filter("product_id__status", 1).Filter("status", 1).Filter("warehouse_id", warehouse).Filter("product_id__category_id__id__in", categoryId).OrderBy("-product_id").All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func GetExportFormWaste(rq *orm.RequestQuery, warehouse int64) (m []*model.Stock, total int64, err error) {

	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Stock))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}
	// get data requested
	var mx []*model.Stock
	if _, err = q.Filter("product_id__status", 1).Filter("status", 1).Filter("warehouse_id", warehouse).OrderBy("-product_id").All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func CheckStockProductTransferSku(productID int64, warehouseID ...string) (*model.Stock, int64, error) {
	var err error
	o := orm.NewOrm()
	o.Using("read_only")

	s := new(model.Stock)

	cond1 := orm.NewCondition()
	cond1 = cond1.And("product_id", productID)

	if len(warehouseID) > 0 {
		cond1 = cond1.And("Warehouse__in", warehouseID)
	}

	q := o.QueryTable(s).SetCond(cond1)

	if total, err := q.All(s); err == nil {
		return s, total, nil
	}

	return nil, 0, err
}

func GetFilterStocksWithProductGroup(rq *orm.RequestQuery, pID int64, orderChannelRestriction string) (m []*model.Stock, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Stock))

	cond := q.GetCond()
	// add order channel restriction clause
	if orderChannelRestriction != "" {
		orderChannelResArr := strings.Split(orderChannelRestriction, ",")

		// validation filter order_channel_restriction
		cond1 := orm.NewCondition()
		cond1 = cond1.And("product__order_channel_restriction__isnull", true)

		cond2 := orm.NewCondition()
		for _, v := range orderChannelResArr {
			valInt, err := strconv.Atoi(v)
			if err != nil {
				return nil, 0, err
			}
			orderChanel := util.IsOrderChannel(valInt)
			if !orderChanel {
				return nil, 0, err
			}

			cond3 := orm.NewCondition()
			cond3 = cond3.Or("product__order_channel_restriction__icontains", ","+v+",").Or("product__order_channel_restriction__istartswith", v+",").Or("product__order_channel_restriction__iendswith", ","+v).Or("product__order_channel_restriction", v)
			cond2 = cond2.OrCond(cond3)
		}

		cond1 = cond1.OrNotCond(cond2)
		cond = cond.AndCond(cond1)
	}

	q = q.SetCond(cond)

	if total, err = q.RelatedSel(2).Exclude("status", 3).Filter("warehouse__status", 1).Exclude("product_id", pID).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Stock
	if _, err = q.RelatedSel(2).Exclude("status", 3).Filter("warehouse__status", 1).Exclude("product_id", pID).All(&mx, rq.Fields...); err != nil {
		return nil, total, err
	}

	for _, v := range mx {

		productGroupItem, _ := GetProductGroupItem("product_id", v.Product.ID)
		if productGroupItem != nil {
			v.ProductGroup = productGroupItem.ProductGroup
		}
	}

	return mx, total, nil
}
