package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Customer struct {
	ID                         int64     `orm:"column(id)" json:"id"`
	Code                       string    `orm:"column(code)" json:"code"`
	ReferralCode               string    `orm:"column(referral_code)" json:"referral_code"`
	Name                       string    `orm:"column(name)" json:"name"`
	Gender                     int8      `orm:"column(gender)" json:"gender"`
	BirthDate                  time.Time `orm:"column(birth_date)" json:"birth_date"`
	PicName                    string    `orm:"column(pic_name)" json:"pic_name"`
	PhoneNumber                string    `orm:"column(phone_number)" json:"phone_number"`
	AltPhoneNumber             string    `orm:"column(alt_phone_number)" json:"alt_phone_number"`
	Email                      string    `orm:"column(email)" json:"email"`
	Password                   string    `orm:"column(password)" json:"-"`
	BillingAddress             string    `orm:"column(billing_address)" json:"billing_address"`
	Note                       string    `orm:"column(note)" json:"note"`
	ReferenceInfo              string    `orm:"column(reference_info)" json:"reference_info"`
	TagCustomer                string    `orm:"column(tag_customer)" json:"tag_customer"`
	Status                     int8      `orm:"column(status)" json:"status"`
	Suspended                  int8      `orm:"column(suspended)" json:"suspended"`
	UpgradeStatus              int8      `orm:"column(upgrade_status)" json:"upgrade_status"`
	CustomerGroup              int8      `orm:"column(customer_group)" json:"customer_group"`
	TagCustomerName            string    `orm:"-" json:"tag_customer_name"`
	ReferrerCode               string    `orm:"column(referrer_code)" json:"referrer_code"`
	CreatedAt                  time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy                  int64     `orm:"column(created_by)" json:"created_by"`
	LastUpdatedAt              time.Time `orm:"column(last_updated_at);type(timestamp);null" json:"last_updated_at"`
	LastUpdatedBy              int64     `orm:"column(last_updated_by)" json:"last_updated_by"`
	TotalPoint                 float64   `orm:"column(total_point);digits(10);decimals(2)" json:"total_point"`
	CustomerTypeCreditLimit    int8      `orm:"column(customer_type_credit_limit)" json:"customer_type_credit_limit"`
	EarnedPoint                float64   `orm:"-" json:"earned_point"`
	RedeemedPoint              float64   `orm:"-" json:"redeemed_point"`
	CustomCreditLimit          int8      `orm:"column(custom_credit_limit)" json:"custom_credit_limit"`
	CreditLimitAmount          float64   `orm:"column(credit_limit_amount)" json:"credit_limit_amount"`
	ProfileCode                string    `orm:"column(profile_code)" json:"profile_code"` // profile id for talon.one
	RemainingCreditLimitAmount float64   `orm:"column(credit_limit_remaining)" json:"remaining_credit_limit_amount"`
	AverageSales               float64   `orm:"-" json:"average_sales"`
	RemainingOutstanding       float64   `orm:"-" json:"remaining_outstanding"`
	OverdueDebt                float64   `orm:"-" json:"overdue_debt"`
	KTPPhotosUrl               string    `orm:"column(ktp_photos_url)" json:"-"`
	CustomerPhotosUrl          string    `orm:"column(customer_photos_url)" json:"-"`
	KTPPhotosUrlArr            []string  `orm:"-" json:"ktp_photos_url"`
	CustomerPhotosUrlArr       []string  `orm:"-" json:"customer_photos_url"`
	MembershipLevelID          int64     `orm:"column(membership_level_id)" json:"-"`
	MembershipCheckpointID     int64     `orm:"column(membership_checkpoint_id)" json:"-"`
	MembershipRewardID         int64     `orm:"column(membership_reward_id)" json:"-"`
	MembershipRewardAmount     float64   `orm:"column(membership_reward_amount)" json:"-"`
	AddressID                  int64     `orm:"column(address_ID)" json:"-"`
	AdmDivisionID              int64     `orm:"column(adm_division_ID)" json:"-"`
	CustomerTypeID             int64     `orm:"column(customer_type_ID)" json:"-"`
}

func init() {
	orm.RegisterModel(new(Customer))
}

func (m *Customer) TableName() string {
	return "customer"
}
