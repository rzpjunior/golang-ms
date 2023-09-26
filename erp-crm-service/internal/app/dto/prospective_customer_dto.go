package dto

import "time"

type ProspectiveCustomerResponse struct {
	ID                         int64                   `json:"id"`
	Code                       string                  `json:"code"`
	Archetype                  *ArchetypeResponse      `json:"archetype"`
	CustomerType               *CustomerTypeResponse   `json:"customer_type"`
	Customer                   *CustomerResponse       `json:"customer"`
	CustomerClass              *CustomerClassResponse  `json:"customer_class"`
	Salesperson                *SalespersonResponse    `json:"salesperson"`
	SalesTerritory             *SalesTerritoryResponse `json:"sales_territory"`
	ShippingMethod             *ShippingMethodResponse `json:"shipping_method"`
	Site                       *SiteResponse           `json:"site"`
	PriceLevel                 *PriceLevelResponse     `json:"price_level"`
	BusinessType               *GlossaryResponse       `json:"business_type"`
	BusinessName               string                  `json:"business_name"`
	BrandName                  string                  `json:"brand_name"`
	CompanyAddressID           int64                   `json:"company_address_id"`
	CompanyAddressName         string                  `json:"company_address_name"`
	CompanyAddressRegion       string                  `json:"company_address_region"`
	CompanyAddressDetail1      string                  `json:"company_address_detail_1"`
	CompanyAddressDetail2      string                  `json:"company_address_detail_2"`
	CompanyAddressDetail3      string                  `json:"company_address_detail_3"`
	CompanyAddressProvince     string                  `json:"company_address_province"`
	CompanyAddressCity         string                  `json:"company_address_city"`
	CompanyAddressDistrict     string                  `json:"company_address_district"`
	CompanyAddressSubDistrict  string                  `json:"company_address_sub_district"`
	CompanyAddressPostalCode   string                  `json:"company_address_postal_code"`
	CompanyAddressNote         string                  `json:"company_address_note"`
	CompanyAddressLatitude     string                  `json:"company_address_latitude"`
	CompanyAddressLongitude    string                  `json:"company_address_longitude"`
	RegStatus                  int8                    `json:"reg_status"`
	RegStatusConvert           string                  `json:"reg_status_convert"`
	Application                *GlossaryResponse       `json:"application"`
	CreatedAt                  time.Time               `json:"created_at"`
	UpdatedAt                  time.Time               `json:"updated_at"`
	ProcessedAt                time.Time               `json:"processed_at"`
	ProcessedBy                *CreatedByResponse      `json:"processed_by"`
	DeclineType                int8                    `json:"decline_type"`
	DeclineNote                string                  `json:"decline_note"`
	ShippingAddressReferTo     int8                    `json:"shipping_address_refer_to"`
	ShippingAddressID          int64                   `json:"shipping_address_id"`
	ShippingAddressName        string                  `json:"shipping_address_name"`
	ShippingAddressRegion      string                  `json:"shipping_address_region"`
	ShippingAddressDetail1     string                  `json:"shipping_address_detail_1"`
	ShippingAddressDetail2     string                  `json:"shipping_address_detail_2"`
	ShippingAddressDetail3     string                  `json:"shipping_address_detail_3"`
	ShippingAddressProvince    string                  `json:"shipping_address_province"`
	ShippingAddressCity        string                  `json:"shipping_address_city"`
	ShippingAddressDistrict    string                  `json:"shipping_address_district"`
	ShippingAddressSubDistrict string                  `json:"shipping_address_sub_district"`
	ShippingAddressPostalCode  string                  `json:"shipping_address_postal_code"`
	ShippingAddressNote        string                  `json:"shipping_address_note"`
	ShippingAddressLatitude    string                  `json:"shipping_address_latitude"`
	ShippingAddressLongitude   string                  `json:"shipping_address_longitude"`
	BillingAddressReferTo      int8                    `json:"billing_address_refer_to"`
	BillingAddressID           int64                   `json:"billing_address_id"`
	BillingAddressName         string                  `json:"billing_address_name"`
	BillingAddressRegion       string                  `json:"billing_address_region"`
	BillingAddressDetail1      string                  `json:"billing_address_detail_1"`
	BillingAddressDetail2      string                  `json:"billing_address_detail_2"`
	BillingAddressDetail3      string                  `json:"billing_address_detail_3"`
	BillingAddressProvince     string                  `json:"billing_address_province"`
	BillingAddressCity         string                  `json:"billing_address_city"`
	BillingAddressDistrict     string                  `json:"billing_address_district"`
	BillingAddressSubDistrict  string                  `json:"billing_address_sub_district"`
	BillingAddressPostalCode   string                  `json:"billing_address_postal_code"`
	BillingAddressNote         string                  `json:"billing_address_note"`
	BillingAddressLatitude     string                  `json:"billing_address_latitude"`
	BillingAddressLongitude    string                  `json:"billing_address_longitude"`
	OutletImage                []string                `json:"outlet_image"`
	TimeConsent                *GlossaryResponse       `json:"time_consent"`
	ReferenceInfo              *GlossaryResponse       `json:"reference_info"`
	ReferrerCode               string                  `json:"referrer_code"`
	OwnerName                  string                  `json:"owner_name"`
	OwnerContact               string                  `json:"owner_contact"`
	OwnerRole                  string                  `json:"owner_role"`
	Email                      string                  `json:"email"`
	PicOperationName           string                  `json:"pic_operation_name"`
	PicOperationContact        string                  `json:"pic_operation_contact"`
	PicOrderName               string                  `json:"pic_order_name"`
	PicOrderContact            string                  `json:"pic_order_contact"`
	PicFinanceName             string                  `json:"pic_finance_name"`
	PicFinanceContact          string                  `json:"pic_finance_contact"`
	IDCardDocName              string                  `json:"id_card_doc_name"`
	IDCardDocNumber            string                  `json:"id_card_doc_number"`
	IDCardDocURL               string                  `json:"id_card_doc_url"`
	TaxpayerDocName            string                  `json:"taxpayer_doc_name"`
	TaxpayerDocNumber          string                  `json:"taxpayer_doc_number"`
	TaxpayerDocURL             string                  `json:"taxpayer_doc_url"`
	CompanyContractDocName     string                  `json:"company_contract_doc_name"`
	CompanyContractDocURL      string                  `json:"company_contract_doc_url"`
	NotarialDeedDocName        string                  `json:"notarial_deed_doc_name"`
	NotarialDeedDocURL         string                  `json:"notarial_deed_doc_url"`
	TaxableEntrepeneurDocName  string                  `json:"taxable_entrepeneur_doc_name"`
	TaxableEntrepeneurDocURL   string                  `json:"taxable_entrepeneur_doc_url"`
	CompanyCertificateRegName  string                  `json:"company_certificate_reg_name"`
	CompanyCertificateRegURL   string                  `json:"company_certificate_reg_url"`
	BusinessLicenseDocName     string                  `json:"business_license_doc_name"`
	BusinessLicenseDocURL      string                  `json:"business_license_doc_url"`
	PaymentTerm                *PaymentTermResponse    `json:"payment_term"`
	ExchangeInvoice            int8                    `json:"exchange_invoice"`
	ExchangeInvoiceTime        string                  `json:"exchange_invoice_time"`
	FinanceEmail               string                  `json:"finance_email"`
	InvoiceTerm                *GlossaryResponse       `json:"invoice_term"`
	Comment1                   string                  `json:"comment_1"`
	Comment2                   string                  `json:"comment_2"`
}

type RegionResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type CustomerResponse struct {
	ID           int64  `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Email        string `json:"email,omitempty"`
	ReferrerCode string `json:"referrerCode,omitempty"`
}

type PaymentTermResponse struct {
	ID                       string `json:"id"`
	Description              string `json:"description"`
	DueType                  string `json:"due_type"`
	PaymentUseFor            int    `json:"payment_usefor"`
	PaymentUseForDescription string `json:"payment_usefor_description"`
}
type ProspectiveCustomerDecineRequest struct {
	DeclineType int8   `json:"decline_type" valid:"required"`
	DeclineNote string `json:"decline_note" valid:"lte:250"`
}

type ProspectiveCustomerCreateRequest struct {
	ProspectiveCustomerID      int64    `json:"prospective_customer_id"`
	CustomerCode               string   `json:"customer_code" valid:"required"`
	BusinessName               string   `json:"business_name" valid:"required"`
	BrandName                  string   `json:"brand_name"`
	CustomerTypeID             string   `json:"customer_type_id" valid:"required"`
	BusinessTypeID             int8     `json:"business_type_id" valid:"required"`
	ArchetypeID                string   `json:"archetype_id" valid:"required"`
	CustomerClassID            string   `json:"customer_class_id" valid:"required"`
	ReferrerCode               string   `json:"referrer_code"`
	TimeConsent                int8     `json:"time_consent"`
	ReferenceInfo              int8     `json:"reference_info"`
	RegistrationChannel        int8     `json:"registration_channel"`
	OutletImage                []string `json:"outlet_image"`
	CompanyAddressID           int64    `json:"company_address_id"`
	CompanyAddressName         string   `json:"company_address_name"`
	CompanyAddressRegion       string   `json:"company_address_region"`
	CompanyAddressDetail1      string   `json:"company_address_detail_1"`
	CompanyAddressDetail2      string   `json:"company_address_detail_2"`
	CompanyAddressDetail3      string   `json:"company_address_detail_3"`
	CompanyAddressProvince     string   `json:"company_address_province"`
	CompanyAddressCity         string   `json:"company_address_city"`
	CompanyAddressDistrict     string   `json:"company_address_district"`
	CompanyAddressSubDistrict  string   `json:"company_address_sub_district"`
	CompanyAddressPostalCode   string   `json:"company_address_postal_code"`
	CompanyAddressNote         string   `json:"company_address_note"`
	CompanyAddressLatitude     string   `json:"company_address_latitude"`
	CompanyAddressLongitude    string   `json:"company_address_longitude"`
	ShippingAddressReferTo     int8     `json:"shipping_address_refer_to"`
	ShippingAddressID          int64    `json:"shipping_address_id"`
	ShippingAddressName        string   `json:"shipping_address_name" valid:"required"`
	ShippingAddressRegion      string   `json:"shipping_address_region" valid:"required"`
	ShippingAddressDetail1     string   `json:"shipping_address_detail_1" valid:"required"`
	ShippingAddressDetail2     string   `json:"shipping_address_detail_2"`
	ShippingAddressDetail3     string   `json:"shipping_address_detail_3"`
	ShippingAddressProvince    string   `json:"shipping_address_province" valid:"required"`
	ShippingAddressCity        string   `json:"shipping_address_city" valid:"required"`
	ShippingAddressDistrict    string   `json:"shipping_address_district" valid:"required"`
	ShippingAddressSubDistrict string   `json:"shipping_address_sub_district" valid:"required"`
	ShippingAddressPostalCode  string   `json:"shipping_address_postal_code" valid:"required"`
	ShippingAddressNote        string   `json:"shipping_address_note"`
	ShippingAddressLatitude    string   `json:"shipping_address_latitude" valid:"required"`
	ShippingAddressLongitude   string   `json:"shipping_address_longitude" valid:"required"`
	SiteID                     string   `json:"site_id"`
	ShippingMethodID           string   `json:"shipping_method_id" valid:"required"`
	PicOrderName               string   `json:"pic_order_name" valid:"required"`
	PicOrderContact            string   `json:"pic_order_contact" valid:"required"`
	SalesTerritoryID           string   `json:"sales_territory_id" valid:"required"`
	SalespersonID              string   `json:"salesperson_id" valid:"required"`
	PriceLevelID               string   `json:"price_level_id" valid:"required"`
	OwnerName                  string   `json:"owner_name"`
	OwnerContact               string   `json:"owner_contact"`
	OwnerRole                  string   `json:"owner_role"`
	Email                      string   `json:"email"`
	IDCardDocNumber            string   `json:"id_card_doc_number"`
	TaxpayerDocNumber          string   `json:"taxpayer_doc_number"`
	PicOperationName           string   `json:"pic_operation_name"`
	PicOperationContact        string   `json:"pic_operation_contact"`
	IDCardDocURL               string   `json:"id_card_doc_url"`
	CompanyContractDocURL      string   `json:"company_contract_doc_url"`
	NotarialDeedDocURL         string   `json:"notarial_deed_doc_url"`
	TaxpayerDocURL             string   `json:"taxpayer_doc_url"`
	TaxableEntrepeneurDocURL   string   `json:"taxable_entrepeneur_doc_url"`
	BusinessLicenseDocURL      string   `json:"business_license_doc_url"`
	CompanyCertificateRegURL   string   `json:"company_certificate_reg_url"`
	PaymentTermID              string   `json:"payment_term_id"`
	PicFinanceName             string   `json:"pic_finance_name"`
	PicFinanceContact          string   `json:"pic_finance_contact"`
	ExchangeInvoice            int8     `json:"exchange_invoice"`
	ExchangeInvoiceTime        string   `json:"exchange_invoice_time"`
	InvoiceTerm                int8     `json:"invoice_term"`
	FinanceEmail               string   `json:"finance_email"`
	BillingAddressReferTo      int8     `json:"billing_address_refer_to"`
	BillingAddressID           int64    `json:"billing_address_id"`
	BillingAddressName         string   `json:"billing_address_name" valid:"required"`
	BillingAddressRegion       string   `json:"billing_address_region" valid:"required"`
	BillingAddressDetail1      string   `json:"billing_address_detail_1" valid:"required"`
	BillingAddressDetail2      string   `json:"billing_address_detail_2"`
	BillingAddressDetail3      string   `json:"billing_address_detail_3"`
	BillingAddressProvince     string   `json:"billing_address_province" valid:"required"`
	BillingAddressCity         string   `json:"billing_address_city" valid:"required"`
	BillingAddressDistrict     string   `json:"billing_address_district" valid:"required"`
	BillingAddressSubDistrict  string   `json:"billing_address_sub_district" valid:"required"`
	BillingAddressPostalCode   string   `json:"billing_address_postal_code" valid:"required"`
	BillingAddressNote         string   `json:"billing_address_note"`
	BillingAddressLatitude     string   `json:"billing_address_latitude" valid:"required"`
	BillingAddressLongitude    string   `json:"billing_address_longitude" valid:"required"`
	Comment1                   string   `json:"comment_1"`
	Comment2                   string   `json:"comment_2"`
}

type SalesTerritoryResponse struct {
	ID            string    `json:"id"`
	Description   string    `json:"description"`
	SalespersonID string    `json:"salesperson_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PriceLevelResponse struct {
	ID             string `json:"id"`
	Description    string `json:"description"`
	CustomerTypeID string `json:"customer_type_id"`
	RegionID       string `json:"region_id"`
}

