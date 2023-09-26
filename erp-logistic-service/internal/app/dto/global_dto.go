// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import (
	"time"
)

type GlobalSalesOrder struct {
	Id            int64      `json:"id,omitempty"`
	Code          string     `json:"code,omitempty"`
	DocNumber     string     `json:"doc_number,omitempty"`
	AddressId     int64      `json:"address_id,omitempty"`
	CustomerId    int64      `json:"customer_id,omitempty"`
	SalespersonId int64      `json:"salesperson_id,omitempty"` // BELUM
	WrtId         int64      `json:"wrt_id,omitempty"`
	OrderTypeId   int64      `json:"order_type_id,omitempty"`
	Application   int32      `json:"application,omitempty"`
	Status        int32      `json:"status,omitempty"`
	OrderDate     *time.Time `json:"order_date,omitempty"`
	Total         *float64   `json:"total,omitempty"`
	CreatedDate   *time.Time `json:"created_date,omitempty"`
	ModifiedDate  *time.Time `json:"modified_date,omitempty"`
	FinishedDate  *time.Time `json:"finished_date,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`

	Address        *GlobalAddress          `json:"address,omitempty"`
	Customer       *GlobalCustomer         `json:"customer,omitempty"`
	Wrt            *GlobalWrt              `json:"wrt,omitempty"`
	OrderType      *GlobalOrderType        `json:"order_type,omitempty"`
	SalesOrderItem []*GlobalSalesOrderItem `json:"sales_order_item,omitempty"`
	DeliveryOrder  *GlobalDeliveryOrder    `json:"delivery_order,omitempty"`
}

type GlobalAddress struct {
	Id               int64      `json:"id,omitempty"`
	Code             string     `json:"code,omitempty"`
	CustomerName     string     `json:"customer_name,omitempty"`
	ArchetypeId      int64      `json:"archetype_id,omitempty"` // BELUM
	AdmDivisionId    int64      `json:"adm_division_id,omitempty"`
	SiteId           int64      `json:"site_id,omitempty"`
	SalespersonId    int64      `json:"salesperson_id,omitempty"` // BELUM
	TerritoryId      int64      `json:"territory_id,omitempty"`   // BELUM
	AddressCode      string     `json:"address_code,omitempty"`
	AddressName      string     `json:"address_name,omitempty"`
	ContactPerson    string     `json:"contact_person,omitempty"`
	City             string     `json:"city,omitempty"`
	State            string     `json:"state,omitempty"`
	ZipCode          string     `json:"zip_code,omitempty"`
	CountryCode      string     `json:"country_code,omitempty"`
	Country          string     `json:"country,omitempty"`
	Latitude         *float64   `json:"latitude,omitempty"`
	Longitude        *float64   `json:"longitude,omitempty"`
	UpsZone          string     `json:"ups_zone,omitempty"`
	ShippingMethod   string     `json:"shipping_method,omitempty"`
	TaxScheduleId    int64      `json:"tax_schedule_id,omitempty"` // BELUM
	PrintPhoneNumber int32      `json:"print_phone_number,omitempty"`
	Phone_1          string     `json:"phone_1,omitempty"`
	Phone_2          string     `json:"phone_2,omitempty"`
	Phone_3          string     `json:"phone_3,omitempty"`
	FaxNumber        string     `json:"fax_number,omitempty"`
	ShippingAddress  string     `json:"shipping_address,omitempty"`
	BcaVa            string     `json:"bca_va,omitempty"`
	OtherVa          string     `json:"other_va,omitempty"`
	Note             string     `json:"note,omitempty"`
	Status           int32      `json:"status,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`

	Site        *GlobalSite       `json:"site,omitempty"`
	AdmDivision *GlobalAdmDivsion `json:"adm_division,omitempty"`
}

type GlobalAdmDivsion struct {
	Id            int64      `json:"id,omitempty"`
	Code          string     `json:"code,omitempty"`
	RegionId      int64      `json:"region_id,omitempty"` // belum
	City          string     `json:"city,omitempty"`
	District      string     `json:"district,omitempty"`
	SubDistrictId int64      `json:"sub_district_id,omitempty"`
	PostalCode    string     `json:"postal_code,omitempty"`
	Status        int32      `json:"status,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`

	SubDistrict *GlobalSubDistrict `json:"sub_district,omitempty"`
}

