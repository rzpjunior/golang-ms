package dto

import (
	"time"

	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"
)

type LoginRequest struct {
	Platform string       `json:"platform" valid:"required"`
	Data     loginRequest `json:"data"`
	Timezone string
}
type loginRequest struct {
	PhoneNumber string `json:"phone_number" valid:"required"`
	OTP         string `json:"otp" valid:"required"`
	FcmToken    string `json:"fcm_token"`
	//Customer     *model.Customer     `json:"-"`
	//UserCustomer *model.UserCustomer `json:"-"`
	IPAddress string `json:"-"`
}
type CheckSessionRequest struct {
	Platform string `json:"platform" valid:"required"`
}

// type SessionDataRequest struct {
// 	// Customer *pb.Customer `json:"customer"`
// 	// Address  *pb.Address  `json:"address"`
// 	Customer *SessionCustomer `json:"customer"`
// 	Address  *SessionAddress  `json:"address"`
// }

type SessionCustomer struct {
	ID                         string    `json:"id,omitempty"`
	Code                       string    `json:"code,omitempty"`
	ReferralCode               string    `json:"referral_code,omitempty"`
	Name                       string    `json:"name,omitempty"`
	Gender                     string    `json:"Gender,omitempty"`
	BirthDate                  time.Time `json:"birth_date,omitempty"`
	PicName                    string    `json:"pic_name,omitempty"`
	PhoneNumber                string    `json:"phone_number,omitempty"`
	AltPhoneNumber             string    `json:"alt_phone_number,omitempty"`
	Email                      string    `json:"email,omitempty"`
	Password                   string    `json:"password,omitempty"`
	BillingAddress             string    `json:"billing_address,omitempty"`
	Note                       string    `json:"Note,omitempty"`
	ReferenceInfo              string    `json:"reference_info,omitempty"`
	TagCustomer                string    `json:"tag_customer,omitempty"`
	Status                     string    `json:"status,omitempty"`
	Suspended                  string    `json:"suspended,omitempty"`
	UpgradeStatus              string    `json:"upgrade_status,omitempty"`
	CustomerGroup              string    `json:"customer_group,omitempty"`
	TagCustomerName            string    `json:"tag_customer_name,omitempty"`
	ReferrerCode               string    `json:"referrer_code,omitempty"`
	CreatedAt                  time.Time `json:"createdAt,omitempty"`
	CreatedBy                  string    `json:"created_by,omitempty"`
	LastUpdatedAt              time.Time `json:"last_updated_at,omitempty"`
	LastUpdatedBy              string    `json:"last_updated_by,omitempty"`
	TotalPoint                 string    `json:"total_point,omitempty"`
	CustomerTypeCreditLimit    string    `json:"customer_type_credit_limit,omitempty"`
	EarnedPoint                string    `json:"earned_point,omitempty"`
	RedeemedPoint              string    `json:"redeemed_point,omitempty"`
	CustomCreditLimit          string    `json:"custom_credit_limit,omitempty"`
	CreditLimitAmount          string    `json:"credit_limit_amount,omitempty"`
	ProfileCode                string    `json:"profile_code,omitempty"`
	RemainingCreditLimitAmount string    `json:"remaining_credit_limit_amount,omitempty"`
	AverageSales               string    `json:"average_sales,omitempty"`
	RemainingOutstanding       string    `json:"remaining_outstanding,omitempty"`
	OverdueDebt                string    `json:"overdueDebt,omitempty"`
	KTPPhotosUrl               string    `json:"ktp_photosUrl,omitempty"`
	CustomerPhotosUrl          string    `json:"customerPhotosUrl,omitempty"`
	KTPPhotosUrlArr            []string  `json:"ktp_photosUrlArr,omitempty"`
	CustomerPhotosUrlArr       []string  `json:"customerPhotosUrlArr,omitempty"`
	MembershipLevelID          string    `json:"membershipLevelID,omitempty"`
	MembershipCheckpointID     string    `json:"membershipCheckpointID,omitempty"`
	MembershipRewardID         string    `json:"membershipRewardID,omitempty"`
	MembershipRewardAmount     string    `json:"membershipRewardAmount,omitempty"`
	TermPaymentSlsId           string    `json:"term_payment_sls_id,omitempty"`
	BirthDateString            string    `json:"birth_date_string,omitempty"`

	UserCustomer   *model.UserCustomer  `json:"user_customer,omitempty"`
	InvoiceTerm    *model.InvoiceTerm   `json:"invoice_term,omitempty"`
	TermPaymentSls string               `json:"term_payment_sls,omitempty"`
	PaymentMethod  *model.PaymentMethod `json:"payment_method,omitempty"`
	//CustomerType         *model.CustomerType         `json:"customer_type,omitempty"`
	FinanceArea          *model.Region                 `json:"finance_area,omitempty"`
	PaymentGroup         *model.PaymentGroup           `json:"payment_group,omitempty"`
	ProspectCustomer     *model.ProspectCustomer       `json:"prospect_customer,omitempty"`
	CustomerPriceSet     []*model.CustomerPriceSet     `json:"price_set_area"`
	CustomerAccNum       []*model.CustomerAccNum       `json:"customer_acc_num"`
	CustomerType         string                        `json:"customer_type"`
	ReferrerCustomer     *model.Customer               `json:"referrer_Customer,omitempty"`
	MembershipLevel      *MembershipLevelResponse      `json:"membership_level,omitempty"`
	MembershipCheckpoint *MembershipCheckpointResponse `json:"membership_checkpoint,omitempty"`
	MembershipReward     *model.MembershipReward       `json:"membership_reward,omitempty"`
}
type SessionAddress struct {
	ID                 string    `orm:"column(id);auto" json:"id,omitempty"`
	Code               string    `orm:"column(code)" json:"code,omitempty"`
	Name               string    `orm:"column(name)" json:"name,omitempty"`
	PicName            string    `orm:"column(pic_name)" json:"pic_name,omitempty"`
	PhoneNumber        string    `orm:"column(phone_number)" json:"phone_number,omitempty"`
	AltPhoneNumber     string    `orm:"column(alt_phone_number)" json:"alt_phone_number,omitempty"`
	AddressName        string    `orm:"column(address_name)" json:"address_name,omitempty"`
	ShippingAddress    string    `orm:"column(shipping_address)" json:"shipping_address,omitempty"`
	Latitude           string    `orm:"column(latitude);null" json:"latitude"`
	Longitude          string    `orm:"column(longitude);null" json:"longitude"`
	Note               string    `orm:"column(note)" json:"note,omitempty"`
	City               string    `orm:"column(city)" json:"city,omitempty"`
	MainBranch         string    `orm:"column(main_branch)" json:"main_branch,omitempty"`
	Status             string    `orm:"column(status)" json:"status"`
	CreatedAt          time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy          string    `orm:"column(created_by)" json:"created_by"`
	LastUpdatedAt      time.Time `orm:"column(last_updated_at);type(timestamp);null" json:"last_updated_at"`
	LastUpdatedBy      string    `orm:"column(last_updated_by)" json:"last_updated_by"`
	PinpointValidation string    `orm:"column(pinpoint_validation)" json:"pinpoint_validation"`
	AdmDivisionId      string    `orm:"column(adm_division_id)" json:"adm_division_id"`
	RegionID           string    `orm:"column(region_id)" json:"region_id"`
	CustomerID         string    `orm:"column(customer_id)" json:"customer_id"`
	ArchetypeID        string    `orm:"column(archetype_id)" json:"archetype_id"`
	PriceSetID         string    `orm:"column(priceset_id)" json:"priceset_id"`
	SiteID             string    `orm:"column(site_id)" json:"site_id"`
	SalesPersonID      string    `orm:"column(sales_person_id)" json:"sales_person_id"`
	SubDistrictID      string    `orm:"column(sub_district_id)" json:"sub_district_id"`

	Customer    *model.Customer    `orm:"column(customer_id);null;rel(fk)" json:"customer,omitempty"`
	Region      *model.Region      `orm:"column(region_id);null;rel(fk)" json:"region,omitempty"`
	Archetype   *model.Archetype   `orm:"column(archetype_id);null;rel(fk)" json:"archetype,omitempty"`
	PriceSet    *model.PriceSet    `orm:"column(price_set_id);null;rel(fk)" json:"price_set,omitempty"`
	Site        *model.Site        `orm:"column(site_id);null;rel(fk)" json:"site,omitempty"`
	Salesperson *model.Staff       `orm:"column(salesperson_id);null;rel(fk)" json:"salesperson,omitempty"`
	AdmDivision *model.AdmDivision `orm:"column(sub_district_id);null;rel(fk)" json:"adm_division,omitempty"`

	StatusConvert string `orm:"-" json:"status_convert"`
}
type LoginResponse struct {
	Token string `json:"token,omitempty"`
}

