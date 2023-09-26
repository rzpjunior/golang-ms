// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/util"
)

func init() {
	orm.RegisterModel(new(ConsolidatedShipment))
}

// ConsolidatedShipment: struct to hold consolidated_shipment model data for database
type ConsolidatedShipment struct {
	ID                int64     `orm:"column(id);auto" json:"-"`
	Code              string    `orm:"column(code);size(50);null" json:"code"`
	DriverName        string    `orm:"column(driver_name);size(100);null" json:"driver_name"`
	VehicleNumber     string    `orm:"column(vehicle_number);size(10);null" json:"vehicle_number"`
	DriverPhoneNumber string    `orm:"column(driver_phone_number);size(15);null" json:"driver_phone_number"`
	DeltaPrint        int8      `orm:"column(delta_print)" json:"delta_print"`
	Status            int8      `orm:"column(status)" json:"status"`
	CreatedAt         time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy         *Staff    `orm:"column(created_by);null;rel(fk)" json:"created_by"`

	ConsolidatedShipmentSignatures []*ConsolidatedShipmentSignature `orm:"reverse(many)" json:"consolidated_shipment_signatures,omitempty"`
	WarehouseName                  string                           `orm:"-" json:"warehouse_name"`
	PurchaseOrders                 []*PurchaseOrder                 `orm:"reverse(many)" json:"purchase_orders,omitempty"`
	SkuSummaries                   []*skuSummary                    `orm:"-" json:"sku_summaries"`
}

type skuSummary struct {
	ProductName    string           `orm:"-" json:"product_name"`
	UomName        string           `orm:"-" json:"uom_name"`
	TotalQty       float64          `orm:"-" json:"total_qty"`
	PurchaseOrders []*purchaseOrder `orm:"-" json:"purchase_orders"`
}

type purchaseOrder struct {
	PurchaseOrderCode string  `orm:"-" json:"purchase_order_code"`
	Qty               float64 `orm:"-" json:"qty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *ConsolidatedShipment) MarshalJSON() ([]byte, error) {
	type Alias ConsolidatedShipment

	return json.Marshal(&struct {
		ID            string `json:"id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		StatusConvert: util.ConvertStatusDoc(m.Status),
		Alias:         (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *ConsolidatedShipment) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *ConsolidatedShipment) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
