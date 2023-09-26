package dto

import (
	"time"
)

type PushNotificationResponse struct {
	ID             int64                `json:"id"`
	Code           string               `json:"code"`
	CampaginName   string               `json:"campaign_name"`
	Regions        []string             `json:"regions"`
	RegionNames    []string             `json:"region_names"`
	Archetypes     []string             `json:"archetypes"`
	ArchetypeNames []string             `json:"archetype_names"`
	RedirectTo     int8                 `json:"redirect_to"`
	RedirectValue  string               `json:"redirect_value"`
	Redirect       *RedirectResponse    `json:"redirect"`
	Title          string               `json:"title"`
	Message        string               `json:"message"`
	PushNow        int8                 `json:"push_now"`
	ScheduledAt    time.Time            `json:"scheduled_at"`
	Status         int8                 `json:"status"`
	StatusConvert  string               `json:"status_convert"`
	SuccessSent    int                  `json:"success_sent"`
	FailedSent     int                  `json:"failed_sent"`
	Opened         int                  `json:"opened"`
	CreatedAt      time.Time            `json:"created_at"`
	CreatedBy      int64                `json:"created_by"`
	UpdatedAt      time.Time            `json:"updated_at,omitempty"`
	Region         []*RegionResponse    `json:"region"`
	Archetype      []*ArchetypeResponse `json:"archetype"`
}

type ArchetypeResponse struct {
	ID             string                `json:"id"`
	Code           string                `json:"code"`
	Description    string                `json:"description"`
	CustomerTypeID string                `json:"customer_type_id"`
	Status         int8                  `json:"status"`
	ConvertStatus  string                `json:"convert_status"`
	CustomerType   *CustomerTypeResponse `json:"customer_type"`
}

type CustomerTypeResponse struct {
	ID            string `json:"id"`
	Code          string `json:"code"`
	Description   string `json:"description"`
	CustomerGroup string `json:"customer_group"`
	Status        int8   `json:"status"`
	ConvertStatus string `json:"convert_status"`
}

type RegionResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type PushNotificationRequestCreate struct {
	CampaignName  string    `json:"campaign_name" valid:"required|lte:100"`
	Regions       []string  `json:"regions" valid:"required"`
	Archetypes    []string  `json:"archetypes" valid:"required"`
	RedirectTo    int8      `json:"redirect_to" valid:"required"`
	RedirectValue string    `json:"redirect_value"`
	Title         string    `json:"title" valid:"required|lte:100"`
	Message       string    `json:"message" valid:"required|lte:150"`
	PushNow       int8      `json:"push_now" valid:"required"`
	ScheduledAt   time.Time `json:"scheduled_at" valid:""`
}

type PushNotificationRequestUpdate struct {
	CampaignName  string    `json:"campaign_name" valid:"required|lte:100"`
	Regions       []string  `json:"regions" valid:"required"`
	Archetypes    []string  `json:"archetypes" valid:"required"`
	RedirectTo    int8      `json:"redirect_to" valid:"required"`
	RedirectValue string    `json:"redirect_value"`
	Title         string    `json:"title" valid:"required|lte:100"`
	Message       string    `json:"message" valid:"required|lte:150"`
	PushNow       int8      `json:"push_now" valid:"required"`
	ScheduledAt   time.Time `json:"scheduled_at" valid:""`
}

type PushNotificationRequestCancel struct {
	Note string `json:"note" valid:"required|lte:255"`
}

type PushNotificationRequestUpdateOpened struct {
	ID     int64  `json:"id"`
	Code   string `json:"code"`
	Opened int    `json:"opened"`
}
