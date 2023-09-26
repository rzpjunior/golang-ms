package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(SalesTerm))
}

// SalesTerm : struct to hold payment term model data for database
type SalesTerm struct {
	ID        int64  `orm:"column(id);auto" json:"-"`
	Code      string `orm:"column(code);size(50);null" json:"code,omitempty"`
	Name      string `orm:"column(name);size(100);null" json:"name,omitempty"`
	DaysValue int64  `orm:"column(days_value);null" json:"days_value"`
	Note      string `orm:"column(note)" json:"note,omitempty"`
	Status    int8   `orm:"column(status);null" json:"status"`
}

// TableName : set table name used by model
func (SalesTerm) TableName() string {
	return "term_payment_sls"
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *SalesTerm) MarshalJSON() ([]byte, error) {
	type Alias SalesTerm

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *SalesTerm) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *SalesTerm) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
