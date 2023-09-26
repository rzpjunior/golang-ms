package dto

import "time"

type CourierResponse struct {
	ID                int64     `json:"id"`
	Code              string    `json:"code"`
	Name              string    `json:"name"`
	PhoneNumber       string    `json:"phone_number"`
	LicensePlate      string    `json:"license_plate"`
	EmergencyMode     int8      `json:"emergency_mode"`
	LastEmergencyTime time.Time `json:"last_emergency_time"`
	Status            int8      `json:"status"`
	StatusConvert     string    `json:"status_convert"`

	RoleID           int64 `json:"role_id"`
	UserID           int64 `json:"user_id"`
	VehicleProfileID int64 `json:"vehicle_profile_id"`
}

type EmergencyModeRequest struct {
	InterID          string `json:"interid"`
	GnlCourierID     string `json:"gnl_courier_id"`
	GnlEmergencyMode int    `json:"gnl_emergencymode"`
}
