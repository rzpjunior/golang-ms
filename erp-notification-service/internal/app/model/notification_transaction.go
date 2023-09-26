package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationTransaction struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CustomerID string             `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	RefID      string             `json:"ref_id,omitempty" bson:"ref_id,omitempty"`
	Type       string             `json:"type,omitempty" bson:"type,omitempty"`
	Title      string             `json:"title,omitempty" bson:"title,omitempty"`
	Message    string             `json:"message,omitempty" bson:"message,omitempty"`
	Read       int64              `json:"read,omitempty" bson:"read,omitempty"`
	CreatedAt  time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Status     int8               `json:"status,omitempty" bson:"status,omitempty"`
}
