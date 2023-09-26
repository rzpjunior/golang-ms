package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type SalesOrder struct {
	ID            int64     `orm:"column(id)" json:"id"`
	Code          string    `orm:"column(code)" json:"code"`
	DocNumber     string    `orm:"column(doc_number)" json:"doc_number"`
	AddressID     int64     `orm:"column(address_id)" json:"address_id"`
	CustomerID    int64     `orm:"column(customer_id)" json:"customer_id"`
	SalespersonID int64     `orm:"column(salesperson_id)" json:"salesperson_id"`
	WrtID         int64     `orm:"column(wrt_id)" json:"wrt_id"`
	OrderTypeID   int64     `orm:"column(order_type_id)" json:"order_type_id"`
	SiteID        int64     `orm:"column(site_id)" json:"site_id"`
	Application   int8      `orm:"column(application)" json:"application"`
	Status        int8      `orm:"column(status)" json:"status"`
	OrderDate     time.Time `orm:"column(order_date)" json:"order_date"`
	Total         float64   `orm:"column(total)" json:"total"`
	CreatedDate   time.Time `orm:"column(created_date)" json:"created_date"`
	ModifiedDate  time.Time `orm:"column(modified_date)" json:"modified_date"`
	FinishedDate  time.Time `orm:"column(finished_date)" json:"finished_date"`
	CreatedAt     time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt     time.Time `orm:"column(updated_at)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(SalesOrder))
}

func (m *SalesOrder) TableName() string {
	return "sales_order"
}
