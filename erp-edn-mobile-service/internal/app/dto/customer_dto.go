package dto

import "time"

type CustomerResponse struct {
	ID                      string              `json:"id"`
	Code                    string              `json:"code"`
	ReferralCode            string              `json:"referral_code"`
	Name                    string              `json:"name"`
	Gender                  int8                `json:"gender"`
	BirthDate               time.Time           `json:"birth_date"`
	PicName                 string              `json:"pic_name"`
	PhoneNumber             string              `json:"phone_number"`
	AltPhoneNumber          string              `json:"alt_phone_number"`
	Email                   string              `json:"email"`
	Password                string              `json:"-"`
	BillingAddress          string              `json:"billing_address"`
	Note                    string              `json:"note"`
	ReferenceInfo           string              `json:"reference_info"`
	TagCustomer             string              `json:"tag_customer"`
	Status                  int8                `json:"status"`
	StatusDescription       string              `json:"status_description"`
	Suspended               int8                `json:"suspended"`
	UpgradeStatus           int8                `json:"upgrade_status"`
	TagCustomerName         string              `json:"tag_customer_name"`
	ReferrerCode            string              `json:"referrer_code"`
	CreatedAt               time.Time           `json:"created_at"`
	CreatedBy               int64               `json:"created_by"`
	LastUpdatedAt           time.Time           `json:"last_updated_at"`
	LastUpdatedBy           int64               `json:"last_updated_by"`
	TotalPoint              float64             `json:"total_point"`
	CustomerTypeCreditLimit int32               `json:"customer_type_credit_limit"`
	EarnedPoint             float64             `json:"earned_point"`
	RedeemedPoint           float64             `json:"redeemed_point"`
	ProfileCode             string              `json:"profile_code"` // profile id for talon.one
	AverageSales            float64             `json:"average_sales"`
	RemainingOutstanding    float64             `json:"remaining_outstanding"`
	OverdueDebt             float64             `json:"overdue_debt"`
	KTPPhotosUrl            string              `json:"-"`
	MerchantPhotosUrl       string              `json:"-"`
	KTPPhotosUrlArr         []string            `json:"ktp_photos_url"`
	MerchantPhotosUrlArr    []string            `json:"merchant_photos_url"`
	MembershipLevelID       int64               `json:"-"`
	MembershipCheckpointID  int64               `json:"-"`
	MembershipRewardID      int64               `json:"-"`
	MembershipRewardAmount  float64             `json:"-"`
	PaymentTerm             *CustomerGPpymtrmid `json:"payment_term"`
	CustomerTypeID          string              `json:"customer_type_id"`
	CustomerTypeDesc        string              `json:"customer_type_desc"`
	SalesPerson             string              `json:"salesperson"`
	SalesPersonName         string              `json:"salesperson_name"`
	SiteCode                string              `json:"site_code"`
	SiteName                string              `json:"site_name"`

	// Additional
	Region            string `json:"region"`
	ShippingAddress   string `json:"shipping_address"`
	CustomerGroup     string `json:"customer_group"`
	PaymentGroup      string `json:"payment_group"`
	PriceLevel        string `json:"price_level"`
	CustomerClass     string `json:"customer_class"`
	CustomerClassDesc string `json:"customer_class_description"`

	// Credit Limit
	CreditLimitType      int32   `json:"credit_limit_type"`
	CreditLimitTypeDesc  string  `json:"credit_limit_type_desc"`
	CreditLimitAmount    float64 `json:"credit_limit_amount"`
	RemainingCreditLimit float64 `json:"remaining_credit_limit"`

	// Adm Division
	Province    string `json:"province"`
	City        string `json:"city"`
	District    string `json:"district"`
	SubDistrict string `json:"sub_district"`
	Zip         string `json:"zip"`

	// Additional 2
	DueDate              string  `json:"due_date"`
	TotalRemainingAmount float64 `json:"total_remaining_amount"`
}

