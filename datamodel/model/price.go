package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(Price))
}

// Price : struct to hold price set model data for database
type Price struct {
	ID             int64   `orm:"column(id);auto" json:"-"`
	UnitPrice      float64 `orm:"column(unit_price)" json:"unit_price"`
	ShadowPrice    float64 `orm:"column(shadow_price)" json:"shadow_price"`
	ShadowPricePct int     `orm:"column(shadow_price_pct)" json:"shadow_price_pct"`

	Product  *Product  `orm:"column(product_id);null;rel(fk)" json:"product,omitempty"`
	PriceSet *PriceSet `orm:"column(price_set_id);null;rel(fk)" json:"price_set,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *Price) MarshalJSON() ([]byte, error) {
	type Alias Price

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *Price) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *Price) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
