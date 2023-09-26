package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(SalesAssignmentObjective))
}

// SalesAssignmentObjective : struct to hold Sales Assignment Objective model data for database
type SalesAssignmentObjective struct {
	ID         int64     `orm:"column(id);auto" json:"-"`
	Code       string    `orm:"column(code)" json:"code"`
	Name       string    `orm:"column(name)" json:"name"`
	Objective  string    `orm:"column(objective)" json:"objective"`
	SurveyLink *string   `orm:"column(surveylink)" json:"surveylink"`
	CreatedAt  time.Time `orm:"column(created_at)" json:"created_at"`
	CreatedBy  *Staff    `orm:"column(created_by);null;rel(fk)" json:"created_by"`
	Status     int8      `orm:"column(status)" json:"status"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *SalesAssignmentObjective) MarshalJSON() ([]byte, error) {
	type Alias SalesAssignmentObjective

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *SalesAssignmentObjective) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *SalesAssignmentObjective) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
