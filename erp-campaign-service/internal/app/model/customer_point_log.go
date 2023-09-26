package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type CustomerPointLog struct {
	ID               int64   `orm:"column(id);auto" json:"-"`
	CustomerID       int64   `orm:"column(customer_id);null" json:"customer_id"`
	SalesOrderID     int64   `orm:"column(sales_order_id);null" json:"sales_oder_id"`
	EPCampaignID     int64   `orm:"column(eden_point_campaign_id);null" json:"eden_point_campaign_id"`
	PointValue       float64 `orm:"column(point_value);null" json:"point_value"`
	RecentPoint      float64 `orm:"column(recent_point);null" json:"recent_point"`
	Status           int8    `orm:"column(status);null" json:"status"`
	CreatedDate      string  `orm:"column(created_date);null" json:"created_date"`
	ExpiredDate      string  `orm:"column(expired_date);null" json:"expired_date"`
	Note             string  `orm:"column(note);null" json:"note"`
	CurrentPointUsed float64 `orm:"column(current_point_used);null" json:"current_point_used"`
	NextPointUsed    float64 `orm:"column(next_point_used);null" json:"next_point_used"`
	TransactionType  int8    `orm:"column(transaction_type);null" json:"transaction_type"`
}

func init() {
	orm.RegisterModel(new(CustomerPointLog))
}

func (m *CustomerPointLog) TableName() string {
	return "customer_point_log"
}

type PointHistoryList struct {
	CreatedDate string  `orm:"column(created_date)" json:"created_date"`
	PointValue  float64 `orm:"column(point_value)" json:"point_value"`
	StatusType  string  `orm:"column(status_type)" json:"status_type"`
	Status      int8    `orm:"column(status)" json:"status"`
}

// ReferralList : struct to hold referral list data
type ReferralList struct {
	Name      string    `orm:"column(name)" json:"name"`
	CreatedAt time.Time `orm:"column(created_at)" json:"created_at"`
}

// ReferralPointList : struct to hold referral point list data
type ReferralPointList struct {
	Name       string    `orm:"column(name)" json:"name"`
	CreatedAt  time.Time `orm:"column(created_date)" json:"created_at"`
	PointValue float64   `orm:"column(point_value)" json:"point"`
}
