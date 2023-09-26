package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type TermCondition struct {
	ID          int64  `orm:"column(id);auto" json:"-"`
	Application int8   `orm:"column(application);null" json:"application"`
	Version     string `orm:"column(version);size(50);null" json:"version"`
	Title       string `orm:"column(title);size(50);null" json:"title"`
	TitleValue  string `orm:"column(title_value);size(50);null" json:"title_value"`
	Content     string `orm:"column(content);null" json:"content"`
}

func init() {
	orm.RegisterModel(new(TermCondition))
}

func (m *TermCondition) TableName() string {
	return "term_condition"
}
