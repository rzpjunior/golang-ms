package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationSite struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	HelperID  string             `json:"helper_id,omitempty" bson:"helper_id,omitempty"`
	RefID     string             `json:"ref_id,omitempty" bson:"ref_id,omitempty"`
	Type      string             `json:"type,omitempty" bson:"type,omitempty"`
	Title     string             `json:"title,omitempty" bson:"title,omitempty"`
	Message   string             `json:"message,omitempty" bson:"message,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
