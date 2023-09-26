// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"

	"git.edenfarm.id/cuxs/dbredis"
	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/cuxs/mongodb"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/tealeg/xlsx"
)

// Save : function to save data requested into database
func Save(r createRequest) (p *model.PickingOrder, e error) {

	//generate codes for document
	code, _ := util.GenerateDocCode("PIO", r.Warehouse.Code, "picking_order")
	o := orm.NewOrm()
	o.Begin()

	if e == nil {
		p = &model.PickingOrder{
			Code:            code,
			Warehouse:       r.Warehouse,
			RecognitionDate: r.RecognitionDateTime,
			Status:          1,
			Note:            r.Note,
		}

		if _, e = o.Insert(p); e == nil {
			for _, row := range r.PickingOrderAssign {
				pickingAssign := &model.PickingOrderAssign{
					SalesOrder:   row.SalesOrder,
					Status:       1,
					Helper:       row.Helper,
					PickingOrder: p,
				}

				if _, e = o.Insert(pickingAssign); e == nil {
					var arrPi []*model.PickingOrderItem
					for _, item := range row.SalesOrderItem {
						pi := &model.PickingOrderItem{
							PickingOrderAssign: pickingAssign,
							Product:            item.Product,
							PickQuantity:       0,
							OrderQuantity:      item.OrderQty,
						}
						arrPi = append(arrPi, pi)
					}

					if _, e := o.InsertMulti(100, &arrPi); e != nil {
						o.Rollback()
						return nil, e
					}

				} else {
					o.Rollback()
					return nil, e
				}

				// flag for picked SO already
				row.SalesOrder.HasPickingAssigned = 1
				if _, e = o.Update(row.SalesOrder, "HasPickingAssigned"); e != nil {
					o.Rollback()
					return nil, e
				}
			}

			e = log.AuditLogByUser(r.Session.Staff, p.ID, "picking_order", "create", "")

		} else {
			o.Rollback()
			return nil, e

		}

	}

	o.Commit()

	return p, e
}

