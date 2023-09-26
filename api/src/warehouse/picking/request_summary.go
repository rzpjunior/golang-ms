// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"math"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type summaryRequest struct {
	WarehouseID      string   `json:"warehouse_id"`
	HubPickingList   bool     `json:"hub_picking_list" valid:"required"`
	HubID            string   `json:"hub_id"`
	DeliveryDate     string   `json:"delivery_date"`
	CityID           string   `json:"city_id"`
	DistrictID       []string `json:"district_id"`
	WrtID            []string `json:"wrt_id"`
	BusinessTypeID   []string `json:"business_type_id"`
	SalesOrderTypeID []string `json:"sales_order_type_id"`

	ArrCity             []int64 `json:"-"`
	ArrSubDistrictID    []int64 `json:"-"`
	ArrDistrictID       []int64 `json:"-"`
	ArrWrtID            []int64 `json:"-"`
	ArrBusinessTypeID   []int64 `json:"-"`
	ArrSalesOrderTypeID []int64 `json:"-"`

	QueryStringDistrict       string `json:"-"`
	QueryStringSubDistrict    string `json:"-"`
	QueryStringWrt            string `json:"-"`
	QueryStringBusinessType   string `json:"-"`
	QueryStringSalesOrderType string `json:"-"`

	DeliveryDateTime time.Time `json:"-"`

	TotalSalesOrder              int64   `json:"-"`
	TotalWeight                  float64 `json:"-"`
	HighestSalesOrderWeight      float64 `json:"-"`
	HighestSalesOrderItemsWeight float64 `json:"-"`

	Staff     *model.Staff     `json:"-"`
	Warehouse *model.Warehouse `json:"-"`
	Hub       *model.Warehouse `json:"-"`
	City      *model.City      `json:"-"`

	Session *auth.SessionData `json:"-"`
}

type pickingRouteSummary struct {
	TotalSalesOrder              int64   `json:"total_sales_order"`
	TotalWeight                  float64 `json:"total_weight"`
	HighestSalesOrderWeight      float64 `json:"highest_sales_order_weight"`
	HighestSalesOrderItemsWeight float64 `json:"highest_sales_order_item_weight"`
}

