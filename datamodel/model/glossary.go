package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(Glossary))
}

// Glossary : struct to hold price set model data for database
type Glossary struct {
	ID        int64  `orm:"column(id);auto" json:"-"`
	Table     string `orm:"column(table);size(30);null" json:"table"`
	Attribute string `orm:"column(attribute);size(30);null" json:"attribute"`
	ValueInt  int8   `orm:"column(value_int);" json:"value_int"`
	ValueName string `orm:"column(value_name);size(50);null" json:"value_name"`
	Note      string `orm:"column(note)" json:"note"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *Glossary) MarshalJSON() ([]byte, error) {
	type Alias Glossary

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *Glossary) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *Glossary) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
