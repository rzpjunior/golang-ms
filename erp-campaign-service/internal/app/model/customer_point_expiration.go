package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type CustomerPointExpiration struct {
	ID                 int64     `orm:"column(id);null" json:"id"`
	CustomerID         int64     `orm:"column(customer_id);null" json:"customer_id"`
	CurrentPeriodPoint float64   `orm:"column(current_period_point);null" json:"current_period_point"`
	NextPeriodPoint    float64   `orm:"column(next_period_point);null" json:"next_period_point"`
	CurrentPeriodDate  time.Time `orm:"column(current_period_date);null" json:"current_period_date"`
	NextPeriodDate     time.Time `orm:"column(next_period_date);null" json:"next_period_date"`
	LastUpdatedAt      time.Time `orm:"column(last_updated_at);null" json:"last_updated_at"`
}

func init() {
	orm.RegisterModel(new(CustomerPointExpiration))
}

func (m *CustomerPointExpiration) TableName() string {
	return "customer_point_expiration"
}
