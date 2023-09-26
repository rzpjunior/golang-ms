package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(Category))
}

// Category : struct to hold category model data for database
type Category struct {
	ID            int64  `orm:"column(id);auto" json:"-"`
	GrandParentID int64  `orm:"column(grandparent_id);" json:"-"`
	ParentID      int64  `orm:"column(parent_id);" json:"-"`
	Code          string `orm:"column(code)" json:"code,omitempty"`
	Name          string `orm:"column(name)" json:"name,omitempty"`
	Note          string `orm:"column(note)" json:"note,omitempty"`
	Status        int8   `orm:"column(status)" json:"status"`

	GrandParent *Category `orm:"-" json:"grand_parent"`
	Parent      *Category `orm:"-" json:"parent"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *Category) MarshalJSON() ([]byte, error) {
	type Alias Category

	return json.Marshal(&struct {
		ID            string `json:"id"`
		GrandParentID string `json:"grand_parent_id"`
		ParentID      string `json:"parent_id"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		GrandParentID: common.Encrypt(m.GrandParentID),
		ParentID:      common.Encrypt(m.ParentID),
		Alias:         (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *Category) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *Category) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
