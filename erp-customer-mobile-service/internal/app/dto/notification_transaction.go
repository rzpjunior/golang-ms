package dto

import "time"

type NotificationTransactionResponse struct {
	ID        string `json:"id"`
	RefId     string `json:"ref_id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	Read      int64  `json:"read"`
	CreateAt  string `json:"created_at"`
	Status    string `json:"status"`
	ValueName string `json:"value_name"`
}

type NotificationTransactionRequestGet struct {
	Platform string `json:"platform" valid:"required"`
	Offset   int64  `json:"offset"`
	Limit    int64  `json:"limit"`

	Session *SessionDataCustomer
}

type NotificationTransactionRequestUpdateRead struct {
	Platform string                            `json:"platform" valid:"required"`
	Data     NotificationTransactionUpdateRead `json:"data" valid:"required"`

	Session *SessionDataCustomer
}

type NotificationTransactionUpdateRead struct {
	RefId string `json:"ref_id" valid:"required"`
}

type NotificationTransactionRequestCountUnread struct {
	Platform string `json:"platform" valid:"required"`

	Session *SessionDataCustomer
}

type NotificationTransactionRequestCreate struct {
	CustomerID string    `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	RefID      string    `json:"ref_id,omitempty" bson:"ref_id,omitempty"`
	Type       string    `json:"type,omitempty" bson:"type,omitempty"`
	Title      string    `json:"title,omitempty" bson:"title,omitempty"`
	Message    string    `json:"message,omitempty" bson:"message,omitempty"`
	Read       int64     `json:"read,omitempty" bson:"read,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type NotificationTransactionCountUnreadResponse struct {
	Unread int64 `json:"unread"`
}
