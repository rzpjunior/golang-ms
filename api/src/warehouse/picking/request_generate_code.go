// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"strings"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/dbredis"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// generateCodePickingRequest : struct to hold picking assign request data
type generateCodePickingRequest struct {
	WarehouseID      string             `json:"warehouse_id" valid:"required"`
	HubPickingList   bool               `json:"hub_picking_list"`
	HubID            string             `json:"hub_id"`
	CityID           string             `json:"city_id"`
	DistrictID       []string           `json:"district_id"`
	WrtID            []string           `json:"wrt_id"`
	BusinessTypeID   []string           `json:"business_type_id"`
	SalesOrderTypeID []string           `json:"sales_order_type_id"`
	JobsID           primitive.ObjectID `json:"jobs_id"`
	DeliveryDate     string             `json:"delivery_date" valid:"required"`
	LimitSalesOrder  int                `json:"limit_sales_order" valid:"required|gt:0"`
	LimitWeight      float64            `json:"limit_weight" valid:"required|gt:0"`
	Note             string             `json:"note"`

	ArrCity             []int64 `json:"-"`
	ArrDistrictID       []int64 `json:"-"`
	ArrSubDistrictID    []int64 `json:"-"`
	ArrWrtID            []int64 `json:"-"`
	ArrBusinessTypeID   []int64 `json:"-"`
	ArrSalesOrderTypeID []int64 `json:"-"`

	QueryStringDistrict       string `json:"-"`
	QueryStringSubDistrict    string `json:"-"`
	QueryStringWrt            string `json:"-"`
	QueryStringBusinessType   string `json:"-"`
	QueryStringSalesOrderType string `json:"-"`

	DeliveryDateTime time.Time `json:"-"`

	TypeRequest string `json:"type"`
	CheckRedis  bool   `json:"check_redis"` // for checking key in the first validation, will not check on kafka consumer validation

	Staff     *model.Staff     `json:"-"`
	Warehouse *model.Warehouse `json:"-"`
	Hub       *model.Warehouse `json:"-"`
	City      *model.City      `json:"-"`

	PickingListFinal map[string]ListPl `json:"picking_list_final"`
	PickingListObj   []*PickingListObj `json:"picking_list_obj"`

	PickingOrder *model.PickingOrder `json:"picking_order"`

	Session *auth.SessionData `json:"-"`
}

type ListPl struct {
	SalesOrderID []int64
	TotalWeight  float64
}

type PickingListObj struct {
	Code           string  `json:"code"`
	TotalWeight    float64 `json:"total_weight"`
	SalesOrderID   []int64 `json:"sales_order_id"`
	PickingRouting int8    `json:"picking_routing"`
}

