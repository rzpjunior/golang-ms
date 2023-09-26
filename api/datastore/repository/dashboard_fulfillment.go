// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/dbredis"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetFulfillmentSummary : function to get fullfillment summary data
func GetFulfillmentSummary(date []string, warehouseID int64) (data *model.Fulfillment, err error) {
	var (
		key       string
		cacheTime time.Duration
		updatedAt time.Time
	)

	data = new(model.Fulfillment)
	startDate, _ := time.Parse("2006-01-02", date[0])
	endDate, _ := time.Parse("2006-01-02", date[1])

	// check if range is a week
	// if it is a week then get data from cache or set data into cache
	if int(startDate.Weekday()) == 1 && endDate.Sub(startDate).Hours()/24 == 5 {
		year, week := startDate.ISOWeek()
		yearWeek := strconv.Itoa(year) + strconv.Itoa(week)
		key = fmt.Sprintf("fulfillment_summary week %s %d", yearWeek, warehouseID)
		if dbredis.Redis.CheckExistKey(key) {
			dbredis.Redis.GetCache(key, &data)
			dbredis.Redis.GetCache(key+" updated at", &updatedAt)
			data.LastUpdatedAt = updatedAt
			return
		}

		if data, err = GetRangedData(startDate, endDate, warehouseID); err == nil && data.TotalSo != 0 {
			updatedAt = time.Now()
			data.LastUpdatedAt = updatedAt

			cacheTime = SetCacheTime(endDate)
			dbredis.Redis.SetCache(key, data, cacheTime)
			dbredis.Redis.SetCache(key+" updated at", updatedAt, cacheTime)

			reportData := &model.ReportFulfillment{StartDate: startDate, EndDate: endDate, WeekNumber: yearWeek, FulfillmentRate: data.FulfillmentRate}
			key = fmt.Sprintf("report_fulfillment week %s %d", yearWeek, warehouseID)
			dbredis.Redis.SetCache(key, reportData, cacheTime)
			dbredis.Redis.SetCache(key+" updated at", updatedAt, cacheTime)
		}

		return
	}

	data, err = GetRangedData(startDate, endDate, warehouseID)

	return
}

