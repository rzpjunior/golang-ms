package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type CodeGenerator struct {
	ID     int64  `orm:"column(id)" json:"-"`
	Code   string `orm:"column(code)" json:"code"`
	Domain string `orm:"column(domain)" json:"domain"`
}

func init() {
	orm.RegisterModel(new(CodeGenerator))
}

func (m *CodeGenerator) TableName() string {
	return "code_generator"
}

type CodeGeneratorReferral struct {
	ID        int64     `orm:"column(id)" json:"-"`
	Code      string    `orm:"column(code)" json:"code"`
	CreatedAt time.Time `orm:"column(created_at)" json:"created_at"`
}

func init() {
	orm.RegisterModel(new(CodeGeneratorReferral))
}

func (m *CodeGeneratorReferral) TableName() string {
	return "code_generator_referral"
}
