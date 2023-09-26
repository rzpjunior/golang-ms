package model

import (
	"encoding/json"

	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(CustomerAcquisition))
}

// CustomerAcquisition : struct to hold category model data for database
type CustomerAcquisition struct {
	ID               int64     `orm:"column(id);auto" json:"-"`
	Task             int8      `orm:"column(task)" json:"-"`
	TaskStr          string    `orm:"-" json:"task"`
	Name             string    `orm:"column(name)" json:"name,omitempty"`
	PhoneNumber      string    `orm:"column(phone_number)" json:"phone_number,omitempty"`
	Latitude         float64   `orm:"column(latitude)" json:"latitude,omitempty"`
	Longitude        float64   `orm:"column(longitude)" json:"longitude,omitempty"`
	AddressName      string    `orm:"column(address_name)" json:"address_name"`
	FoodApp          int8      `orm:"column(food_app)" json:"food_app"`
	PotentialRevenue float64   `orm:"column(potential_revenue)" json:"potential_revenue"`
	TaskPhoto        string    `orm:"column(task_photo)" json:"-"`
	TaskPhotoArr     []string  `orm:"-" json:"-"`
	TaskPhotoList    []string  `orm:"-" json:"task_photo_list"`
	SubmitDate       time.Time `orm:"column(submit_date);type(timestamp);null" json:"submit_date"`
	FinishDate       time.Time `orm:"column(finish_date);type(timestamp);null" json:"finish_date"`
	Status           int8      `orm:"column(status)" json:"status"`

	Salesperson             *Staff                     `orm:"column(salesperson_id);null;rel(fk)" json:"salesperson"`
	Salesgroup              *SalesGroup                `orm:"column(sales_group_id);null;rel(fk)" json:"sales_group"`
	CustomerAcquisitionItem []*CustomerAcquisitionItem `orm:"-" json:"customer_acquisition_items"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *CustomerAcquisition) MarshalJSON() ([]byte, error) {
	type Alias CustomerAcquisition

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *CustomerAcquisition) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *CustomerAcquisition) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
