package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type CustomerType struct {
	ID           int64     `orm:"column(id)" json:"id"`
	Code         string    `orm:"column(code)" json:"code"`
	Description  string    `orm:"column(description)" json:"description"`
	GroupType    string    `orm:"column(group_type)" json:"group_type"`
	Abbreviation string    `orm:"column(abbreviation)" json:"abbreviation"`
	Status       int8      `orm:"column(status)" json:"status"`
	CreatedAt    time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt    time.Time `orm:"column(updated_at)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(CustomerType))
}

func (m *CustomerType) TableName() string {
	return "customer_type"
}
