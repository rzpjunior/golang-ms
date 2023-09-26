package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(Wrt))
}

// Wrt : struct to hold model data for database
type Wrt struct {
	ID     int64  `orm:"column(id);auto" json:"-"`
	Code   string `orm:"column(code)" json:"code"`
	Name   string `orm:"column(name)" json:"name"`
	Note   string `orm:"column(note)" json:"note"`
	Status int8   `orm:"column(status)" json:"status"`
	Type   int8   `orm:"column(type)" json:"type"`

	Area *Area `orm:"column(area_id);null;rel(fk)" json:"area_id,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *Wrt) MarshalJSON() ([]byte, error) {
	type Alias Wrt

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *Wrt) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *Wrt) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
