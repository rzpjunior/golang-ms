package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(CodeGeneratorReferral))
}

// CodeGeneratorReferral : struct to hold price set model data for database
type CodeGeneratorReferral struct {
	ID        int64     `orm:"column(id);auto" json:"-"`
	Code      string    `orm:"column(code);size(30);null" json:"table"`
	CreatedAt time.Time `orm:"column(created_at);size(30);null" json:"created_at"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *CodeGeneratorReferral) MarshalJSON() ([]byte, error) {
	type Alias CodeGeneratorReferral

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *CodeGeneratorReferral) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *CodeGeneratorReferral) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