// GetDailyData : function to get daily summary data based on date and warehouse
func GetDailyData(dateStr string, warehouseID int64) (data *model.Fulfillment, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	data = new(model.Fulfillment)

	qTotalSO := "select count(distinct(tab2.`SO Code`)) as 'Total SO' " +
		"from " +
		"(select tab.Date, tab.City, tab.`SO Code`, " +
		"tab.Merchant_Id, tab.Merchant_Code, tab.Merchant, " +
		"tab.Item_Id, tab.`Item Code`, tab.Item, " +
		"tab.`Category`, tab.UOM, tab.`SO Qty`, tab.`Deliv Qty`, " +
		"tab.`status`, tab.cancel_type " +
		"from " +
		"(select so.delivery_date as Date, " +
		"case when w2.id in (11,15,17,18,20) then 'Jakarta 2' " +
		"else a2.`name` end as City, " +
		"so.`code` as 'SO Code', " +
		"m2.id as 'Merchant_Id', m2.`code` as 'Merchant_Code', m2.`name` as Merchant, " +
		"p.id as 'Item_Id', p.`code` as 'Item Code', p.`name` as 'Item', " +
		"c.C0_Name as 'Category', u2.`name` as 'UOM', sum(soi.order_qty) as 'SO Qty', " +
		"case when sum(doi2.deliver_qty) is NULL then 0 " +
		"when sum(doi2.deliver_qty) > sum(soi.order_qty) then sum(soi.order_qty) " +
		"else sum(doi2.deliver_qty) end as 'Deliv Qty', " +
		"so.`status`, so.cancel_type, concat(so.`status`,so.cancel_type) as note " +
		"from " +
		"sales_order_item soi  " +
		"left join sales_order so on so.id = soi.sales_order_id  " +
		"left join delivery_order do2 on do2.sales_order_id = so.id AND do2.status not in (3,4) " +
		"left join delivery_order_item doi2 ON doi2.sales_order_item_id = soi.id and doi2.delivery_order_id = do2.id  " +
		"left join warehouse w2 on w2.id = so.warehouse_id " +
		"left join area a2 on a2.id = so.area_id  " +
		"left join branch b2 on b2.id = so.branch_id  " +
		"left join merchant m2 on m2.id = b2.merchant_id  " +
		"left join archetype a3 on a3.id = so.archetype_id  " +
		"left join product p ON p.id = soi.product_id  " +
		"left join (select c2.id,  " +
		"c2.`code` as 'C2_Code',c2.`name` as 'C2_Name', " +
		"c1.`code` as 'C1_code', c1.`name` as 'C1_Name', " +
		"c0.`code` as 'C0_Code', c0.`name` as 'C0_Name' " +
		"from " +
		"category c2 " +
		"left join category c1 on c1.id = c2.parent_id and c1.parent_id = 0 " +
		"left join category c0 on c0.id = c2.grandparent_id and c0.parent_id = 0 and c0.grandparent_id = 0 " +
		"where " +
		"c2.`status` = 1 and c2.grandparent_id != 0 and c2.parent_id != 0) c on c.id = p.category_id  " +
		"left join uom u2 on u2.id = p.uom_id  " +
		"left join business_type bt on bt.id = a3.business_type_id  " +
		"where " +
		"so.delivery_date = ? " +
		"and w2.id = ? " +
		"and so.`status` not in (4) " +
		"and bt.id not in (7,12,14) " +
		"group by Date, City, `SO Code`, Merchant_Id, Merchant_Code, Merchant, Item_Id, `Item Code`, " +
		"Item, Category, `UOM`, `status`, cancel_type, note) tab " +
		"where tab.note not in (32,30)) tab2"

	err = o.Raw(qTotalSO, dateStr, warehouseID).QueryRow(&data.TotalSo)

	qTotalSOUnfulfilled := "select count(distinct(tab2.`SO Code`)) as 'SO Unfulfilled' " +
		"from " +
		"(select tab.Date, tab.City, tab.`SO Code`, " +
		"tab.Merchant_Id, tab.Merchant_Code, tab.Merchant, " +
		"tab.Item_Id, tab.`Item Code`, tab.Item, " +
		"tab.`Category`, tab.UOM, tab.`SO Qty`, tab.`Deliv Qty`, " +
		"tab.`status`, tab.cancel_type " +
		"from " +
		"(select so.delivery_date as Date, " +
		"case when w2.id in (11,15,17,18,20) then 'Jakarta 2' " +
		"else a2.`name` end as City, " +
		"so.`code` as 'SO Code', " +
		"m2.id as 'Merchant_Id', m2.`code` as 'Merchant_Code', m2.`name` as Merchant, " +
		"p.id as 'Item_Id', p.`code` as 'Item Code', p.`name` as 'Item', " +
		"c.C0_Name as 'Category', u2.`name` as 'UOM', sum(soi.order_qty) as 'SO Qty', " +
		"case when sum(doi2.deliver_qty) is NULL then 0 " +
		"when sum(doi2.deliver_qty) > sum(soi.order_qty) then sum(soi.order_qty) " +
		"else sum(doi2.deliver_qty) end as 'Deliv Qty', " +
		"so.`status`, so.cancel_type, concat(so.`status`,so.cancel_type) as note " +
		"from " +
		"sales_order_item soi  " +
		"left join sales_order so on so.id = soi.sales_order_id  " +
		"left join delivery_order do2 on do2.sales_order_id = so.id AND do2.status not in (3,4) " +
		"left join delivery_order_item doi2 ON doi2.sales_order_item_id = soi.id and doi2.delivery_order_id = do2.id  " +
		"left join warehouse w2 on w2.id = so.warehouse_id " +
		"left join area a2 on a2.id = so.area_id  " +
		"left join branch b2 on b2.id = so.branch_id  " +
		"left join merchant m2 on m2.id = b2.merchant_id  " +
		"left join archetype a3 on a3.id = so.archetype_id  " +
		"left join product p ON p.id = soi.product_id  " +
		"left join (select c2.id,  " +
		"c2.`code` as 'C2_Code',c2.`name` as 'C2_Name', " +
		"c1.`code` as 'C1_code', c1.`name` as 'C1_Name', " +
		"c0.`code` as 'C0_Code', c0.`name` as 'C0_Name' " +
		"from " +
		"category c2 " +
		"left join category c1 on c1.id = c2.parent_id and c1.parent_id = 0 " +
		"left join category c0 on c0.id = c2.grandparent_id and c0.parent_id = 0 and c0.grandparent_id = 0 " +
		"where " +
		"c2.`status` = 1 and c2.grandparent_id != 0 and c2.parent_id != 0) c on c.id = p.category_id  " +
		"left join uom u2 on u2.id = p.uom_id  " +
		"left join business_type bt on bt.id = a3.business_type_id  " +
		"where " +
		"so.delivery_date = ? " +
		"and w2.id = ? " +
		"and so.`status` not in (4) " +
		"and bt.id not in (7,12,14) " +
		"group by Date, City, `SO Code`, Merchant_Id, Merchant_Code, Merchant, Item_Id, `Item Code`, " +
		"Item, Category, `UOM`, `status`, cancel_type, note) tab " +
		"where tab.note not in (32,30)) tab2 " +
		"where tab2.`SO Qty` > tab2.`Deliv Qty`"
	err = o.Raw(qTotalSOUnfulfilled, dateStr, warehouseID).QueryRow(&data.TotalSoUnfulfilled)

	qTotalCust := "select count(distinct(tab2.Merchant_Id)) as 'Total Customer' " +
		"from " +
		"(select tab.Date, tab.City, tab.`SO Code`, " +
		"tab.Merchant_Id, tab.Merchant_Code, tab.Merchant, " +
		"tab.Item_Id, tab.`Item Code`, tab.Item, " +
		"tab.`Category`, tab.UOM, tab.`SO Qty`, tab.`Deliv Qty`, " +
		"tab.`status`, tab.cancel_type " +
		"from " +
		"(select so.delivery_date as Date, " +
		"case when w2.id in (11,15,17,18,20) then 'Jakarta 2' " +
		"else a2.`name` end as City, " +
		" so.`code` as 'SO Code', " +
		"m2.id as 'Merchant_Id', m2.`code` as 'Merchant_Code', m2.`name` as Merchant, " +
		"p.id as 'Item_Id', p.`code` as 'Item Code', p.`name` as 'Item', " +
		"c.C0_Name as 'Category', u2.`name` as 'UOM', sum(soi.order_qty) as 'SO Qty', " +
		"case when sum(doi2.deliver_qty) is NULL then 0 " +
		"when sum(doi2.deliver_qty) > sum(soi.order_qty) then sum(soi.order_qty) " +
		"else sum(doi2.deliver_qty) end as 'Deliv Qty', " +
		"so.`status`, so.cancel_type, concat(so.`status`,so.cancel_type) as note " +
		"from " +
		"sales_order_item soi  " +
		"left join sales_order so on so.id = soi.sales_order_id  " +
		"left join delivery_order do2 on do2.sales_order_id = so.id AND do2.status not in (3,4) " +
		"left join delivery_order_item doi2 ON doi2.sales_order_item_id = soi.id and doi2.delivery_order_id = do2.id  " +
		"left join warehouse w2 on w2.id = so.warehouse_id " +
		"left join area a2 on a2.id = so.area_id  " +
		"left join branch b2 on b2.id = so.branch_id  " +
		"left join merchant m2 on m2.id = b2.merchant_id  " +
		"left join archetype a3 on a3.id = so.archetype_id  " +
		"left join product p ON p.id = soi.product_id  " +
		"left join (select c2.id,  " +
		"c2.`code` as 'C2_Code',c2.`name` as 'C2_Name', " +
		"c1.`code` as 'C1_code', c1.`name` as 'C1_Name', " +
		"c0.`code` as 'C0_Code', c0.`name` as 'C0_Name' " +
		"from " +
		"category c2 " +
		"left join category c1 on c1.id = c2.parent_id and c1.parent_id = 0 " +
		"left join category c0 on c0.id = c2.grandparent_id and c0.parent_id = 0 and c0.grandparent_id = 0 " +
		"where " +
		"c2.`status` = 1 and c2.grandparent_id != 0 and c2.parent_id != 0) c on c.id = p.category_id  " +
		"left join uom u2 on u2.id = p.uom_id  " +
		"left join business_type bt on bt.id = a3.business_type_id  " +
		"where " +
		"so.delivery_date = ? " +
		"and w2.id = ? " +
		"and so.`status` not in (4) " +
		"and bt.id not in (7,12,14) " +
		"group by Date, City, `SO Code`, Merchant_Id, Merchant_Code, Merchant, Item_Id, `Item Code`, " +
		"Item, Category, `UOM`, `status`, cancel_type, note) tab " +
		"where tab.note not in (32,30)) tab2 "
	err = o.Raw(qTotalCust, dateStr, warehouseID).QueryRow(&data.TotalCust)

	qTotalCustUnfulfilled := "select count(distinct(tab2.Merchant_Id)) as 'Customer Unfulfilled' " +
		"from " +
		"(select tab.Date, tab.City, tab.`SO Code`, " +
		"tab.Merchant_Id, tab.Merchant_Code, tab.Merchant, " +
		"tab.Item_Id, tab.`Item Code`, tab.Item, " +
		"tab.`Category`, tab.UOM, tab.`SO Qty`, tab.`Deliv Qty`, " +
		"tab.`status`, tab.cancel_type " +
		"from " +
		"(select so.delivery_date as Date, " +
		"case when w2.id in (11,15,17,18,20) then 'Jakarta 2' " +
		"else a2.`name` end as City, " +
		"so.`code` as 'SO Code', " +
		"m2.id as 'Merchant_Id', m2.`code` as 'Merchant_Code', m2.`name` as Merchant, " +
		"p.id as 'Item_Id', p.`code` as 'Item Code', p.`name` as 'Item', " +
		"c.C0_Name as 'Category', u2.`name` as 'UOM', sum(soi.order_qty) as 'SO Qty', " +
		"case when sum(doi2.deliver_qty) is NULL then 0 " +
		"when sum(doi2.deliver_qty) > sum(soi.order_qty) then sum(soi.order_qty) " +
		"else sum(doi2.deliver_qty) end as 'Deliv Qty', " +
		"so.`status`, so.cancel_type, concat(so.`status`,so.cancel_type) as note " +
		"from " +
		"sales_order_item soi  " +
		"left join sales_order so on so.id = soi.sales_order_id  " +
		"left join delivery_order do2 on do2.sales_order_id = so.id AND do2.status not in (3,4) " +
		"left join delivery_order_item doi2 ON doi2.sales_order_item_id = soi.id and doi2.delivery_order_id = do2.id " +
		"left join warehouse w2 on w2.id = so.warehouse_id " +
		"left join area a2 on a2.id = so.area_id  " +
		"left join branch b2 on b2.id = so.branch_id  " +
		"left join merchant m2 on m2.id = b2.merchant_id  " +
		"left join archetype a3 on a3.id = so.archetype_id  " +
		"left join product p ON p.id = soi.product_id  " +
		"left join (select c2.id,  " +
		"c2.`code` as 'C2_Code',c2.`name` as 'C2_Name', " +
		"c1.`code` as 'C1_code', c1.`name` as 'C1_Name', " +
		"c0.`code` as 'C0_Code', c0.`name` as 'C0_Name' " +
		"from " +
		"category c2 " +
		"left join category c1 on c1.id = c2.parent_id and c1.parent_id = 0 " +
		"left join category c0 on c0.id = c2.grandparent_id and c0.parent_id = 0 and c0.grandparent_id = 0 " +
		"where " +
		"c2.`status` = 1 and c2.grandparent_id != 0 and c2.parent_id != 0) c on c.id = p.category_id  " +
		"left join uom u2 on u2.id = p.uom_id  " +
		"left join business_type bt on bt.id = a3.business_type_id  " +
		"where " +
		"so.delivery_date = ? " +
		"and w2.id = ? " +
		"and so.`status` not in (4) " +
		"and bt.id not in (7,12,14) " +
		"group by Date, City, `SO Code`, Merchant_Id, Merchant_Code, Merchant, Item_Id, `Item Code`, " +
		"Item, Category, `UOM`, `status`, cancel_type, note) tab " +
		"where tab.note not in (32,30)) tab2 " +
		"where tab2.`SO Qty` > tab2.`Deliv Qty`"
	err = o.Raw(qTotalCustUnfulfilled, dateStr, warehouseID).QueryRow(&data.TotalCustUnfulfilled)

	qTotalProductUnfulfilled := "select count(distinct(tab2.Item_Id)) as 'Product Unfulfilled' " +
		"from " +
		"(select tab.Date, tab.City, tab.`SO Code`, " +
		"tab.Merchant_Id, tab.Merchant_Code, tab.Merchant, " +
		"tab.Item_Id, tab.`Item Code`, tab.Item, " +
		"tab.`Category`, tab.UOM, tab.`SO Qty`, tab.`Deliv Qty`, " +
		"tab.`status`, tab.cancel_type " +
		"from " +
		"(select so.delivery_date as Date, " +
		"case when w2.id in (11,15,17,18,20) then 'Jakarta 2' " +
		"else a2.`name` end as City, " +
		"so.`code` as 'SO Code', " +
		"m2.id as 'Merchant_Id', m2.`code` as 'Merchant_Code', m2.`name` as Merchant, " +
		"p.id as 'Item_Id', p.`code` as 'Item Code', p.`name` as 'Item', " +
		"c.C0_Name as 'Category', u2.`name` as 'UOM', sum(soi.order_qty) as 'SO Qty', " +
		"case when sum(doi2.deliver_qty) is NULL then 0 " +
		"when sum(doi2.deliver_qty) > sum(soi.order_qty) then sum(soi.order_qty) " +
		"else sum(doi2.deliver_qty) end as 'Deliv Qty', " +
		"so.`status`, so.cancel_type, concat(so.`status`,so.cancel_type) as note " +
		"from " +
		"sales_order_item soi  " +
		"left join sales_order so on so.id = soi.sales_order_id  " +
		"left join delivery_order do2 on do2.sales_order_id = so.id AND do2.status not in (3,4) " +
		"left join delivery_order_item doi2 ON doi2.sales_order_item_id = soi.id and doi2.delivery_order_id = do2.id  " +
		"left join warehouse w2 on w2.id = so.warehouse_id " +
		"left join area a2 on a2.id = so.area_id  " +
		"left join branch b2 on b2.id = so.branch_id  " +
		"left join merchant m2 on m2.id = b2.merchant_id  " +
		"left join archetype a3 on a3.id = so.archetype_id  " +
		"left join product p ON p.id = soi.product_id  " +
		"left join (select c2.id,  " +
		"c2.`code` as 'C2_Code',c2.`name` as 'C2_Name', " +
		"c1.`code` as 'C1_code', c1.`name` as 'C1_Name', " +
		"c0.`code` as 'C0_Code', c0.`name` as 'C0_Name' " +
		"from " +
		"category c2 " +
		"left join category c1 on c1.id = c2.parent_id and c1.parent_id = 0 " +
		"left join category c0 on c0.id = c2.grandparent_id and c0.parent_id = 0 and c0.grandparent_id = 0 " +
		"where " +
		"c2.`status` = 1 and c2.grandparent_id != 0 and c2.parent_id != 0) c on c.id = p.category_id  " +
		"left join uom u2 on u2.id = p.uom_id  " +
		"left join business_type bt on bt.id = a3.business_type_id  " +
		"where " +
		"so.delivery_date = ? " +
		"and w2.id = ? " +
		"and so.`status` not in (4) " +
		"and bt.id not in (7,12,14) " +
		"group by Date, City, `SO Code`, Merchant_Id, Merchant_Code, Merchant, Item_Id, `Item Code`, " +
		"Item, Category, `UOM`, `status`, cancel_type, note) tab " +
		"where tab.note not in (32,30)) tab2 " +
		"where tab2.`SO Qty` > tab2.`Deliv Qty`"
	err = o.Raw(qTotalProductUnfulfilled, dateStr, warehouseID).QueryRow(&data.TotalProductUnfulfilled)

	data.UnfulfillmentRate = math.Round(float64(data.TotalSoUnfulfilled)/float64(data.TotalSo)*100*100) / 100
	data.FulfillmentRate = math.Round((100-data.UnfulfillmentRate)*100) / 100
	data.CustFulfillmentRate = 100 - math.Round(float64(data.TotalCustUnfulfilled)/float64(data.TotalCust)*100*100)/100

	return
}