// Validate : function to validate picking assign request data
func (r *generateCodePickingRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var (
		err              error
		warehouseID      int64
		districtID       int64
		wrtID            int64
		cityID           int64
		businessTypeID   int64
		salesOrderTypeID int64
		hubFilter        bool
	)

	if warehouseID, err = common.Decrypt(r.WarehouseID); err != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}
	if r.Warehouse, err = repository.ValidWarehouse(warehouseID); err != nil {
		o.Failure("warehouse.invalid", util.ErrorInvalidData("warehouse"))
		return o
	}

	// check if there's an on going picking list creation for the associated warehouse
	if r.CheckRedis == false {
		// if key exist
		if dbredis.Redis.CheckExistKey("picking_list_" + r.WarehouseID) {
			o.Failure("picking_list.invalid", util.ErrorCreationInProgress("picking list"))
			return o
		}

		dbredis.Redis.SetCache("picking_list_"+r.WarehouseID, true, 15*time.Second)

		// set to true so consumer will not check on the key anymore
		r.CheckRedis = true
	}

	if r.DeliveryDateTime, err = time.Parse("2006-01-02", r.DeliveryDate); err != nil {
		o.Failure("delivery_date.invalid", util.ErrorInputRequired("delivery_date date"))
	}

	if r.HubPickingList {
		if len(r.HubID) > 0 {
			if warehouseID, err = common.Decrypt(r.HubID); err != nil {
				o.Failure("warehouse_hub_id.invalid", util.ErrorInvalidData("warehouse hub"))
			}
			if r.Hub, err = repository.ValidWarehouse(warehouseID); err != nil {
				o.Failure("warehouse_hub.invalid", util.ErrorInvalidData("warehouse hub"))
				return o
			}
			hubFilter = true
		}

		if hubFilter {
			_, err = o1.Raw("SELECT sub_district_id FROM warehouse_coverage WHERE warehouse_id = ?", r.Hub.ID).QueryRows(&r.ArrSubDistrictID)
			if len(r.ArrSubDistrictID) == 0 {
				o.Failure("warehouse_coverage.invalid", util.ErrorWarehouseCoverage())
			}

		} else {
			_, err = o1.Raw("SELECT sub_district_id FROM warehouse_coverage WHERE warehouse_id = ? AND sub_district_id NOT IN (SELECT sub_district_id FROM warehouse_coverage WHERE parent_warehouse_id = ?)", r.Warehouse.ID, r.Warehouse.ID).QueryRows(&r.ArrSubDistrictID)
		}

		for range r.ArrSubDistrictID {
			r.QueryStringSubDistrict += "?,"
		}
		r.QueryStringSubDistrict = strings.TrimSuffix(r.QueryStringSubDistrict, ",")

	} else {
		if len(r.CityID) > 0 {
			if cityID, err = common.Decrypt(r.CityID); err != nil {
				o.Failure("city_id.invalid", util.ErrorInvalidData("city"))
			}
			if r.City, err = repository.ValidCity(cityID); err != nil {
				o.Failure("city.invalid", util.ErrorInvalidData("city"))
				return o
			}

			r.ArrCity = append(r.ArrCity, cityID)
		}

		for _, v := range r.DistrictID {
			if districtID, err = common.Decrypt(v); err != nil {
				o.Failure("district_id.invalid", util.ErrorInvalidData("district"))
			}
			r.ArrDistrictID = append(r.ArrDistrictID, districtID)

			r.QueryStringDistrict += "?,"
		}
		r.QueryStringDistrict = strings.TrimSuffix(r.QueryStringDistrict, ",")
	}

	for _, v := range r.WrtID {
		if wrtID, err = common.Decrypt(v); err != nil {
			o.Failure("wrt_id.invalid", util.ErrorInvalidData("wrt"))
		}
		r.ArrWrtID = append(r.ArrWrtID, wrtID)

		r.QueryStringWrt += "?,"
	}
	r.QueryStringWrt = strings.TrimSuffix(r.QueryStringWrt, ",")

	// business type filter
	for _, v := range r.BusinessTypeID {
		if businessTypeID, err = common.Decrypt(v); err != nil {
			o.Failure("business_type_id.invalid", util.ErrorInvalidData("business type"))
		}
		r.ArrBusinessTypeID = append(r.ArrBusinessTypeID, businessTypeID)

		r.QueryStringBusinessType += "?,"
	}
	r.QueryStringBusinessType = strings.TrimSuffix(r.QueryStringBusinessType, ",")

	// sales order type filter
	for _, v := range r.SalesOrderTypeID {
		if salesOrderTypeID, err = common.Decrypt(v); err != nil {
			o.Failure("business_type_id.invalid", util.ErrorInvalidData("business type"))
		}
		r.ArrSalesOrderTypeID = append(r.ArrSalesOrderTypeID, salesOrderTypeID)

		r.QueryStringSalesOrderType += "?,"
	}
	r.QueryStringSalesOrderType = strings.TrimSuffix(r.QueryStringSalesOrderType, ",")

	r.PickingListFinal = make(map[string]ListPl)

	err = o1.Raw("SELECT id, code, warehouse_id, recognition_date, note, status FROM picking_order where warehouse_id = ? and recognition_date = ?", r.Warehouse.ID, r.DeliveryDate).QueryRow(&r.PickingOrder)

	return o
}

// Messages : function to return error validation messages
func (r *generateCodePickingRequest) Messages() map[string]string {
	messages := map[string]string{
		"delivery_date.required":     util.ErrorInputRequired("delivery date"),
		"warehouse_id.required":      util.ErrorInputRequired("warehouse"),
		"limit_sales_order.required": util.ErrorGreater("limit sales order", "0"),
		"limit_weight.required":      util.ErrorGreater("limit weight", "0"),
	}

	return messages
}
