package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Log struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID      int64              `json:"user_id,omitempty" bson:"user_id,omitempty"`
	UserIdGp    string             `json:"user_id_gp,omitempty" bson:"user_id_gp,omitempty"`
	RoleID      int64              `json:"role_id,omitempty" bson:"role_id,omitempty"`
	ReferenceID string             `json:"reference_id,omitempty" bson:"reference_id,omitempty"`
	MainRole    int8               `json:"main_role,omitempty" bson:"main_role,omitempty"`
	Type        string             `json:"type,omitempty" bson:"type,omitempty"`
	Function    string             `json:"function,omitempty" bson:"function,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Note        string             `json:"note,omitempty" bson:"note,omitempty"`
}
