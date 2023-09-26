// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(ConsolidatedPurchaseDeliverSignature))
}

// ConsolidatedPurchaseDeliverSignature: struct to hold signature data for consolidated_purchase_deliver
type ConsolidatedPurchaseDeliverSignature struct {
	ID                          int64                        `orm:"column(id);auto" json:"-"`
	ConsolidatedPurchaseDeliver *ConsolidatedPurchaseDeliver `orm:"column(consolidated_purchase_deliver_id);null;rel(fk)" json:"consolidated_purchase_deliver,omitempty"`
	Role                        string                       `orm:"column(role)" json:"role"`
	Name                        string                       `orm:"column(name)" json:"name"`
	Signature                   string                       `orm:"column(signature)" json:"signature"`
	CreatedAt                   time.Time                    `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy                   int64                        `orm:"column(created_by)" json:"-"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *ConsolidatedPurchaseDeliverSignature) MarshalJSON() ([]byte, error) {
	type Alias ConsolidatedPurchaseDeliverSignature

	return json.Marshal(&struct {
		ID        string `json:"id"`
		CreatedBy string `json:"created_by"`
		*Alias
	}{
		ID:        common.Encrypt(m.ID),
		CreatedBy: common.Encrypt(m.CreatedBy),
		Alias:     (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *ConsolidatedPurchaseDeliverSignature) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *ConsolidatedPurchaseDeliverSignature) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
