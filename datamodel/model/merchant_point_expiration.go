package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(MerchantPointExpiration))
}

// MerchantPointExpiration : struct to hold MerchantPointExpiration data for database
type MerchantPointExpiration struct {
	ID                 int64     `orm:"column(merchant_id);null" json:"merchant_id"`
	CurrentPeriodPoint float64   `orm:"column(current_period_point);null" json:"current_period_point"`
	NextPeriodPoint    float64   `orm:"column(next_period_point);null" json:"next_period_point"`
	CurrentPeriodDate  time.Time `orm:"column(current_period_date);null" json:"current_period_date"`
	NextPeriodDate     time.Time `orm:"column(next_period_date);null" json:"next_period_date"`
	LastUpdatedAt      time.Time `orm:"column(last_updated_at);null" json:"last_updated_at"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *MerchantPointExpiration) MarshalJSON() ([]byte, error) {
	type Alias MerchantPointExpiration

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Read execute select based on data struct that already
// assigned.
func (m *MerchantPointExpiration) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}

// Save inserting or updating User struct into user table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *MerchantPointExpiration) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}
