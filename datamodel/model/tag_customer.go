package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(TagCustomer))
}

// CustomerTag : struct to hold tag model data for database
type TagCustomer struct {
	ID     int64  `orm:"column(id);auto" json:"-"`
	Code   string `orm:"column(code)" json:"code"`
	Name   string `orm:"column(name)" json:"name"`
	Note   string `orm:"column(note)" json:"note"`
	Status int8   `orm:"column(status)" json:"status"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *TagCustomer) MarshalJSON() ([]byte, error) {
	type Alias TagCustomer

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *TagCustomer) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *TagCustomer) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
