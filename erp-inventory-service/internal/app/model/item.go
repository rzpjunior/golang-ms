package model

import (
	"encoding/json"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Item struct {
	ID                      int64        `orm:"column(id)" json:"id"`
	Name                    string       `orm:"column(name)" json:"name"`
	ItemIDGP                string       `orm:"column(item_id_gp)" json:"item_id_gp"`
	ItemCategoryID          string       `orm:"column(item_category_id)" json:"item_category_id,omitempty"`
	Note                    string       `orm:"column(note)" json:"note,omitempty"`
	ExcludeArchetype        string       `orm:"column(exclude_archetype)" json:"exclude_archetype,omitempty"`
	MaxDayDeliveryDate      int8         `orm:"column(max_day_delivery_date)" json:"max_day_delivery_date,omitempty"`
	OrderChannelRestriction string       `orm:"column(order_channel_restriction)" json:"order_channel_restriction,omitempty"`
	Packability             string       `orm:"column(packability)" json:"packability"`
	FragileGoods            string       `orm:"column(fragile_goods)" json:"fragile_goods"`
	ItemImage               []*ItemImage `orm:"-" json:"item_image"`
	ItemCategoryNameArr     []string     `orm:"-" json:"item_category_name_arr"`
}

func init() {
	orm.RegisterModel(new(Item))
}

func (m *Item) TableName() string {
	return "item"
}

func (m *Item) MarshalJSON() ([]byte, error) {
	type Alias Item

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
