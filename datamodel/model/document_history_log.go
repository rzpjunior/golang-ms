// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocumentHistoryLog struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AuditLogID string             `json:"audit_log_id,omitempty" bson:"audit_log_id,omitempty"`
	RefID      string             `json:"ref_id,omitempty" bson:"ref_id,omitempty"`
	Type       string             `json:"type,omitempty" bson:"type,omitempty"`
	ChangesLog ChangesLog         `json:"changes_log,omitempty" bson:"changes_log,omitempty"`
}

type ChangesLog struct {
	PreviousData []Data `json:"previous_data,omitempty" bson:"previous_data,omitempty"`
	AfterData    []Data `json:"after_data,omitempty" bson:"after_data,omitempty"`
}

type Data struct {
	FieldName string `json:"field_name,omitempty" bson:"field_name,omitempty"`
	Value     string `json:"value,omitempty" bson:"value,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *DocumentHistoryLog) MarshalJSON() ([]byte, error) {
	type Alias DocumentHistoryLog

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
