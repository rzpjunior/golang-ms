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
	orm.RegisterModel(new(Merchant))
}

// Merchant : struct to hold model data for database
type Merchant struct {
	ID                         int64     `orm:"column(id);auto" json:"-"`
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
	BusinessTypeCreditLimit    int8      `orm:"column(business_type_credit_limit)" json:"business_type_credit_limit"`
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
	MerchantPhotosUrl          string    `orm:"column(merchant_photos_url)" json:"-"`
	KTPPhotosUrlArr            []string  `orm:"-" json:"ktp_photos_url"`
	MerchantPhotosUrlArr       []string  `orm:"-" json:"merchant_photos_url"`
	MembershipLevelID          int64     `orm:"column(membership_level_id)" json:"-"`
	MembershipCheckpointID     int64     `orm:"column(membership_checkpoint_id)" json:"-"`
	MembershipRewardID         int64     `orm:"column(membership_reward_id)" json:"-"`
	MembershipRewardAmount     float64   `orm:"column(membership_reward_amount)" json:"-"`

	CreditLimit          *CreditLimit          `orm:"-" json:"credit_limit,omitempty"`
	UserMerchant         *UserMerchant         `orm:"column(user_merchant_id);null;rel(fk)" json:"user_merchant,omitempty"`
	InvoiceTerm          *InvoiceTerm          `orm:"column(term_invoice_sls_id);null;rel(fk)" json:"invoice_term,omitempty"`
	PaymentTerm          *SalesTerm            `orm:"column(term_payment_sls_id);null;rel(fk)" json:"payment_term,omitempty"`
	PaymentMethod        *PaymentMethod        `orm:"column(payment_method_id);null;rel(fk)" json:"payment_method,omitempty"`
	BusinessType         *BusinessType         `orm:"column(business_type_id);null;rel(fk)" json:"business_type,omitempty"`
	FinanceArea          *Area                 `orm:"column(finance_area_id);null;rel(fk)" json:"finance_area,omitempty"`
	PaymentGroup         *PaymentGroup         `orm:"column(payment_group_sls_id);null;rel(fk)" json:"payment_group,omitempty"`
	ProspectCustomer     *ProspectCustomer     `orm:"column(prospect_customer_id);null;rel(fk)" json:"prospect_customer,omitempty"`
	Referrer             *Merchant             `orm:"column(referrer_id);null;rel(fk)" json:"referrer,omitempty"`
	MerchantPriceSet     []*MerchantPriceSet   `orm:"reverse(many)" json:"price_set_area"`
	MerchantAccNum       []*MerchantAccNum     `orm:"reverse(many)" json:"merchant_acc_num"`
	MembershipLevel      *MembershipLevel      `orm:"-" json:"membership_level,omitempty"`
	MembershipCheckpoint *MembershipCheckpoint `orm:"-" json:"membership_checkpoint,omitempty"`
	MembershipReward     *MembershipReward     `orm:"-" json:"membership_reward,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *Merchant) MarshalJSON() ([]byte, error) {
	type Alias Merchant

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *Merchant) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *Merchant) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}

type MerchantOrderPerformance struct {
	ProductId    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	QtySell      float64 `json:"qty_sell"`
	AverageSales float64 `json:"avg_sales"`
	OrderTotal   int64   `json:"order_total"`
}

type MerchantPaymentPerformance struct {
	CreditLimitAmount                   float64 `json:"credit_limit_amount"`
	CreditLimitRemaining                float64 `json:"credit_limit_remaining"`
	RemainingOutstanding                float64 `json:"remaining_outstanding"`
	CreditLimitUsageRemainingPercentage float64 `json:"credit_limit_usage_remaining_percentage"`
	OverdueDebtAmount                   float64 `json:"overdue_debt_amount"`
	OverdueDebtRemainingPercentage      float64 `json:"overdue_debt_remaining_percentage"`
	AveragePaymentAmount                float64 `json:"average_payment_amount"`
	AveragePaymentPercentage            float64 `json:"average_payment_percentage"`
	AveragePaymentPeriod                int     `json:"average_payment_period"`
}
