package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(MerchantPointSummary))
}

// MerchantPointSummary : struct to hold MerchantPointSummary data for database
type MerchantPointSummary struct {
	ID            int64     `orm:"column(id);auto" json:"-"`
	EarnedPoint   float64   `orm:"column(earned_point);null" json:"earned_point"`
	RedeemedPoint float64   `orm:"column(redeemed_point);null" json:"redeemed_point"`
	SummaryDate   time.Time `orm:"column(summary_date);null" json:"summary_date"`

	Merchant *Merchant `orm:"column(merchant_id);null;rel(fk)" json:"merchant"`
}

// MerchantPointList : struct to return list of eden point data
type MerchantPointList struct {
	WidgetPoint          *WidgetPoint            `json:"points"`
	MerchantPointSummary []*MerchantPointSummary `json:"merchant_point_summary"`
}

// WidgetPoint : struct to return points data for widget
type WidgetPoint struct {
	LastUpdated        time.Time `json:"last_updated"`
	TotalPoint         float64   `json:"total_point"`
	TotalEarnPoint     float64   `json:"total_earned_point"`
	TotalRedeemedPoint float64   `json:"total_redeemed_point"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *MerchantPointSummary) MarshalJSON() ([]byte, error) {
	type Alias MerchantPointSummary

	return json.Marshal(&struct {
		ID          string `json:"id"`
		SummaryDate string `json:"summary_date"`
		*Alias
	}{
		ID:          common.Encrypt(m.ID),
		SummaryDate: m.SummaryDate.Format("2006-01-02"),
		Alias:       (*Alias)(m),
	})
}

// Read execute select based on data struct that already
// assigned.
func (m *MerchantPointSummary) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}

// Save inserting or updating User struct into user table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *MerchantPointSummary) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}
