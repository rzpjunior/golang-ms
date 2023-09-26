package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Item struct {
	ID                      int64     `orm:"column(id)" json:"id"`
	Code                    string    `orm:"column(code)" json:"code"`
	UomID                   int64     `orm:"column(uom_id)" json:"uom_id"`
	ClassID                 int64     `orm:"column(class_id)" json:"class_id"`
	ItemCategoryID          int64     `orm:"column(item_category_id)" json:"item_category"`
	Description             string    `orm:"column(description)" json:"description"`
	UnitWeightConversion    float64   `orm:"column(unit_weight_conversion)" json:"unit_weight_conversion"`
	OrderMinQty             float64   `orm:"column(order_min_qty)" json:"order_min_qty"`
	OrderMaxQty             float64   `orm:"column(order_max_qty)" json:"order_max_qty"`
	ItemType                string    `orm:"column(item_type)" json:"item_type"`
	Packability             string    `orm:"column(packability)" json:"packability"`
	Capitalize              string    `orm:"column(capitalize)" json:"capitalize"`
	Note                    string    `orm:"column(note)" json:"note"`
	ExcludeArchetype        string    `orm:"column(exclude_archetype)" json:"exclude_archetype"`
	MaxDayDeliveryDate      int8      `orm:"column(max_day_delivery_date)" json:"max_day_delivery_date"`
	FragileGoods            string    `orm:"column(fragile_goods)" json:"fragile_goods"`
	Taxable                 string    `orm:"column(taxable)" json:"taxable"`
	OrderChannelRestriction string    `orm:"column(order_channel_restriction)" json:"order_channel_restriction"`
	Status                  int8      `orm:"column(status)" json:"status"`
	CreatedAt               time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt               time.Time `orm:"column(updated_at)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(Item))
}

func (m *Item) TableName() string {
	return "item"
}
