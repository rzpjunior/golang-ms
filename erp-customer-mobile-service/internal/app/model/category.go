package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

func init() {
	orm.RegisterModel(new(Category))
}

// Category : struct to hold category model data for database
type Category struct {
	ID     int64  `orm:"column(id);auto" json:"-"`
	Code   string `orm:"column(code)" json:"code,omitempty"`
	Name   string `orm:"column(name)" json:"name,omitempty"`
	Note   string `orm:"column(note)" json:"note,omitempty"`
	Status int8   `orm:"column(status)" json:"status"`

	GrandParentID int64     `orm:"column(grandparent_id);" json:"-"`
	ParentID      int64     `orm:"column(parent_id);" json:"-"`
	GrandParent   *Category `orm:"-" json:"grand_parent"`
	Parent        *Category `orm:"-" json:"parent"`
}
