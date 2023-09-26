package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(Stock))
}

// Stock : struct to hold stock model data for database
type Stock struct {
	ID                int64   `orm:"column(id);auto" json:"-"`
	AvailableStock    float64 `orm:"column(available_stock)" json:"available_stock"`
	WasteStock        float64 `orm:"column(waste_stock)" json:"waste_stock"`
	SafetyStock       float64 `orm:"column(safety_stock)" json:"safety_stock"`
	CommitedInStock   float64 `orm:"column(commited_in_stock)" json:"committed_in_stock"`
	CommitedOutStock  float64 `orm:"column(commited_out_stock)" json:"committed_out_stock"`
	ExpectedQty       float64 `orm:"column(expected_qty)" json:"expected_qty"`
	IntransitQty      float64 `orm:"column(intransit_qty)" json:"intransit_qty"`
	ReceivedQty       float64 `orm:"column(received_qty)" json:"received_qty"`
	IntransitWasteQty float64 `orm:"column(intransit_waste_qty)" json:"intransit_waste_qty"`
	Salable           int8    `orm:"column(salable)" json:"salable"`
	Purchasable       int8    `orm:"column(purchasable)" json:"purchasable"`
	Status            int8    `orm:"column(status)" json:"status"`

	Product      *Product      `orm:"column(product_id);null;rel(fk)" json:"product,omitempty"`
	Warehouse    *Warehouse    `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse,omitempty"`
	Bin          *Bin          `orm:"column(bin_id);null;rel(fk)" json:"bin,omitempty"`
	ProductGroup *ProductGroup `orm:"-" json:"product_group"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *Stock) MarshalJSON() ([]byte, error) {
	type Alias Stock

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *Stock) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *Stock) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
