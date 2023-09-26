package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

type VehicleProfile struct {
	ID                  int64   `orm:"column(id)" json:"id"`
	Code                string  `orm:"column(code)" json:"code"`
	Name                string  `orm:"column(name)" json:"name"`
	MaxKoli             float64 `orm:"column(max_koli)" json:"max_koli"`
	MaxWeight           float64 `orm:"column(max_weight)" json:"max_weight"`
	MaxFragile          float64 `orm:"column(max_fragile)" json:"max_fragile"`
	SpeedFactor         float64 `orm:"column(speed_factor)" json:"speed_factor"`
	RoutingProfile      int8    `orm:"column(routing_profile)" json:"routing_profile"`
	Skills              string  `orm:"column(skills)" json:"skills"`
	InitialCost         float64 `orm:"column(initial_cost)" json:"initial_cost"`
	SubsequentCost      float64 `orm:"column(subsequent_cost)" json:"subsequent_cost"`
	MaxAvailableVehicle int64   `orm:"column(max_available_vehicle)" json:"max_available_vehicle"`
	Status              int8    `orm:"column(status)" json:"status"`

	CourierVendorID int64 `orm:"column(courier_vendor_id)" json:"courier_vendor"`
}

func init() {
	orm.RegisterModel(new(VehicleProfile))
}

func (m *VehicleProfile) TableName() string {
	return "vehicle_profile"
}
