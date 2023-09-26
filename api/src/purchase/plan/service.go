// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package plan

import (
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/orm"

	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to insert data requested into database
func Save(r createRequest) (pp *model.PurchasePlan, err error) {
	o := orm.NewOrm()
	o.Begin()

	r.Code, err = util.GenerateDocCode(r.Code, r.SupplierOrganization.Code, "purchase_plan")
	if err != nil {
		o.Rollback()
		return nil, err
	}

	pp = &model.PurchasePlan{
		Code:                 r.Code,
		SupplierOrganization: r.SupplierOrganization,
		Warehouse:            r.Warehouse,
		RecognitionDate:      r.RecognitionDate,
		EtaDate:              r.EtaDate,
		EtaTime:              r.EtaTime,
		TotalPrice:           r.TotalPrice,
		TotalWeight:          r.TotalWeight,
		Note:                 r.Note,
		Status:               1,
		TotalPurchasePlanQty: r.TotalPurchasePlanQty,
		CreatedAt:            time.Now(),
		CreatedBy:            r.Session.Staff,
	}

	if r.FieldPurchaserID != "" {
		pp.AssignedTo = r.FieldPurchaser
		pp.AssignedBy = r.Session.Staff
		pp.AssignedAt = time.Now()
	}

	if _, err = o.Insert(pp); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, pp.ID, "purchase_plan", "create", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	for _, v := range r.PurchasePlanItems {
		ppi := &model.PurchasePlanItem{
			PurchasePlan:    pp,
			Product:         v.Product,
			PurchasePlanQty: v.PurchasePlanQty,
			UnitPrice:       v.UnitPrice,
			Subtotal:        v.Subtotal,
			Weight:          v.PurchasePlanQty * v.Product.UnitWeight,
		}

		if _, err = o.Insert(ppi); err != nil {
			o.Rollback()
			return
		}

		if err = log.AuditLogByUser(r.Session.Staff, ppi.ID, "purchase_plan_item", "create", strconv.FormatFloat(ppi.UnitPrice, 'f', 2, 64)); err != nil {
			o.Rollback()
			return nil, err
		}
	}

	if r.FieldPurchaserID != "" {
		r.MessageNotifPurchaser.Message = util.ReplaceCodeString(r.MessageNotifPurchaser.Message, map[string]interface{}{"#purchase_plan_code#": r.Code})

		mnp := &util.ModelPurchaserNotification{
			SendTo:    r.FieldPurchaser.User.PurchaserNotifToken,
			Title:     r.MessageNotifPurchaser.Title,
			Message:   r.MessageNotifPurchaser.Message,
			Type:      "6",
			RefID:     pp.ID,
			StaffID:   r.FieldPurchaser.ID,
			ServerKey: util.PurchaserServerKeyFireBase,
		}
		util.PostPurchaserModelNotification(mnp)
	}

	for _, v := range r.PurchasingManagers {
		if r.FieldPurchaserID == "" {
			r.MessageNotifManager.Message = util.ReplaceCodeString(r.MessageNotifManager.Message, map[string]interface{}{"#purchase_plan_code#": r.Code})

			v.User.Read("ID")

			mnm := &util.ModelPurchaserNotification{
				SendTo:    v.User.PurchaserNotifToken,
				Title:     r.MessageNotifManager.Title,
				Message:   r.MessageNotifManager.Message,
				Type:      "6",
				RefID:     pp.ID,
				StaffID:   v.ID,
				ServerKey: util.PurchaserServerKeyFireBase,
			}
			util.PostPurchaserModelNotification(mnm)
		} else {
			r.MessageNotifManager.Message = util.ReplaceCodeString(r.MessageNotifManager.Message, map[string]interface{}{"#purchase_plan_code#": r.Code, "#field_purchaser_name#": r.FieldPurchaser.DisplayName})

			v.User.Read("ID")

			mnm := &util.ModelPurchaserNotification{
				SendTo:    v.User.PurchaserNotifToken,
				Title:     r.MessageNotifManager.Title,
				Message:   r.MessageNotifManager.Message,
				Type:      "6",
				RefID:     pp.ID,
				StaffID:   v.ID,
				ServerKey: util.PurchaserServerKeyFireBase,
			}
			util.PostPurchaserModelNotification(mnm)
		}
	}

	o.Commit()
	return pp, err
}

// Assign : function to assign purchase plan to field purchaser
func Assign(r assignRequest) (m *model.PurchasePlan, err error) {
	o := orm.NewOrm()
	o.Begin()

	m = &model.PurchasePlan{
		ID:         r.ID,
		AssignedTo: r.FieldPurchaser,
		AssignedBy: r.Session.Staff,
		AssignedAt: time.Now(),
	}

	if _, err = o.Update(m, "AssignedTo", "AssignedBy", "AssignedAt"); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, m.ID, "purchase_plan", "assign", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	if r.PreviousAssignee != nil {
		if r.FieldPurchaser.ID != r.PreviousAssignee.ID {
			mnpa := &util.ModelPurchaserNotification{
				SendTo:    r.PreviousAssignee.User.PurchaserNotifToken,
				Title:     r.MessageNotifPreviousAssignee.Title,
				Message:   r.MessageNotifPreviousAssignee.Message,
				Type:      "6",
				RefID:     r.PurchasePlan.ID,
				StaffID:   r.PreviousAssignee.ID,
				ServerKey: util.PurchaserServerKeyFireBase,
			}
			util.PostPurchaserModelNotification(mnpa)

			mnna := &util.ModelPurchaserNotification{
				SendTo:    r.FieldPurchaser.User.PurchaserNotifToken,
				Title:     r.MessageNotifNewAssignee.Title,
				Message:   r.MessageNotifNewAssignee.Message,
				Type:      "6",
				RefID:     r.PurchasePlan.ID,
				StaffID:   r.FieldPurchaser.ID,
				ServerKey: util.PurchaserServerKeyFireBase,
			}
			util.PostPurchaserModelNotification(mnna)
		}
	} else {
		mnna := &util.ModelPurchaserNotification{
			SendTo:    r.FieldPurchaser.User.PurchaserNotifToken,
			Title:     r.MessageNotifNewAssignee.Title,
			Message:   r.MessageNotifNewAssignee.Message,
			Type:      "6",
			RefID:     r.PurchasePlan.ID,
			StaffID:   r.FieldPurchaser.ID,
			ServerKey: util.PurchaserServerKeyFireBase,
		}
		util.PostPurchaserModelNotification(mnna)
	}
	o.Commit()

	return m, err
}

// Update : function to update data in database
func Update(r updateRequest) (pp *model.PurchasePlan, err error) {
	o := orm.NewOrm()
	o.Begin()

	var keepItemsId []int64
	var isItemCreated bool

	pp = &model.PurchasePlan{
		ID:                   r.PurchasePlan.ID,
		RecognitionDate:      r.RecognitionDate,
		EtaDate:              r.EtaDate,
		EtaTime:              r.EtaTime,
		TotalPrice:           r.TotalPrice,
		TotalWeight:          r.TotalWeight,
		Note:                 r.Note,
		TotalPurchasePlanQty: r.TotalPurchasePlanQty,
	}

	if r.FieldPurchaserID != "" {
		pp.AssignedTo = r.FieldPurchaser
		pp.AssignedBy = r.PurchasePlan.AssignedBy
		pp.AssignedAt = r.PurchasePlan.AssignedAt

		if r.PurchasePlan.AssignedTo == nil || r.PurchasePlan.AssignedTo.ID != r.FieldPurchaser.ID {
			pp.AssignedBy = r.Session.Staff
			pp.AssignedAt = time.Now()
		}
	} else {
		if r.PurchasePlan.AssignedTo != nil {
			pp.AssignedTo = nil
			pp.AssignedBy = nil
			pp.AssignedAt = time.Time{}
		}
	}

	if _, err = o.Update(pp, "RecognitionDate", "EtaDate", "EtaTime", "TotalPrice", "TotalWeight", "Note", "TotalPurchasePlanQty", "AssignedTo", "AssignedBy", "AssignedAt"); err != nil {
		o.Rollback()
		return nil, err
	}

	for _, v := range r.PurchasePlanItems {

		ppi := &model.PurchasePlanItem{
			PurchasePlan:    &model.PurchasePlan{ID: pp.ID},
			Product:         v.Product,
			PurchasePlanQty: v.PurchasePlanQty,
			UnitPrice:       v.UnitPrice,
			Subtotal:        v.Subtotal,
			Weight:          v.PurchasePlanQty * v.Product.UnitWeight,
		}

		if isItemCreated, ppi.ID, err = o.ReadOrCreate(ppi, "PurchasePlan", "Product"); err != nil {
			o.Rollback()
			return nil, err
		}

		if !isItemCreated {
			ppi.PurchasePlanQty = v.PurchasePlanQty
			ppi.UnitPrice = v.UnitPrice
			ppi.Subtotal = v.Subtotal
			ppi.Weight = v.PurchasePlanQty * v.Product.UnitWeight

			if _, err = o.Update(ppi, "PurchasePlanQty", "UnitPrice", "Subtotal", "Weight"); err != nil {
				o.Rollback()
				return nil, err
			}

		}

		keepItemsId = append(keepItemsId, ppi.ID)
	}

	if _, e := o.QueryTable(new(model.PurchasePlanItem)).Filter("purchase_plan_id", pp.ID).Exclude("ID__in", keepItemsId).Delete(); e != nil {
		o.Rollback()
		return nil, e
	}

	if err = log.AuditLogByUser(r.Session.Staff, pp.ID, "purchase_plan", "update", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	if r.PreviousAssignee != nil {
		if r.FieldPurchaserID == "" {
			mnpa := &util.ModelPurchaserNotification{
				SendTo:    r.PreviousAssignee.User.PurchaserNotifToken,
				Title:     r.MessageNotifPreviousAssignee.Title,
				Message:   r.MessageNotifPreviousAssignee.Message,
				Type:      "6",
				RefID:     pp.ID,
				StaffID:   r.PreviousAssignee.ID,
				ServerKey: util.PurchaserServerKeyFireBase,
			}
			util.PostPurchaserModelNotification(mnpa)

			for _, v := range r.PurchasingManagers {

				v.User.Read("ID")

				mnm := &util.ModelPurchaserNotification{
					SendTo:    v.User.PurchaserNotifToken,
					Title:     r.MessageNotifManager.Title,
					Message:   r.MessageNotifManager.Message,
					Type:      "6",
					RefID:     pp.ID,
					StaffID:   v.ID,
					ServerKey: util.PurchaserServerKeyFireBase,
				}
				util.PostPurchaserModelNotification(mnm)
			}
		} else if r.FieldPurchaser.ID != r.PreviousAssignee.ID {
			mnpa := &util.ModelPurchaserNotification{
				SendTo:    r.PreviousAssignee.User.PurchaserNotifToken,
				Title:     r.MessageNotifPreviousAssignee.Title,
				Message:   r.MessageNotifPreviousAssignee.Message,
				Type:      "6",
				RefID:     r.PurchasePlan.ID,
				StaffID:   r.PreviousAssignee.ID,
				ServerKey: util.PurchaserServerKeyFireBase,
			}
			util.PostPurchaserModelNotification(mnpa)

			mnna := &util.ModelPurchaserNotification{
				SendTo:    r.FieldPurchaser.User.PurchaserNotifToken,
				Title:     r.MessageNotifNewAssignee.Title,
				Message:   r.MessageNotifNewAssignee.Message,
				Type:      "6",
				RefID:     r.PurchasePlan.ID,
				StaffID:   r.FieldPurchaser.ID,
				ServerKey: util.PurchaserServerKeyFireBase,
			}
			util.PostPurchaserModelNotification(mnna)
		}
	} else {
		if r.FieldPurchaserID != "" {
			mnna := &util.ModelPurchaserNotification{
				SendTo:    r.FieldPurchaser.User.PurchaserNotifToken,
				Title:     r.MessageNotifNewAssignee.Title,
				Message:   r.MessageNotifNewAssignee.Message,
				Type:      "6",
				RefID:     r.PurchasePlan.ID,
				StaffID:   r.FieldPurchaser.ID,
				ServerKey: util.PurchaserServerKeyFireBase,
			}
			util.PostPurchaserModelNotification(mnna)
		}
	}

	o.Commit()
	return pp, err
}

// Cancel : function to cancel data
func Cancel(r cancelRequest) (pp *model.PurchasePlan, e error) {
	o := orm.NewOrm()
	o.Begin()

	pp = &model.PurchasePlan{
		ID: r.PurchasePlan.ID,
	}

	pp.Status = 3

	if _, e = o.Update(pp, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, pp.ID, "purchase_plan", "cancel", r.Note); e != nil {
		o.Rollback()
		return nil, e
	}

	if r.PurchasePlan.AssignedTo != nil {

		mnp := &util.ModelPurchaserNotification{
			SendTo:    r.FieldPurchaser.User.PurchaserNotifToken,
			Title:     r.MessageNotifPurchaser.Title,
			Message:   r.MessageNotifPurchaser.Message,
			Type:      "6",
			RefID:     pp.ID,
			StaffID:   r.FieldPurchaser.ID,
			ServerKey: util.PurchaserServerKeyFireBase,
		}
		util.PostPurchaserModelNotification(mnp)
	}

	for _, v := range r.PurchasingManagers {
		v.User.Read("ID")

		mnm := &util.ModelPurchaserNotification{
			SendTo:    v.User.PurchaserNotifToken,
			Title:     r.MessageNotifManager.Title,
			Message:   r.MessageNotifManager.Message,
			Type:      "6",
			RefID:     pp.ID,
			StaffID:   v.ID,
			ServerKey: util.PurchaserServerKeyFireBase,
		}
		util.PostPurchaserModelNotification(mnm)
	}

	o.Commit()

	return r.PurchasePlan, e
}

// Confirm : function to finish data
func Confirm(r confirmRequest) (pp *model.PurchasePlan, e error) {
	o := orm.NewOrm()
	o.Begin()

	pp = &model.PurchasePlan{
		ID: r.PurchasePlan.ID,
	}

	pp.Status = 2

	if _, e = o.Update(pp, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, pp.ID, "purchase_plan", "confirm", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return r.PurchasePlan, e
}

// CancelAssignment : function to cancel purchase plan assignment
func CancelAssignment(r cancelAssignmentRequest) (pp *model.PurchasePlan, e error) {
	o := orm.NewOrm()
	o.Begin()

	pp = &model.PurchasePlan{
		ID: r.PurchasePlan.ID,
	}

	pp.AssignedTo = nil
	pp.AssignedBy = nil
	pp.AssignedAt = time.Time{}

	if _, e = o.Update(pp, "AssignedTo", "AssignedBy", "AssignedAt"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, pp.ID, "purchase_plan", "cancel assignment", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	mnp := &util.ModelPurchaserNotification{
		SendTo:    r.PreviousAssignee.User.PurchaserNotifToken,
		Title:     r.MessageNotifPurchaser.Title,
		Message:   r.MessageNotifPurchaser.Message,
		Type:      "6",
		RefID:     pp.ID,
		StaffID:   r.PreviousAssignee.ID,
		ServerKey: util.PurchaserServerKeyFireBase,
	}
	util.PostPurchaserModelNotification(mnp)

	mnm := &util.ModelPurchaserNotification{
		SendTo:    r.Session.Staff.User.PurchaserNotifToken,
		Title:     r.MessageNotifManager.Title,
		Message:   r.MessageNotifManager.Message,
		Type:      "6",
		RefID:     pp.ID,
		StaffID:   r.Session.Staff.ID,
		ServerKey: util.PurchaserServerKeyFireBase,
	}
	util.PostPurchaserModelNotification(mnm)

	o.Commit()

	return r.PurchasePlan, e
}