type GlobalWrt struct {
	Id        int64  `json:"id,omitempty"`
	RegionId  int64  `json:"region_id,omitempty"` // belum
	Code      string `json:"code,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
}

type GlobalSubDistrict struct {
	Id          int64      `json:"id,omitempty"`
	Code        string     `json:"code,omitempty"`
	Description string     `json:"description,omitempty"`
	Status      int32      `json:"status,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type GlobalOrderType struct {
	Id          int64      `json:"id,omitempty"`
	Code        string     `json:"code,omitempty"`
	Description string     `json:"description,omitempty"`
	Status      int32      `json:"status,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type GlobalCustomer struct {
	Id                         int64      `json:"id,omitempty"`
	Code                       string     `json:"code,omitempty"`
	ReferralCode               string     `json:"referral_code,omitempty"`
	Name                       string     `json:"name,omitempty"`
	Gender                     int32      `json:"Gender,omitempty"`
	BirthDate                  *time.Time `json:"birth_date,omitempty"`
	PicName                    string     `json:"pic_name,omitempty"`
	PhoneNumber                string     `json:"phone_number,omitempty"`
	AltPhoneNumber             string     `json:"alt_phone_number,omitempty"`
	Email                      string     `json:"email,omitempty"`
	Password                   string     `json:"password,omitempty"`
	BillingAddress             string     `json:"billing_address,omitempty"`
	Note                       string     `json:"Note,omitempty"`
	ReferenceInfo              string     `json:"reference_info,omitempty"`
	TagCustomer                string     `json:"tag_customer,omitempty"`
	Status                     int32      `json:"status,omitempty"`
	Suspended                  int32      `json:"suspended,omitempty"`
	UpgradeStatus              int32      `json:"upgrade_status,omitempty"`
	CustomerGroup              int32      `json:"customer_group,omitempty"`
	TagCustomerName            string     `json:"tag_customer_name,omitempty"`
	ReferrerCode               string     `json:"referrer_code,omitempty"`
	CreatedAt                  *time.Time `json:"createdAt,omitempty"`
	CreatedBy                  int64      `json:"created_by,omitempty"`
	LastUpdatedAt              *time.Time `json:"last_updated_at,omitempty"`
	LastUpdatedBy              int64      `json:"last_updated_by,omitempty"`
	TotalPoint                 *float64   `json:"total_point,omitempty"`
	BusinessTypeCreditLimit    int32      `json:"business_type_credit_limit,omitempty"`
	EarnedPoint                *float64   `json:"earned_point,omitempty"`
	RedeemedPoint              *float64   `json:"redeemed_point,omitempty"`
	CustomCreditLimit          int32      `json:"custom_credit_limit,omitempty"`
	CreditLimitAmount          *float64   `json:"credit_limit_amount,omitempty"`
	ProfileCode                string     `json:"profile_code,omitempty"`
	RemainingCreditLimitAmount *float64   `json:"remaining_credit_limit_amount,omitempty"`
	AverageSales               *float64   `json:"average_sales,omitempty"`
	RemainingOutstanding       *float64   `json:"remaining_outstanding,omitempty"`
	OverdueDebt                *float64   `json:"OverdueDebt,omitempty"`
	KTPPhotosUrl               string     `json:"KTPPhotosUrl,omitempty"`
	MerchantPhotosUrl          string     `json:"MerchantPhotosUrl,omitempty"`
	KTPPhotosUrlArr            []string   `json:"KTPPhotosUrlArr,omitempty"`
	MerchantPhotosUrlArr       []string   `json:"MerchantPhotosUrlArr,omitempty"`
	MembershipLevelID          int64      `json:"MembershipLevelID,omitempty"`
	MembershipCheckpointID     int64      `json:"MembershipCheckpointID,omitempty"`
	MembershipRewardID         int64      `json:"MembershipRewardID,omitempty"`
	MembershipRewardAmount     *float64   `json:"MembershipRewardAmount,omitempty"`
}

type GlobalSalesOrderItem struct {
	Id            int64      `json:"id,omitempty"`
	SalesOrderId  int64      `json:"sales_order_id,omitempty"`
	ItemId        int64      `json:"item_id,omitempty"`
	OrderQty      *float64   `json:"order_qty,omitempty"`
	DefaultPrice  *float64   `json:"default_price,omitempty"`
	UnitPrice     *float64   `json:"unit_price,omitempty"`
	TaxableItem   int32      `json:"taxable_item,omitempty"`
	TaxPercentage *float64   `json:"tax_percentage,omitempty"`
	ShadowPrice   *float64   `json:"shadow_price,omitempty"`
	Subtotal      *float64   `json:"subtotal,omitempty"`
	Weight        *float64   `json:"weight,omitempty"`
	Note          string     `json:"note,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`

	SalesOrder        *GlobalSalesOrder        `json:"sales_order,omitempty"`
	Item              *GlobalItem              `json:"item,omitempty"`
	DeliveryOrderItem *GlobalDeliveryOrderItem `json:"delivery_order_item"`
}

