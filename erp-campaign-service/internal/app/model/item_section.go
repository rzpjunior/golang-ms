package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type ItemSection struct {
	ID              int64     `orm:"column(id)" json:"-"`
	Code            string    `orm:"column(code)" json:"code"`
	Name            string    `orm:"column(name)" json:"name"`
	BackgroundImage string    `orm:"column(background_image)" json:"background_image"`
	StartAt         time.Time `orm:"column(start_at)" json:"start_at"`
	FinishAt        time.Time `orm:"column(finish_at)" json:"finish_at"`
	Regions         string    `orm:"column(regions)" json:"regions"`
	Archetypes      string    `orm:"column(archetypes)" json:"archetypes"`
	Items           string    `orm:"column(items)" json:"items"`
	Sequence        int       `orm:"column(sequence)" json:"sequence"`
	Note            string    `orm:"column(note)" json:"note"`
	Status          int8      `orm:"column(status)" json:"status"`
	CreatedAt       time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt       time.Time `orm:"column(updated_at)" json:"updated_at"`
	Type            int8      `orm:"column(type)" json:"type"`
}

func init() {
	orm.RegisterModel(new(ItemSection))
}

func (m *ItemSection) TableName() string {
	return "item_section"
}