type CustomerListResponse struct {
	Data []*CustomerResponse `json:"data"`
}

type CustomerListRequest struct {
	Limit          int32  `json:"limit"`
	Offset         int32  `json:"offset"`
	Status         int32  `json:"status"`
	Search         string `json:"search"`
	OrderBy        string `json:"order_by"`
	CustomerTypeId int64  `json:"customer_type_id"`
}

type CustomerDetailRequest struct {
	Id string `json:"id"`
}

type CustomerGP struct {
	Custnmbr                 string                `json:"custnmbr"`
	Custclas                 string                `json:"custclas"`
	Custname                 string                `json:"custname"`
	Cprcstnm                 string                `json:"cprcstnm"`
	Cntcprsn                 string                `json:"cntcprsn"`
	Stmtname                 string                `json:"stmtname"`
	Shrtname                 string                `json:"shrtname"`
	Upszone                  string                `json:"upszone"`
	Shipmthd                 string                `json:"shipmthd"`
	Taxschid                 string                `json:"taxschid"`
	AddresS1                 string                `json:"addresS1"`
	AddresS2                 string                `json:"addresS2"`
	AddresS3                 string                `json:"addresS3"`
	Country                  string                `json:"country"`
	City                     string                `json:"city"`
	State                    string                `json:"state"`
	Zip                      string                `json:"zip"`
	PhonE1                   string                `json:"phonE1"`
	PhonE2                   string                `json:"phonE2"`
	PhonE3                   string                `json:"phonE3"`
	Fax                      string                `json:"fax"`
	Prbtadcd                 string                `json:"prbtadcd"`
	Prstadcd                 string                `json:"prstadcd"`
	Staddrcd                 string                `json:"staddrcd"`
	Slprsnid                 string                `json:"slprsnid"`
	Chekbkid                 string                `json:"chekbkid"`
	Pymtrmid                 []*CustomerGPpymtrmid `json:"pymtrmid"`
	Crlmttyp                 int32                 `json:"crlmttyp"`
	CreditLimitTypeDesc      string                `json:"credit_limit_type_desc"`
	Crlmtamt                 float64               `json:"crlmtamt"`
	Curncyid                 string                `json:"curncyid"`
	Ratetpid                 string                `json:"ratetpid"`
	Custdisc                 string                `json:"custdisc"`
	Prclevel                 string                `json:"prclevel"`
	Minpytyp                 string                `json:"minpytyp"`
	MinimumPaymentTypeDesc   string                `json:"minimum_payment_type_desc"`
	Minpydlr                 string                `json:"minpydlr"`
	Minpypct                 string                `json:"minpypct"`
	Fnchatyp                 string                `json:"fnchatyp"`
	FinanceChargeAmtTypeDesc string                `json:"finance_charge_amt_type_desc"`
	Fnchpcnt                 string                `json:"fnchpcnt"`
	Finchdlr                 string                `json:"finchdlr"`
	Mxwoftyp                 string                `json:"mxwoftyp"`
	MaximumWriteoffTypeDesc  string                `json:"maximum_writeoff_type_desc"`
	Mxwrofam                 string                `json:"mxwrofam"`
	CommenT1                 string                `json:"commenT1"`
	CommenT2                 string                `json:"commenT2"`
	UserdeF1                 string                `json:"userdeF1"`
	UserdeF2                 string                `json:"userdeF2"`
	TaxexmT1                 string                `json:"taxexmT1"`
	TaxexmT2                 string                `json:"taxexmT2"`
	TaxexmT3                 string                `json:"taxexmT3"`
	Txrgnnum                 string                `json:"txrgnnum"`
	Balnctyp                 string                `json:"balnctyp"`
	BalanceTypeDesc          string                `json:"balance_type_desc"`
	Stmtcycl                 string                `json:"stmtcycl"`
	StatementCycleDesc       string                `json:"statement_cycle_desc"`
	Bankname                 string                `json:"bankname"`
	Bnkbrnch                 string                `json:"bnkbrnch"`
	Salsterr                 string                `json:"salsterr"`
	Inactive                 int32                 `json:"inactive"`
	Hold                     string                `json:"hold"`
	Crcardid                 string                `json:"crcardid"`
	Crcrdnum                 string                `json:"crcrdnum"`
	Ccrdxpdt                 string                `json:"ccrdxpdt"`

	PaymentTermGP *SalesPaymentTermGPResponse `json:"payment_term_gp"`
	CustomerType  *CustomerTypeResponse       `json:"customer_type"`
	Region        *RegionResponse             `json:"region"`
	Site          *SiteResponse               `json:"site"`
}

