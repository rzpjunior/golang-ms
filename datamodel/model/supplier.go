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
	orm.RegisterModel(new(Supplier))
}

// Supplier : struct to hold supplier model data for database
type Supplier struct {
	ID             int64     `orm:"column(id);auto" json:"-"`
	Code           string    `orm:"column(code);size(50);null" json:"code,omitempty"`
	Name           string    `orm:"column(name);size(100);null" json:"name,omitempty"`
	Email          string    `orm:"column(email);size(100);null" json:"email,omitempty"`
	PhoneNumber    string    `orm:"column(phone_number);size(15);null" json:"phone_number,omitempty"`
	AltPhoneNumber string    `orm:"column(alt_phone_number);size(15);null" json:"alt_phone_number,omitempty"`
	PicName        string    `orm:"column(pic_name);size(100);null" json:"pic_name,omitempty"`
	Address        string    `orm:"column(address);size(350);null" json:"address,omitempty"`
	Note           string    `orm:"column(note)" json:"note,omitempty"`
	Status         int8      `orm:"column(status);null" json:"status"`
	Returnable     int8      `orm:"column(returnable)" json:"returnable"`
	Rejectable     int8      `orm:"column(rejectable)" json:"rejectable"`
	BlockNumber    string    `orm:"column(block_number);size(10);null" json:"block_number"`
	CreatedAt      time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy      *Staff    `orm:"column(created_by);null;rel(fk)" json:"created_by"`

	SupplierType         *SupplierType         `orm:"column(supplier_type_id);null;rel(fk)" json:"supplier_type,omitempty"`
	PaymentTerm          *PurchaseTerm         `orm:"column(term_payment_pur_id);null;rel(fk)" json:"purchase_term,omitempty"`
	PaymentMethod        *PaymentMethod        `orm:"column(payment_method_id);null;rel(fk)" json:"payment_method,omitempty"`
	SubDistrict          *SubDistrict          `orm:"column(sub_district_id);null;rel(fk)" json:"sub_district,omitempty"`
	SupplierBadge        *SupplierBadge        `orm:"column(supplier_badge_id);null;rel(fk)" json:"supplier_badge"`
	SupplierCommodity    *SupplierCommodity    `orm:"column(supplier_commodity_id);null;rel(fk)" json:"supplier_commodity"`
	SupplierOrganization *SupplierOrganization `orm:"column(supplier_organization_id);null;rel(fk)" json:"supplier_organization"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *Supplier) MarshalJSON() ([]byte, error) {
	type Alias Supplier

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *Supplier) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *Supplier) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
