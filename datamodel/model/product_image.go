package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(ProductImage))
}

// ProductImage : struct to hold model data for database
type ProductImage struct {
	ID        int64  `orm:"column(id);auto" json:"-"`
	ImageUrl  string `orm:"column(image_url)" json:"image_url"`
	MainImage int8   `orm:"column(main_image)" json:"main_image"`

	Product *Product `orm:"column(product_id);null;rel(fk)" json:"product,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *ProductImage) MarshalJSON() ([]byte, error) {
	type Alias ProductImage

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *ProductImage) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *ProductImage) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
