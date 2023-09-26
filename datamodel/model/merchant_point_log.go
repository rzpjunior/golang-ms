package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(MerchantPointLog))
}

// MerchantPointLog : struct to hold MerchantPointLog data for database
type MerchantPointLog struct {
	ID               int64     `orm:"column(id);auto" json:"-"`
	PointValue       float64   `orm:"column(point_value);null" json:"point_value"`
	RecentPoint      float64   `orm:"column(recent_point);null" json:"recent_point"`
	Status           int8      `orm:"column(status);null" json:"status"`
	CreatedDate      time.Time `orm:"column(created_date);null" json:"created_date"`
	ExpiredDate      time.Time `orm:"column(expired_date);null" json:"expired_date"`
	Note             string    `orm:"column(note);null" json:"note"`
	CurrentPointUsed float64   `orm:"column(current_point_used);null" json:"current_point_used"`
	NextPointUsed    float64   `orm:"column(next_point_used);null" json:"next_point_used"`
	TransactionType  int8      `orm:"column(transaction_type);null" json:"transaction_type"`

	Merchant   *Merchant          `orm:"column(merchant_id);null;rel(fk)" json:"merchant"`
	SalesOrder *SalesOrder        `orm:"column(sales_order_id);null;rel(fk)" json:"sales_oder"`
	EPCampaign *EdenPointCampaign `orm:"column(eden_point_campaign_id);null;rel(fk)" json:"eden_point_campaign"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *MerchantPointLog) MarshalJSON() ([]byte, error) {
	type Alias MerchantPointLog

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
func (m *MerchantPointLog) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}

// Save inserting or updating User struct into user table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *MerchantPointLog) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}
