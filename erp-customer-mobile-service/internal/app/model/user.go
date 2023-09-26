package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type User struct {
	ID           int64     `orm:"column(id)" json:"-"`
	Name         string    `orm:"column(name)" json:"name"`
	Nickname     string    `orm:"column(nickname)" json:"nickname"`
	Email        string    `orm:"column(email)" json:"email"`
	Password     string    `orm:"column(password)" json:"password"`
	RegionID     int64     `orm:"column(region_id)" json:"region_id"`
	ParentID     int64     `orm:"column(parent_id)" json:"parent_id"`
	SiteID       int64     `orm:"column(site_id)" json:"site_id"`
	TerritoryID  int64     `orm:"column(territory_id)" json:"territory_id"`
	EmployeeCode string    `orm:"column(employee_code)" json:"employee_code"`
	PhoneNumber  string    `orm:"column(phone_number)" json:"phone_number"`
	CreatedAt    time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt    time.Time `orm:"column(updated_at)" json:"updated_at"`
	Status       int8      `orm:"column(status)" json:"status"`
}

func init() {
	orm.RegisterModel(new(User))
}

func (m *User) TableName() string {
	return "user"
}