type GlobalItem struct {
	Id                      int64      `json:"id,omitempty"`
	Code                    string     `json:"code,omitempty"`
	UomId                   int64      `json:"uom_id,omitempty"`
	ClassId                 int64      `json:"class_id,omitempty"`         // BELUM
	ItemCategoryId          int64      `json:"item_category_id,omitempty"` // BELUM
	Description             string     `json:"description,omitempty"`
	UnitWeightConversion    *float64   `json:"unit_weight_conversion,omitempty"`
	OrderMinQty             *float64   `json:"order_min_qty,omitempty"`
	OrderMaxQty             *float64   `json:"order_max_qty,omitempty"`
	ItemType                string     `json:"item_type,omitempty"`
	Packability             string     `json:"packability,omitempty"`
	Capitalize              string     `json:"capitalize,omitempty"`
	Note                    string     `json:"note,omitempty"`
	ExcludeArchetype        string     `json:"exclude_archetype,omitempty"`
	MaxDayDeliveryDate      int32      `json:"max_day_delivery_date,omitempty"`
	FragileGoods            string     `json:"fragile_goods,omitempty"`
	Taxable                 string     `json:"taxable,omitempty"`
	OrderChannelRestriction string     `json:"order_channel_restriction,omitempty"`
	Status                  int32      `json:"status,omitempty"`
	CreatedAt               *time.Time `json:"created_at,omitempty"`
	UpdatedAt               *time.Time `json:"updated_at,omitempty"`

	Uom *GlobalUom `json:"uom,omitempty"`
}

type GlobalUom struct {
	Id          int64      `json:"id,omitempty"`
	Code        string     `json:"code,omitempty"`
	Description string     `json:"description,omitempty"`
	Status      int32      `json:"status,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type GlobalSite struct {
	Id          int64      `json:"id,omitempty"`
	Code        string     `json:"code,omitempty"`
	Description string     `json:"description,omitempty"`
	Status      int32      `json:"status,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type GlobalDeliveryOrder struct {
	Id                int64                      `json:"id,omitempty"`
	DeliveryOrderItem []*GlobalDeliveryOrderItem `json:"delivery_order_item,omitempty"`
}

type GlobalDeliveryOrderItem struct {
	Id             int64                 `json:"id,omitempty"`
	DeliverQty     *float64              `json:"deliver_qty,omitempty"`
	SalesOrderItem *GlobalSalesOrderItem `json:"sales_order_item,omitempty"`
}

type GlobalCourier struct {
	Id                int64      `json:"id,omitempty"`
	RoleId            int64      `json:"role_id,omitempty"`
	UserId            int64      `json:"user_id,omitempty"`
	Code              string     `json:"code,omitempty"`
	Name              string     `json:"name,omitempty"`
	PhoneNumber       string     `json:"phone_number,omitempty"`
	VehicleProfileId  int64      `json:"vehicle_profile_id,omitempty"`
	LicensePlate      string     `json:"license_plate,omitempty"`
	EmergencyMode     int32      `json:"emergency_mode,omitempty"`
	LastEmergencyTime *time.Time `json:"last_emergency_time,omitempty"`
	Status            int32      `json:"status,omitempty"`

	VehicleProfile *GlobalVehicleProfile `json:"vehicle_profile,omitempty"`
}

type GlobalVehicleProfile struct {
	Id                  int64    `json:"id,omitempty"`
	CourierVendorId     int64    `json:"courier_vendor_id,omitempty"`
	Code                string   `json:"code,omitempty"`
	Name                string   `json:"name,omitempty"`
	MaxKoli             *float64 `json:"max_koli,omitempty"`
	MaxWeight           *float64 `json:"max_weight,omitempty"`
	MaxFragile          *float64 `json:"max_fragile,omitempty"`
	SpeedFactor         *float64 `json:"speed_factor,omitempty"`
	RoutingProfile      int32    `json:"routing_profile,omitempty"`
	Status              int32    `json:"status,omitempty"`
	Skills              string   `json:"skills,omitempty"`
	InitialCost         *float64 `json:"initial_cost,omitempty"`
	SubsequentCost      *float64 `json:"subsequent_cost,omitempty"`
	MaxAvailableVehicle int64    `json:"max_available_vehicle,omitempty"`
}

