// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

func init() {
	orm.RegisterModel(new(ProspectCustomer))
}

// ProspectiveCustomer model for prospective_customer table.
type ProspectCustomer struct {
	ID                 int64     `orm:"column(id);auto" json:"-"`
	Code               string    `orm:"column(code);size(50);null" json:"code"`
	Name               string    `orm:"column(name);size(100);null" json:"name"`
	Gender             int8      `orm:"column(gender);null" json:"gender"`
	BirthDate          string    `orm:"column(birth_date);null" json:"birth_date"`
	Email              string    `orm:"column(email);null" json:"email"`
	BusinessTypeName   string    `orm:"column(business_type_name);null" json:"business_type_name"`
	PicName            string    `orm:"column(pic_name);size(100);null" json:"pic_name"`
	PhoneNumber        string    `orm:"column(phone_number);size(15);null" json:"phone_number"`
	AltPhoneNumber     string    `orm:"column(alt_phone_number);size(15);null" json:"alt_phone_number"`
	StreetAddress      string    `orm:"column(street_address);null" json:"street_address"`
	TimeConsent        int8      `orm:"column(time_consent);" json:"time_consent"`
	ReferenceInfo      string    `orm:"column(reference_info);null" json:"reference_info"`
	RegStatus          int8      `orm:"column(reg_status);null" json:"reg_status"`
	RegChannel         int8      `orm:"column(reg_channel);null" json:"reg_channel"`
	ReferrerCode       string    `orm:"column(referrer_code);null" json:"referrer_code"`
	CreatedAt          time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	ProcessedAt        time.Time `orm:"column(processed_at);type(timestamp);null" json:"processed_at"`
	ProcessedBy        int64     `orm:"column(processed_by)" json:"processed_by"`
	PicFinanceName     string    `orm:"column(pic_finance_name);size(100)" json:"pic_finance_name"`
	PicFinanceContact  string    `orm:"column(pic_finance_contact);size(15)" json:"pic_finance_contact"`
	PicBusinessName    string    `orm:"column(pic_business_name);size(100);null" json:"pic_business_name"`
	PicBusinessContact string    `orm:"column(pic_business_contact);size(15);null" json:"pic_business_contact"`
	IDCardNumber       string    `orm:"column(id_card_number);size(16);null" json:"id_card_number"`
	IDCardImage        string    `orm:"column(id_card_image);size(300);null" json:"id_card_image"`
	SelfieImage        string    `orm:"column(selfie_image);size(300);null" json:"selfie_image"`
	TaxpayerNumber     string    `orm:"column(taxpayer_number);size(20);null" json:"taxpayer_number"`
	TaxpayerImage      string    `orm:"column(taxpayer_image);size(300);null" json:"taxpayer_image"`
	TermPaymentSlsId   int64     `orm:"column(term_payment_sls_id);null" json:"term_payment_sls_id"`
	OutletPhoto        string    `orm:"column(outlet_photo);null" json:"outlet_photo"`
	PaymentGroupSlsId  int64     `orm:"column(payment_group_sls_id);null" json:"-"`
	TermInvoiceSlsId   int64     `orm:"column(term_invoice_sls_id);null" json:"term_invoice_sls_id"`
	BillingAddress     string    `orm:"column(billing_address);size(350);null" json:"billing_address"`
	Note               string    `orm:"column(note);size(250);null" json:"note"`

	Archetype   *Archetype   `orm:"-" json:"archetype,omitempty"`
	SubDistrict *AdmDivision `orm:"-" json:"adm_division,omitempty"`
	Customer    *Customer    `orm:"-" json:"customer,omitempty"`
}
