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
	orm.RegisterModel(new(PurchasePlan))
}

// PurchasePlan: struct to hold model data for database
type PurchasePlan struct {
	ID                   int64                 `orm:"column(id);auto" json:"-"`
	Code                 string                `orm:"column(code);size(50);null" json:"code"`
	SupplierOrganization *SupplierOrganization `orm:"column(supplier_organization_id);null;rel(fk)" json:"supplier_organization,omitempty"`
	Warehouse            *Warehouse            `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse,omitempty"`
	RecognitionDate      time.Time             `orm:"column(recognition_date)" json:"recognition_date"`
	EtaDate              time.Time             `orm:"column(eta_date)" json:"eta_date"`
	EtaTime              string                `orm:"column(eta_time)" json:"eta_time"`
	TotalPrice           float64               `orm:"column(total_price)" json:"total_price"`
	TotalWeight          float64               `orm:"column(total_weight)" json:"total_weight"`
	Note                 string                `orm:"column(note)" json:"note"`
	Status               int8                  `orm:"column(status)" json:"status"`
	TotalPurchasePlanQty float64               `orm:"column(total_purchase_plan_qty)" json:"total_purchase_plan_qty"`
	TotalPurchaseQty     float64               `orm:"column(total_purchase_qty)" json:"total_purchase_qty"`

	PurchasePlanItems []*PurchasePlanItem `orm:"reverse(many)" json:"purchase_plan_items,omitempty"`

	CreatedAt  time.Time `orm:"column(created_at)" json:"created_at"`
	CreatedBy  *Staff    `orm:"column(created_by);null;rel(fk)" json:"created_by"`
	AssignedTo *Staff    `orm:"column(assigned_to);null;rel(fk)" json:"assigned_to"`
	AssignedBy *Staff    `orm:"column(assigned_by);null;rel(fk)" json:"assigned_by"`
	AssignedAt time.Time `orm:"column(assigned_at)" json:"assigned_at"`

	TotalSku            int                    `orm:"-" json:"total_sku"`
	TonnagePurchasePlan []*tonnagePurchasePlan `orm:"-" json:"tonnage"`
}

// struct for hold calculation of Total weight per UoM for Get summary of purchase plan item in field purchaser app
type tonnagePurchasePlan struct {
	UomName     string  `orm:"-" json:"uom_name"`
	TotalWeight float64 `orm:"-" json:"total_weight"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PurchasePlan) MarshalJSON() ([]byte, error) {
	type Alias PurchasePlan

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
func (m *PurchasePlan) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PurchasePlan) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