type GlobalDeliveryRunSheet struct {
	ID                int64      `json:"id,omitempty"`
	Code              string     `json:"code,omitempty"`
	DeliveryDate      *time.Time `json:"delivery_date,omitempty"`
	StartedAt         *time.Time `json:"started_at,omitempty"`
	FinishedAt        *time.Time `json:"finished_at,omitempty"`
	StartingLatitude  *float64   `json:"starting_latitude,omitempty"`
	StartingLongitude *float64   `json:"starting_longitude,omitempty"`
	FinishedLatitude  *float64   `json:"finished_latitude,omitempty"`
	FinishedLongitude *float64   `json:"finished_longitude,omitempty"`
	Status            int8       `json:"status,omitempty"`

	CourierId int64 `json:"courier_id,omitempty"`

	Courier *GlobalCourier `json:"courier,omitempty"`
}

type GlobalDeliveryRunSheetItem struct {
	Id                          int64      `json:"id,omitempty"`
	StepType                    int8       `json:"step_type,omitempty"`
	Latitude                    *float64   `json:"latitude,omitempty"`
	Longitude                   *float64   `json:"longitude,omitempty"`
	Status                      int8       `json:"status,omitempty"`
	Note                        string     `json:"note,omitempty"`
	RecipientName               string     `json:"recipient_name,omitempty"`
	MoneyReceived               *float64   `json:"money_received,omitempty"`
	DeliveryEvidenceImageURL    string     `json:"delivery_evidence_image_url,omitempty"`
	TransactionEvidenceImageURL string     `json:"transaction_evidence_image_url,omitempty"`
	ArrivalTime                 *time.Time `json:"arrival_time,omitempty"`
	UnpunctualReason            int8       `json:"unpunctual_reason,omitempty"`
	UnpunctualDetail            int8       `json:"unpunctual_detail,omitempty"`
	FarDeliveryReason           string     `json:"far_delivery_reason,omitempty"`
	CreatedAt                   *time.Time `json:"created_at,omitempty"`
	StartedAt                   *time.Time `json:"started_at,omitempty"`
	FinishedAt                  *time.Time `json:"finished_at,omitempty"`

	DeliveryRunSheetID int64 `json:"delivery_run_sheet_id,omitempty"`
	CourierID          int64 `json:"courier_id,omitempty"`
	SalesOrderID       int64 `json:"sales_order_id,omitempty"`

	SalesOrder *GlobalSalesOrder `json:"sales_order,omitempty"`
	Courier    *GlobalCourier    `json:"courier,omitempty"`
}

type GlobalDeliveryRunReturn struct {
	ID          int64      `json:"id,omitempty"`
	Code        string     `json:"code,omitempty"`
	TotalPrice  *float64   `json:"total_price,omitempty"`
	TotalCharge *float64   `json:"total_charge,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`

	DeliveryRunSheetItemId int64 `json:"delivery_run_sheet_item_id,omitempty"`

	DeliveryRunSheetItem  *GlobalDeliveryRunSheetItem    `json:"delivery_run_sheet_item,omitempty"`
	DeliveryRunReturnItem []*GlobalDeliveryRunReturnItem `json:"delivery_run_return_item,omitempty"`
}

type GlobalDeliveryRunReturnItem struct {
	ID             int64    `json:"id,omitempty"`
	ReceiveQty     *float64 `json:"receive_qty,omitempty"`
	ReturnReason   int8     `json:"return_reason,omitempty"`
	ReturnEvidence string   `json:"return_evidence,omitempty"`
	Subtotal       *float64 `json:"subtotal,omitempty"`

	DeliveryRunReturnId int64 `json:"delivery_run_return_id,omitempty"`
	DeliveryOrderItemId int64 `json:"delivery_order_item_id,omitempty"`

	DeliveryRunReturn *GlobalDeliveryRunReturn `json:"delivery_run_return,omitempty"`
	DeliveryOrderItem *GlobalDeliveryOrderItem `json:"delivery_order_item,omitempty"`
}
