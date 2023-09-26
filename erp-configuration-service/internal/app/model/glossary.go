package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Glossary struct {
	ID        int64  `orm:"column(id)" json:"-"`
	Table     string `orm:"column(table)" json:"table,omitempty"`
	Attribute string `orm:"column(attribute)" json:"attribute,omitempty"`
	ValueInt  int8   `orm:"column(value_int)" json:"value_int"`
	ValueName string `orm:"column(value_name)" json:"value_name"`
	Note      string `orm:"column(note)" json:"note"`
}

func init() {
	orm.RegisterModel(new(Glossary))
}

func (m *Glossary) TableName() string {
	return "glossary"
}
