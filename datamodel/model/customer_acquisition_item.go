package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(CustomerAcquisitionItem))
}

// CustomerAcquisitionItem : struct to hold category model data for database
type CustomerAcquisitionItem struct {
	ID    int64 `orm:"column(id);auto" json:"-"`
	IsTop int8  `orm:"column(is_top)" json:"is_top"`

	CustomerAcquisition *CustomerAcquisition `orm:"column(customer_acquisition_id);null;rel(fk)" json:"customer_acquisition"`
	Product             *Product             `orm:"column(product_id);null;rel(fk)" json:"product"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *CustomerAcquisitionItem) MarshalJSON() ([]byte, error) {
	type Alias CustomerAcquisitionItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *CustomerAcquisitionItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *CustomerAcquisitionItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
