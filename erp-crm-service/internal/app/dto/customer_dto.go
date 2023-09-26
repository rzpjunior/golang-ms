package dto

import "time"

type CustomerResponseGet struct {
	ID                         int64     `json:"id"`
	Code                       string    `json:"customer_code"`
	Name                       string    `json:"name"`
	ProspectiveCustomerID      int64     `json:"prospective_customer_id"`
	MembershipLevelID          int64     `json:"membership_level_id"`
	MembershipCheckpointID     int64     `json:"membership_checkpoint_id"`
	TotalPoint                 int64     `json:"total_point"`
	ProfileCode                string    `json:"profile_code"`
	Email                      string    `json:"email"`
	ReferenceInfo              string    `json:"reference_info"`
	UpgradeStatus              int8      `json:"upgrade_status"`
	KtpPhotosUrl               string    `json:"ktp_photos_url"`
	CustomerPhotosUrl          string    `json:"customer_photos_url"`
	CustomerSelfieUrl          string    `json:"customer_selfie_url"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
	MembershipRewardID         int64     `json:"membership_reward_id"`
	MembershipRewardAmmount    float64   `json:"membership_reward_amount"`
	CorporateCustomerNumber    string    `json:"corporate_customer_number"`
	ContactPerson              string    `json:"contact_person"`
	StatementName              string    `json:"statement_name"`
	ShortName                  string    `json:"short_name"`
	Upszone                    string    `json:"upszone"`
	TaxScheduleID              string    `json:"tax_schedule_id"`
	AddresS1                   string    `json:"address1"`
	AddresS2                   string    `json:"address2"`
	AddresS3                   string    `json:"address3"`
	Country                    string    `json:"country"`
	City                       string    `json:"city"`
	State                      string    `json:"state"`
	Zip                        string    `json:"zip"`
	PhonE1                     string    `json:"phone1"`
	PhonE2                     string    `json:"phone2"`
	PhonE3                     string    `json:"phone3"`
	Fax                        string    `json:"fax"`
	PrimaryAddressCode         string    `json:"primary_address_code"`
	WarehouseAddressCode       string    `json:"warehouse_address_code"`
	StreetAddressCode          string    `json:"street_address_code"`
	SalesPersonID              string    `json:"sale_person_id"`
	CheckbookID                string    `json:"checkbook_id"`
	PaymentTermID              string    `json:"payment_term_id"`
	CreditLimitType            string    `json:"credit_limit_type"`
	CreditLimitTypeDescription string    `json:"credit_limit_type_desc"`
	CreditLimitAmount          string    `json:"credit_limit_amount"`
	CurrencyID                 string    `json:"currency_id"`
	RateTypeID                 string    `json:"rate_type_id"`
	CustomerDiscount           string    `json:"customer_discount"`
	MinimumPaymentType         string    `json:"minimum_payment_type"`
	MinimumPaymentTypeDesc     string    `json:"minimum_payment_type_desc"`
	MinimumPaymentDollarAmount string    `json:"minimum_payment_dollar_amount"`
	MinimumPaymentPercent      string    `json:"minimum_payment_percent"`
	FinanceChargeType          string    `json:"financial_charge_type"`
	FinanceChargeTypeDesc      string    `json:"finance_charge_amt_type_desc"`
	FinanceChargePercent       string    `json:"finance_charge_percent"`
	FinanceChargeDollarAmount  string    `json:"financial_charge_dollar_amount"`
	MaximumWriteoffType        string    `json:"maximum_writeoff_type"`
	MaximumWriteoffTypeDesc    string    `json:"maximum_writeoff_type_desc"`
	MaximumWriteoffAmount      string    `json:"maximum_writeoff_amount"`
	Comment1                   string    `json:"comment1"`
	Comment2                   string    `json:"comment2"`
	UserDefined1               string    `json:"user_defined1"`
	UserDefined2               string    `json:"user_defined2"`
	TaxExempt1                 string    `json:"tax_exempt1"`
	TaxExempt2                 string    `json:"tax_exempt2"`
	TaxRegistrationNumber      string    `json:"tax_registration_number"`
	BalanceType                string    `json:"balance_type"`
	BalanceTypeDesc            string    `json:"balance_type_desc"`
	StatementCycle             string    `json:"statement_cycle"`
	StatementCycleDesc         string    `json:"statement_cycle_desc"`
	BankName                   string    `json:"bank_name"`
	BankBranch                 string    `json:"bank_branch"`
	Inactive                   int32     `json:"inactive"`
	Hold                       string    `json:"hold"`
	CreditCardID               string    `json:"credit_card_id"`
	CreditCardNumber           string    `json:"credit_card_number"`
	CreditCardExpDate          string    `json:"credit_card_example"`
	ReferrerID                 int64     `json:"referrer_id"`
	ReferrerCode               string    `json:"referrer_code"`
	ReferralCode               string    `json:"referral_code"`

	CompanyAddress *AddressCustomerResponse `json:"company_address"`
	ShipToAddress  *AddressCustomerResponse `json:"ship_to_address"`
	BillToAddress  *AddressCustomerResponse `json:"bill_to_address"`
	Archetype      *ArchetypeResponse       `json:"archetype"`
	CustomerType   *CustomerTypeResponse    `json:"customer_type"`
	CustomerClass  *CustomerClassResponse   `json:"customer_class"`
	Salesperson    *SalespersonResponse     `json:"salesperson"`
	SalesTerritory *SalesTerritoryResponse  `json:"sales_territory"`
	ShippingMethod *ShippingMethodResponse  `json:"shipping_method"`
	Site           *SiteResponse            `json:"site"`
	PriceLevel     *PriceLevelResponse      `json:"price_level"`
	PaymentTerm    *PaymentTermResponse     `json:"payment_term"`
	BusinessType   *GlossaryResponse        `json:"business_type"`
}

type CustomerRequestGetDetail struct {
	ID           int64  `json:"id"`
	CustomerIDGP string `json:"customer_id_gp"`
	Email        string `json:"email"`
	ReferrerCode string `json:"referrer_code"`
}

type CustomerRequestUpdate struct {
	ID                     int64    `json:"id"`
	CustomerIDGP           string   `json:"customer_id_gp"`
	ProspectiveCustomerID  int64    `json:"prospective_customer_id"`
	MembershipLevelID      int64    `json:"membership_level_id"`
	MembershipCheckpointID int64    `json:"membership_checkpoint_id"`
	TotalPoint             int64    `json:"total_point"`
	ProfileCode            string   `json:"profile_code"`
	ReferenceInfo          string   `json:"reference_info"`
	UpgradeStatus          int8     `json:"upgrade_status"`
	FieldUpdate            []string `json:"field_update"`
}

type CustomerGetListRequest struct {
	Limit        int    `json:"limit"`
	Offset       int    `json:"offset"`
	Search       string `json:"search"`
	Status       int8   `json:"status"`
	CustomerType string `json:"customer_type"`
}

type AddressCustomerResponse struct {
	ID            string  `json:"id"`
	AdmDivisionID string  `json:"adm_division_id"`
	AddressName   string  `json:"address_name"`
	AddressType   string  `json:"address_type"`
	Address1      string  `json:"address_1"`
	Address2      string  `json:"address_2"`
	Address3      string  `json:"address_3"`
	Region        string  `json:"region"`
	Province      string  `json:"province"`
	City          string  `json:"city"`
	District      string  `json:"district"`
	SubDistrict   string  `json:"sub_district"`
	PostalCode    string  `json:"postal_code"`
	Note          string  `json:"note"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
}