// UpdateAssign : function to save data requested into database
func UpdateAssign(r updateRequestAssign) (u *model.PickingOrderAssign, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.PickingOrderAssign.Status = 3
	if _, e = o.Update(r.PickingOrderAssign, "Status"); e != nil {
		o.Rollback()
		return nil, e
	} else {
		r.PickingOrderAssign.PickingOrder.Status = 3
		if _, e = o.Update(r.PickingOrderAssign.PickingOrder, "Status"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	for _, row := range r.PickingOrderItems {

		row.PickingOrderItem.PickQuantity = row.PickOrderQty
		if _, e = o.Update(row.PickingOrderItem, "PickQuantity"); e != nil {
			o.Rollback()
			return nil, e
		}
	}
	e = log.AuditLogByUser(r.Session.Staff, r.PickingOrderAssign.ID, "picking order", "update pick qty", "")

	o.Commit()

	return r.PickingOrderAssign, e
}

// UpdateChecker : function to save data requested into database
func UpdateChecker(r updateRequestAssign) (u *model.PickingOrderAssign, e error) {
	o := orm.NewOrm()
	o.Begin()

	for _, row := range r.PickingOrderItems {

		row.PickingOrderItem.CheckQuantity = row.CheckOrderQty
		if _, e = o.Update(row.PickingOrderItem, "CheckQuantity"); e != nil {
			o.Rollback()
			return nil, e
		}
	}
	e = log.AuditLogByUser(r.Session.Staff, r.PickingOrderAssign.ID, "picking order", "update check qty", "")

	o.Commit()

	return r.PickingOrderAssign, e
}

// CheckoutAssign : function to save data requested into database
func CheckoutAssign(r checkoutRequestAssign) (u *model.PickingOrderAssign, e error) {
	o := orm.NewOrm()
	o1 := orm.NewOrm()
	o1.Using("read_only")
	o.Begin()

	r.PickingOrderAssign.Status = 5
	r.PickingOrderAssign.TotalKoli = r.TotalColly
	r.PickingOrderAssign.CheckoutTimestamp = time.Now()
	if _, e = o.Update(r.PickingOrderAssign, "Status", "CheckoutTimestamp", "TotalKoli"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = insertToDeliveryKolies(r.DeliveryKolies); e != nil {
		o.Rollback()
		return nil, e
	}

	var idx float64
	var arrDeliveryKoliIncrement []*model.DeliveryKoliIncrement
	for idx = 0; idx < r.TotalColly; idx++ {
		dki := &model.DeliveryKoliIncrement{
			SalesOrder: r.SalesOrder,
			Increment:  idx + 1,
			IsRead:     0,
		}
		arrDeliveryKoliIncrement = append(arrDeliveryKoliIncrement, dki)
	}

	if _, e := o.InsertMulti(100, &arrDeliveryKoliIncrement); e != nil {
		o.Rollback()
	}

	for _, row := range r.PickingOrderItems {

		row.PickingOrderItem.UnfullfillNote = row.UnfullfillNote
		if _, e = o.Update(row.PickingOrderItem, "UnfullfillNote"); e != nil {
			o.Rollback()
			return nil, e
		}
	}
	o.Commit()
	var isExist bool
	if isExist = o.QueryTable(new(model.PickingOrderAssign)).Filter("picking_list_id", r.PickingOrderAssign.PickingList.ID).Exclude("status__in", 5, 7).Exist(); !isExist {
		pl := orm.Params{
			"status": int8(2),
		}
		if _, e = o.QueryTable(new(model.PickingList)).Filter("id", r.PickingOrderAssign.PickingList.ID).Update(pl); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	e = log.AuditLogByUser(r.Session.Staff, r.PickingOrderAssign.ID, "picking order", "checkout", "")

	return r.PickingOrderAssign, e
}

// CheckinAssign : function to save data requested into database
func CheckinAssign(r checkinRequestAssign) (u *model.PickingOrderAssign, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.PickingOrderAssign.Status = 3
	r.PickingOrderAssign.CheckinTimestamp = time.Now()
	if _, e = o.Update(r.PickingOrderAssign, "Status", "CheckinTimestamp"); e != nil {
		o.Rollback()
		return nil, e
	} else {
		r.PickingOrderAssign.PickingOrder.Status = 3
		if _, e = o.Update(r.PickingOrderAssign.PickingOrder, "Status"); e != nil {
			o.Rollback()
			return nil, e
		} else {
			e = log.AuditLogByUser(r.Session.Staff, r.PickingOrderAssign.ID, "picking order", "checkin", "")
		}
	}
	o.Commit()

	return r.PickingOrderAssign, e
}

// CheckinBulkAssign : function to save data requested into database
func CheckinBulkAssign(r checkinBulkRequestAssign) (u *model.PickingOrder, e error) {
	o := orm.NewOrm()
	o.Begin()

	switch r.TypeRequest {
	case "rollback":
		for _, v := range r.PickingOrderAssign {
			v.Status = 1
			v.CheckinTimestamp = time.Time{}

			if _, e = o.Update(v, "Status", "CheckinTimestamp"); e != nil {
				o.Rollback()
				return nil, e
			}
		}
	default:
		for _, v := range r.PickingOrderAssign {
			v.Status = 3
			v.CheckinTimestamp = time.Now()

			if _, e = o.Update(v, "Status", "CheckinTimestamp"); e != nil {
				o.Rollback()
				return nil, e
			}
		}
	}

	pickingList := orm.Params{
		"status": 3,
	}
	if _, e = o.QueryTable(new(model.PickingList)).Filter("id", r.PickingOrderAssign[0].PickingList.ID).Update(pickingList); e != nil {
		o.Rollback()
		return nil, e
	}

	r.PickingOrder.Status = 3
	if _, e = o.Update(r.PickingOrder, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return r.PickingOrder, e
}

// CheckinChecker : function to save data requested into database
func CheckinChecker(r checkinRequestChecker) (u []*model.PickingOrderItem, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.PickingOrderAssign.CheckerInTimestamp = time.Now()
	r.PickingOrderAssign.CheckedBy = r.Session.Staff
	r.PickingOrderAssign.Status = 6
	if _, e = o.Update(r.PickingOrderAssign, "Status", "CheckerInTimestamp", "CheckedBy"); e != nil {
		o.Rollback()
		return nil, e
	}

	r.PickingOrderAssign.SalesOrder.IsLocked = 1
	r.PickingOrderAssign.SalesOrder.LockedBy = r.Session.Staff.ID
	if _, e = o.Update(r.PickingOrderAssign.SalesOrder, "IsLocked", "LockedBy"); e != nil {
		o.Rollback()
		return nil, e
	}

	// For Getting no change, updated, deleted item
	// FlagOrder -> 1: New, 2: Updated, 3: No changes, 4: Deleted
loop:
	for _, poi := range r.PickingOrderAssign.PickingOrderItem {
		for _, soi := range r.SalesOrderItems {
			if poi.Product.ID == soi.Product.ID {
				// ---no changes item---
				if poi.OrderQuantity == soi.OrderQty {
					poi.FlagOrder = 3
					u = append(u, poi)

					if _, e = o.Update(poi, "FlagOrder"); e != nil {
						o.Rollback()
						return nil, e
					}
				} else {
					// ---updated item---
					poi.OrderQuantity = soi.OrderQty
					poi.FlagOrder = 2
					u = append(u, poi)

					if _, e = o.Update(poi, "OrderQuantity", "FlagOrder"); e != nil {
						o.Rollback()
						return nil, e
					}
				}
				continue loop
			} else {
				// ---deleted item---
				poi.FlagOrder = 4
				u = append(u, poi)

				if _, e = o.Update(poi, "FlagOrder"); e != nil {
					o.Rollback()
					return nil, e
				}
			}
		}
	}

	// New Item for picking order item
	// Item save to db with new item also using orm beego
	// give back response to front with new flag

	var newPoi *model.PickingOrderItem
	var isExist bool

loop2:
	for _, soi2 := range r.SalesOrderItems {
		for _, poi2 := range r.PickingOrderAssign.PickingOrderItem {
			if isExist, _ = repository.GetItemIfExistByProductId(poi2.PickingOrderAssign.ID, soi2.Product.ID); isExist {
				continue loop2
			}

			if soi2.Product.ID != poi2.Product.ID {
				newPoi = &model.PickingOrderItem{
					PickingOrderAssign: r.PickingOrderAssign,
					OrderQuantity:      soi2.OrderQty,
					Product:            soi2.Product,
					CheckQuantity:      0,
					PickQuantity:       0,
					FlagOrder:          1,
				}
				u = append(u, newPoi)
			}
			if _, e = o.Insert(newPoi); e != nil {
				o.Rollback()
				return nil, e
			}
			o.Commit()
			time.Sleep(1 * time.Second)
		}
	}

	e = log.AuditLogByUser(r.Session.Staff, r.PickingOrderAssign.ID, "picking order", "checking", "")

	o.Commit()

	return u, e
}

// RequestApproval : function to save data requested into database
func RequestApproval(r updateRequestApprovalAssign) (u *model.PickingOrderAssign, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.PickingOrderAssign.TotalKoli = r.TotalColly
	r.PickingOrderAssign.Status = 4
	if _, e = o.Update(r.PickingOrderAssign, "Status", "TotalKoli"); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, row := range r.PickingOrderItems {

		row.PickingOrderItem.PickQuantity = row.PickOrderQty
		row.PickingOrderItem.UnfullfillNote = row.UnfullfillNote
		if _, e = o.Update(row.PickingOrderItem, "PickQuantity", "UnfullfillNote"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if e = insertToDeliveryKolies(r.DeliveryKolies); e != nil {
		o.Rollback()
		return nil, e
	}

	e = log.AuditLogByUser(r.Session.Staff, r.PickingOrderAssign.ID, "picking order", "need approval", "")

	o.Commit()

	return r.PickingOrderAssign, e
}

// Approve : function to save data requested into database
func Approve(r approveRequestAssign) (u *model.PickingOrderAssign, e error) {
	o := orm.NewOrm()
	o.Begin()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	r.PickingOrderAssign.Status = 5
	r.PickingOrderAssign.CheckoutTimestamp = time.Now()
	if _, e = o.Update(r.PickingOrderAssign, "Status", "CheckoutTimestamp"); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, row := range r.PickingOrderItems {

		row.PickingOrderItem.UnfullfillNote = row.UnfullfillNote
		row.PickingOrderItem.PickingFlag = 2
		if _, e = o.Update(row.PickingOrderItem, "UnfullfillNote", "PickingFlag"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	e = log.AuditLogByUser(r.Session.Staff, r.PickingOrderAssign.ID, "picking order", "approve", "")

	messageNotif := &util.MessageNotification{}
	notificationCode := "NOT0011"

	if e = orSelect.Raw("SELECT message, title FROM notification WHERE code= ?", notificationCode).QueryRow(&messageNotif); e != nil {
		o.Rollback()
		return nil, e
	}

	r.PickingOrderAssign.Helper.Read("ID")
	r.PickingOrderAssign.Helper.User.Read("ID")

	mn := &util.ModelPickingNotification{
		SendTo:    r.PickingOrderAssign.Helper.User.PickingNotifToken,
		Title:     messageNotif.Title,
		Message:   messageNotif.Message,
		Type:      "4",
		RefID:     r.PickingOrderAssign.PickingOrder.ID,
		StaffID:   r.Session.Staff.ID,
		ServerKey: util.PickingServerKeyFireBase,
	}
	util.PostPickingModelNotification(mn)

	o.Commit()

	return r.PickingOrderAssign, e
}

// Reject : function to save data requested into database
func Reject(r rejectRequestAssign) (u *model.PickingOrderAssign, e error) {
	o := orm.NewOrm()
	o.Begin()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	r.PickingOrderAssign.Status = 8
	r.PickingOrderAssign.BeenRejected = 1
	r.PickingOrderAssign.Note = r.Note
	if _, e = o.Update(r.PickingOrderAssign, "Status", "BeenRejected", "Note"); e != nil {
		o.Rollback()
		return nil, e
	} else {
		r.PickingOrderAssign.PickingOrder.Status = 3
		if _, e = o.Update(r.PickingOrderAssign.PickingOrder, "Status"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	r.PickingOrderAssign.PickingList.Status = 4
	if _, e = o.Update(r.PickingOrderAssign.PickingList, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, row := range r.PickingOrderItems {

		row.PickingOrderItem.UnfullfillNote = row.UnfullfillNote
		row.PickingOrderItem.PickingFlag = row.POIFlagging

		if _, e = o.Update(row.PickingOrderItem, "UnfullfillNote", "PickingFlag"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	e = log.AuditLogByUser(r.Session.Staff, r.PickingOrderAssign.ID, "picking order", "reject", "")

	messageNotif := &util.MessageNotification{}
	notificationCode := "NOT0010"

	if e = orSelect.Raw("SELECT message, title FROM notification WHERE code= ?", notificationCode).QueryRow(&messageNotif); e != nil {
		o.Rollback()
		return nil, e
	}

	r.PickingOrderAssign.Helper.Read("ID")
	r.PickingOrderAssign.Helper.User.Read("ID")

	mn := &util.ModelPickingNotification{
		SendTo:    r.PickingOrderAssign.Helper.User.PickingNotifToken,
		Title:     messageNotif.Title,
		Message:   messageNotif.Message,
		Type:      "4",
		RefID:     r.PickingOrderAssign.PickingOrder.ID,
		StaffID:   r.Session.Staff.ID,
		ServerKey: util.PickingServerKeyFireBase,
	}
	util.PostPickingModelNotification(mn)

	o.Commit()

	return r.PickingOrderAssign, e
}

// ApproveByChecker : function to save data requested into database
func ApproveByChecker(r approveRequestChecker) (u *model.PickingOrderAssign, e error) {
	o := orm.NewOrm()
	o.Begin()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	r.PickingOrderAssign.Status = 2
	r.PickingOrderAssign.TotalKoli = r.TotalColly
	r.PickingOrderAssign.CheckedAt = time.Now()
	r.PickingOrderAssign.CheckerOutTimestamp = time.Now()
	r.PickingOrderAssign.CheckedBy = r.Session.Staff
	if _, e = o.Update(r.PickingOrderAssign, "Status", "TotalKoli", "CheckedAt", "CheckerOutTimestamp", "CheckedBy"); e != nil {
		o.Rollback()
		return nil, e
	} else {
		if e = o.Commit(); e != nil {
			o.Rollback()
			return nil, e
		} else {
			if isExist := o.QueryTable(new(model.PickingOrderAssign)).Filter("picking_order_id", r.PickingOrderAssign.PickingOrder.ID).Exclude("status", 2).Exist(); !isExist {
				r.IsFinished = true
			}
		}
	}

	// condition if all of sales order in picking is finished, the picking order status will be finished
	if r.IsFinished {
		r.PickingOrderAssign.PickingOrder.Status = 2
		if _, e = o.Update(r.PickingOrderAssign.PickingOrder, "Status"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	for _, row := range r.PickingOrderItems {

		row.PickingOrderItem.CheckQuantity = row.CheckOrderQty
		row.PickingOrderItem.PickingFlag = 2
		if _, e = o.Update(row.PickingOrderItem, "CheckQuantity", "PickingFlag"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	// condition for add kolies when checker finished picking
	if e = insertToDeliveryKolies(r.DeliveryKolies); e != nil {
		o.Rollback()
		return nil, e
	}

	r.PickingOrderAssign.SalesOrder.IsLocked = 2
	r.PickingOrderAssign.SalesOrder.LockedBy = 0
	if _, e = o.Update(r.PickingOrderAssign.SalesOrder, "IsLocked", "LockedBy"); e != nil {
		o.Rollback()
		return nil, e
	}

	if _, e := o.QueryTable(new(model.DeliveryKoliIncrement)).Filter("sales_order_id", r.PickingOrderAssign.SalesOrder.ID).Delete(); e != nil {
		o.Rollback()
		return nil, e
	}

	var idx float64
	var arrDeliveryKoliIncrement []*model.DeliveryKoliIncrement
	for idx = 0; idx < r.TotalColly; idx++ {
		dki := &model.DeliveryKoliIncrement{
			SalesOrder: r.PickingOrderAssign.SalesOrder,
			Increment:  idx + 1,
			IsRead:     0,
		}
		arrDeliveryKoliIncrement = append(arrDeliveryKoliIncrement, dki)
	}

	if _, e := o.InsertMulti(100, &arrDeliveryKoliIncrement); e != nil {
		o.Rollback()
	}

	e = log.AuditLogByUser(r.Session.Staff, r.PickingOrderAssign.ID, "picking order", "approve_checker", "")

	messageNotif := &util.MessageNotification{}
	notificationCode := "NOT0016"

	if e = orSelect.Raw("SELECT message, title FROM notification WHERE code= ?", notificationCode).QueryRow(&messageNotif); e != nil {
		o.Rollback()
		return nil, e
	}

	r.PickingOrderAssign.Helper.Read("ID")
	r.PickingOrderAssign.Helper.User.Read("ID")

	mn := &util.ModelPickingNotification{
		SendTo:    r.PickingOrderAssign.Helper.User.PickingNotifToken,
		Title:     messageNotif.Title,
		Message:   messageNotif.Message,
		Type:      "4",
		RefID:     r.PickingOrderAssign.PickingOrder.ID,
		StaffID:   r.Session.Staff.ID,
		ServerKey: util.PickingServerKeyFireBase,
	}
	util.PostPickingModelNotification(mn)

	o.Commit()

	return r.PickingOrderAssign, e
}

// UpdateCheckQtyByScan: function for update check qty
func UpdateCheckQtyByScan(r scanRequestChecker) (poi *model.PickingOrderItem, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.PickingOrderItem.CheckQuantity = r.CheckOrderQty
	if _, e = o.Update(r.PickingOrderItem, "CheckQuantity"); e != nil {
		o.Rollback()
		return nil, e
	}
	e = log.AuditLogByUser(r.Session.Staff, r.ID, "picking_order_item", "update check qty", "")

	o.Commit()

	return r.PickingOrderItem, e
}

// RejectByChecker : function to save data requested into database
func RejectByChecker(r rejectRequestChecker) (u *model.PickingOrderAssign, e error) {
	o := orm.NewOrm()
	o.Begin()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	r.PickingOrderAssign.Status = 8
	r.PickingOrderAssign.BeenRejected = 1
	r.PickingOrderAssign.Note = r.Note
	r.PickingOrderAssign.CheckedAt = time.Now()
	r.PickingOrderAssign.CheckedBy = r.Session.Staff
	if _, e = o.Update(r.PickingOrderAssign, "Status", "BeenRejected", "Note", "CheckedAt", "CheckedBy"); e != nil {
		o.Rollback()
		return nil, e
	} else {
		r.PickingOrderAssign.PickingOrder.Status = 3
		if _, e = o.Update(r.PickingOrderAssign.PickingOrder, "Status"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	for _, row := range r.PickingOrderItems {

		row.PickingOrderItem.CheckQuantity = row.CheckOrderQty
		row.PickingOrderItem.PickingFlag = row.POIFlagging
		if _, e = o.Update(row.PickingOrderItem, "CheckQuantity", "PickingFlag"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	r.PickingOrderAssign.PickingList.Status = 4
	if _, e = o.Update(r.PickingOrderAssign.PickingList, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	// deleted picking_order_item with conditions where its belongs to and flag_order = deleted
	for _, v := range r.DeletedPOIPickingRoutingStep {
		v.PickingOrderItem = nil
		if _, e = o.Update(v, "PickingOrderItem"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	flagOrder := 4
	if _, e = o.Raw("delete from picking_order_item where picking_order_assign_id = ? and flag_order = ?", r.PickingOrderAssign.ID, flagOrder).Exec(); e != nil {
		o.Rollback()
		return nil, e
	}

	r.PickingOrderAssign.SalesOrder.IsLocked = 2
	r.PickingOrderAssign.SalesOrder.LockedBy = 0
	if _, e = o.Update(r.PickingOrderAssign.SalesOrder, "IsLocked", "LockedBy"); e != nil {
		o.Rollback()
		return nil, e
	}

	e = log.AuditLogByUser(r.Session.Staff, r.PickingOrderAssign.ID, "picking order", "reject checker", "")

	messageNotif := &util.MessageNotification{}
	notificationCode := "NOT0015"

	if e = orSelect.Raw("SELECT message, title FROM notification WHERE code= ?", notificationCode).QueryRow(&messageNotif); e != nil {
		o.Rollback()
		return nil, e
	}

	r.PickingOrderAssign.Helper.Read("ID")
	r.PickingOrderAssign.Helper.User.Read("ID")

	mn := &util.ModelPickingNotification{
		SendTo:    r.PickingOrderAssign.Helper.User.PickingNotifToken,
		Title:     messageNotif.Title,
		Message:   messageNotif.Message,
		Type:      "4",
		RefID:     r.PickingOrderAssign.PickingOrder.ID,
		StaffID:   r.Session.Staff.ID,
		ServerKey: util.PickingServerKeyFireBase,
	}
	util.PostPickingModelNotification(mn)

	o.Commit()

	return r.PickingOrderAssign, e
}

func GetPickingOrderXls(date time.Time, r []*templatePickingOrder, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	dir := util.ExportDirectory

	filename := fmt.Sprintf("TemplatePickingOrder_%s_%s_%s.xlsx", date.Format("2006-01-02"), util.ReplaceSpace(warehouse.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "Business_Type"
		row.AddCell().Value = "Order_Code"
		row.AddCell().Value = "Order_Type"
		row.AddCell().Value = "Merchant_Name"
		row.AddCell().Value = "Order_Status"
		row.AddCell().Value = "Shipping_Address"
		row.AddCell().Value = "Province"
		row.AddCell().Value = "City"
		row.AddCell().Value = "District"
		row.AddCell().Value = "Sub District"
		row.AddCell().Value = "Postal_Code"
		row.AddCell().Value = "WRT"
		row.AddCell().Value = "Order_Weight"
		row.AddCell().Value = "Delivery_Date"
		row.AddCell().Value = "Payment_Term"
		row.AddCell().Value = "Picker"
		row.AddCell().Value = "Vendor"
		row.AddCell().Value = "Planning"

		for i, v := range r {

			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.Warehouse                       // Warehouse
			row.AddCell().Value = v.BusinessType                    // Business Type
			row.AddCell().Value = v.OrderCode                       // Order Code
			row.AddCell().Value = v.OrderType                       // Order Type
			row.AddCell().Value = v.MerchantName                    // Merchant Name
			row.AddCell().Value = v.OrderStatus                     // Order Status
			row.AddCell().Value = v.ShippingAddress                 // Shipping Address
			row.AddCell().Value = v.Province                        // Province
			row.AddCell().Value = v.City                            // City
			row.AddCell().Value = v.District                        // District
			row.AddCell().Value = v.SubDistrict                     // Sub District
			row.AddCell().Value = v.PostalCode                      // Postal Code
			row.AddCell().Value = v.Wrt                             // WRT
			row.AddCell().SetFloatWithFormat(v.OrderWeight, "0.00") // Order Weight
			row.AddCell().Value = v.DeliveryDate                    // Delivery Date
			row.AddCell().Value = v.PaymentTerm                     // Payment Term
			row.AddCell().Value = v.Picker                          // Picker
			row.AddCell().Value = v.Vendor                          // Vendor
			row.AddCell().Value = v.Planning                        // Planning
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// UploadAssign : function to save data requested into database
func UploadAssign(r uploadAssignRequest) (p *model.PickingOrder, e error) {

	o1 := orm.NewOrm()
	o1.Using("read_only")
	o1.Raw("select id,code,warehouse_id,recognition_date,note,status from picking_order po where po.warehouse_id= ? and po.recognition_date = ?", r.Warehouse, r.RecognitionDate).QueryRow(&p)

	if p != nil {
		o := orm.NewOrm()
		o.Begin()

		for _, row := range r.ItemUploadAssign {
			if row.Helper == nil && row.CourierVendor == nil {
				continue
			}

			var poa *model.PickingOrderAssign
			o1.Raw("select poa.id, poa.courier_vendor_id, poa.staff_id ,poa.planning_vendor from picking_order_assign poa where poa.picking_order_id = ? and poa.sales_order_id = ?", p.ID, row.SalesOrder.ID).QueryRow(&poa)
			if poa != nil {
				poa.Helper = row.Helper
				poa.CourierVendor = row.CourierVendor
				poa.PlanningVendor = row.PlanningStr
				poa.AssignTimestamp = time.Now()
				if _, e = o.Update(poa, "Helper", "CourierVendor", "PlanningVendor", "AssignTimestamp"); e != nil {
					o.Rollback()
					return nil, e
				}
			} else {
				pl := &model.PickingList{
					ID: 1,
				}
				poa = &model.PickingOrderAssign{
					SalesOrder:      row.SalesOrder,
					Status:          1,
					Helper:          row.Helper,
					CourierVendor:   row.CourierVendor,
					PlanningVendor:  row.PlanningStr,
					PickingOrder:    p,
					DispatchStatus:  1,
					AssignTimestamp: time.Now(),
					PickingList:     pl,
				}

				if _, e = o.Insert(poa); e != nil {
					o.Rollback()
					return nil, e
				}

				var arrPi []*model.PickingOrderItem
				for _, item := range row.SalesOrderItem {
					pi := &model.PickingOrderItem{
						PickingOrderAssign: poa,
						Product:            item.Product,
						PickQuantity:       0,
						OrderQuantity:      item.OrderQty,
						FlagOrder:          3,
					}
					arrPi = append(arrPi, pi)
				}

				if _, e := o.InsertMulti(100, &arrPi); e != nil {
					o.Rollback()
					return nil, e
				}
			}
		}
		o.Commit()

		e = log.AuditLogByUser(r.Session.Staff, p.ID, "picking order", "update", "")

		return p, nil
	} else {

		//generate codes for document
		code, _ := util.GenerateDocCode("PIO", r.Warehouse.Code, "picking_order")
		o := orm.NewOrm()
		o.Begin()

		p = &model.PickingOrder{
			Code:            code,
			Warehouse:       r.Warehouse,
			RecognitionDate: r.RecognitionDateTime,
			Status:          1,
			Note:            r.Note,
		}

		if _, e = o.Insert(p); e == nil {
			for _, row := range r.ItemUploadAssign {
				if row.Helper == nil {
					continue
				}

				pl := &model.PickingList{
					ID: 1,
				}

				pickingAssign := &model.PickingOrderAssign{
					SalesOrder:      row.SalesOrder,
					Status:          1,
					Helper:          row.Helper,
					CourierVendor:   row.CourierVendor,
					PlanningVendor:  row.PlanningStr,
					PickingOrder:    p,
					DispatchStatus:  1,
					AssignTimestamp: time.Now(),
					PickingList:     pl,
				}

				if _, e = o.Insert(pickingAssign); e == nil {
					var arrPi []*model.PickingOrderItem
					for _, item := range row.SalesOrderItem {
						pi := &model.PickingOrderItem{
							PickingOrderAssign: pickingAssign,
							Product:            item.Product,
							PickQuantity:       0,
							OrderQuantity:      item.OrderQty,
							FlagOrder:          3,
						}
						arrPi = append(arrPi, pi)
					}

					if _, e := o.InsertMulti(100, &arrPi); e != nil {
						o.Rollback()
						return nil, e
					}

				} else {
					o.Rollback()
					return nil, e
				}

				// flag for picked SO already
				row.SalesOrder.HasPickingAssigned = 1
				if _, e = o.Update(row.SalesOrder, "HasPickingAssigned"); e != nil {
					o.Rollback()
					return nil, e
				}
			}

			e = log.AuditLogByUser(r.Session.Staff, p.ID, "picking order", "create", "")

		} else {
			o.Rollback()
			return nil, e

		}

		o.Commit()

		return p, e
	}
}

func InsertPickingList(code string, dt time.Time, wh *model.Warehouse, so *model.SalesOrder, hp *model.Staff) (e error) {
	o := orm.NewOrm()
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var pl *model.PickingList
	var poa1 *model.PickingOrderAssign
	o1.Raw("SELECT id, code, warehouse_id, delivery_date, status FROM eden_v2.picking_list where code = ?", code).QueryRow(&pl)
	if pl == nil {
		pl = &model.PickingList{
			Code:         code,
			DeliveryDate: dt,
			Warehouse:    wh,
			Status:       1,
		}

		if _, e = o.Insert(pl); e != nil {
			o.Rollback()
			return e
		}
	} else {
		o1.Raw("select poa.id from picking_order_assign poa where poa.sales_order_id = ? and poa.staff_id = ?", so.ID, hp.ID).QueryRow(&poa1)
		poa1.PickingList = pl
		if _, e = o.Update(poa1, "PickingList"); e != nil {
			return e
		}
	}
	return e
}

func RegeneratePickingOrderXls(date time.Time, r uploadAssignRequest, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	dir := util.ExportDirectory

	filename := fmt.Sprintf("TemplatePickingOrder_%s_%s_%s.xlsx", date.Format("2006-01-02"), util.ReplaceSpace(warehouse.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "Business_Type"
		row.AddCell().Value = "Order_Code"
		row.AddCell().Value = "Order_Type"
		row.AddCell().Value = "Merchant_Name"
		row.AddCell().Value = "Order_Status"
		row.AddCell().Value = "Shipping_Address"
		row.AddCell().Value = "Province"
		row.AddCell().Value = "City"
		row.AddCell().Value = "District"
		row.AddCell().Value = "Sub_District"
		row.AddCell().Value = "Postal_Code"
		row.AddCell().Value = "WRT"
		row.AddCell().Value = "Order_Weight"
		row.AddCell().Value = "Delivery_Date"
		row.AddCell().Value = "Payment_Term"
		row.AddCell().Value = "Picker"
		row.AddCell().Value = "Vendor"
		row.AddCell().Value = "Planning"
		row.AddCell().Value = "Error"

		for i, v := range r.ItemUploadAssign {

			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = r.Warehouse.Name     // Warehouse
			row.AddCell().Value = v.BusinessTypeStr    // Business Type
			row.AddCell().Value = v.SalesOrderCode     // Order Code
			row.AddCell().Value = v.OrderTypeStr       // Order Type
			row.AddCell().Value = v.MerchantNameStr    // Merchant Name
			row.AddCell().Value = v.OrderStatusStr     // Order Status
			row.AddCell().Value = v.ShippingAddressStr // Shipping Address
			row.AddCell().Value = v.ProvinceStr        // Province
			row.AddCell().Value = v.CityStr            // City
			row.AddCell().Value = v.DistrictStr        // District
			row.AddCell().Value = v.SubDistrictStr     // Sub District
			row.AddCell().Value = v.PostalCodeStr      // Postal Code
			row.AddCell().Value = v.WRTStr             // WRT
			row.AddCell().Value = v.OrderWeightStr     // Order Weight
			row.AddCell().Value = r.RecognitionDate    // Delivery Date
			row.AddCell().Value = v.PaymentTermStr     // Payment Term
			row.AddCell().Value = v.HelperStr          // Picker
			row.AddCell().Value = v.VendorStr          // Vendor
			row.AddCell().Value = v.PlanningStr        // Planning
			if r.DataCorrection[i] != "" {
				row.AddCell().Value = r.DataCorrection[i] // Error
			}
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

func insertToDeliveryKolies(r []*DeliveryKoli) error {
	o := orm.NewOrm()
	var e error

	var dKoli *model.DeliveryKoli

	var isCreated bool
	var deliveryKoliID int64
	var arrDeliveryKoli []int64
	var qString []string
	var queryString string

	for _, row2 := range r {
		if row2.Quantity == 0 {
			continue
		}

		dKoli = &model.DeliveryKoli{
			SalesOrder: row2.SalesOrder,
			Koli:       row2.Koli,
			Quantity:   row2.Quantity,
			Note:       row2.Note,
		}

		if isCreated, deliveryKoliID, e = o.ReadOrCreate(dKoli, "SalesOrder", "Koli"); e != nil {
			return e
		}
		arrDeliveryKoli = append(arrDeliveryKoli, deliveryKoliID)
		qString = append(qString, "?")
		queryString = strings.Join(qString, ",")

		// update quantity if different with existing delivery koli
		if !isCreated {
			dKoli.Read("ID")
			if dKoli.Quantity != row2.Quantity {
				dKoli.Quantity = row2.Quantity
				if _, e = o.Update(dKoli, "Quantity"); e != nil {
					return e
				}
			}
		}
	}

	o.Raw("delete from delivery_koli where id not in ("+queryString+") and sales_order_id = ?", arrDeliveryKoli, r[0].SalesOrder.ID).Exec()

	return e
}

// InsertPickingListGenerateCode : function to save data requested into database
func InsertPickingListGenerateCode(r generateCodePickingRequest) (p *model.PickingOrder, t []*PickingListObj, e error) {

	o1 := orm.NewOrm()
	o1.Using("read_only")
	o := orm.NewOrm()
	o.Begin()

	for k, v := range r.PickingListFinal {
		var pickingListObj = new(PickingListObj)
		pickingListObj.Code = k
		pickingListObj.TotalWeight = v.TotalWeight
		pickingListObj.SalesOrderID = v.SalesOrderID

		r.PickingListObj = append(r.PickingListObj, pickingListObj)
	}

	var jobs *model.Jobs
	if r.PickingOrder != nil {

		p = r.PickingOrder
		p.Jobs = jobs

	} else {
		code, e := util.GenerateDocCode("PIO", r.Warehouse.Code, "picking_order")

		if e != nil {
			o.Rollback()
			return nil, nil, e
		}
		p = &model.PickingOrder{
			Code:            code,
			Warehouse:       r.Warehouse,
			RecognitionDate: r.DeliveryDateTime,
			Status:          1,
		}

		if _, e := o.Insert(p); e != nil {
			o.Rollback()
			return nil, nil, e
		}

		p.Jobs = jobs
	}
	o.Commit()
	return p, r.PickingListObj, nil
}

func SaveOrderIntoAssign(r generateCodePickingRequest, p *model.PickingOrder) (po *model.PickingOrder, e error) {
	o := orm.NewOrm()
	o.Begin()

	o1 := orm.NewOrm()
	o1.Using("read_only")

	for _, v1 := range r.PickingListObj {
		pl := &model.PickingList{
			Code:         v1.Code,
			DeliveryDate: r.DeliveryDateTime,
			Warehouse:    r.Warehouse,
			Note:         r.Note,
			Status:       1,
		}

		if _, e := o.Insert(pl); e != nil {
			o.Rollback()
			return nil, e
		}
		for _, v2 := range v1.SalesOrderID {
			var salesOrder *model.SalesOrder
			if e = o1.Raw("SELECT id, branch_id, term_payment_sls_id, term_invoice_sls_id, salesperson_id, sales_group_id, sub_district_id, warehouse_id, wrt_id, area_id, voucher_id, price_set_id, payment_group_sls_id, archetype_id, order_type_sls_id, order_channel, code, status, recognition_date, delivery_date, billing_address, shipping_address, shipping_address_note, delivery_fee, vou_redeem_code, vou_disc_amount, point_redeem_amount, point_redeem_id, total_price, total_charge, total_weight, note, reload_packing, payment_reminder, is_locked, has_ext_invoice, has_picking_assigned, cancel_type, created_at, created_by, last_updated_at, last_updated_by, finished_at, locked_by "+
				"FROM eden_v2.sales_order where id = ?", v2).QueryRow(&salesOrder); e != nil {
				o.Rollback()
				return nil, e
			}

			pickingOrderAssign := &model.PickingOrderAssign{
				SalesOrder:     salesOrder,
				Status:         1,
				PickingOrder:   p,
				PickingList:    pl,
				DispatchStatus: 1,
			}

			if _, e = o.Insert(pickingOrderAssign); e != nil {
				o.Rollback()
				return nil, e
			}
			var salesOrderItem []*model.SalesOrderItem
			if _, e = o1.Raw("SELECT id, sales_order_id, product_id, order_qty, forecast_qty, unit_price, shadow_price, subtotal, weight, note "+
				"FROM eden_v2.sales_order_item where sales_order_id = ?", v2).QueryRows(&salesOrderItem); e != nil {
				return nil, e
			}
			var arrPi []*model.PickingOrderItem
			for _, item := range salesOrderItem {
				pi := &model.PickingOrderItem{
					PickingOrderAssign: pickingOrderAssign,
					Product:            item.Product,
					PickQuantity:       0,
					OrderQuantity:      item.OrderQty,
					FlagOrder:          3,
					PickingFlag:        1,
				}
				arrPi = append(arrPi, pi)
			}

			if _, e := o.InsertMulti(100, &arrPi); e != nil {
				o.Rollback()
				return nil, e
			}
		}
	}

	// delete key that blocks picking list creation for associated warehouse
	dbredis.Redis.DeleteCache("picking_list_" + r.WarehouseID)

	o.Commit()
	return
}

func AssignLeadPicker(r assignLeadPickerRequest) (pl *model.PickingList, e error) {
	o := orm.NewOrm()
	o.Begin()

	var pickingOrderAssign map[string]interface{}
	if r.Helper == nil {
		pickingOrderAssign = orm.Params{
			"staff_id":         nil,
			"assign_timestamp": r.AssignTimeStamp,
		}
	} else {
		pickingOrderAssign = orm.Params{
			"staff_id":         r.Helper.ID,
			"assign_timestamp": r.AssignTimeStamp,
		}
	}

	if _, e = o.QueryTable(new(model.PickingOrderAssign)).Filter("picking_list_id", r.PickingList.ID).Update(pickingOrderAssign); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return r.PickingList, e
}

func UpdateBulkQty(r updateBulkQtyRequest) (pl *model.PickingList, e error) {
	o := orm.NewOrm()
	o.Begin()

	pickingList := orm.Params{
		"status": 3,
	}
	if _, e = o.QueryTable(new(model.PickingList)).Filter("id", r.PickingList.ID).Update(pickingList); e != nil {
		o.Rollback()
		return nil, e
	}

	var pickingOrderItem map[string]interface{}
	for _, v := range r.Items {
		if v.SalesOrder.Status == 3 || v.SalesOrder.Status == 4 {
			pickingOrderItem = orm.Params{
				"flag_saved_pick": 1,
				"picking_flag":    2,
			}

			if _, e = o.QueryTable(new(model.PickingOrderItem)).Filter("product_id", r.Product.ID).Filter("picking_order_assign_id", v.PickingOrderAssign.ID).Update(pickingOrderItem); e != nil {
				o.Rollback()
				return nil, e
			}
			continue
		}
		pickingOrderItem = orm.Params{
			"pick_qty":        v.PickOrderQty,
			"unfullfill_note": v.UnfullfillNote,
			"flag_saved_pick": 1,
			"picking_flag":    v.PoiFlagging,
		}

		if _, e = o.QueryTable(new(model.PickingOrderItem)).Filter("product_id", r.Product.ID).Filter("picking_order_assign_id", v.PickingOrderAssign.ID).Update(pickingOrderItem); e != nil {
			o.Rollback()
			return nil, e
		}

		v.PickingOrderAssign.Status = 3

		if _, e = o.Update(v.PickingOrderAssign, "Status"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	o.Commit()
	return r.PickingList, e
}

// PickerAction : function to save data requested into database
func PickerAction(r requestPickerAction) (u *model.PickingRoutingStep, e error) {
	o := orm.NewOrm()
	o.Begin()

	if r.ActionType.ValueName == "start" {
		r.PickingRoutingStep.WalkingStartTime = time.Now()
		r.PickingRoutingStep.WalkingFinishTime = time.Now()
		r.PickingRoutingStep.PickingStartTime = time.Now()
		r.PickingRoutingStep.PickingFinishTime = time.Now()
		r.PickingRoutingStep.StatusStep = 3
		if _, e = o.Update(r.PickingRoutingStep, "WalkingStartTime", "WalkingFinishTime", "PickingStartTime", "PickingFinishTime", "StatusStep"); e != nil {
			o.Rollback()
			return nil, e
		}

		r.NextPickingRoutingStep.WalkingStartTime = time.Now()
		if _, e = o.Update(r.NextPickingRoutingStep, "WalkingStartTime"); e != nil {
			o.Rollback()
			return nil, e
		}
	} else if r.ActionType.ValueName == "pickup" {
		if r.IsPicking == true {
			if r.PickingOrderItem.PickingOrderAssign.SalesOrder.Status == 3 || r.PickingOrderItem.PickingOrderAssign.SalesOrder.Status == 4 {
				r.PickingOrderItem.FlagSavePick = 1
				r.PickingOrderItem.PickingFlag = 2
				if _, e = o.Update(r.PickingOrderItem, "FlagSavePick", "PickingFlag"); e != nil {
					o.Rollback()
					return nil, e
				}
			} else {
				r.PickingOrderItem.PickQuantity = r.PickOrderQty
				r.PickingOrderItem.UnfullfillNote = r.UnfulfillNote
				r.PickingOrderItem.FlagSavePick = 1
				r.PickingOrderItem.PickingFlag = r.PickingFlag
				if _, e = o.Update(r.PickingOrderItem, "PickQuantity", "UnfullfillNote", "FlagSavePick", "PickingFlag"); e != nil {
					o.Rollback()
					return nil, e
				}
			}

			r.PickingRoutingStep.PickingFinishTime = time.Now()
			r.PickingRoutingStep.StatusStep = 3
			if _, e = o.Update(r.PickingRoutingStep, "PickingFinishTime", "StatusStep"); e != nil {
				o.Rollback()
				return nil, e
			}

			r.NextPickingRoutingStep.WalkingStartTime = time.Now()
			if _, e = o.Update(r.NextPickingRoutingStep, "WalkingStartTime"); e != nil {
				o.Rollback()
				return nil, e
			}
		} else {
			r.PickingRoutingStep.WalkingFinishTime = time.Now()
			r.PickingRoutingStep.PickingStartTime = time.Now()
			if _, e = o.Update(r.PickingRoutingStep, "WalkingFinishTime", "PickingStartTime"); e != nil {
				o.Rollback()
				return nil, e
			}
		}
	} else if r.ActionType.ValueName == "delivery" {
		r.PickingRoutingStep.WalkingFinishTime = time.Now()
		r.PickingRoutingStep.PickingStartTime = time.Now()
		r.PickingRoutingStep.PickingFinishTime = time.Now()
		r.PickingRoutingStep.StatusStep = 3
		if _, e = o.Update(r.PickingRoutingStep, "WalkingFinishTime", "PickingStartTime", "PickingFinishTime", "StatusStep"); e != nil {
			o.Rollback()
			return nil, e
		}

		r.NextPickingRoutingStep.WalkingStartTime = time.Now()
		if _, e = o.Update(r.NextPickingRoutingStep, "WalkingStartTime"); e != nil {
			o.Rollback()
			return nil, e
		}
	} else if r.ActionType.ValueName == "end" {
		r.PickingRoutingStep.WalkingFinishTime = time.Now()
		r.PickingRoutingStep.PickingStartTime = time.Now()
		r.PickingRoutingStep.PickingFinishTime = time.Now()
		r.PickingRoutingStep.StatusStep = 3
		if _, e = o.Update(r.PickingRoutingStep, "WalkingFinishTime", "PickingStartTime", "PickingFinishTime", "StatusStep"); e != nil {
			o.Rollback()
			return nil, e
		}
		for _, v := range r.AllPickingRoutingStep {
			v.StatusStep = 4
			if _, e = o.Update(v, "StatusStep"); e != nil {
				o.Rollback()
				return nil, e
			}
		}
	}

	o.Commit()

	return r.PickingRoutingStep, e
}

func PostVroom(r *model.VroomRequest, session *auth.SessionData) (res *model.VroomResponse, err error) {
	var client = &http.Client{}
	vroomUrl := util.VroomUrl

	jsonReq, _ := json.Marshal(r)
	request, err := http.NewRequest("POST", vroomUrl, bytes.NewBuffer(jsonReq))
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		// update row routing to status 4 (failed)
		_ = UpdatePickingListRoutingStatusError(r.Code, err.Error())
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		o := orm.NewOrm()
		o.Using("read_only")
		w := orm.NewOrm()

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(bodyBytes, &res)
		if err != nil {
			return nil, err
		}

		// insert to mongo
		mongoData := model.RoutingMongoModel{
			Code:          r.Code,
			VroomResponse: res,
		}

		m := mongodb.NewMongo()
		m.CreateIndex("picking_list_routing", "code", false)
		ID, err := m.InsertID("picking_list_routing")
		if err != nil {
			m.DisconnectMongoClient()
			return nil, err
		}
		mongoData.ID = ID.(int64)

		_, err = m.InsertOneData("picking_list_routing", mongoData)
		if err != nil {
			m.DisconnectMongoClient()
			return nil, err
		}

		m.DisconnectMongoClient()

		if res.Summary.Unassigned > 0 {
			err = UpdatePickingListRoutingStatus(r.Code, "There's a sales order item that cannot be routed")
			if err != nil {
				return nil, err
			}

		} else {
			w.Begin()
			err = RoutingSuccessStatus(r.Code)
			if err != nil {
				return nil, err
			}

			var pickingRoutingStepArr []*model.PickingRoutingStep
			pickingList := &model.PickingList{Code: r.Code}
			if err = pickingList.Read("code"); err != nil {
				UpdatePickingListRoutingStatus(r.Code, "There's a problem with the picking list")
				return nil, err
			}

			for a := 0; a < len(res.Routes); a++ {
				picker := &model.Staff{ID: res.Routes[a].Vehicle}
				if err = picker.Read("id"); err != nil {
					UpdatePickingListRoutingStatus(r.Code, "There's a problem with the picker")
					return nil, err
				}

				for b := 0; b < len(res.Routes[a].Steps); b++ {
					pickingRoutingStep := &model.PickingRoutingStep{
						PickingList: pickingList,
						Staff:       picker,
						Sequence:    int64(b + 1),
						CreatedBy:   session.Staff,
						CreatedAt:   time.Now(),
						StatusStep:  2,
					}

					if res.Routes[a].Steps[b].Type == "start" {
						pickingRoutingStep.StepType = 1
					} else if res.Routes[a].Steps[b].Type == "pickup" {
						walkingDuration := (res.Routes[a].Steps[b].Duration - res.Routes[a].Steps[b-1].Duration)
						serviceTime := res.Routes[a].Steps[b].Service

						pickingOrderItem := &model.PickingOrderItem{ID: res.Routes[a].Steps[b].ID}
						if err = pickingOrderItem.Read("id"); err != nil {
							UpdatePickingListRoutingStatus(r.Code, "There's a problem with the picking order item")
							return nil, err
						}

						bin := &model.Bin{
							Warehouse: pickingList.Warehouse,
							Product:   pickingOrderItem.Product,
						}
						if err = o.Read(bin, "warehouse", "product"); err != nil {
							UpdatePickingListRoutingStatus(r.Code, "There's a problem with the bin")
							return nil, err
						}

						pickingRoutingStep.PickingOrderItem = pickingOrderItem
						pickingRoutingStep.ExpectedWalkingDuration = walkingDuration
						pickingRoutingStep.ExpectedServiceDuration = serviceTime
						pickingRoutingStep.Bin = bin
						pickingRoutingStep.StepType = 2
					} else if res.Routes[a].Steps[b].Type == "delivery" {
						walkingDuration := (res.Routes[a].Steps[b].Duration - res.Routes[a].Steps[b-1].Duration)
						serviceTime := res.Routes[a].Steps[b].Service

						pickingOrderItem := &model.PickingOrderItem{ID: res.Routes[a].Steps[b].ID}
						if err = pickingOrderItem.Read("id"); err != nil {
							UpdatePickingListRoutingStatus(r.Code, "There's a problem with the picking order item")
							return nil, err
						}

						bin := &model.Bin{
							Warehouse: pickingList.Warehouse,
							Product:   pickingOrderItem.Product,
						}
						if err = o.Read(bin, "warehouse", "product"); err != nil {
							UpdatePickingListRoutingStatus(r.Code, "There's a problem with the bin")
							return nil, err
						}

						pickingRoutingStep.PickingOrderItem = pickingOrderItem
						pickingRoutingStep.ExpectedWalkingDuration = walkingDuration
						pickingRoutingStep.ExpectedServiceDuration = serviceTime
						pickingRoutingStep.Bin = bin
						pickingRoutingStep.StepType = 3
					} else if res.Routes[a].Steps[b].Type == "end" {
						walkingDuration := (res.Routes[a].Steps[b].Duration - res.Routes[a].Steps[b-1].Duration)
						serviceTime := res.Routes[a].Steps[b].Service

						pickingRoutingStep.ExpectedWalkingDuration = walkingDuration
						pickingRoutingStep.ExpectedServiceDuration = serviceTime
						pickingRoutingStep.StepType = 4
					}

					pickingRoutingStepArr = append(pickingRoutingStepArr, pickingRoutingStep)
				}
			}

			if _, err = w.InsertMulti(100, &pickingRoutingStepArr); err != nil {
				w.Rollback()
				return nil, err
			}
			w.Commit()
		}

		return res, err
	} else if response.StatusCode == http.StatusBadRequest || response.StatusCode == http.StatusInternalServerError {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(bodyBytes, &res)
		if err != nil {
			return nil, err
		}

		// update row routing to status 4 (failed)
		err = UpdatePickingListRoutingStatusError(r.Code, res.Error)
		if err != nil {
			return nil, err
		}

		// insert to mongo
		mongoData := model.RoutingMongoModel{
			Code:          r.Code,
			VroomResponse: res,
		}

		m := mongodb.NewMongo()
		m.CreateIndex("picking_list_routing", "code", false)
		ID, err := m.InsertID("picking_list_routing")
		if err != nil {
			m.DisconnectMongoClient()
			return nil, err
		}
		mongoData.ID = ID.(int64)

		_, err = m.InsertOneData("picking_list_routing", mongoData)
		if err != nil {
			m.DisconnectMongoClient()
			return nil, err
		}

		m.DisconnectMongoClient()
	}

	return res, err
}

// UpdatePickingListRoutingStatusError : function to update data requested into database
func UpdatePickingListRoutingStatusError(code string, errorRes ...string) error {
	o := orm.NewOrm()
	o.Begin()
	var err error

	if _, err = o.Raw(
		"UPDATE picking_list SET routing_note = ? WHERE code = ? ", errorRes, code).Exec(); err != nil {
		o.Rollback()
		return err
	}
	o.Commit()

	return nil
}

// UpdatePickingListRoutingStatus : function to update data requested into database
func UpdatePickingListRoutingStatus(code string, reason string) error {
	o := orm.NewOrm()
	o.Begin()
	var err error

	if _, err = o.Raw(
		"UPDATE picking_list SET routing_note = ? WHERE code = ? ", reason, code).Exec(); err != nil {
		o.Rollback()
		return err
	}
	o.Commit()

	return nil
}

// RoutingSuccessStatus : function to update data requested into database
func RoutingSuccessStatus(code string) error {
	o := orm.NewOrm()
	o.Begin()
	var err error

	if _, err = o.Raw(
		"UPDATE picking_list SET routing_note='' WHERE code = ? ", code).Exec(); err != nil {
		o.Rollback()
		return err
	}
	o.Commit()

	return nil
}

// StartAssignRouting : function to save data requested into database
func StartAssignRouting(r startRoutingAssignment) (u *model.PickingList, e error) {
	// not using read only because the data of the picking list is not the latest and the code won't work as how it supposed to be
	w := orm.NewOrm()
	w.Begin()

	bodyReq, e := json.Marshal(r.VroomRequest)
	if e != nil {
		w.Rollback()
		return nil, e
	}
	var client = &http.Client{}
	baseURL := env.GetString("SERVER_HOST", "")

	request, _ := http.NewRequest("POST", "http://"+baseURL+"/v1/warehouse/picking_order/assign/generateroute", bytes.NewBuffer(bodyReq))

	request.Header.Set("Authorization", "Bearer "+r.Session.Token)
	request.Header.Set("Content-Type", "application/json")
	resp, e := client.Do(request)
	if e != nil {
		w.Rollback()
		return nil, e
	}
	defer resp.Body.Close()

	for _, v := range r.PickingOrderAssign {
		v.Status = 3
		v.CheckinTimestamp = time.Now()
		v.SubPicker = r.PickerString
		v.PickerCapacity = r.PickerCapacity / 1000
		if _, e = w.Update(v, "Status", "CheckinTimestamp", "subpicker", "PickerCapacity"); e != nil {
			w.Rollback()
			return nil, e
		}

		e = log.AuditLogByUser(r.Session.Staff, v.ID, "picking order", "started picking routing", r.PickingList.Code)
	}

	r.PickingList.Status = 3
	if _, e = w.Update(r.PickingList, "Status"); e != nil {
		w.Rollback()
		return nil, e
	}

	if e = w.Raw("SELECT * from picking_list where id=?", r.PickingList.ID).QueryRow(&r.PickingList); e != nil {
		w.Rollback()
		return nil, e
	}

	notificationCode, e := repository.GetConfigApp("attribute", "notif_picker")
	if e != nil {
		w.Rollback()
		return nil, e
	}

	notification := model.Notification{
		Code: notificationCode.Value,
	}
	notification.Read("code")

	for _, v := range r.PickerArr {
		staff := &model.Staff{ID: v}
		if e = staff.Read("id"); e != nil {
			w.Rollback()
			return nil, e
		}

		user := &model.User{ID: staff.User.ID}
		if e = user.Read("id"); e != nil {
			w.Rollback()
			return nil, e
		}

		mn := &util.ModelPickingNotification{
			SendTo:    user.PickingNotifToken,
			Title:     notification.Title,
			Message:   notification.Message,
			Type:      "4",
			RefID:     r.PickingList.ID,
			StaffID:   v,
			ServerKey: util.PickingServerKeyFireBase,
		}
		util.PostPickingModelNotification(mn)
	}

	w.Commit()

	return r.PickingList, e
}

// CancelAssignRouting : function to save data requested into database
func CancelAssignRouting(r cancelRoutingAssignment) (u *model.PickingList, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.PickingList.RoutingNote = "Canceled Picking Routing"
	if _, e = o.Update(r.PickingList, "RoutingNote"); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, v := range r.PickingRoutingStep {
		v.StatusStep = 5
		if _, e = o.Update(v, "StatusStep"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	messageNotif := &util.MessageNotification{}
	notificationCode, _ := repository.GetConfigApp("attribute", "cancel_notif_picker")

	if e = o.Raw("SELECT message, title FROM notification WHERE code= ?", notificationCode.Value).QueryRow(&messageNotif); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, v := range r.Picker {
		staff := &model.Staff{ID: v}
		if e = staff.Read("id"); e != nil {
			o.Rollback()
			return nil, e
		}

		user := &model.User{ID: staff.User.ID}
		if e = user.Read("id"); e != nil {
			o.Rollback()
			return nil, e
		}

		mn := &util.ModelPickingNotification{
			SendTo:    user.PickingNotifToken,
			Title:     messageNotif.Title,
			Message:   messageNotif.Message,
			Type:      "4",
			RefID:     r.PickingList.ID,
			StaffID:   v,
			ServerKey: util.PickingServerKeyFireBase,
		}
		util.PostPickingModelNotification(mn)
	}

	o.Commit()

	return r.PickingList, e
}

// StaffUsedInformation : function to read history of the staff with the associated Picking List
func StaffUsedInformation(pickingListID int64, staffs []*model.Staff) (staffsTempered []*model.Staff, e error) {
	o := orm.NewOrm()
	o.Using("read_only")

	var (
		pickerStr       string
		pickerArr       []string
		filter, exclude map[string]interface{}
	)

	mapStaff := map[int64]bool{}

	o.Raw("select sub_picker_id from picking_order_assign poa where poa.picking_list_id = ? LIMIT 1", pickingListID).QueryRow(&pickerStr)
	pickerArr = strings.Split(pickerStr, ",")

	for _, v := range pickerArr {
		pickerID, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		mapStaff[int64(pickerID)] = true
	}

	for _, v := range staffs {
		filter = map[string]interface{}{"staff_id": v.ID, "status_step__in": []int{2, 3}}
		_, total, err := repository.CheckPickingRoutingStepData(filter, exclude)
		if err != nil {
			return
		}

		if total != 0 {
			v.IsBusy = true
		}

		if mapStaff[v.ID] {
			v.UsedStaff = true
		}

	}

	return staffs, nil
}
