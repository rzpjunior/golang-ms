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
