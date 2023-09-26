package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type CustomerAcquisition struct {
	ID               int64     `orm:"column(id)" json:"id"`
	Code             string    `orm:"column(code)" json:"code"`
	Task             int8      `orm:"column(task)" json:"task"`
	Name             string    `orm:"column(name)" json:"name"`
	PhoneNumber      string    `orm:"column(phone_number)" json:"phone_number"`
	Latitude         float64   `orm:"column(latitude)" json:"latitude"`
	Longitude        float64   `orm:"column(longitude)" json:"longitude"`
	AddressName      string    `orm:"column(address_name)" json:"address_name"`
	FoodApp          int8      `orm:"column(food_app)" json:"food_app"`
	PotentialRevenue float64   `orm:"column(potential_revenue)" json:"potential_revenue"`
	TaskImageUrl     string    `orm:"column(task_image_url)" json:"task_image_url"`
	SalespersonID    int64     `orm:"column(salesperson_id)" json:"salesperson_id"`
	SalespersonIDGP  string    `orm:"column(salesperson_id_gp)" json:"salesperson_id_gp"`
	TerritoryID      int64     `orm:"column(territory_id)" json:"territory_id"`
	TerritoryIDGP    string    `orm:"column(territory_id_gp)" json:"territory_id_gp"`
	FinishDate       time.Time `orm:"column(finish_date)" json:"finish_date"`
	SubmitDate       time.Time `orm:"column(submit_date)" json:"submit_date"`
	CreatedAt        time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt        time.Time `orm:"column(updated_at)" json:"updated_at"`
	Status           int8      `orm:"column(status)" json:"status"`
}

func init() {
	orm.RegisterModel(new(CustomerAcquisition))
}

func (m *CustomerAcquisition) TableName() string {
	return "customer_acquisition"
}
