// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Jobs struct {
	ID             primitive.ObjectID `orm:"column(_id)" json:"_id" bson:"_id,omitempty" `
	EndpointUrl    string             `orm:"column(endpoint_url)" json:"endpoint_url" bson:"endpoint_url,omitempty"`
	Topic          string             `orm:"column(topic);null" json:"topic" bson:"topic,omitempty"`
	EndpointMethod string             `orm:"column(endpoint_method);null" json:"endpoint_method" bson:"endpoint_method,omitempty"`
	RequestBody    string             `orm:"column(request_body);null" json:"request_body" bson:"request_body,omitempty"`
	ResponseBody   string             `orm:"column(response_body);null" json:"response_body" bson:"response_body,omitempty"`
	ResponseCode   int8               `orm:"column(response_code);null" json:"response_code" bson:"response_code,omitempty"`
	Status         int8               `orm:"column(status);null" json:"status" bson:"status,omitempty"`
	StartedAt      time.Time          `orm:"column(started_at);type(timestamp);null" json:"started_at" bson:"started_at,omitempty"`
	CompletedAt    time.Time          `orm:"column(completed_at);type(timestamp);null" json:"completed_at" bson:"completed_at,omitempty"`
	CreatedAt      time.Time          `orm:"column(created_at);type(timestamp);null" json:"created_at" bson:"created_at,omitempty"`
	CreatedBy      int64              `orm:"column(created_by)" json:"created_by" bson:"created_by,omitempty"`
	RetryCount     int8               `orm:"column(retry_count);null" json:"retry_count" bson:"retry_count,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *Jobs) MarshalJSON() ([]byte, error) {
	type Alias Jobs

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
