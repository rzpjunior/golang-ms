package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type DayOff struct {
	ID      int64     `orm:"column(id)" json:"-"`
	OffDate time.Time `orm:"column(off_date)" json:"off_date"`
	Note    string    `orm:"column(note)" json:"note"`
	Status  int8      `orm:"column(status)" json:"status"`
}

func init() {
	orm.RegisterModel(new(DayOff))
}

func (m *DayOff) TableName() string {
	return "day_off"
}
