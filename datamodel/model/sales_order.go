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
	orm.RegisterModel(new(SalesOrder))
}

// Sales Order: struct to hold model data for database
type SalesOrder struct {
	ID                       int64              `orm:"column(id);auto" json:"-"`
	Code                     string             `orm:"column(code);size(50);null" json:"code"`
	RecognitionDate          time.Time          `orm:"column(recognition_date)" json:"recognition_date"`
	DeliveryDate             time.Time          `orm:"column(delivery_date)" json:"delivery_date"`
	BillingAddress           string             `orm:"column(billing_address)" json:"billing_address"`
	ShippingAddress          string             `orm:"column(shipping_address)" json:"shipping_address"`
	DeliveryFee              float64            `orm:"column(delivery_fee)" json:"delivery_fee"`
	VouRedeemCode            string             `orm:"column(vou_redeem_code)" json:"vou_redeem_code"`
	VouDiscAmount            float64            `orm:"column(vou_disc_amount)" json:"vou_disc_amount"`
	TotalPrice               float64            `orm:"column(total_price)" json:"total_price"`
	TotalCharge              float64            `orm:"column(total_charge)" json:"total_charge"`
	TotalWeight              float64            `orm:"column(total_weight)" json:"total_weight"`
	Note                     string             `orm:"column(note)" json:"note"`
	Status                   int8               `orm:"column(status)" json:"status"`
	OrderChannel             int8               `orm:"column(order_channel)" json:"order_channel"`
	HasExtInvoice            int8               `orm:"column(has_ext_invoice)" json:"has_ext_invoice"`
	CancellationNote         string             `orm:"-" json:"cancellation_note,omitempty"`
	ReloadPacking            int8               `orm:"column(reload_packing)" json:"reload_packing"`
	HasPickingAssigned       int8               `orm:"column(has_picking_assigned)" json:"has_picking_assigned"`
	SalesGroupID             int64              `orm:"column(sales_group_id);null" json:"sales_group_id"`
	CreatedAt                time.Time          `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy                int64              `orm:"column(created_by)" json:"created_by"`
	LastUpdatedAt            time.Time          `orm:"column(last_updated_at);type(timestamp);null" json:"last_updated_at"`
	LastUpdatedBy            int64              `orm:"column(last_updated_by)" json:"last_updated_by"`
	IsLocked                 int8               `orm:"column(is_locked)" json:"is_locked"`
	LockedBy                 int64              `orm:"column(locked_by)" json:"locked_by"`
	LockedByName             string             `orm:"-" json:"locked_by_name,omitempty"`
	CancelType               int8               `orm:"column(cancel_type)" json:"cancel_type"`
	StatusPickingOrderAssign int8               `orm:"-" json:"status_picking_order_assign"`
	TotalItem                int8               `orm:"-" json:"total_item"`
	TotalSkuDiscAmount       float64            `orm:"column(total_sku_disc_amount)" json:"total_sku_disc_amount"`
	EPCampaignID             int64              `orm:"column(eden_point_campaign_id)" json:"-"`
	EPCampaign               *EdenPointCampaign `orm:"-" json:"eden_point_campaign"`
	EstimateTimeDeparture    time.Time          `orm:"column(estimate_time_departure);type(timestamp);null" json:"estimate_time_departure"`
	IntegrationCode          string             `orm:"column(integration_code)" json:"integration_code"` // integration id for talon.one

	// Eden Point
	FinishedAt        time.Time `orm:"column(finished_at);type(timestamp);null" json:"finished_at"`
	PointRedeemAmount float64   `orm:"column(point_redeem_amount);null" json:"point_redeem_amount"`
	PointRedeemID     int64     `orm:"column(point_redeem_id);null" json:"point_redeem_id"`

	Branch          *Branch           `orm:"column(branch_id);null;rel(fk)" json:"branch"`
	SalesTerm       *SalesTerm        `orm:"column(term_payment_sls_id);null;rel(fk)" json:"term_payment_sls"`
	InvoiceTerm     *InvoiceTerm      `orm:"column(term_invoice_sls_id);null;rel(fk)" json:"term_invoice_sls"`
	Salesperson     *Staff            `orm:"column(salesperson_id);null;rel(fk)" json:"salesperson"`
	SubDistrict     *SubDistrict      `orm:"column(sub_district_id);null;rel(fk)" json:"sub_district"`
	Warehouse       *Warehouse        `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse"`
	Wrt             *Wrt              `orm:"column(wrt_id);null;rel(fk)" json:"wrt"`
	Area            *Area             `orm:"column(area_id);null;rel(fk)" json:"area"`
	Voucher         *Voucher          `orm:"column(voucher_id);null;rel(fk)" json:"voucher"`
	PriceSet        *PriceSet         `orm:"column(price_set_id);null;rel(fk)" json:"price_set"`
	PaymentGroup    *PaymentGroup     `orm:"column(payment_group_sls_id);null;rel(fk)" json:"payment_group"`
	Archetype       *Archetype        `orm:"column(archetype_id);null;rel(fk)" json:"archetype"`
	OrderType       *OrderType        `orm:"column(order_type_sls_id);null;rel(fk)" json:"order_type"`
	SalesOrderItems []*SalesOrderItem `orm:"reverse(many)" json:"sales_order_items,omitempty"`
	DeliveryOrder   []*DeliveryOrder  `orm:"reverse(many)" json:"delivery_order,omitempty"`
	SalesInvoice    []*SalesInvoice   `orm:"reverse(many)" json:"sales_invoice,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *SalesOrder) MarshalJSON() ([]byte, error) {
	type Alias SalesOrder

	return json.Marshal(&struct {
		ID                   string `json:"id"`
		StatusConvert        string `json:"status_convert"`
		StatusPickingConvert string `json:"status_picking_order_assign_convert"`
		*Alias
		PointRedeemID string `json:"point_redeem_id"`
		SalesGroupID  string `json:"sales_group_id"`
	}{
		ID:                   common.Encrypt(m.ID),
		StatusConvert:        util.ConvertStatusDoc(m.Status),
		StatusPickingConvert: util.ConvertStatusPicking(m.StatusPickingOrderAssign),
		Alias:                (*Alias)(m),
		PointRedeemID:        common.Encrypt(m.PointRedeemID),
		SalesGroupID:         common.Encrypt(m.SalesGroupID),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *SalesOrder) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *SalesOrder) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
