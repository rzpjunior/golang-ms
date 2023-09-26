package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type DeliveryOrder struct {
	ID              int64     `orm:"column(id)" json:"id"`
	Code            string    `orm:"column(code)" json:"code"`
	CustomerID      int64     `orm:"column(customer_id)" json:"customer_id"`
	SalesOrderID    int64     `orm:"column(sales_order_id)" json:"sales_order_id"`
	WrtID           int64     `orm:"column(wrt_id)" json:"wrt_id"`
	SiteID          int64     `orm:"column(site_id)" json:"site_id"`
	Status          int8      `orm:"column(status)" json:"status"`
	RecognitionDate time.Time `orm:"column(recognition_date)" json:"recognition_date"`
	CreatedDate     time.Time `orm:"column(created_date)" json:"created_date"`
}

func init() {
	orm.RegisterModel(new(DeliveryOrder))
}

func (m *DeliveryOrder) TableName() string {
	return "delivery_order"
}
