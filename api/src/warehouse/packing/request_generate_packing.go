// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packing

import (
	"sort"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// generatePackingRequest: struct generate picking
type generatePackingRequest struct {
	WarehouseID           string                                  `json:"warehouse_id" valid:"required"`
	DeliveryDate          string                                  `json:"delivery_date" valid:"required"`
	Note                  string                                  `json:"note"`
	KeysMapProductPack    []int                                   `json:"-"`
	MapProductPack        map[int64]map[float64]float64           `json:"-"`
	MapResProductPackType map[int64]map[float64]float64           `json:"-"`
	MapExistData          map[int64]map[float64]float64           `json:"-"`
	MapOldExistData       map[int64]map[float64]float64           `json:"-"`
	MapPackingSalesOrder  map[int64]map[float64]map[int64]float64 `json:"-"`

	SalesOrderItems []*SalesOrderItemByProduct `json:"-"`
	ResponseData    []*ResponseData            `json:"-"`
	PackType        []float64                  `json:"-"`

	DeliveryDateTime time.Time           `json:"-"`
	PackingOrder     *model.PackingOrder `json:"-"`
	Warehouse        *model.Warehouse    `json:"-"`
	Session          *auth.SessionData   `json:"-"`
}

type SalesOrderItemByProduct struct {
	SalesOrderID    int64   `orm:"column(id)" json:"-"`
	OrderMinimalQty float64 `orm:"column(order_min_qty)" json:"-"`
	UomID           int64   `orm:"column(uom_id)" json:"-"`
	ProductID       int64   `orm:"column(product_id)" json:"-"`
	OrderQty        float64 `orm:"column(order_qty)" json:"-"`
}

type ResponseData struct {
	ID                primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
	PackingOrderID    int64               `bson:"packing_order_id" json:"packing_order_id"`
	ProductID         int64               `bson:"product_id" json:"product_id"`
	PackType          float64             `bson:"pack_type" json:"pack_type"`
	ExpectedTotalPack float64             `bson:"expected_total_pack" json:"expected_total_pack"`
	ActualTotalPack   float64             `bson:"actual_total_pack" json:"actual_total_pack"`
	WeightPack        float64             `bson:"weight_pack" json:"weight_pack"`
	Product           *model.Product      `json:"product,omitempty"`
	PackingOrder      *model.PackingOrder `json:"packing_order,omitempty"`
}

// Validate : function to validate packing generate request data
func (c *generatePackingRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var err error
	var warehouseID int64

	if c.DeliveryDateTime, err = time.Parse("2006-01-02", c.DeliveryDate); err != nil {
		o.Failure("delivery_date.invalid", util.ErrorInputRequired("delivery_date date"))
	}

	if warehouseID, err = common.Decrypt(c.WarehouseID); err == nil {
		if c.Warehouse, err = repository.ValidWarehouse(warehouseID); err != nil {
			o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
		}
	} else {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}

	o1.Raw("select p.order_min_qty,p.uom_id,so.id, soi.product_id, soi.order_qty order_qty "+
		"FROM sales_order_item soi "+
		"JOIN sales_order so ON so.id = soi.sales_order_id "+
		"JOIN product p ON p.id = soi.product_id and p.packability = 1 "+
		"where so.status != 3 and so.order_type_sls_id NOT IN (5,10) and so.delivery_date = ? and so.warehouse_id = ? ", c.DeliveryDate, c.Warehouse.ID).QueryRows(&c.SalesOrderItems)

	o1.Raw("select * from packing_order po where po.delivery_date = ? and po.warehouse_id  = ?", c.DeliveryDate, c.Warehouse.ID).QueryRow(&c.PackingOrder)
	o1.Raw("SELECT value_name FROM glossary WHERE `table` = 'packing_order' and `attribute` = 'pack_size'").QueryRows(&c.PackType)

	c.MapResProductPackType = make(map[int64]map[float64]float64)
	c.MapProductPack = make(map[int64]map[float64]float64)
	c.MapExistData = make(map[int64]map[float64]float64)
	c.MapOldExistData = make(map[int64]map[float64]float64)
	c.MapPackingSalesOrder = make(map[int64]map[float64]map[int64]float64)

	// region for define product based on configuration pack type
	for _, v := range c.SalesOrderItems {
		c.MapProductPack[v.ProductID] = make(map[float64]float64)
		c.MapResProductPackType[v.ProductID] = make(map[float64]float64)
		c.MapExistData[v.ProductID] = make(map[float64]float64)
		c.MapOldExistData[v.ProductID] = make(map[float64]float64)
		c.MapPackingSalesOrder[v.ProductID] = make(map[float64]map[int64]float64)
		c.MapPackingSalesOrder[v.ProductID][-1] = make(map[int64]float64)
		for _, v2 := range c.PackType {
			c.MapPackingSalesOrder[v.ProductID][v2] = make(map[int64]float64)
			c.MapProductPack[v.ProductID][v2] = v2
		}

	}

	// endregion
	sort.Sort(sort.Reverse(sort.Float64Slice(c.PackType)))
	return o
}

// Messages : function to return error validation messages
func (c *generatePackingRequest) Messages() map[string]string {
	messages := map[string]string{
		"warehouse_id.required":  util.ErrorInputRequired("warehouse"),
		"delivery_date.required": util.ErrorInputRequired("delivery date"),
	}

	return messages
}
