package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

func init() {
	orm.RegisterModel(new(SalesOrderFeedback))
}

// SalesOrderFeedback : struct to hold payment term model data for database
type SalesOrderFeedback struct {
	ID             int64     `orm:"column(id);auto" json:"-"`
	SalesOrderCode string    `orm:"column(sales_order_code);size(50);null" json:"sales_order_code,omitempty"`
	DeliveryDate   string    `orm:"column(delivery_date);null" json:"delivery_date,omitempty"`
	RatingScore    int8      `orm:"column(rating_score)" json:"rating_score"`
	Tags           string    `orm:"column(tags);size(100);null" json:"tags,omitempty"`
	Description    string    `orm:"column(description);size(250)" json:"description,omitempty"`
	ToBeContacted  int8      `orm:"column(to_be_contacted)" json:"-"`
	CreatedAt      time.Time `orm:"column(created_at);null" json:"-"`
	TotalCharge    float64   `orm:"-" json:"total_charge"`
	SalesOrder     int64     `orm:"column(sales_order_id)" json:"-"`
	Customer       int64     `orm:"column(customer_id)" json:"-"`
}

// TableName : set table name used by model
func (SalesOrderFeedback) TableName() string {
	return "sales_order_feedback"
}
