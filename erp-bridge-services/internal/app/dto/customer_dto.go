package dto

import "time"

type CustomerResponse struct {
	ID                         int64     `json:"id"`
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
	MerchantPhotosUrl          string    `orm:"column(merchant_photos_url)" json:"-"`
	KTPPhotosUrlArr            []string  `orm:"-" json:"ktp_photos_url"`
	MerchantPhotosUrlArr       []string  `orm:"-" json:"merchant_photos_url"`
	MembershipLevelID          int64     `orm:"column(membership_level_id)" json:"-"`
	MembershipCheckpointID     int64     `orm:"column(membership_checkpoint_id)" json:"-"`
	MembershipRewardID         int64     `orm:"column(membership_reward_id)" json:"-"`
	MembershipRewardAmount     float64   `orm:"column(membership_reward_amount)" json:"-"`
	SalesPaymentTermID         int64     `json:"payment_term_sls_id"`
	BirthDateString            string    `orm:"-" json:"birth_date_string"`
	CustomerTypeId             string    `orm:"-" json:"customer_type_id"`
}

type CreateCustomerGPRequest struct {
	InterID         string `json:"interid,omitempty"`
	CustNmbr        string `json:"custnmbr"`
	CustName        string `json:"custname"`
	CustClas        string `json:"custclas"`
	CustPriority    string `json:"custpriority"`
	CprCstNm        string `json:"cprcstnm"`
	StmtName        string `json:"stmtname"`
	ShrtName        string `json:"shrtname"`
	UPSZone         string `json:"upszone"`
	ShipMthd        string `json:"shipmthd"`
	TaxSchID        string `json:"taxschid"`
	PrbtAdcd        string `json:"prbtadcd"`
	PrstAdcd        string `json:"prstadcd"`
	StAdrcd         string `json:"staddrcd"`
	SlprsnID        string `json:"slprsnid"`
	PymtrmID        string `json:"pymtrmid"`
	Salsterr        string `json:"salsterr"`
	UserDef1        string `json:"userdef1"`
	UserDef2        string `json:"userdef2"`
	DeclID          string `json:"declid"`
	Comment1        string `json:"comment1"`
	Comment2        string `json:"comment2"`
	CustDisc        int32  `json:"custdisc"`
	DisGrper        int32  `json:"disgrper"`
	DueGrper        int32  `json:"duegrper"`
	PrcLevel        string `json:"prclevel"`
	GnlCustTypeID   string `json:"gnl_cust_type_id"`
	GnlReferrerCode string `json:"gnl_referrer_code"`
	GnlReferralCode string `json:"gnl_referral_code"`
	GnlBusinessType int32  `json:"gnl_business_type"`
	GnlSocialSecNum string `json:"gnl_social_sec_num"`
	Hold            int32  `json:"hold"`
	Inactive        int32  `json:"inactive"`
	ShipComplete    int32  `json:"shipcomplete"`

	// Credit Limit
	CreditLimitAmount float64 `json:"crlmtamt"`
	CreditLimitType   int32   `json:"crlmttyp"`
	CreditLimitDesc   string  `json:"credit_limit_desc"`

	Address *CreateOrUpdateAddressGpRequest `json:"address"`
}

type CreateOrUpdateAddressGpRequest struct {
	AdrsCode string `json:"adrscode,omitempty"`
	CntcPrsn string `json:"cntcprsn,omitempty"`
	Address1 string `json:"address1,omitempty"`
	Address2 string `json:"address2,omitempty"`
	Address3 string `json:"address3,omitempty"`
	City     string `json:"city,omitempty"`
	State    string `json:"state,omitempty"`
	Zip      string `json:"zip,omitempty"`
	CCode    string `json:"ccode,omitempty"`
	Country  string `json:"country,omitempty"`
	Phone1   string `json:"phone1,omitempty"`
	Phone2   string `json:"phone2,omitempty"`
	Phone3   string `json:"phone3,omitempty"`
	Fax      string `json:"fax,omitempty"`
	UserDef1 string `json:"userdef1,omitempty"`
	UserDef2 string `json:"userdef2,omitempty"`
}

type UpdateCustomerGPRequest struct {
	InterID         string `json:"interid,omitempty"`
	CustNmbr        string `json:"custnmbr,omitempty"`
	CustName        string `json:"custname,omitempty"`
	AdrsCode        string `json:"adrscode,omitempty"`
	CustClas        string `json:"custclas,omitempty"`
	CustPriority    string `json:"custpriority,omitempty"`
	CprCstNm        string `json:"cprcstnm,omitempty"`
	StmtName        string `json:"stmtname,omitempty"`
	ShrtName        string `json:"shrtname,omitempty"`
	UPSZone         string `json:"upszone,omitempty"`
	ShipMthd        string `json:"shipmthd,omitempty"`
	TaxSchID        string `json:"taxschid,omitempty"`
	PrbtAdcd        string `json:"prbtadcd,omitempty"`
	PrstAdcd        string `json:"prstadcd,omitempty"`
	StAdrcd         string `json:"staddrcd,omitempty"`
	SlprsnID        string `json:"slprsnid,omitempty"`
	PymtrmID        string `json:"pymtrmid,omitempty"`
	Salsterr        string `json:"salsterr,omitempty"`
	UserDef1        string `json:"userdef1,omitempty"`
	UserDef2        string `json:"userdef2,omitempty"`
	DeclID          string `json:"declid,omitempty"`
	Comment1        string `json:"comment1,omitempty"`
	Comment2        string `json:"comment2,omitempty"`
	CustDisc        int32  `json:"custdisc,omitempty"`
	DisGrper        int32  `json:"disgrper,omitempty"`
	DueGrper        int32  `json:"duegrper,omitempty"`
	PrcLevel        string `json:"prclevel,omitempty"`
	GnlCustTypeID   string `json:"gnl_cust_type_id,omitempty"`
	GnlReferrerCode string `json:"gnl_referrer_code,omitempty"`
	GnlReferralCode string `json:"gnl_referral_code,omitempty"`
	GnlBusinessType int32  `json:"gnl_business_type,omitempty"`
	GnlSocialSecNum string `json:"gnl_social_sec_num,omitempty"`
	Hold            int32  `json:"hold,omitempty"`
	Inactive        int32  `json:"inactive,omitempty"`
	ShipComplete    string `json:"shipcomplete,omitempty"`

	Address *CreateOrUpdateAddressGpRequest `json:"address"`
}

type CreateCustomerGPResponse struct {
	Code            int    `json:"code"`
	CustNmbr        string `json:"custnmbr"`
	GnlReferralCode string `json:"gnL_referral_code"`
	GnlReferrerCode string `json:"gnl_referrer_code"`
	ShipComplete    int    `json:"shipcomplete"`
	Message         string `json:"message"`
}

type UpdateFixedVaRequest struct {
	InterID  string `json:"interid,omitempty"`
	CustNmbr string `json:"custnmbr,omitempty"`
	UserDef1 string `json:"userdef1,omitempty"`
	UserDef2 string `json:"userdef2,omitempty"`
}
