package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type ConsolidatedShipment struct {
	ID                int64     `orm:"column(id)" json:"id"`
	Code              string    `orm:"column(code)" json:"code"`
	DriverName        string    `orm:"column(driver_name)" json:"driver_name"`
	VehicleNumber     string    `orm:"column(vehicle_number)" json:"vehicle_number"`
	DriverPhoneNumber string    `orm:"column(driver_phone_number)" json:"driver_phone_number"`
	DeltaPrint        int8      `orm:"column(delta_print)" json:"delta_print"`
	Status            int8      `orm:"column(status)" json:"status"`
	CreatedAt         time.Time `orm:"column(created_at)" json:"created_at"`
	CreatedBy         int64     `orm:"column(created_by)" json:"created_by"`
	SiteName          string    `orm:"column(site_name)" json:"site_name"`
}

func init() {
	orm.RegisterModel(new(ConsolidatedShipment))
}

func (m *ConsolidatedShipment) TableName() string {
	return "consolidated_shipment"
}
