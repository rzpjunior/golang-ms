package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type ItemImage struct {
	ID        int64     `orm:"column(id)" json:"id"`
	ItemID    int64     `orm:"column(item_id)" json:"item_id"`
	ImageUrl  string    `orm:"column(image_url)" json:"image_url"`
	MainImage int8      `orm:"column(main_image)" json:"main_image"`
	CreatedAt time.Time `orm:"column(created_at)" json:"created_at"`
}

func init() {
	orm.RegisterModel(new(ItemImage))
}

func (m *ItemImage) TableName() string {
	return "item_image"
}

func (m *ItemImage) MarshalJSON() ([]byte, error) {
	type Alias ItemImage

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