// GetRangedData : function to get summary data based on range date and warehouse
func GetRangedData(startDate, endDate time.Time, warehouseID int64) (data *model.Fulfillment, err error) {
	var (
		key                                                                                   string
		cacheTime                                                                             time.Duration
		totalSo, totalSoUnfulfilled, totalCust, totalCustUnfulfilled, totalProductUnfulfilled float64
		updatedAt, lastUpdatedAt                                                              time.Time
	)

	data = new(model.Fulfillment)
	for iDate := startDate; !(iDate.After(endDate)); iDate = iDate.AddDate(0, 0, 1) {
		dailyData := new(model.Fulfillment)
		dateStr := iDate.Format("2006-01-02")
		key = fmt.Sprintf("fulfillment_summary day %s %d", dateStr, warehouseID)
		if dbredis.Redis.CheckExistKey(key) {
			dbredis.Redis.GetCache(key, &dailyData)
			dbredis.Redis.GetCache(key+" updated at", &updatedAt)

			if updatedAt.After(lastUpdatedAt) {
				lastUpdatedAt = updatedAt
			}
		} else {
			if dailyData, err = GetDailyData(dateStr, warehouseID); dailyData == nil {
				continue
			}

			updatedAt = time.Now()
			cacheTime = SetCacheTime(iDate)
			dbredis.Redis.SetCache(key, dailyData, cacheTime)
			dbredis.Redis.SetCache(key+" updated at", updatedAt, cacheTime)

			lastUpdatedAt = updatedAt
		}

		totalSo += dailyData.TotalSo
		totalSoUnfulfilled += dailyData.TotalSoUnfulfilled
		totalCust += dailyData.TotalCust
		totalCustUnfulfilled += dailyData.TotalCustUnfulfilled
		totalProductUnfulfilled += dailyData.TotalProductUnfulfilled
	}

	if totalSo == 0 {
		return
	}

	data.TotalSo = totalSo
	data.TotalCust = totalCust
	data.TotalSoUnfulfilled = totalSoUnfulfilled
	data.TotalCustUnfulfilled = totalCustUnfulfilled
	data.TotalProductUnfulfilled = totalProductUnfulfilled
	data.UnfulfillmentRate = math.Round(data.TotalSoUnfulfilled/data.TotalSo*100*100) / 100
	data.FulfillmentRate = math.Round((100-data.UnfulfillmentRate)*100) / 100
	data.CustFulfillmentRate = 100 - math.Round(data.TotalCustUnfulfilled/data.TotalCust*100*100)/100
	data.LastUpdatedAt = lastUpdatedAt

	return
}