type ShippingMethodResponse struct {
	ID              string `json:"id"`
	Description     string `json:"description"`
	Type            int8   `json:"type"`
	TypeDescription string `json:"type_description"`
}

type SiteResponse struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type GlossaryResponse struct {
	ID        int64  `json:"id,omitempty"`
	Table     string `json:"table,omitempty"`
	Attribute string `json:"attribute,omitempty"`
	ValueInt  int8   `json:"value_int"`
	ValueName string `json:"value_name"`
	Note      string `json:"note"`
}

type ProspectiveCustomerUpgradeRequest struct {
	ProspectiveCustomerID      int64    `json:"prospective_customer_id"`
	Type                       int8     `json:"type"`
	CustomerCode               string   `json:"customer_code" valid:"required"`
	BusinessName               string   `json:"business_name" valid:"required"`
	BrandName                  string   `json:"brand_name"`
	CustomerTypeID             string   `json:"customer_type_id" valid:"required"`
	BusinessTypeID             int8     `json:"business_type_id" valid:"required"`
	ArchetypeID                string   `json:"archetype_id" valid:"required"`
	CustomerClassID            string   `json:"customer_class_id" valid:"required"`
	ReferrerCode               string   `json:"referrer_code"`
	TimeConsent                int8     `json:"time_consent" valid:"required"`
	ReferenceInfo              int8     `json:"reference_info" valid:"required"`
	RegistrationChannel        int8     `json:"registration_channel"`
	OutletImage                []string `json:"outlet_image" valid:"required"`
	CompanyAddressID           int64    `json:"company_address_id"`
	CompanyAddressName         string   `json:"company_address_name"`
	CompanyAddressRegion       string   `json:"company_address_region"`
	CompanyAddressDetail1      string   `json:"company_address_detail_1"`
	CompanyAddressDetail2      string   `json:"company_address_detail_2"`
	CompanyAddressDetail3      string   `json:"company_address_detail_3"`
	CompanyAddressProvince     string   `json:"company_address_province"`
	CompanyAddressCity         string   `json:"company_address_city"`
	CompanyAddressDistrict     string   `json:"company_address_district"`
	CompanyAddressSubDistrict  string   `json:"company_address_sub_district"`
	CompanyAddressPostalCode   string   `json:"company_address_postal_code"`
	CompanyAddressNote         string   `json:"company_address_note"`
	CompanyAddressLatitude     string   `json:"company_address_latitude"`
	CompanyAddressLongitude    string   `json:"company_address_longitude"`
	ShippingAddressReferTo     int8     `json:"shipping_address_refer_to"`
	ShippingAddressID          int64    `json:"shipping_address_id"`
	ShippingAddressName        string   `json:"shipping_address_name" valid:"required"`
	ShippingAddressRegion      string   `json:"shipping_address_region" valid:"required"`
	ShippingAddressDetail1     string   `json:"shipping_address_detail_1" valid:"required"`
	ShippingAddressDetail2     string   `json:"shipping_address_detail_2"`
	ShippingAddressDetail3     string   `json:"shipping_address_detail_3"`
	ShippingAddressProvince    string   `json:"shipping_address_province" valid:"required"`
	ShippingAddressCity        string   `json:"shipping_address_city" valid:"required"`
	ShippingAddressDistrict    string   `json:"shipping_address_district" valid:"required"`
	ShippingAddressSubDistrict string   `json:"shipping_address_sub_district" valid:"required"`
	ShippingAddressPostalCode  string   `json:"shipping_address_postal_code" valid:"required"`
	ShippingAddressNote        string   `json:"shipping_address_note"`
	ShippingAddressLatitude    string   `json:"shipping_address_latitude" valid:"required"`
	ShippingAddressLongitude   string   `json:"shipping_address_longitude" valid:"required"`
	ShippingMethodID           string   `json:"shipping_method_id" valid:"required"`
	PicOrderName               string   `json:"pic_order_name" valid:"required"`
	PicOrderContact            string   `json:"pic_order_contact" valid:"required"`
	SalesTerritoryID           string   `json:"sales_territory_id" valid:"required"`
	SalespersonID              string   `json:"salesperson_id" valid:"required"`
	PriceLevelID               string   `json:"price_level_id" valid:"required"`
	OwnerName                  string   `json:"owner_name" valid:"required"`
	OwnerContact               string   `json:"owner_contact"`
	OwnerRole                  string   `json:"owner_role"`
	Email                      string   `json:"email" valid:"required"`
	IDCardDocNumber            string   `json:"id_card_doc_number" valid:"required"`
	TaxpayerDocNumber          string   `json:"taxpayer_doc_number" valid:"required"`
	PicOperationName           string   `json:"pic_operation_name" valid:"required"`
	PicOperationContact        string   `json:"pic_operation_contact" valid:"required"`
	IDCardDocURL               string   `json:"id_card_doc_url" valid:"required"`
	CompanyContractDocURL      string   `json:"company_contract_doc_url"`
	NotarialDeedDocURL         string   `json:"notarial_deed_doc_url"`
	TaxpayerDocURL             string   `json:"taxpayer_doc_url" valid:"required"`
	TaxableEntrepeneurDocURL   string   `json:"taxable_entrepeneur_doc_url"`
	BusinessLicenseDocURL      string   `json:"business_license_doc_url"`
	CompanyCertificateRegURL   string   `json:"company_certificate_reg_url"`
	PicFinanceName             string   `json:"pic_finance_name" valid:"required"`
	PicFinanceContact          string   `json:"pic_finance_contact" valid:"required"`
	PaymentTermID              string   `json:"payment_term_id" valid:"required"`
	ExchangeInvoice            int8     `json:"exchange_invoice" valid:"required"`
	ExchangeInvoiceTime        string   `json:"exchange_invoice_time"`
	InvoiceTerm                int8     `json:"invoice_term"`
	FinanceEmail               string   `json:"finance_email"`
	BillToRefersTo             int8     `json:"bill_to_refers_to"`
	BillingAddressReferTo      int8     `json:"billing_address_refer_to"`
	BillingAddressID           int64    `json:"billing_address_id"`
	BillingAddressName         string   `json:"billing_address_name" valid:"required"`
	BillingAddressRegion       string   `json:"billing_address_region" valid:"required"`
	BillingAddressDetail1      string   `json:"billing_address_detail_1" valid:"required"`
	BillingAddressDetail2      string   `json:"billing_address_detail_2"`
	BillingAddressDetail3      string   `json:"billing_address_detail_3"`
	BillingAddressProvince     string   `json:"billing_address_province" valid:"required"`
	BillingAddressCity         string   `json:"billing_address_city" valid:"required"`
	BillingAddressDistrict     string   `json:"billing_address_district" valid:"required"`
	BillingAddressSubDistrict  string   `json:"billing_address_sub_district" valid:"required"`
	BillingAddressPostalCode   string   `json:"billing_address_postal_code" valid:"required"`
	BillingAddressNote         string   `json:"billing_address_note"`
	BillingAddressLatitude     string   `json:"billing_address_latitude" valid:"required"`
	BillingAddressLongitude    string   `json:"billing_address_longitude" valid:"required"`
	Comment1                   string   `json:"comment_1"`
	Comment2                   string   `json:"comment_2"`
}

type ProspectiveCustomerGetRequest struct {
	Offset         int64
	Limit          int64
	CustomerID     string
	Search         string
	OrderBy        string
	ArchetypeID    string
	CustomerTypeID string
	RegionID       string
	SalesPersonID  string
	RequestBy      string
	Status         int8
}
type ProspectiveCustomerGetDetailRequest struct {
	ID           int64
	Code         string
	CustomerIDGP string
}
