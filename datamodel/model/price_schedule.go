package model

import (
	"encoding/json"

	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(PriceSchedule))
}

// PriceSchedule : struct to hold price set schedule model data for database
type PriceSchedule struct {
	ID           int64     `orm:"column(id);auto" json:"-"`
	Status       int8      `orm:"column(status)" json:"status"`
	ScheduleDate string    `orm:"column(schedule_date)" json:"schedule_date"`
	ScheduleTime string    `orm:"column(schedule_time)" json:"schedule_time"`
	CreatedAt    time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	Note         string    `orm:"column(note)" json:"note"`

	CreatedBy          *Staff               `orm:"column(created_by);null;rel(fk)" json:"created_by,omitempty"`
	PriceSet           *PriceSet            `orm:"column(price_set_id);null;rel(fk)" json:"price_set,omitempty"`
	PriceScheduleDumps []*PriceScheduleDump `orm:"reverse(many)" json:"price_schedule_dumps,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PriceSchedule) MarshalJSON() ([]byte, error) {
	type Alias PriceSchedule

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PriceSchedule) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PriceSchedule) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
