// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
)

// PostponeDeliveryLog model for postpone_delivery_log table IN MONGO.
// this table hold the postpone log of a delivery run sheet item
type PostponeDeliveryLog struct {
	ID               int64  `json:"id" bson:"id"`
	PostponeReason   string `json:"postpone_reason" bson:"postpone_reason"`
	StartedAtUnix    int64  `json:"started_at_unix,omitempty" bson:"started_at_unix"`
	PostponedAtUnix  int64  `json:"postponed_at_unix,omitempty" bson:"postponed_at_unix"`
	PostponeEvidence string `json:"postpone_evidence" bson:"postpone_evidence"`

	StartedAt   time.Time `json:"started_at,omitempty" bson:"-"`
	PostponedAt time.Time `json:"postponed_at,omitempty" bson:"-"`

	DeliveryRunSheetItem *DeliveryRunSheetItem `json:"delivery_run_sheet_item" bson:"delivery_run_sheet_item"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PostponeDeliveryLog) MarshalJSON() ([]byte, error) {
	type Alias PostponeDeliveryLog

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}
