package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationCampaign struct {
	ID                     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	NotificationCampaignID string             `json:"notification_campaign_id,omitempty" bson:"notification_campaign_id,omitempty"`
	CustomerID             string             `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	UserCustomerID         string             `json:"user_customer_id,omitempty" bson:"user_customer_id,omitempty"`
	FirebaseToken          string             `json:"firebase_token,omitempty" bson:"firebase_token,omitempty"`
	RedirectTo             int64              `json:"redirect_to,omitempty" bson:"redirect_to,omitempty"`
	RedirectToName         string             `json:"redirect_to_name,omitempty" bson:"redirect_to_name,omitempty"`
	RedirectValue          string             `json:"redirect_value,omitempty" bson:"redirect_value,omitempty"`
	RedirectValueName      string             `json:"redirect_value_name,omitempty" bson:"redirect_value_name,omitempty"`
	Sent                   int64              `json:"sent,omitempty" bson:"sent,omitempty"`
	Opened                 int64              `json:"opened,omitempty" bson:"opened,omitempty"`
	Conversion             int64              `json:"conversion,omitempty" bson:"conversion,omitempty"`
	CreatedAt              time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt              time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	RetryCount             int8               `json:"retry_count,omitempty" bson:"retry_count,omitempty"`
	FcmResultStatus        string             `json:"fcm_result_status,omitempty" bson:"fcm_result_status,omitempty"`
}
