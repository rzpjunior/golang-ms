// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import "time"

type PostponeDeliveryLogResponse struct {
	ID               int64     `json:"id"`
	PostponeReason   string    `json:"postpone_reason"`
	StartedAt        time.Time `json:"started_at"`
	PostponedAt      time.Time `json:"postponed_at"`
	PostponeEvidence string    `json:"postpone_evidence"`

	DeliveryRunSheetItemID int64 `json:"delivery_run_sheet_item_id"`
}
