package dto

import (
	"time"
)

type NotificationCampaignResponse struct {
	NotificationCampaignID   string    `json:"notification_campaign_id"`
	NotificationCampaignName string    `json:"notification_campaign_name"`
	Title                    string    `json:"title"`
	Message                  string    `json:"message"`
	RedirectTo               int64     `json:"redirect_to,omitempty"`
	RedirectToName           string    `json:"redirect_to_name,omitempty"`
	RedirectValue            string    `json:"redirect_value,omitempty"`
	RedirectValueName        string    `json:"redirect_value_name,omitempty"`
	Sent                     int64     `json:"sent,omitempty"`
	Opened                   int64     `json:"opened,omitempty"`
	Conversion               int64     `json:"conversion,omitempty"`
	CreatedAt                time.Time `json:"created_at,omitempty"`
}

type NotificationCampaignRequestGet struct {
	Platform string `json:"platform" valid:"required"`
	Offset   int64  `json:"offset"`
	Limit    int64  `json:"limit"`

	Session *SessionDataCustomer
}

type NotificationCampaignRequestUpdateRead struct {
	Platform string                         `json:"platform" valid:"required"`
	Data     NotificationCampaignUpdateRead `json:"data" valid:"required"`

	Session *SessionDataCustomer
}

type NotificationCampaignUpdateRead struct {
	NotificationCampaignID string `json:"notification_campaign_id" valid:"required"`
}

type NotificationCampaignRequestCountUnread struct {
	Platform string `json:"platform" valid:"required"`

	Session *SessionDataCustomer
}

type NotificationCampaignRequestCreate struct {
	NotificationCampaignID string    `json:"notification_campaign_id,omitempty" bson:"notification_campaign_id,omitempty"`
	CustomerID             string    `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	UserCustomerID         string    `json:"user_customer_id,omitempty" bson:"user_customer_id,omitempty"`
	FirebaseToken          string    `json:"firebase_token,omitempty" bson:"firebase_token,omitempty"`
	RedirectTo             int64     `json:"redirect_to,omitempty" bson:"redirect_to,omitempty"`
	RedirectToName         string    `json:"redirect_to_name,omitempty" bson:"redirect_to_name,omitempty"`
	RedirectValue          string    `json:"redirect_value,omitempty" bson:"redirect_value,omitempty"`
	RedirectValueName      string    `json:"redirect_value_name,omitempty" bson:"redirect_value_name,omitempty"`
	Sent                   int64     `json:"sent,omitempty" bson:"sent,omitempty"`
	Opened                 int64     `json:"opened,omitempty" bson:"opened,omitempty"`
	Conversion             int64     `json:"conversion,omitempty" bson:"conversion,omitempty"`
	CreatedAt              time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt              time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	RetryCount             int8      `json:"retry_count,omitempty" bson:"retry_count,omitempty"`
	FcmResultStatus        string    `json:"fcm_result_status,omitempty" bson:"fcm_result_status,omitempty"`
}

type NotificationCampaignCountUnreadResponse struct {
	Unread int64 `json:"unread"`
}
