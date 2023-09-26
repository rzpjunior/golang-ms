package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Courier struct {
	ID                int64     `orm:"column(id)" json:"id"`
	Code              string    `orm:"column(code)" json:"code"`
	Name              string    `orm:"column(name)" json:"name"`
	PhoneNumber       string    `orm:"column(phone_number)" json:"phone_number"`
	LicensePlate      string    `orm:"column(license_plate)" json:"license_plate"`
	EmergencyMode     int8      `orm:"column(emergency_mode)" json:"emergency_mode"`
	LastEmergencyTime time.Time `orm:"column(last_emergency_time)" json:"last_emergency_time"`
	Status            int8      `orm:"column(status)" json:"status"`

	RoleID           int64 `orm:"column(role_id)" json:"role_id"`
	UserID           int64 `orm:"column(user_id)" json:"user_id"`
	VehicleProfileID int64 `orm:"column(vehicle_profile_id)" json:"vehicle_profile_id"`
}

func init() {
	orm.RegisterModel(new(Courier))
}

func (m *Courier) TableName() string {
	return "courier"
}
