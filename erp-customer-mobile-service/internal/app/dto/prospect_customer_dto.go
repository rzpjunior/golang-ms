package dto

type RequestUpgradeBusiness struct {
	Platform string              `json:"platform" valid:"required"`
	Data     dataUpgradeBusiness `json:"data" valid:"required"`

	Session *SessionDataCustomer
}

type dataUpgradeBusiness struct {
	Name               string `json:"name" valid:"required"`
	BusinessTypeID     string `json:"business_type_id" valid:"required"`
	PicName            string `json:"pic_name" valid:"required"`
	Email              string `json:"email" valid:"required"`
	PhoneNumber        string `json:"phone_number" valid:"required"`
	AltPhoneNumber     string `json:"alt_phone_number"`
	ProvinceID         string `json:"province_id" valid:"required"`
	CityID             string `json:"city_id" valid:"required"`
	DistrictID         string `json:"district_id" valid:"required"`
	SubDistrictID      string `json:"sub_district_id" valid:"required"`
	StreetAddress      string `json:"street_address" valid:"required"`
	TimeConsent        int8   `json:"time_consent" valid:"required"`
	ReferenceInfo      string `json:"reference_info"`
	ReferrerCode       string `json:"referrer_code"`
	Referral           string `json:"referral"`
	PicFinanceName     string `json:"pic_finance_name"`
	PicFinanceContact  string `json:"pic_finance_contact"`
	PicBusinessName    string `json:"pic_business_name"`
	PicBusinessContact string `json:"pic_business_contact"`
	IDCardNumber       string `json:"id_card_number"`
	IDCardImage        string `json:"id_card_image"`
	SelfieImage        string `json:"selfie_image"`
	TaxpayerNumber     string `json:"taxpayer_number"`
	TaxpayerImage      string `json:"taxpayer_image"`
	OutletPhoto        string `json:"outlet_photo,omitempty"`
	OutletPhotoFront   string `json:"outlet_photo_front"`
	OutletPhotoSide    string `json:"outlet_photo_side"`
	OutletPhotoInside  string `json:"outlet_photo_inside"`
	TermPaymentSlsID   string `json:"term_payment_sls_id"`
	TermInvoiceSlsID   string `json:"term_invoice_sls_id"`
	BillingAddress     string `json:"billing_address"`
	Note               string `json:"note"`
	TNCDataIsRight     bool   `json:"tnc_data_is_right"`
	IsCheck            bool   `json:"is_check"`

	RegChannel   int8
	DataResponse responseUpgradeBusiness
	// PaymentGroupComb *model.PaymentGroupComb
	// TermInvoiceSls   *model.InvoiceTerm
	// TermPaymentSls   *model.SalesTerm
	// SubDistrict      *model.SubDistrict
	// Archetype        *model.Archetype
	// BusinessType     *model.BusinessType
	ReferrerID int64
	// MerchantReferral *model.Merchant
}

type responseUpgradeBusiness struct {
	// ProspectCustomer *model.ProspectCustomer

	// ID                 int64     `orm:"column(id);auto" json:"-"`
	// Code               string    `orm:"column(code);size(50);null" json:"code"`
	// Name               string    `orm:"column(name);size(100);null" json:"name"`
	// Gender             int8      `orm:"column(gender);null" json:"gender"`
	// BirthDate          string    `orm:"column(birth_date);null" json:"birth_date"`
	// Email              string    `orm:"column(email);null" json:"email"`
	// BusinessTypeName   string    `orm:"column(business_type_name);null" json:"business_type_name"`
	// PicName            string    `orm:"column(pic_name);size(100);null" json:"pic_name"`
	// PhoneNumber        string    `orm:"column(phone_number);size(15);null" json:"phone_number"`
	// AltPhoneNumber     string    `orm:"column(alt_phone_number);size(15);null" json:"alt_phone_number"`
	// StreetAddress      string    `orm:"column(street_address);null" json:"street_address"`
	// TimeConsent        int8      `orm:"column(time_consent);" json:"time_consent"`
	// ReferenceInfo      string    `orm:"column(reference_info);null" json:"reference_info"`
	// RegStatus          int8      `orm:"column(reg_status);null" json:"reg_status"`
	// RegChannel         int8      `orm:"column(reg_channel);null" json:"reg_channel"`
	// ReferrerCode       string    `orm:"column(referrer_code);null" json:"referrer_code"`
	// CreatedAt          time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	// ProcessedAt        time.Time `orm:"column(processed_at);type(timestamp);null" json:"processed_at"`
	// ProcessedBy        int64     `orm:"column(processed_by)" json:"processed_by"`
	// PicFinanceName     string    `orm:"column(pic_finance_name);size(100)" json:"pic_finance_name"`
	// PicFinanceContact  string    `orm:"column(pic_finance_contact);size(15)" json:"pic_finance_contact"`
	// PicBusinessName    string    `orm:"column(pic_business_name);size(100);null" json:"pic_business_name"`
	// PicBusinessContact string    `orm:"column(pic_business_contact);size(15);null" json:"pic_business_contact"`
	// IDCardNumber       string    `orm:"column(id_card_number);size(16);null" json:"id_card_number"`
	// IDCardImage        string    `orm:"column(id_card_image);size(300);null" json:"id_card_image"`
	// SelfieImage        string    `orm:"column(selfie_image);size(300);null" json:"selfie_image"`
	// TaxpayerNumber     string    `orm:"column(taxpayer_number);size(20);null" json:"taxpayer_number"`
	// TaxpayerImage      string    `orm:"column(taxpayer_image);size(300);null" json:"taxpayer_image"`
	// TermPaymentSlsId   int64     `orm:"column(term_payment_sls_id);null" json:"term_payment_sls_id"`
	// OutletPhoto        string    `orm:"column(outlet_photo);null" json:"outlet_photo"`
	// PaymentGroupSlsId  int64     `orm:"column(payment_group_sls_id);null" json:"-"`
	// TermInvoiceSlsId   int64     `orm:"column(term_invoice_sls_id);null" json:"term_invoice_sls_id"`
	// BillingAddress     string    `orm:"column(billing_address);size(350);null" json:"billing_address"`
	// Note               string    `orm:"column(note);size(250);null" json:"note"`

	// Archetype   *Archetype   `orm:"column(archetype_id);null;rel(fk)" json:"archetype,omitempty"`
	// SubDistrict *SubDistrict `orm:"column(sub_district_id);null;rel(fk)" json:"sub_district,omitempty"`
	// Merchant    *Merchant    `orm:"column(merchant_id);null;rel(fk)" json:"merchant,omitempty"`
}
