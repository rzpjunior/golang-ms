package dto

import "time"

type CourierResponse struct {
	ID                string     `json:"id"`
	Name              string     `json:"name"`
	PhoneNumber       string     `json:"phone_number"`
	VehicleProfileId  string     `json:"vehicle_profile_id"`
	LicensePlate      string     `json:"license_plate"`
	EmergencyMode     int32      `json:"emergency_mode"`
	LastEmergencyTime *time.Time `json:"last_emergency_time"`
	Status            int32      `json:"status"`
}

// Get Courier
type GetCourierRequest struct {
	Limit           int
	Offset          int
	Name            string
	CourierVendorID string
}