// GetReportFulfillment : function to get data for report fulfillment
func GetReportFulfillment(year string, warehouseID int64) (data []*model.ReportFulfillment, lastUpdatedAt time.Time, err error) {
	var (
		key                           string
		i, yearInt                    int
		startDate, endDate, updatedAt time.Time
		cacheTime                     time.Duration
	)

	key = fmt.Sprintf("report_fulfillment year %s %d", year, warehouseID)
	if dbredis.Redis.CheckExistKey(key) {
		dbredis.Redis.GetCache(key, &data)
		dbredis.Redis.GetCache(key+" updated at", &lastUpdatedAt)

		return
	}

	yearInt, _ = strconv.Atoi(year)

	// get first week number of a year
	startDate = util.GetWeekStart(yearInt, 1)
	endDate = startDate.AddDate(0, 0, 6)
	firstWeekOfAYear := 0
	if yearInt == 2021 {
		firstWeekOfAYear = 46
	} else {
		if startDate.Year() != endDate.Year() || endDate.Day() <= 7 {
			firstWeekOfAYear = 1
		}
	}

	// get last week number of a year
	startDate = util.GetWeekStart(yearInt, 52)
	endDate = startDate.AddDate(0, 0, 6)
	lastWeekOfAYear := 53
	if startDate.Year() != endDate.Year() || endDate.Day() == 31 {
		lastWeekOfAYear = 52
	}

	// loop through weeks in a year
	for i = firstWeekOfAYear; i <= lastWeekOfAYear; i++ {
		startDate = util.GetWeekStart(yearInt, i)
		endDate = startDate.AddDate(0, 0, 6)

		yearWeek := fmt.Sprintf("%s%02d", year, i)
		key = fmt.Sprintf("report_fulfillment week %s %d", yearWeek, warehouseID)

		if dbredis.Redis.CheckExistKey(key) {
			reportData := new(model.ReportFulfillment)
			dbredis.Redis.GetCache(key, &reportData)
			dbredis.Redis.GetCache(key+" updated at", &updatedAt)

			if updatedAt.After(lastUpdatedAt) {
				lastUpdatedAt = updatedAt
			}

			data = append(data, reportData)

			continue
		}

		summaryData := new(model.Fulfillment)
		summaryData, err = GetRangedData(startDate, endDate, warehouseID)
		if summaryData.TotalSo == 0 {
			continue
		}

		cacheTime = SetCacheTime(endDate)
		dbredis.Redis.SetCache(key, summaryData, cacheTime)
		dbredis.Redis.SetCache(key+" updated at", time.Now(), cacheTime)
		lastUpdatedAt = time.Now()

		data = append(data, &model.ReportFulfillment{WeekNumber: yearWeek, StartDate: startDate, EndDate: endDate, FulfillmentRate: summaryData.FulfillmentRate})
	}

	if len(data) > 0 {
		cacheTime = SetCacheTime(lastUpdatedAt)
		key = fmt.Sprintf("report_fulfillment year %s %d", year, warehouseID)
		dbredis.Redis.SetCache(key, data, cacheTime)
		dbredis.Redis.SetCache(key+" updated at", lastUpdatedAt, cacheTime)
	}

	return
}

