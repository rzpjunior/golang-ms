package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Salesperson struct {
	ID         int64     `orm:"column(id)" json:"id"`
	Code       string    `orm:"column(code)" json:"code"`
	FirstName  string    `orm:"column(firstname)" json:"firstname"`
	MiddleName string    `orm:"column(middlename)" json:"namemiddle"`
	LastName   string    `orm:"column(lastname)" json:"lastname"`
	Status     int8      `orm:"column(status)" json:"status"`
	CreatedAt  time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt  time.Time `orm:"column(updated_at)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(Salesperson))
}

func (m *Salesperson) TableName() string {
	return "salesperson"
}
