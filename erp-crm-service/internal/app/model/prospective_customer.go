package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type ProspectiveCustomer struct {
	ID                        int64     `orm:"column(id)" json:"id"`
	Code                      string    `orm:"column(code)" json:"code"`
	ArchetypeIDGP             string    `orm:"column(archetype_id_gp)" json:"archetype_id_gp"`
	CustomerTypeIDGP          string    `orm:"column(customer_type_id_gp)" json:"customer_type_id_gp"`
	CustomerIDGP              string    `orm:"column(customer_id_gp)" json:"customer_id_gp"`
	CustomerClassIDGP         string    `orm:"column(customer_class_id_gp)" json:"customer_class_id_gp"`
	SalespersonIDGP           string    `orm:"column(salesperson_id_gp)" json:"salesperson_id_gp"`
	SalesTerritoryIDGP        string    `orm:"column(sales_territory_id_gp)" json:"sales_territory_id_gp"`
	SalesPriceLevelIDGP       string    `orm:"column(sales_price_level_id_gp)" json:"sales_price_level_id_gp"`
	ShippingMethodIDGP        string    `orm:"column(shipping_method_id_gp)" json:"shipping_method_id_gp"`
	RegionIDGP                string    `orm:"column(region_id_gp)" json:"region_id_gp"`
	BusinessName              string    `orm:"column(business_name)" json:"business_name"`
	BrandName                 string    `orm:"column(brand_name)" json:"brand_name"`
	RegStatus                 int8      `orm:"column(reg_status)" json:"reg_status"`
	Application               int8      `orm:"column(application)" json:"application"`
	CreatedAt                 time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt                 time.Time `orm:"column(updated_at)" json:"updated_at"`
	ProcessedAt               time.Time `orm:"column(processed_at)" json:"processed_at"`
	ProcessedBy               int64     `orm:"column(processed_by)" json:"processed_by"`
	DeclineType               int8      `orm:"column(decline_type)" json:"decline_type"`
	DeclineNote               string    `orm:"column(decline_note)" json:"decline_note"`
	SiteIDGP                  string    `orm:"column(site_id_gp)" json:"site_id_gp"`
	OutletImage               string    `orm:"column(outlet_image)" json:"outlet_image"`
	TimeConsent               int8      `orm:"column(time_consent)" json:"time_consent"`
	ReferenceInfo             int8      `orm:"column(reference_info)" json:"reference_info"`
	ReferrerCode              string    `orm:"column(referrer_code)" json:"referrer_code"`
	OwnerName                 string    `orm:"column(owner_name)" json:"owner_name"`
	OwnerContact              string    `orm:"column(owner_contact)" json:"owner_contact"`
	OwnerRole                 string    `orm:"column(owner_role)" json:"owner_role"`
	Email                     string    `orm:"column(email)" json:"email"`
	PicOperationName          string    `orm:"column(pic_operation_name)" json:"pic_operation_name"`
	PicOperationContact       string    `orm:"column(pic_operation_contact)" json:"pic_operation_contact"`
	BusinessTypeIDGP          int8      `orm:"column(business_type_id_gp)" json:"business_type_id_gp"`
	PicOrderName              string    `orm:"column(pic_order_name)" json:"pic_order_name"`
	PicOrderContact           string    `orm:"column(pic_order_contact)" json:"pic_order_contact"`
	PicFinanceName            string    `orm:"column(pic_finance_name)" json:"pic_finance_name"`
	PicFinanceContact         string    `orm:"column(pic_finance_contact)" json:"pic_finance_contact"`
	IDCardDocName             string    `orm:"column(id_card_doc_name)" json:"id_card_doc_name"`
	IDCardDocNumber           string    `orm:"column(id_card_doc_number)" json:"id_card_doc_number"`
	IDCardDocURL              string    `orm:"column(id_card_doc_url)" json:"id_card_doc_url"`
	TaxpayerDocName           string    `orm:"column(taxpayer_doc_name)" json:"taxpayer_doc_name"`
	TaxpayerDocNumber         string    `orm:"column(taxpayer_doc_number)" json:"taxpayer_doc_number"`
	TaxpayerDocURL            string    `orm:"column(taxpayer_doc_url)" json:"taxpayer_doc_url"`
	CompanyContractDocName    string    `orm:"column(company_contract_doc_name)" json:"company_contract_doc_name"`
	CompanyContractDocURL     string    `orm:"column(company_contract_doc_url)" json:"company_contract_doc_url"`
	NotarialDeedDocName       string    `orm:"column(notarial_deed_doc_name)" json:"notarial_deed_doc_name"`
	NotarialDeedDocURL        string    `orm:"column(notarial_deed_doc_url)" json:"notarial_deed_doc_url"`
	TaxableEntrepeneurDocName string    `orm:"column(taxable_entrepeneur_doc_name)" json:"taxable_entrepeneur_doc_name"`
	TaxableEntrepeneurDocURL  string    `orm:"column(taxable_entrepeneur_doc_url)" json:"taxable_entrepeneur_doc_url"`
	CompanyCertificateRegName string    `orm:"column(company_certificate_reg_name)" json:"company_certificate_reg_name"`
	CompanyCertificateRegURL  string    `orm:"column(company_certificate_reg_url)" json:"company_certificate_reg_url"`
	BusinessLicenseDocName    string    `orm:"column(business_license_doc_name)" json:"business_license_doc_name"`
	BusinessLicenseDocURL     string    `orm:"column(business_license_doc_url)" json:"business_license_doc_url"`
	PaymentTermIDGP           string    `orm:"column(payment_term_id_gp)" json:"payment_term_id_gp"`
	ExchangeInvoice           int8      `orm:"column(exchange_invoice)" json:"exchange_invoice"`
	ExchangeInvoiceTime       string    `orm:"column(exchange_invoice_time)" json:"exchange_invoice_time"`
	FinanceEmail              string    `orm:"column(finance_email)" json:"finance_email"`
	InvoiceTerm               int8      `orm:"column(invoice_term)" json:"invoice_term"`
	Comment1                  string    `orm:"column(comment_1)" json:"comment_1"`
	Comment2                  string    `orm:"column(comment_2)" json:"comment_2"`
}

func init() {
	orm.RegisterModel(new(ProspectiveCustomer))
}

func (m *ProspectiveCustomer) TableName() string {
	return "prospective_customer"
}
