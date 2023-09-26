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
	orm.RegisterModel(new(PurchaseOrder))
}

// PurchaseOrder: struct to hold model data for database
type PurchaseOrder struct {
	ID               int64     `orm:"column(id);auto" json:"-"`
	Code             string    `orm:"column(code);size(50);null" json:"code"`
	RecognitionDate  time.Time `orm:"column(recognition_date)" json:"recognition_date"`
	EtaDate          time.Time `orm:"column(eta_date)" json:"eta_date"`
	EtaTime          string    `orm:"column(eta_time)" json:"eta_time"`
	TaxPct           float64   `orm:"column(tax_pct)" json:"tax_pct"`
	DeliveryFee      float64   `orm:"column(delivery_fee)" json:"delivery_fee"`
	TotalPrice       float64   `orm:"column(total_price)" json:"total_price"`
	TaxAmount        float64   `orm:"column(tax_amount)" json:"tax_amount"`
	TotalCharge      float64   `orm:"column(total_charge)" json:"total_charge"`
	TotalInvoice     float64   `orm:"column(total_invoice)" json:"total_invoice"`
	TotalWeight      float64   `orm:"column(total_weight)" json:"total_weight"`
	Note             string    `orm:"column(note)" json:"note"`
	Status           int8      `orm:"column(status);null" json:"status"`
	CancellationNote string    `orm:"-" json:"cancellation_note,omitempty"`
	Tax              float64   `orm:"-" json:"tax"`
	WarehouseAddress string    `orm:"column(warehouse_address)" json:"warehouse_address"`
	CreatedFrom      int8      `orm:"column(created_from)"`
	HasFinishedGr    int8      `orm:"column(has_finished_gr)" json:"has_finished_gr"`
	Locked           int8      `orm:"column(locked)" json:"locked"`
	LockedBy         int64     `orm:"column(locked_by)" json:"-"`
	LockedByObj      *Staff    `orm:"-" json:"locked_by"`
	UpdatedAt        time.Time `orm:"column(updated_at)" json:"updated_at"`
	UpdatedBy        int64     `orm:"column(updated_by)" json:"-"`
	UpdatedByObj     *Staff    `orm:"-" json:"updated_by"`
	DeltaPrint       int8      `orm:"column(delta_print)" json:"delta_print"`
	Latitude         float64   `orm:"column(latitude)" json:"latitude"`
	Longitude        float64   `orm:"column(longitude)" json:"longitude"`

	Supplier           *Supplier            `orm:"column(supplier_id);null;rel(fk)" json:"supplier"`
	Warehouse          *Warehouse           `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse"`
	PurchaseOrderItems []*PurchaseOrderItem `orm:"reverse(many)" json:"purchase_order_items,omitempty"`
	TermPaymentPur     *PurchaseTerm        `orm:"column(term_payment_pur_id);null;rel(fk)" json:"term_payment_pur"`
	GoodsReceipt       []*GoodsReceipt      `orm:"reverse(many)" json:"goods_receipt"`
	PurchaseInvoice    []*PurchaseInvoice   `orm:"reverse(many)" json:"purchase_invoice,omitempty"`
	SupplierBadge      *SupplierBadge       `orm:"column(supplier_badge_id);null;rel(fk)" json:"supplier_badge"`

	CreatedAt            time.Time             `orm:"column(created_at)" json:"created_at"`
	CreatedBy            *Staff                `orm:"column(created_by);null;rel(fk)" json:"created_by"`
	CommittedAt          time.Time             `orm:"column(committed_at)" json:"committed_at"`
	CommittedBy          *Staff                `orm:"column(committed_by);null;rel(fk)" json:"committed_by"`
	AssignedTo           *Staff                `orm:"column(assigned_to);null;rel(fk)" json:"assigned_to"`
	AssignedBy           *Staff                `orm:"column(assigned_by);null;rel(fk)" json:"assigned_by"`
	AssignedAt           time.Time             `orm:"column(assigned_at)" json:"assigned_at"`
	PurchasePlan         *PurchasePlan         `orm:"column(purchase_plan_id);null;rel(fk)" json:"purchase_plan,omitempty"`
	ConsolidatedShipment *ConsolidatedShipment `orm:"column(consolidated_shipment_id);null;rel(fk)" json:"consolidated_shipment,omitempty"`

	TotalSku               int8                      `orm:"-" json:"total_sku"`
	TonasePurchaseOrder    []*TonasePurchaseOrder    `orm:"-" json:"tonase_purchase_order"`
	PurchaseOrderSignature []*PurchaseOrderSignature `orm:"reverse(many)" json:"signature,omitempty"`
	PurchaseOrderImage     []*PurchaseOrderImage     `orm:"reverse(many)" json:"images,omitempty"`
}

// struct for hold calculation of Total weight per UoM for Get summary of purchase order item in field purchaser app
type TonasePurchaseOrder struct {
	UomName     string  `orm:"-" json:"uom_name"`
	TotalWeight float64 `orm:"-" json:"total_weight"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PurchaseOrder) MarshalJSON() ([]byte, error) {
	type Alias PurchaseOrder

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
func (m *PurchaseOrder) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PurchaseOrder) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
