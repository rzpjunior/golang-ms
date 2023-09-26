// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"math"
	"sort"
	"time"

	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetSalesRecap find a single data of sales recap using field and value condition.
func GetSalesRecap(field string, values ...interface{}) (*model.SalesRecap, error) {
	m := new(model.SalesRecap)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetSalesRecaps : function to get data from database based on parameters
func GetSalesRecaps(deliveryDate time.Time, warehouseID, categoryID int64) (salesRecap []*model.SalesRecap, total int64, err error) {
	var (
		o     orm.Ormer
		q     orm.QuerySeter
		query string
		s     []*model.Stock
	)

	o = orm.NewOrm()
	o.Using("read_only")

	q = o.QueryTable(new(model.Stock))
	q = q.RelatedSel("Product").RelatedSel("Warehouse").RelatedSel("Product__Category").RelatedSel("Product__Uom").Filter("Warehouse", warehouseID).Exclude("CommitedOutStock", 0)

	if categoryID != 0 {
		q = q.Filter("Product__Category__id", categoryID)
	}

	if total, err = q.All(&s); err != nil {
		return nil, 0, err
	}

	for _, v := range s {
		query = "select sum(soi.order_qty) so_qty " +
			"from sales_order so " +
			"join sales_order_item soi on so.id = soi.sales_order_id " +
			"join order_type_sls ots on so.order_type_sls_id = ots.id " +
			"where (so.status in (9,12) or (so.status = 1 and so.payment_group_sls_id != 1)) and ots.value != \"draft\" and so.delivery_date = ? and so.warehouse_id = ? and soi.product_id = ? "

		week1Demand := float64(0)
		o.Raw(query, deliveryDate.AddDate(0, 0, -7).Format("2006-01-02"), warehouseID, v.Product.ID).QueryRow(&week1Demand)

		week2Demand := float64(0)
		o.Raw(query, deliveryDate.AddDate(0, 0, -14).Format("2006-01-02"), warehouseID, v.Product.ID).QueryRow(&week2Demand)

		sr := new(model.SalesRecap)
		sr.Product = v.Product
		sr.SumSoQty = v.CommitedOutStock
		sr.SumPoQty = v.CommitedInStock
		sr.AvailableStock = v.AvailableStock
		sr.ExpectedRemainingStock = math.Round((v.CommitedInStock+v.AvailableStock-v.CommitedOutStock)*100) / 100
		sr.SpareQty = math.Round((v.Product.SparePercentage*v.CommitedOutStock/100)*100) / 100
		sr.Week1Demand = week1Demand
		sr.Week2Demand = week2Demand
		sr.WeekAvgDemand = (week1Demand + week2Demand) / 2

		salesRecap = append(salesRecap, sr)
	}

	// sort asc sales recap data by it's remaining stock
	sort.SliceStable(salesRecap[:], func(i, j int) bool {
		return salesRecap[i].ExpectedRemainingStock < salesRecap[j].ExpectedRemainingStock
	})

	return salesRecap, total, nil
}