// Validate : function to validate routing request data
func (c *summaryRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	f := orm.NewOrm()
	f.Using("read_only")

	var (
		err             error
		errorValidation bool

		maxWeightSOI float64
		maxWeightSO  float64

		warehouseID      int64
		districtID       int64
		wrtID            int64
		cityID           int64
		businessTypeID   int64
		salesOrderTypeID int64
		hubFilter        bool
	)

	salesOrderIDMap := map[int64]bool{}

	if len(c.WarehouseID) > 0 {
		if warehouseID, err = common.Decrypt(c.WarehouseID); err != nil {
			o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
			errorValidation = true
		}
		if c.Warehouse, err = repository.ValidWarehouse(warehouseID); err != nil {
			o.Failure("warehouse.invalid", util.ErrorInvalidData("warehouse"))
			errorValidation = true
		}

		if c.HubPickingList {
			if len(c.HubID) > 0 {
				if warehouseID, err = common.Decrypt(c.HubID); err != nil {
					o.Failure("warehouse_hub_id.invalid", util.ErrorInvalidData("warehouse hub"))
					errorValidation = true
				}
				if c.Hub, err = repository.ValidWarehouse(warehouseID); err != nil {
					o.Failure("warehouse_hub.invalid", util.ErrorInvalidData("warehouse hub"))
					errorValidation = true
				}
				hubFilter = true
			}

			if hubFilter {
				_, err = f.Raw("SELECT sub_district_id FROM warehouse_coverage WHERE warehouse_id = ?", c.Hub.ID).QueryRows(&c.ArrSubDistrictID)
				if len(c.ArrSubDistrictID) == 0 {
					o.Failure("warehouse_coverage.invalid", util.ErrorWarehouseCoverage())
					errorValidation = true
				}

			} else {
				_, err = f.Raw("SELECT sub_district_id FROM warehouse_coverage WHERE warehouse_id = ? AND sub_district_id NOT IN (SELECT sub_district_id FROM warehouse_coverage WHERE parent_warehouse_id = ?)", c.Warehouse.ID, c.Warehouse.ID).QueryRows(&c.ArrSubDistrictID)
			}

			for range c.ArrSubDistrictID {
				c.QueryStringSubDistrict += "?,"
			}
			c.QueryStringSubDistrict = strings.TrimSuffix(c.QueryStringSubDistrict, ",")

		} else {
			if len(c.CityID) > 0 {
				if cityID, err = common.Decrypt(c.CityID); err != nil {
					o.Failure("city_id.invalid", util.ErrorInvalidData("city"))
					errorValidation = true
				}
				if c.City, err = repository.ValidCity(cityID); err != nil {
					o.Failure("city.invalid", util.ErrorInvalidData("city"))
					errorValidation = true
				}

				c.ArrCity = append(c.ArrCity, cityID)
			}

			for _, v := range c.DistrictID {
				if districtID, err = common.Decrypt(v); err != nil {
					o.Failure("district_id.invalid", util.ErrorInvalidData("district"))
					errorValidation = true
				}
				c.ArrDistrictID = append(c.ArrDistrictID, districtID)

				c.QueryStringDistrict += "?,"
			}
			c.QueryStringDistrict = strings.TrimSuffix(c.QueryStringDistrict, ",")
		}
	}

	if len(c.DeliveryDate) > 0 {
		if c.DeliveryDateTime, err = time.Parse("2006-01-02", c.DeliveryDate); err != nil {
			o.Failure("delivery_date.invalid", util.ErrorInputRequired("delivery_date date"))
		}
	}

	for _, v := range c.WrtID {
		if wrtID, err = common.Decrypt(v); err != nil {
			o.Failure("wrt_id.invalid", util.ErrorInvalidData("wrt"))
			errorValidation = true
		}
		c.ArrWrtID = append(c.ArrWrtID, wrtID)

		c.QueryStringWrt += "?,"
	}
	c.QueryStringWrt = strings.TrimSuffix(c.QueryStringWrt, ",")

	// business type filter
	for _, v := range c.BusinessTypeID {
		if businessTypeID, err = common.Decrypt(v); err != nil {
			o.Failure("business_type_id.invalid", util.ErrorInvalidData("business type"))
			errorValidation = true
		}
		c.ArrBusinessTypeID = append(c.ArrBusinessTypeID, businessTypeID)

		c.QueryStringBusinessType += "?,"
	}
	c.QueryStringBusinessType = strings.TrimSuffix(c.QueryStringBusinessType, ",")

	// sales order type filter
	for _, v := range c.SalesOrderTypeID {
		if salesOrderTypeID, err = common.Decrypt(v); err != nil {
			o.Failure("business_type_id.invalid", util.ErrorInvalidData("business type"))
			errorValidation = true
		}
		c.ArrSalesOrderTypeID = append(c.ArrSalesOrderTypeID, salesOrderTypeID)

		c.QueryStringSalesOrderType += "?,"
	}
	c.QueryStringSalesOrderType = strings.TrimSuffix(c.QueryStringSalesOrderType, ",")

	if len(c.DeliveryDate) > 0 && len(c.WarehouseID) > 0 {
		//CHECK IF THERE'S AN ERROR WITHIN THE VALIDATION PHASE
		if errorValidation {
			return o
		}

		var filterDistrict string
		if len(c.ArrDistrictID) != 0 {
			filterDistrict = "AND ad.district_id IN (" + c.QueryStringDistrict + ") "
		}

		var filterSubDistrict string
		if len(c.ArrSubDistrictID) != 0 {
			filterDistrict = "AND ad.sub_district_id IN (" + c.QueryStringSubDistrict + ") "
		}

		var filterWrt string
		if len(c.ArrWrtID) != 0 {
			filterWrt = "AND w.id IN (" + c.QueryStringWrt + ") "
		}

		var filterBusinessType string
		if len(c.ArrBusinessTypeID) != 0 {
			filterBusinessType = "AND bt.id IN (" + c.QueryStringBusinessType + ") "
		}

		var filterSalesOrderType string
		if len(c.ArrSalesOrderTypeID) != 0 {
			filterSalesOrderType = "AND so.order_type_sls_id IN (" + c.QueryStringSalesOrderType + ") "
		}

		var filterCity string
		if len(c.ArrCity) != 0 {
			filterCity = "AND ad.city_id = ? "
		}

		var salesOrderInformation []GenerateCodePickingList

		q := "SELECT so2.id, so2.code 'so_code', so2.so_total, w.name 'wrt' , p.name 'product_name', soi2.weight 'weight_item', soi2.order_qty 'order_item' " +
			"FROM sales_order_item soi2 " +
			"JOIN(SELECT so.id, so.wrt_id, so.status , so.code , so.branch_id , SUM(soi.weight) 'so_total' " +
			"FROM sales_order so " +
			"JOIN sales_order_item soi ON soi.sales_order_id = so.id " +
			"LEFT JOIN picking_order_assign poa ON poa.sales_order_id = so.id " +
			"LEFT JOIN archetype a ON a.id = so.archetype_id " +
			"LEFT JOIN business_type bt ON bt.id = a.business_type_id " +
			"WHERE so.delivery_date = ? " + filterBusinessType + filterSalesOrderType +
			"and so.status IN (1,9,12) and so.order_type_sls_id != 10 and so.id not in (SELECT id from sales_order so2 WHERE so2.status =1 and so2.term_payment_sls_id = 11) " +
			"and so.warehouse_id = ? and poa.id is NULL " +
			"GROUP BY so.id) so2 on so2.id = soi2.sales_order_id " +
			"JOIN product p ON p.id = soi2.product_id " +
			"JOIN wrt w ON w.id = so2.wrt_id " + filterWrt + " " +
			"JOIN branch b ON b.id = so2.branch_id " +
			"JOIN adm_division ad ON ad.sub_district_id = b.sub_district_id " + filterDistrict + filterSubDistrict + filterCity + " ORDER BY ad.district_id, w.name, so2.code"

		if _, err = f.Raw(q, c.DeliveryDate, c.ArrBusinessTypeID, c.ArrSalesOrderTypeID, c.Warehouse, c.ArrWrtID, c.ArrDistrictID, c.ArrSubDistrictID, c.ArrCity).QueryRows(&salesOrderInformation); err != nil {
			o.Failure("sales_order.not_found", util.ErrorNotFound("sales order"))
		} else {
			for _, v := range salesOrderInformation {
				if salesOrderIDMap[v.SalesOrderID] != true {
					salesOrderIDMap[v.SalesOrderID] = true
				}

				if v.TotalWeight > maxWeightSO {
					maxWeightSO = v.TotalWeight
				}

				if v.WeightItem > maxWeightSOI {
					maxWeightSOI = v.WeightItem
				}

				c.TotalWeight += v.WeightItem
			}
		}

		c.TotalSalesOrder = int64(len(salesOrderIDMap))
		c.TotalWeight = math.Round(c.TotalWeight*100) / 100
		c.HighestSalesOrderWeight = math.Round(maxWeightSO*100) / 100
		c.HighestSalesOrderItemsWeight = math.Round(maxWeightSOI*100) / 100
	}

	return o
}

// Messages : function to return error validation messages
func (c *summaryRequest) Messages() map[string]string {
	messages := map[string]string{}
	return messages
}
