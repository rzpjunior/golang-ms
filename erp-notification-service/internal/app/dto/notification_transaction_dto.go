package dto

import (
	"time"
)

type NotificationTransactionResponse struct {
	ID         string    `json:"id,omitempty" bson:"_id,omitempty"`
	CustomerID string    `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	RefID      string    `json:"ref_id,omitempty" bson:"ref_id,omitempty"`
	Type       string    `json:"type,omitempty" bson:"type,omitempty"`
	Title      string    `json:"title,omitempty" bson:"title,omitempty"`
	Message    string    `json:"message,omitempty" bson:"message,omitempty"`
	Read       int64     `json:"read,omitempty" bson:"read,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type GetNotificationTransactionRequest struct {
	Offset     int64  `json:"offset"`
	Limit      int64  `json:"limit"`
	CustomerID string `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
}

type SendNotificationTransactionRequest struct {
	CustomerID string `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	RefID      string `json:"ref_id,omitempty" bson:"ref_id,omitempty"`
	Type       string `json:"type,omitempty" bson:"type,omitempty"`
	Status     int8   `json:"status,omitempty" bson:"status,omitempty"`
	SendTo     string `json:"send_to" valid:"required"`
	NotifCode  string `json:"notif_code" bson:"notif_code"`
	RefCode    string `json:"ref_code" bson:"ref_code"`
}

type UpdateReadNotificationTransactionRequest struct {
	RefID      string `json:"ref_id,omitempty" bson:"ref_id,omitempty"`
	CustomerID string `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
}

type CountUnreadNotificationTransactionRequest struct {
	CustomerID string `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
}
