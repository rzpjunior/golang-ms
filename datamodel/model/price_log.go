package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(PriceLog))
}

type PriceLog struct {
	ID        int64     `orm:"column(id);auto" json:"-"`
	PriceID   int64     `orm:"column(price_id);null;" json:"price_id"`
	UnitPrice float64   `orm:"column(unit_price);null" json:"unit_price"`
	CreatedAt time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy *Staff    `orm:"column(created_by);null;rel(fk)" json:"created_by"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *PriceLog) MarshalJSON() ([]byte, error) {
	type Alias PriceLog

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save inserting or updating ProcureOrder struct into ProcureOrder table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to ProcureOrder.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *PriceLog) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting ProcureOrder data
// this also will truncated all data from all table
// that have relation with this ProcureOrder.
func (m *PriceLog) Delete() (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		var i int64
		if i, err = o.Delete(m); i == 0 && err == nil {
			err = orm.ErrNoAffected
		}
		return
	}
	return orm.ErrNoRows
}

// Read execute select based on data struct that already
// assigned.
func (m *PriceLog) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
