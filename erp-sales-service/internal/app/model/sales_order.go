// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

func init() {
	orm.RegisterModel(new(SalesOrder))
}

func (m *SalesOrder) TableName() string {
	return "sales_order"
}

// Purchase Order: struct to hold model data for database
type SalesOrder struct {
	ID                  int64     `orm:"column(id)" json:"id"`
	AddressIDGP         string    `orm:"column(address_id_gp)" json:"address_id_gp"`
	CustomerIDGP        string    `orm:"column(customer_id_gp)" json:"customer_id_gp"`
	TermPaymentSlsIDGP  string    `orm:"column(term_payment_sls_id_gp)" json:"term_payment_sls_id_gp"`
	SubDistrictIDGP     string    `orm:"column(sub_district_id_gp)" json:"sub_district_id_gp"`
	SiteIDGP            string    `orm:"column(site_id_gp)" json:"site_id_gp"`
	WrtIDGP             string    `orm:"column(wrt_id_gp)" json:"wrt_id_gp"`
	RegionIDGP          string    `orm:"column(region_id_gp)" json:"region_id_gp"`
	PriceLevelIDGP      string    `orm:"column(price_level_gp)" json:"price_level_id_gp"`
	PaymentGroupSlsID   int32     `orm:"column(payment_group_sls_id)" json:"payment_group_sls_id"`
	ArchetypeIDGP       string    `orm:"column(archetype_id_gp)" json:"archetype_id_gp"`
	SalesOrderNumber    string    `orm:"column(sales_order_number)" json:"sales_order_number"`
	IntegrationCode     string    `orm:"column(integration_code)" json:"integration_code"`
	SalesOrderNumberGP  string    `orm:"column(sales_order_number_gp)" json:"sales_order_number_gp"`
	Status              int8      `orm:"column(status)" json:"status"`
	RecognitionDate     time.Time `orm:"column(recognition_date)" json:"recognition_date"`
	RequestsShipDate    time.Time `orm:"column(requests_ship_date)" json:"requests_ship_date"`
	BillingAddress      string    `orm:"column(billing_address)" json:"billing_address"`
	ShippingAddress     string    `orm:"column(shipping_address)" json:"shipping_address"`
	ShippingAddressNote string    `orm:"column(shipping_address_note)" json:"shipping_address_note"`
	DeliveryFee         float64   `orm:"column(delivery_fee)" json:"delivery_fee"`
	VouDiscAmount       float64   `orm:"column(vou_disc_amount)" json:"vou_disc_amount"`
	CustomerPointLogID  int64     `orm:"column(customer_point_log_id)" json:"customer_point_log_id"`
	EdenPointCampaignID int64     `orm:"column(eden_point_campaign_id)" json:"eden_point_campaign_id"`
	TotalPrice          float64   `orm:"column(total_price)" json:"total_price"`
	TotalCharge         float64   `orm:"column(total_charge)" json:"total_charge"`
	TotalWeight         float64   `orm:"column(total_weight)" json:"total_weight"`
	Note                string    `orm:"column(note)" json:"note"`
	PaymentReminder     int8      `orm:"column(payment_reminder)" json:"payment_reminder"`
	CancelType          int8      `orm:"column(cancel_type)" json:"cancel_type"`
	CreatedAt           time.Time `orm:"column(created_at)" json:"created_at"`
	CreatedBy           int64     `orm:"column(created_by)" json:"created_by"`
	ShippingMethodIDGP  string    `orm:"column(shipping_method_id_gp)" json:"shipping_method_id_gp"`
	CustomerNameGP      string    `orm:"column(customer_name_gp)" json:"customer_name_gp"`
}
