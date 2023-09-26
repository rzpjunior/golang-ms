package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type WhiteListLogin struct {
	ID          int64  `orm:"column(id)" json:"-"` // id not set
	PhoneNumber string `orm:"column(phone_number);" json:"phone_number"`
	OTP         string `orm:"column(otp)" json:"otp"`
}

func init() {
	orm.RegisterModel(new(WhiteListLogin))
}

func (m *WhiteListLogin) TableName() string {
	return "white_list_login"
}
