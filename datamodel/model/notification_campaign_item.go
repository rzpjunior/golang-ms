// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationCampaignItem struct {
	ID                     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	NotificationCampaignID string             `json:"notification_campaign_id,omitempty" bson:"notification_campaign_id,omitempty"`
	MerchantID             string             `json:"merchant_id,omitempty" bson:"merchant_id,omitempty"`
	RedirectTo             int64              `json:"redirect_to,omitempty" bson:"redirect_to,omitempty"`
	RedirectToName         string             `json:"redirect_to_name,omitempty" bson:"redirect_to_name,omitempty"`
	RedirectValue          string             `json:"redirect_value,omitempty" bson:"redirect_value,omitempty"`
	Sent                   int64              `json:"sent,omitempty" bson:"sent,omitempty"`
	Opened                 int64              `json:"opened,omitempty" bson:"opened,omitempty"`
	Conversion             int64              `json:"conversion,omitempty" bson:"conversion,omitempty"`
	CreatedAt              time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt              time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	RetryCount             int8               `json:"retry_count,omitempty" bson:"retry_count,omitempty"`
	FcmResultStatus        string             `json:"fcm_result_status,omitempty" bson:"fcm_result_status,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *NotificationCampaignItem) MarshalJSON() ([]byte, error) {
	type Alias NotificationCampaignItem

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