// GetUnfulfillmentProduct : function to get unfullfilled product list
func GetUnfulfillmentProduct(date []string, warehouseID int64) (data []*model.UnfulfilledProduct, err error) {
	var (
		key       string
		cacheTime time.Duration
	)

	startDate, _ := time.Parse("2006-01-02", date[0])
	endDate, _ := time.Parse("2006-01-02", date[1])

	// check if range is a week
	// if it is a week then get data from cache or set data into cache
	if int(startDate.Weekday()) == 1 && endDate.Sub(startDate).Hours()/24 == 5 {
		year, week := startDate.ISOWeek()
		yearWeek := strconv.Itoa(year) + strconv.Itoa(week)
		key = fmt.Sprintf("unfulfilled_product week %s %d", yearWeek, warehouseID)
		if dbredis.Redis.CheckExistKey(key) {
			dbredis.Redis.GetCache(key, &data)
			return
		}

		if data, err = GetRangedProductData(startDate, endDate, warehouseID); err == nil {
			if len(data) > 0 {
				cacheTime = SetCacheTime(endDate)
				dbredis.Redis.SetCache(key, data, cacheTime)
			}

			return
		}
	}

	data, err = GetRangedProductData(startDate, endDate, warehouseID)

	return
}

// GetDailyProductData : function to get unfulfilled product based on date and warehouse
func GetDailyProductData(dateStr string, warehouseID int64) (data []*model.UnfulfilledProduct, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	q := "SELECT tab2.Item_Id as product_id, tab2.Item as 'product', tab2.UOM as 'uom', " +
		"count(distinct(tab2.`SO Code`)) as 'count_so', " +
		"count(distinct(tab2.Merchant_Id)) as 'count_cust', " +
		"sum(tab2.`SO Qty`)-sum(tab2.`Deliv Qty`) as  'unfulfilled_qty' " +
		"from " +
		"(select tab.Date, tab.City, tab.`SO Code`, " +
		"tab.Merchant_Id, tab.Merchant_Code, tab.Merchant, " +
		"tab.Item_Id, tab.`Item Code`, tab.Item, " +
		"tab.`Category`, tab.UOM, tab.`SO Qty`, tab.`Deliv Qty`, " +
		"tab.`status`, tab.cancel_type " +
		"from " +
		"(select so.delivery_date as Date, " +
		"case when w2.id in (11,15,17,18,20) then 'Jakarta 2' " +
		"else a2.`name` end as City, " +
		"so.`code` as 'SO Code', " +
		"m2.id as 'Merchant_Id', m2.`code` as 'Merchant_Code', m2.`name` as Merchant, " +
		"p.id as 'Item_Id', p.`code` as 'Item Code', p.`name` as 'Item', " +
		"c.C0_Name as 'Category', u2.`name` as 'UOM', sum(soi.order_qty) as 'SO Qty', " +
		"case when sum(doi2.deliver_qty) is NULL then 0 " +
		"when sum(doi2.deliver_qty) > sum(soi.order_qty) then sum(soi.order_qty) " +
		"else sum(doi2.deliver_qty) end as 'Deliv Qty', " +
		"so.`status`, so.cancel_type, concat(so.`status`,so.cancel_type) as note " +
		"from " +
		"sales_order_item soi  " +
		"left join sales_order so on so.id = soi.sales_order_id  " +
		"left join delivery_order do2 on do2.sales_order_id = so.id AND do2.status not in (3,4) " +
		"left join delivery_order_item doi2 ON doi2.sales_order_item_id = soi.id and doi2.delivery_order_id = do2.id  " +
		"left join warehouse w2 on w2.id = so.warehouse_id " +
		"left join area a2 on a2.id = so.area_id  " +
		"left join branch b2 on b2.id = so.branch_id  " +
		"left join merchant m2 on m2.id = b2.merchant_id  " +
		"left join archetype a3 on a3.id = so.archetype_id  " +
		"left join product p ON p.id = soi.product_id  " +
		"left join (select c2.id,  " +
		"c2.`code` as 'C2_Code',c2.`name` as 'C2_Name', " +
		"c1.`code` as 'C1_code', c1.`name` as 'C1_Name', " +
		"c0.`code` as 'C0_Code', c0.`name` as 'C0_Name' " +
		"from " +
		"category c2 " +
		"left join category c1 on c1.id = c2.parent_id and c1.parent_id = 0 " +
		"left join category c0 on c0.id = c2.grandparent_id and c0.parent_id = 0 and c0.grandparent_id = 0 " +
		"where " +
		"c2.`status` = 1 and c2.grandparent_id != 0 and c2.parent_id != 0) c on c.id = p.category_id  " +
		"left join uom u2 on u2.id = p.uom_id  " +
		"left join business_type bt on bt.id = a3.business_type_id  " +
		"where " +
		"so.delivery_date = ? " +
		"and w2.id = ? " +
		"and so.`status` not in (4) " +
		"and bt.id not in (7,12,14) " +
		"group by Date, City, `SO Code`, Merchant_Id, Merchant_Code, Merchant, Item_Id, `Item Code`, " +
		"Item, Category, `UOM`, `status`, cancel_type, note) tab " +
		"where tab.note not in (32,30)) tab2 " +
		"where tab2.`SO Qty` > tab2.`Deliv Qty` " +
		"group by 1,2,3"
	_, err = o.Raw(q, dateStr, warehouseID).QueryRows(&data)

	return
}

