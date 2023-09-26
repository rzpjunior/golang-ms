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
	orm.RegisterModel(new(PickingRoutingStep))
}

type PickingRoutingStep struct {
	ID                      int64             `orm:"column(id);auto" json:"-"`
	PickingList             *PickingList      `orm:"column(picking_list_id);null;rel(fk)" json:"picking_list,omitempty"`
	Staff                   *Staff            `orm:"column(staff_id);null;rel(fk)" json:"staff,omitempty"`
	PickingOrderItem        *PickingOrderItem `orm:"column(picking_order_item_id);null;rel(fk)" json:"picking_order_item,omitempty"`
	Bin                     *Bin              `orm:"column(bin_id);null;rel(fk)" json:"bin_id"`
	StepType                int64             `orm:"column(step_type)" json:"step_type"`
	Sequence                int64             `orm:"column(sequence)" json:"sequence"`
	ExpectedWalkingDuration int64             `orm:"column(expected_walking_duration)" json:"expected_walking_duration"`
	ExpectedServiceDuration int64             `orm:"column(expected_service_duration)" json:"expected_service_duration"`
	WalkingStartTime        time.Time         `orm:"column(walking_start_time)" json:"walking_start_time"`
	WalkingFinishTime       time.Time         `orm:"column(walking_finish_time)" json:"walking_finish_time"`
	PickingStartTime        time.Time         `orm:"column(picking_start_time)" json:"picking_start_time"`
	PickingFinishTime       time.Time         `orm:"column(picking_finish_time)" json:"picking_finish_time"`
	CreatedAt               time.Time         `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy               *Staff            `orm:"column(created_by);rel(fk)" json:"created_by"`
	StatusStep              int8              `orm:"column(status_step)" json:"status_step"`

	LeadPicker         *Staff                `orm:"-" json:"lead_picker"`
	SalesOrders        []*SalesOrder         `orm:"-" json:"sales_orders"`
	PackRecommendation []*PackRecommendation `orm:"-" json:"pack_recommendation"`
	SalesOrderItemNote string                `orm:"-" json:"sales_order_item_note"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *PickingRoutingStep) MarshalJSON() ([]byte, error) {
	type Alias PickingRoutingStep

	alias := &struct {
		ID            string `json:"id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PickingRoutingStep) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PickingRoutingStep) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