type CustomerGPpymtrmid struct {
	Pymtrmid              string `json:"id"`
	Code                  string `json:"code"`
	CalculateDateFromDays int64  `json:"days_value"`
}

type GetCustomerGPResponse struct {
	PageNumber   int32         `json:"pageNumber"`
	PageSize     int32         `json:"pageSize"`
	TotalPages   int32         `json:"totalPages"`
	TotalRecords int32         `json:"totalRecords"`
	Data         []*CustomerGP `json:"data"`
	Succeeded    bool          `json:"succeeded"`
	Errors       []string      `json:"errors"`
	Message      string        `json:"message"`
}

type GetCustomerGPListRequest struct {
	Limit   int32  `query:"limit"`
	Offset  int32  `query:"offset"`
	Status  int32  `query:"status"`
	Search  string `query:"search"`
	OrderBy string `query:"orderBy"`
}

type CreateCustomerGPRequest struct {
	CustNmbr                string   `json:"customer_number"`
	CustName                string   `json:"customer_name" valid:"required"`
	CntcPrsn                string   `json:"customer_pic_name"`
	AdrsCode                string   `json:"adrscode"`
	ShipMthd                string   `json:"shipmthd"`
	Phone1                  string   `json:"customer_phone_number" valid:"required"`
	Phone2                  string   `json:"customer_alt_phone_number"`
	PymtrmID                string   `json:"payment_term" valid:"required"`
	PrcLevel                string   `json:"price_level"`
	CustTypeID              string   `json:"customer_type_id" valid:"required"`
	ReferrerCode            string   `json:"referrer_code"`
	ReferralCode            string   `json:"referral_code"`
	Region                  string   `json:"region" valid:"required"`
	Archetype               string   `json:"archetype" valid:"required"`
	Email                   string   `json:"email"`
	ImageKtp                string   `json:"image_ktp"`
	ImageAddress            []string `json:"image_address"`
	CustomerNote            string   `json:"customer_note"`
	BusinessTypeCreditLimit int32    `json:"business_type_credit_limit"`

	// add info address
	AddressName    string `json:"adress_name"`
	CodeAddress    string `json:"address_code"`
	SalesPerson    string `json:"salesperson"`
	PICName        string `json:"address_pic_name"`
	PhoneNumber    string `json:"address_phone_number"`
	AltPhoneNumber string `json:"address_alt_phone_number"`
	State          string `json:"province" valid:"required"`
	City           string `json:"city" valid:"required"`
	District       string `json:"district" valid:"required"`
	SubDistrict    string `json:"sub_district" valid:"required"`
	Zip            string `json:"zip"`
	Site           string `json:"site"`
	AddressAddr1   string `json:"address_addr1" valid:"required"`
	AddressAddr2   string `json:"address_addr2"`
	AddressAddr3   string `json:"address_addr3"`
	AddressNote    string `json:"address_note"`
	CustomerClass  string `json:"customer_class"`

	// Credit Limit
	CreditLimitType     int32   `json:"credit_limit_type"`
	CreditLimitTypeDesc string  `json:"credit_limit_type_desc"`
	CreditLimitAmount   float64 `json:"credit_limit_amount"`
}