package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(Cogs))
}

// Cogs : struct to hold model data for database
type Cogs struct {
	ID            int64     `orm:"column(id);auto" json:"-"`
	EtaDate       time.Time `orm:"column(eta_date)" json:"eta_date"`
	TotalQty      float64   `orm:"column(total_qty)" json:"total_qty"`
	TotalSubtotal float64   `orm:"column(total_subtotal)" json:"total_subtotal"`
	TotalAvg      float64   `orm:"column(total_avg);" json:"total_avg"`

	Product   *Product   `orm:"column(product_id);null;rel(fk)" json:"product_id"`
	Warehouse *Warehouse `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse_id"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *Cogs) MarshalJSON() ([]byte, error) {
	type Alias Cogs

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *Cogs) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *Cogs) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
