// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

// PostponeDeliveryLog model for postpone_delivery_log table IN MONGO.
// this table hold the postpone log of a delivery run sheet item
type PostponeDeliveryLog struct {
	ID               int64  `json:"id" bson:"id"`
	PostponeReason   string `json:"postpone_reason" bson:"postpone_reason"`
	StartedAtUnix    int64  `json:"started_at_unix" bson:"started_at_unix"`
	PostponedAtUnix  int64  `json:"postponed_at_unix" bson:"postponed_at_unix"`
	PostponeEvidence string `json:"postpone_evidence" bson:"postpone_evidence"`

	DeliveryRunSheetItemID int64 `json:"delivery_run_sheet_item_id" bson:"delivery_run_sheet_item_id"`
}
