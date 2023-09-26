package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type ItemCategoryImage struct {
	ID             int64     `orm:"column(id)" json:"id"`
	ItemCategoryID int64     `orm:"column(item_category_id)" json:"item_category_id"`
	ImageUrl       string    `orm:"column(image_url)" json:"image_url"`
	CreatedAt      time.Time `orm:"column(created_at)" json:"created_at"`
}

func init() {
	orm.RegisterModel(new(ItemCategoryImage))
}

func (m *ItemCategoryImage) TableName() string {
	return "item_category_image"
}

func (m *ItemCategoryImage) MarshalJSON() ([]byte, error) {
	type Alias ItemCategoryImage

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