type RequestGetOTP struct {
	Platform string     `json:"platform" valid:"required"`
	Data     otpRequest `json:"data"`
}
type otpRequest struct {
	PhoneNumber         string                `json:"phone_number" valid:"required"`
	Type                string                `json:"type" valid:"required"`
	Customer            *model.Customer       `json:"-"`
	WhiteListLogin      *model.WhiteListLogin `json:"-"`
	RegistrationRequest createRequest         `json:"registration_request"`
	CreatedAt           time.Time
	IPAddress           string `json:"-"`
	OtpType             string `json:"otp_type"`
}

type createRequest struct {
	CodeUserCustomer string
	CodeCustomer     string
	CodeBranch       string
	CodeReferral     string

	CustomerName           string `json:"customer_name"`
	CustomerPhoneNumber    string `json:"customer_phone_number"`
	CustomerAltPhoneNumber string `json:"customer_alt_phone_number" `
	CustomerEmail          string `json:"customer_email"`
	CustomerBirthDate      string `json:"customer_birth_date"`
	CustomerGender         int    `json:"customer_gender"`

	ShippingAddress string `json:"shipping_address" `
	AdmDivisionID   string `json:"adm_division_id"`
	ReferenceInfo   string `json:"reference_info"`
	ReferrerCode    string `json:"referrer_code"`

	BirthDateAt time.Time `json:"-"`

	Customer *model.Customer
	//CustomerBusinessType *model.BusinessType

	//BranchPriceSet    *model.PriceSet
	//SubDistrict       *model.SubDistrict
	//WarehouseCoverage *model.WarehouseCoverage
	//PriceSet          *model.PriceSet
	// Area              *model.Area
}

