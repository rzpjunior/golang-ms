// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/util"
)

func init() {
	orm.RegisterModel(new(PickingOrderAssign))
}

// PickingOrderAssign model for picking order assign table.
type PickingOrderAssign struct {
	ID                  int64     `orm:"column(id);auto" json:"-"`
	Status              int8      `orm:"column(status);null" json:"status"`
	DispatchStatus      int8      `orm:"column(dispatch_status);null" json:"dispatch_status"`
	DispatchTimestamp   time.Time `orm:"column(dispatch_timestamp)" json:"dispatch_timestamp"`
	PlanningVendor      string    `orm:"column(planning_vendor)" json:"planning_vendor"`
	BeenRejected        int8      `orm:"column(been_rejected);null" json:"been_rejected"`
	Note                string    `orm:"column(note);null" json:"note"`
	CheckinTimestamp    time.Time `orm:"column(checkin_timestamp)" json:"checkin_timestamp"`
	CheckoutTimestamp   time.Time `orm:"column(checkout_timestamp)" json:"checkout_timestamp"`
	CheckerInTimestamp  time.Time `orm:"column(checker_in_timestamp)" json:"checker_in_timestamp"`
	CheckerOutTimestamp time.Time `orm:"column(checker_out_timestamp)" json:"checker_out_timestamp"`
	AssignTimestamp     time.Time `orm:"column(assign_timestamp)" json:"assign_timestamp"`
	TotalKoli           float64   `orm:"column(total_koli)" json:"total_koli"`
	TotalScanDispatch   int       `orm:"column(total_scan_dispatch)" json:"total_scan_dispatch"`
	CheckedAt           time.Time `orm:"column(checked_at);type(timestamp);null" json:"checked_at"`
	CheckedBy           *Staff    `orm:"column(checked_by);null;rel(fk)" json:"checked_by"`
	StatusConvert       string    `orm:"-" json:"status_convert"`

	PickingOrder     *PickingOrder       `orm:"column(picking_order_id);null;rel(fk)" json:"picking_order,omitempty"`
	SalesOrder       *SalesOrder         `orm:"column(sales_order_id);null;rel(fk)" json:"sales_order,omitempty"`
	Helper           *Staff              `orm:"column(staff_id);null;rel(fk)" json:"helper"`
	Courier          *Courier            `orm:"column(courier_id);null;rel(fk)" json:"courier"`
	CourierVendor    *CourierVendor      `orm:"column(courier_vendor_id);null;rel(fk)" json:"courier_vendor"`
	Dispatcher       *Staff              `orm:"column(dispatcher_id);null;rel(fk)" json:"dispatcher"`
	PickingList      *PickingList        `orm:"column(picking_list_id);null;rel(fk)" json:"picking_list"`
	PickingOrderItem []*PickingOrderItem `orm:"reverse(many)" json:"picking_order_item,omitempty"`

	SubPicker      string   `orm:"column(sub_picker_id)" json:"-"`
	PickerCapacity int64    `orm:"column(picker_capacity)" json:"picker_capacity"`
	SubPickers     []*Staff `orm:"-" json:"sub_pickers,omitempty"`

	TotalItemSO             float64 `orm:"-" json:"total_item_so"`
	TotalItemOnProgress     float64 `orm:"-" json:"total_item_on_progress"`
	DeltaPrintSalesInvoice  int     `orm:"-" json:"delta_print_sales_invoice"`
	DeltaPrintDeliveryOrder int     `orm:"-" json:"delta_print_delivery_order"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *PickingOrderAssign) MarshalJSON() ([]byte, error) {
	type Alias PickingOrderAssign

	var checkinTimestampStrTemp string
	var checkoutTimestampStrTemp string
	var checkerInTimestampStrTemp string
	var checkerOutTimestampStrTemp string
	var dispatchTimestampStrTemp string
	var assignTimestampStrTemp string

	courierVendor := &CourierVendor{}
	courier := &Courier{}
	dispatcher := &Staff{}
	checker := &Staff{}
	helper := &Staff{}

	if m.CheckinTimestamp.Format("2006-01-02 15:04:05") == "0001-01-01 00:00:00" {
		checkinTimestampStrTemp = ""
	} else {
		checkinTimestampStrTemp = m.CheckinTimestamp.Format("2006-01-02 15:04:05")
	}

	if m.CheckoutTimestamp.Format("2006-01-02 15:04:05") == "0001-01-01 00:00:00" {
		checkoutTimestampStrTemp = ""
	} else {
		checkoutTimestampStrTemp = m.CheckoutTimestamp.Format("2006-01-02 15:04:05")
	}

	if m.CheckerInTimestamp.Format("2006-01-02 15:04:05") == "0001-01-01 00:00:00" {
		checkerInTimestampStrTemp = ""
	} else {
		checkerInTimestampStrTemp = m.CheckerInTimestamp.Format("2006-01-02 15:04:05")
	}

	if m.CheckerOutTimestamp.Format("2006-01-02 15:04:05") == "0001-01-01 00:00:00" {
		checkerOutTimestampStrTemp = ""
	} else {
		checkerOutTimestampStrTemp = m.CheckerOutTimestamp.Format("2006-01-02 15:04:05")
	}

	if m.DispatchTimestamp.Format("2006-01-02 15:04:05") == "0001-01-01 00:00:00" {
		dispatchTimestampStrTemp = ""
	} else {
		dispatchTimestampStrTemp = m.DispatchTimestamp.Format("2006-01-02 15:04:05")
	}

	if m.AssignTimestamp.Format("2006-01-02 15:04:05") == "0001-01-01 00:00:00" {
		assignTimestampStrTemp = ""
	} else {
		assignTimestampStrTemp = m.AssignTimestamp.Format("2006-01-02 15:04:05")
	}

	if m.CourierVendor == nil {
		m.CourierVendor = courierVendor
	}

	if m.Courier == nil {
		m.Courier = courier
	}

	if m.Dispatcher == nil {
		m.Dispatcher = dispatcher
	}

	if m.CheckedBy == nil {
		m.CheckedBy = checker
	}

	if m.Helper == nil {
		m.Helper = helper
	}

	return json.Marshal(&struct {
		ID                     string `json:"id"`
		StatusConvert          string `json:"status_convert"`
		DispatchStatusConvert  string `json:"dispatch_status_convert"`
		CheckinTimestampStr    string `json:"checkin_timestamp_str"`
		CheckoutTimestampStr   string `json:"checkout_timestamp_str"`
		CheckerInTimestampStr  string `json:"checker_in_timestamp_str"`
		CheckerOutTimestampStr string `json:"checker_out_timestamp_str"`
		DispatchTimestampStr   string `json:"dispatch_timestamp_str"`
		AssignTimestampStr     string `json:"assign_timestamp_str"`
		*Alias
	}{
		ID:                     common.Encrypt(m.ID),
		StatusConvert:          util.ConvertStatusPicking(m.Status),
		DispatchStatusConvert:  util.ConvertStatusPicking(m.DispatchStatus),
		CheckinTimestampStr:    checkinTimestampStrTemp,
		CheckoutTimestampStr:   checkoutTimestampStrTemp,
		CheckerInTimestampStr:  checkerInTimestampStrTemp,
		CheckerOutTimestampStr: checkerOutTimestampStrTemp,
		DispatchTimestampStr:   dispatchTimestampStrTemp,
		AssignTimestampStr:     assignTimestampStrTemp,
		Alias:                  (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PickingOrderAssign) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PickingOrderAssign) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
