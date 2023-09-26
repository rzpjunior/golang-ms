// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(Product))
}

// Product : struct to hold model data for database
type Product struct {
	ID                      int64   `orm:"column(id);auto" json:"-"`
	Code                    string  `orm:"column(code)" json:"code"`
	Name                    string  `orm:"column(name)" json:"name"`
	Note                    string  `orm:"column(note)" json:"note"`
	Description             string  `orm:"column(description)" json:"description"`
	Status                  int8    `orm:"column(status)" json:"status"`
	UnitWeight              float64 `orm:"column(unit_weight)" json:"unit_weight"`
	OrderMinQty             float64 `orm:"column(order_min_qty)" json:"order_min_qty"`
	OrderMaxQty             float64 `orm:"column(order_max_qty)" json:"order_max_qty"`
	WarehouseStoStr         string  `orm:"column(warehouse_sto)" json:"-"`
	WarehousePurStr         string  `orm:"column(warehouse_pur)" json:"-"`
	WarehouseSalStr         string  `orm:"column(warehouse_sal)" json:"-"`
	TagProduct              string  `orm:"column(tag_product)" json:"tag_product"`
	OrderChannelRestriction string  `orm:"column(order_channel_restriction)" json:"-"`
	Storability             int8    `orm:"column(storability)" json:"storability"`
	Purchasability          int8    `orm:"column(purchasability)" json:"purchasability"`
	Salability              int8    `orm:"column(salability)" json:"salability"`
	UnivProductCode         string  `orm:"column(up_code)" json:"up_code"`
	Packability             int8    `orm:"column(packability)" json:"packability"`
	SparePercentage         float64 `orm:"column(spare_percentage)" json:"spare_percentage"`
	Taxable                 int8    `orm:"column(taxable)" json:"taxable"`
	TaxPercentage           float64 `orm:"column(tax_percentage)" json:"tax_percentage"`
	ExcludeArchetype        string  `orm:"column(exclude_archetype)" json:"exclude_archetype"`
	MaxDayDeliveryDate      int64   `orm:"column(max_day_delivery_date)" json:"max_day_delivery_date"`
	FragileGoods            int8    `orm:"column(fragile_goods)" json:"fragile_goods"`

	ExcludeArchetypeStr string       `orm:"-" json:"exclude_archetype_str"`
	TagProductStr       string       `orm:"-" json:"tag_product_str"`
	WarehouseSto        []*Warehouse `orm:"-" json:"warehouse_sto"`
	WarehousePur        []*Warehouse `orm:"-" json:"warehouse_pur"`
	WarehouseSal        []*Warehouse `orm:"-" json:"warehouse_sal"`
	WarehouseStoArr     []string     `orm:"-" json:"warehouse_sto_arr"`
	WarehousePurArr     []string     `orm:"-" json:"warehouse_pur_arr"`
	WarehouseSalArr     []string     `orm:"-" json:"warehouse_sal_arr"`

	Uom                      *Uom            `orm:"column(uom_id);null;rel(fk)" json:"uom"`
	Category                 *Category       `orm:"column(category_id);null;rel(fk)" json:"category"`
	GrandParent              *Category       `orm:"-" json:"grand_parent"`
	Parent                   *Category       `orm:"-" json:"parent"`
	ProductImage             []*ProductImage `orm:"reverse(many)" json:"product_image"`
	OrderChannelsRestriction []*Glossary     `orm:"-" json:"order_channel_restriction"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *Product) MarshalJSON() ([]byte, error) {
	type Alias Product

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *Product) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *Product) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
