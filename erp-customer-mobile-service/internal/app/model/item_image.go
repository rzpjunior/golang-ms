package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

func init() {
	orm.RegisterModel(new(ItemImage))
}

// ItemImage : struct to hold model data for database
type ItemImage struct {
	ID        int64  `orm:"column(id);auto" json:"-"`
	ImageUrl  string `orm:"column(image_url)" json:"image_url"`
	MainImage int8   `orm:"column(main_image)" json:"main_image"`

	Item *Item `orm:"-" json:"item,omitempty"`
}
