package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type ItemCategory struct {
	ID        int64     `orm:"column(id)" json:"id"`
	Code      string    `orm:"column(code)" json:"code"`
	Regions   string    `orm:"column(regions)" json:"regions"`
	Name      string    `orm:"column(name)" json:"name"`
	Status    int8      `orm:"column(status)" json:"status"`
	CreatedAt time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt time.Time `orm:"column(updated_at)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(ItemCategory))
}

func (m *ItemCategory) TableName() string {
	return "item_category"
}

func (m *ItemCategory) MarshalJSON() ([]byte, error) {
	type Alias ItemCategory

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
