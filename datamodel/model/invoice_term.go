package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(InvoiceTerm))
}

// Invoice Term : struct to hold payment term model data for database
type InvoiceTerm struct {
	ID     int64  `orm:"column(id);auto" json:"-"`
	Code   string `orm:"column(code);size(50);null" json:"code"`
	Name   string `orm:"column(name);size(100);null" json:"name"`
	NameID string `orm:"column(name_id);size(100);null" json:"name_id"`
	Note   string `orm:"column(note)" json:"note"`
	Status int8   `orm:"column(status);null" json:"status"`
}

// TableName : set table name used by model
func (InvoiceTerm) TableName() string {
	return "term_invoice_sls"
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *InvoiceTerm) MarshalJSON() ([]byte, error) {
	type Alias InvoiceTerm

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *InvoiceTerm) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *InvoiceTerm) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