// GetRangedProductData : function to get unfulfilled product data based on date range and warehouse
func GetRangedProductData(startDate, endDate time.Time, warehouseID int64) (data []*model.UnfulfilledProduct, err error) {
	var (
		dateStr, key string
		cacheTime    time.Duration
		dailyData    []*model.UnfulfilledProduct
	)

	type ProductData struct {
		SoCount    int64
		CustCount  int64
		ProductQty float64
		Product    string
		Uom        string
	}

	productData := make(map[int64]*ProductData)

	// loop through dates
	for iDate := startDate; !(iDate.After(endDate)); iDate = iDate.AddDate(0, 0, 1) {
		dateStr = iDate.Format("2006-01-02")
		key = fmt.Sprintf("unfulfilled_product day %s %d", dateStr, warehouseID)
		if dbredis.Redis.CheckExistKey(key) {
			dbredis.Redis.GetCache(key, &data)

			// sort desc grouped data by it's quantity
			sort.SliceStable(data[:], func(i, j int) bool {
				return data[i].UnfulfilledQty > data[j].UnfulfilledQty
			})
			return
		} else {
			if dailyData, err = GetDailyProductData(dateStr, warehouseID); dailyData == nil {
				continue
			}

			cacheTime = SetCacheTime(iDate)
			dbredis.Redis.SetCache(key, dailyData, cacheTime)
		}

		// loop through dailyData to group each data based on it's product id
		for _, v := range dailyData {
			if _, isExists := productData[v.ID]; isExists {
				productData[v.ID].SoCount += v.UnfulfilledSO
				productData[v.ID].CustCount += v.UnfulfilledCust
				productData[v.ID].ProductQty += v.UnfulfilledQty
			} else {
				productData[v.ID] = &ProductData{SoCount: v.UnfulfilledSO, CustCount: v.UnfulfilledCust, ProductQty: v.UnfulfilledQty, Product: v.Product, Uom: v.Uom}
			}
		}
	}

	// append grouped data
	for i, v := range productData {
		data = append(data, &model.UnfulfilledProduct{ID: i, Product: v.Product, Uom: v.Uom, UnfulfilledSO: v.SoCount, UnfulfilledCust: v.CustCount, UnfulfilledQty: v.ProductQty})
	}

	// sort desc grouped data by it's quantity
	sort.SliceStable(data[:], func(i, j int) bool {
		return data[i].UnfulfilledQty > data[j].UnfulfilledQty
	})

	return
}

// SetCacheTime : function to set cache time
func SetCacheTime(date time.Time) (cacheTime time.Duration) {
	// get current date
	currentDate := time.Now()
	// count hours difference between now and date parameter
	totalHours := currentDate.Sub(date).Hours()

	switch {
	case totalHours <= 48: // case <= 2 days
		cacheTime = 2 * time.Hour // cache for 2 hours
	case totalHours <= 168: // case <= 7 days
		cacheTime = 24 * time.Hour // cache for 1 day
	case totalHours <= 2160: // case <= 90 days
		cacheTime = 720 * time.Hour // cache for 30 days
	default: // else case
		cacheTime = 2160 * time.Hour // cache for 90 days
	}

	return
}
