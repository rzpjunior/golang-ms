package dto

import (
	"time"
)

type SaveRegistrationRequest struct {
	CodeUserCustomer string
	CodeCustomer     string
	CodeAddress      string
	CodeReferral     string

	CustomerName           string `json:"customer_name" valid:"required"`
	CustomerPhoneNumber    string `json:"customer_phone_number" valid:"required"`
	CustomerAltPhoneNumber string `json:"customer_alt_phone_number"`
	CustomerEmail          string `json:"customer_email" valid:"required"`
	CustomerBirthDate      string `json:"customer_birth_date"`
	CustomerGender         int    `json:"customer_gender"`
	OTP                    string `json:"otp"`
	Platform               string `json:"-"`
	AppVersion             string `json:"-"`

	// ShippingAddress string `json:"shipping_address" valid:"required"`
	Address1      string `json:"address_1" valid:"required"`
	Address2      string `json:"address_2" `
	Address3      string `json:"address_3" `
	AddressNote   string `json:"address_note" `
	SubDistrictID string `json:"sub_district_id" valid:"required"`
	ReferenceInfo string `json:"reference_info" valid:"required"`
	ReferrerCode  string `json:"referrer_code"`
	ReferrerID    int64  `json:"-"`
	FirebaseToken string `json:"fcm_token"`

	BirthDateAt time.Time `json:"-"`

	Latitude          *float64 `json:"latitude"`
	Longitude         *float64 `json:"longitude"`
	IsValidCoordinate int8     `json:"is_valid_coordinate"`

	// Customer             *model.Customer
	// Branch               *model.Branch
	// CustomerBusinessType *model.BusinessType
	// InvoiceTerm          *model.InvoiceTerm
	// PaymentTerm          *model.SalesTerm
	// PaymentMethod        *model.PaymentMethod
	// BusinessType         *model.BusinessType
	// PaymentGroup         *model.PaymentGroup

	// BranchPriceSet    *model.PriceSet
	// SubDistrict       *model.SubDistrict
	// WarehouseCoverage *model.WarehouseCoverage
	LoginToken string `json:"login_token"`

	// PriceSet     *model.PriceSet
	// PriceSetArea []*model.AreaPolicy
	// SalesPerson  *model.Staff
}
