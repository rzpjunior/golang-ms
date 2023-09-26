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
	ReferrerCode       string    `orm:"column(referrer_code);null" json:"referrer_code"`
	RegStatus          int8      `orm:"column(reg_status);null" json:"reg_status"`
	RegChannel         int8      `orm:"column(reg_channel);null" json:"reg_channel"`
	RegChannelName     string    `orm:"-"  json:"reg_channel_name"`
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
	BillingAddress     string    `orm:"column(billing_address);size(350);null" json:"billing_address,omitempty"`
	Note               string    `orm:"column(note);size(250);null" json:"note,omitempty"`
	PaymentGroupID     int64     `orm:"column(payment_group_sls_id);null;" json:"payment_group_id,omitempty"`
	OutletPhoto        string    `orm:"column(outlet_photo);null" json:"-"`
	OutletPhotoArr     []string  `orm:"-" json:"-"`
	OutletPhotoList    []string  `orm:"-" json:"outlet_photo_list"`
	SalespersonID      int64     `orm:"column(salesperson_id);null;" json:"salesperson_id,omitempty"`
	Salesperson        string    `orm:"-" json:"salesperson,omitempty"`
	DeclineTypeID      int64     `orm:"column(decline_type);null;" json:"-"`
	DeclineType        string    `orm:"-" json:"decline_type,omitempty"`
	DeclineNote        string    `orm:"column(decline_note);null" json:"decline_note,omitempty"`

	InvoiceTerm *InvoiceTerm `orm:"column(term_invoice_sls_id);null;rel(fk)" json:"invoice_term,omitempty"`
	PaymentTerm *SalesTerm   `orm:"column(term_payment_sls_id);null;rel(fk)" json:"payment_term,omitempty"`
	Merchant    *Merchant    `orm:"column(merchant_id);null;rel(fk)" json:"merchant,omitempty"`
	Archetype   *Archetype   `orm:"column(archetype_id);null;rel(fk)" json:"archetype,omitempty"`
	SubDistrict *SubDistrict `orm:"column(sub_district_id);null;rel(fk)" json:"sub_district,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *ProspectCustomer) MarshalJSON() ([]byte, error) {
	type Alias ProspectCustomer

	return json.Marshal(&struct {
		ID            string `json:"id"`
		SalespersonID string `json:"salesperson_id"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		SalespersonID: common.Encrypt(m.SalespersonID),
		Alias:         (*Alias)(m),
	})
}

// Save inserting or updating Promotion struct into prospect_customer table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to prospect_customer.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *ProspectCustomer) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting prospect_customer data
// this also will truncated all data from all table
// that have relation with this promotion.
func (m *ProspectCustomer) Delete() (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		var i int64
		if i, err = o.Delete(m); i == 0 && err == nil {
			err = orm.ErrNoAffected
		}
		return
	}
	return orm.ErrNoRows
}

// Read execute select based on data struct that already
// assigned.
func (m *ProspectCustomer) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
