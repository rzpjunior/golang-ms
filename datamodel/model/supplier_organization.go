// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/util"
)

func init() {
	orm.RegisterModel(new(SupplierOrganization))
}

// SupplierOrganization : struct to hold supplier organization model data for database
type SupplierOrganization struct {
	ID                int64              `orm:"column(id);auto" json:"-"`
	SupplierCommodity *SupplierCommodity `orm:"column(supplier_commodity_id);null;rel(fk)" json:"supplier_commodity,omitempty"`
	SupplierBadge     *SupplierBadge     `orm:"column(supplier_badge_id);null;rel(fk)" json:"supplier_badge,omitempty"`
	SupplierType      *SupplierType      `orm:"column(supplier_type_id);null;rel(fk)" json:"supplier_type,omitempty"`
	TermPaymentPur    *PurchaseTerm      `orm:"colum(term_payment_pur_id);null;rel(fk)" json:"term_payment_pur,omitempty"`
	SubDistrict       *SubDistrict       `orm:"column(sub_district_id);null;rel(fk)" json:"sub_district,omitempty"`
	Code              string             `orm:"column(code)" json:"code"`
	Name              string             `orm:"column(name)" json:"name"`
	Address           string             `orm:"column(address);size(350);null" json:"address"`
	Note              string             `orm:"column(note)" json:"note"`
	Status            int8               `orm:"column(status)" json:"status"`
	CreatedAt         time.Time          `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	UpdatedAt         time.Time          `orm:"column(updated_at);type(timestamp);null" json:"updated_at"`
	CreatedBy         *Staff             `orm:"column(created_by);null;rel(fk)" json:"created_by"`
	UpdatedBy         *Staff             `orm:"column(updated_by);null;rel(fk)" json:"updated_by,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *SupplierOrganization) MarshalJSON() ([]byte, error) {
	type Alias SupplierOrganization

	return json.Marshal(&struct {
		ID            string `json:"id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		StatusConvert: util.ConvertStatusMaster(m.Status),
		Alias:         (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *SupplierOrganization) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *SupplierOrganization) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
