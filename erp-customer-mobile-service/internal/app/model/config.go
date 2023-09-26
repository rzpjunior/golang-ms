package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type ApplicationConfig struct {
	ID          int64  `orm:"column(id)" json:"-"`
	Application int8   `orm:"column(application)" json:"application"`
	Field       string `orm:"column(field)" json:"field"`
	Attribute   string `orm:"column(attribute)" json:"attribute"`
	Value       string `orm:"column(value)" json:"value"`
}

func init() {
	orm.RegisterModel(new(ApplicationConfig))
}

func (m *ApplicationConfig) TableName() string {
	return "config_app"
}

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

type DeliveryFee struct {
	ID          string  `json:"id"`
	MinOrder    float64 `json:"min_order"`
	DeliveryFee float64 `json:"delivery_fee"`
}
