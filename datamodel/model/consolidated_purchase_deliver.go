// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(ConsolidatedPurchaseDeliver))
}

// ConsolidatedPurchaseDeliver: struct to hold consolidated_purchase_deliver model data for database
type ConsolidatedPurchaseDeliver struct {
	ID                int64     `orm:"column(id);auto" json:"-"`
	Code              string    `orm:"column(code);size(50);null" json:"code"`
	DriverName        string    `orm:"column(driver_name);size(100);null" json:"driver_name"`
	VehicleNumber     string    `orm:"column(vehicle_number);size(10);null" json:"vehicle_number"`
	DriverPhoneNumber string    `orm:"column(driver_phone_number);size(15);null" json:"driver_phone_number"`
	DeltaPrint        int8      `orm:"column(delta_print)" json:"delta_print"`
	CreatedAt         time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy         *Staff    `orm:"column(created_by);null;rel(fk)" json:"created_by"`

	ConsolidatedPurchaseDeliverSignature []*ConsolidatedPurchaseDeliverSignature `orm:"reverse(many)" json:"consolidated_purchase_deliver_signature,omitempty"`
	SupplierName                         string                                  `orm:"-" json:"supplier_name"`
	TotalProduct                         int                                     `orm:"-" json:"total_product,omitempty"`
	TotalPurchaseDeliver                 int                                     `orm:"-" json:"total_purchase_deliver,omitempty"`
	PurchaseDelivers                     []*PurchaseDeliver                      `orm:"reverse(many)" json:"purchase_delivers,omitempty"`
	Products                             []*product                              `orm:"-" json:"products,omitempty"`
}

type product struct {
	ProductName    string            `orm:"-" json:"product_name"`
	UomName        string            `orm:"-" json:"uom_name"`
	PurchaseOrders []*purchaseOrders `orm:"-" json:"purchase_orders"`
}

type purchaseOrders struct {
	PurchaseOrderCode string  `orm:"-" json:"purchase_order_code"`
	Items             []*item `orm:"-" json:"items"`
}

type item struct {
	PurchaseDeliverCode string  `orm:"-" json:"purchase_deliver_code"`
	PurchaseQty         float64 `orm:"-" json:"purchase_qty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *ConsolidatedPurchaseDeliver) MarshalJSON() ([]byte, error) {
	type Alias ConsolidatedPurchaseDeliver

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *ConsolidatedPurchaseDeliver) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *ConsolidatedPurchaseDeliver) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
