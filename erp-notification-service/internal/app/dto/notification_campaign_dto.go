package dto

import (
	"time"
)

type NotificationCampaignResponse struct {
	ID                     string    `json:"_id,omitempty" bson:"_id,omitempty"`
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

type GetNotificationCampaignRequest struct {
	Offset     int64  `json:"offset"`
	Limit      int64  `json:"limit"`
	CustomerID string `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
}

type SendNotificationCampaignRequest struct {
	NotificationCampaignID   string          `json:"notification_campaign_id,omitempty" bson:"notification_campaign_id,omitempty"`
	NotificationCampaignCode string          `json:"notification_campaign_code,omitempty" bson:"notification_campaign_code,omitempty"`
	NotificationCampaignName string          `json:"notification_campaign_name,omitempty" bson:"notification_campaign_name,omitempty"`
	Title                    string          `json:"title,omitempty" bson:"title,omitempty"`
	Message                  string          `json:"message,omitempty" bson:"message,omitempty"`
	RedirectTo               int64           `json:"redirect_to,omitempty" bson:"redirect_to,omitempty"`
	RedirectToName           string          `json:"redirect_to_name,omitempty" bson:"redirect_to_name,omitempty"`
	RedirectValue            string          `json:"redirect_value,omitempty" bson:"redirect_value,omitempty"`
	RedirectValueName        string          `json:"redirect_value_name,omitempty" bson:"redirect_value_name,omitempty"`
	Sent                     int64           `json:"sent,omitempty" bson:"sent,omitempty"`
	Opened                   int64           `json:"opened,omitempty" bson:"opened,omitempty"`
	Conversion               int64           `json:"conversion,omitempty" bson:"conversion,omitempty"`
	CreatedAt                time.Time       `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt                time.Time       `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	RetryCount               int8            `json:"retry_count,omitempty" bson:"retry_count,omitempty"`
	FcmResultStatus          string          `json:"fcm_result_status,omitempty" bson:"fcm_result_status,omitempty"`
	UserCustomer             []*UserCustomer `json:"user_customer,omitempty" bson:"user_customer,omitempty"`
}

type UserCustomer struct {
	CustomerID     int64  `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	UserCustomerID int64  `json:"user_customer_id,omitempty" bson:"user_customer_id,omitempty"`
	FirebaseToken  string `json:"firebase_token,omitempty" bson:"firebase_token,omitempty"`
}

type NotificationStatus struct {
	SuccessSent int64 `json:"success_sent,omitempty" bson:"success_sent,omitempty"`
	FailedSent  int64 `json:"failed_sent,omitempty" bson:"failed_sent,omitempty"`
}

type UpdateReadNotificationCampaignRequest struct {
	NotificationCampaignID string `json:"notification_campaign_id,omitempty" bson:"notification_campaign_id,omitempty"`
	CustomerID             string `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
}

type CountUnreadNotificationCampaignRequest struct {
	CustomerID string `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
}

type LarkBotMessage struct {
	MsgType string             `json:"msg_type"`
	Card    LarkBotMessageCard `json:"card"`
}

type LarkBotMessageCard struct {
	Config   LarkBotMessageConfig     `json:"config"`
	Elements []LarkBotMessageElements `json:"elements"`
	Header   LarkBotMessageHeader     `json:"header"`
}

type LarkBotMessageConfig struct {
	WideScreenMode bool `json:"wide_screen_mode"`
	EnableForward  bool `json:"enable_forward"`
}

type LarkBotMessageHeader struct {
	Template string                    `json:"template"`
	Title    LarkBotMessageHeaderTitle `json:"title"`
}

type LarkBotMessageHeaderTitle struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

type LarkBotMessageElements struct {
	Tag    string                 `json:"tag"`
	Text   LarkBotMessageText     `json:"text,omitempty"`
	Action []LarkBotMessageAction `json:"actions,omitempty"`
}

type LarkBotMessageAction struct {
	Tag  string             `json:"tag"`
	Text LarkBotMessageText `json:"text"`
	Type string             `json:"type"`
	URL  string             `json:"url"`
}

type LarkBotMessageText struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type LarkMessage struct {
	MsgType string `json:"msg_type"`
	Card    struct {
		Config struct {
			WideScreenMode bool `json:"wide_screen_mode"`
		} `json:"config"`
		LarkBotMessageElements
	}
}