type ResponseSendOtp struct {
	OtpName string `json:"otp_name"`
}

type SignOutRequest struct {
	Session *SessionDataCustomer
}

type SessionDataCustomer struct {
	Customer *SessionCustomer `json:"customer"`
	Address  *SessionAddress  `json:"address"`
}

type CheckPhoneNumberRequest struct {
	Platform string                       `json:"platform" valid:"required"`
	Data     CheckPhoneNumberLoginRequest `json:"data"`
}

type CheckPhoneNumberLoginRequest struct {
	PhoneNumber string `json:"phone_number" valid:"required"`
}

type SendOtpWASociomile struct {
	WAID       string      `json:"wa_id"`
	TemplateID string      `json:"template_id"`
	Components []Component `json:"components"`
}
type Component struct {
	Type       string      `json:"type"`
	Parameters []Parameter `json:"parameters"`
}
type Parameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type VerifyRegistRequest struct {
	Platform string       `json:"platform" valid:"required"`
	Data     verifyRegist `json:"data"`
}
type verifyRegist struct {
	PhoneNumber string `json:"phone_number" valid:"required"`
	OTP         string `json:"otp" valid:"required"`
	IPAddress   string `json:"-"`
}

type RequestDeleteAccount struct {
	Data    dataDeleteAccount `json:"data" valid:"required"`
	Session *SessionDataCustomer
}

type dataDeleteAccount struct {
	Reason string `json:"reason" valid:"required"`
}

type RequestGetPostSession struct {
	Platform string               `json:"platform" valid:"required"`
	Session  *SessionDataCustomer `json:"-"`
}
