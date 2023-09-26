// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/cuxs/orm"
	"github.com/tealeg/xlsx"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to insert data requested into database
func Save(r createRequest) (po *model.PurchaseOrder, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.Code, e = util.GenerateDocCode("PO", r.Supplier.Code, "purchase_order")
	if e != nil {
		o.Rollback()
		return nil, e
	}

	po = &model.PurchaseOrder{
		Code:             r.Code,
		Supplier:         r.Supplier,
		Warehouse:        r.Warehouse,
		TermPaymentPur:   r.PurchaseTerm,
		Status:           5,
		RecognitionDate:  r.RecognitionDate,
		EtaDate:          r.EtaDate,
		WarehouseAddress: r.WarehouseAddress,
		EtaTime:          r.EtaTime,
		TaxPct:           r.TaxPct,
		DeliveryFee:      r.DeliveryFee,
		TotalPrice:       r.TotalPrice,
		TaxAmount:        r.TaxAmount,
		TotalCharge:      common.Rounder(r.TotalCharge, 0.5, 2),
		TotalInvoice:     float64(0),
		TotalWeight:      r.TotalWeight,
		Note:             r.Note,
		SupplierBadge:    r.Supplier.SupplierBadge,
		CreatedAt:        time.Now(),
		CreatedBy:        r.Session.Staff,
		CreatedFrom:      r.CreatedFrom,
		HasFinishedGr:    2,
		Locked:           2,
		Latitude:         r.Latitude,
		Longitude:        r.Longitude,
		DeltaPrint:       0,
	}

	if r.PurchasePlanID != "" {
		po.Status = 1
		po.CreatedFrom = 2
		po.PurchasePlan = r.PurchasePlan
		po.TaxPct = 0
		po.CommittedAt = time.Now()
		po.CommittedBy = r.Session.Staff
	}

	if _, e = o.Insert(po); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, po.ID, "purchase_order", "create", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, v := range r.PurchaseOrderItems {
		poi := &model.PurchaseOrderItem{
			PurchaseOrder:     po,
			Product:           v.Product,
			OrderQty:          v.OrderQty,
			UnitPrice:         v.UnitPrice,
			UnitPriceTax:      v.UnitPriceTax,
			IncludeTax:        v.IncludeTax,
			TaxableItem:       v.TaxableItem,
			TaxPercentage:     v.TaxPercentage,
			TaxAmount:         v.TaxAmount,
			Subtotal:          v.Subtotal,
			Weight:            v.OrderQty * v.Product.UnitWeight,
			Note:              v.Note,
			MarketPurchaseStr: "[]",
		}

		if v.PurchasePlanItemID != "" {
			poi.PurchaseQty = v.OrderQty
			poi.IncludeTax = 2
			poi.TaxableItem = 2
			poi.TaxPercentage = 0
			poi.PurchasePlanItem = v.PurchasePlanItem

			v.PurchasePlanItem.PurchaseQty += v.OrderQty
			r.PurchasePlan.TotalPurchaseQty += v.OrderQty

			if _, e = o.Update(v.PurchasePlanItem, "PurchaseQty"); e != nil {
				o.Rollback()
				return nil, e
			}
		}

		if _, e = o.Insert(poi); e != nil {
			o.Rollback()
			return
		}

		if e = log.AuditLogByUser(r.Session.Staff, poi.ID, "purchase_order_item", "create", strconv.FormatFloat(poi.UnitPrice, 'f', 2, 64)); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if r.PurchasePlanID != "" {
		if _, e = o.Update(r.PurchasePlan, "TotalPurchaseQty"); e != nil {
			o.Rollback()
			return nil, e
		}

		if e = log.AuditLogByUser(r.Session.Staff, po.ID, "purchase_order", "commit", "auto commit by field purchaser apps"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if len(r.Images) > 0 {
		for _, v := range r.Images {
			image := &model.PurchaseOrderImage{
				PurchaseOrder: po,
				ImageURL:      v,
				CreatedAt:     time.Now(),
				CreatedBy:     r.Session.Staff.ID,
			}

			if _, e = o.Insert(image); e != nil {
				o.Rollback()
				return
			}
		}
	}

	o.Commit()
	return po, e
}

// Update : function to update data in database
func Update(r updateRequest) (po *model.PurchaseOrder, e error) {
	o := orm.NewOrm()
	o.Begin()

	var keepItemsId []int64
	var isItemCreated bool

	po = &model.PurchaseOrder{
		ID:              r.PurchaseOrder.ID,
		Warehouse:       r.Warehouse,
		RecognitionDate: r.RecognitionDate,
		EtaDate:         r.EtaDate,
		EtaTime:         r.EtaTime,
		TaxPct:          r.TaxPct,
		DeliveryFee:     r.DeliveryFee,
		TotalPrice:      r.TotalPrice,
		TaxAmount:       r.TaxAmount,
		TotalCharge:     common.Rounder(r.TotalCharge, 0.5, 2),
		TotalWeight:     r.TotalWeight,
		Note:            r.Note,
		Status:          5,
	}

	if _, e = o.Update(po, "RecognitionDate", "EtaDate", "EtaTime", "TaxPct", "DeliveryFee", "TotalPrice", "TaxAmount", "TotalCharge", "TotalWeight", "Note"); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, v := range r.PurchaseOrderItems {

		if v.ID != "" {
			poiID := &model.PurchaseOrderItem{PurchaseOrder: r.PurchaseOrder, Product: v.Product}
			poiID.Read("PurchaseOrder", "Product")

			if e = poiID.Read("PurchaseOrder", "Product"); e != nil {
				o.Rollback()
				return nil, e
			}
		}

		poi := &model.PurchaseOrderItem{
			PurchaseOrder:     &model.PurchaseOrder{ID: po.ID},
			Product:           v.Product,
			OrderQty:          v.OrderQty,
			UnitPrice:         v.UnitPrice,
			TaxableItem:       v.TaxableItem,
			IncludeTax:        v.IncludeTax,
			TaxPercentage:     v.TaxPercentage,
			TaxAmount:         v.TaxAmount,
			UnitPriceTax:      v.UnitPriceTax,
			Subtotal:          v.Subtotal,
			Weight:            v.OrderQty * v.Product.UnitWeight,
			Note:              v.Note,
			MarketPurchaseStr: "[]",
		}

		if isItemCreated, poi.ID, e = o.ReadOrCreate(poi, "PurchaseOrder", "Product"); e != nil {
			o.Rollback()
			return nil, e
		}

		if !isItemCreated {
			poi.OrderQty = v.OrderQty
			poi.UnitPrice = v.UnitPrice
			poi.Subtotal = v.Subtotal
			poi.Note = v.Note
			poi.Weight = v.OrderQty * v.Product.UnitWeight
			poi.TaxableItem = v.TaxableItem
			poi.IncludeTax = v.IncludeTax
			poi.TaxPercentage = v.TaxPercentage
			poi.TaxAmount = v.TaxAmount
			poi.UnitPriceTax = v.UnitPriceTax

			if _, e = o.Update(poi, "OrderQty", "UnitPrice", "TaxableItem", "IncludeTax", "TaxPercentage", "TaxAmount", "UnitPriceTax", "Subtotal", "Weight", "Note"); e != nil {
				o.Rollback()
				return nil, e
			}

		}

		keepItemsId = append(keepItemsId, poi.ID)
	}

	if _, e := o.QueryTable(new(model.PurchaseOrderItem)).Filter("purchase_order_id", po.ID).Exclude("ID__in", keepItemsId).Delete(); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, po.ID, "purchase_order", "update", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return po, e
}

// Cancel : function to cancel data
func (r *cancelRequest) Cancel() (po *model.PurchaseOrder, e error) {
	o := orm.NewOrm()
	o.Begin()
	var totalPurchaseQty float64

	po = &model.PurchaseOrder{
		ID: r.PurchaseOrder.ID,
	}

	po.Status = 3

	if _, e = o.Update(po, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, po.ID, "purchase_order", "cancel", r.Note); e != nil {
		o.Rollback()
		return nil, e
	}

	if r.PurchaseOrder.PurchasePlan != nil {
		for _, v := range r.PurchaseOrderItems {
			ppi, err := repository.ValidPurchasePlanItem(v.PurchasePlanItem.ID)
			if err != nil {
				o.Rollback()
				return nil, e
			}

			ppi.PurchaseQty -= v.OrderQty

			if _, e = o.Update(ppi, "PurchaseQty"); e != nil {
				o.Rollback()
				return nil, e
			}

			totalPurchaseQty += v.OrderQty
		}

		r.PurchasePlan.TotalPurchaseQty -= totalPurchaseQty

		if _, e = o.Update(r.PurchasePlan, "TotalPurchaseQty"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	o.Commit()

	return r.PurchaseOrder, e
}

// Commit : function to change status data into active
func (r *commitRequest) Commit() (po *model.PurchaseOrder, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.PurchaseOrder.Status = 1
	r.PurchaseOrder.CommittedAt = time.Now()
	r.PurchaseOrder.CommittedBy = r.Session.Staff

	if _, e = o.Update(r.PurchaseOrder, "Status", "CommittedAt", "CommittedBy"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.PurchaseOrder.ID, "purchase_order", "commit", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return r.PurchaseOrder, e
}

// UpdateProduct : function to change data of product quantity
func UpdateProduct(r updateProductRequest) (po *model.PurchaseOrder, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.PurchaseOrder.TotalPrice = r.TotalPrice
	r.PurchaseOrder.TotalWeight = r.TotalWeight
	r.PurchaseOrder.TotalCharge = r.TotalCharge
	r.PurchaseOrder.UpdatedBy = r.Session.Staff.ID
	r.PurchaseOrder.UpdatedAt = time.Now()

	if _, e = o.Update(r.PurchaseOrder, "TotalPrice", "TotalWeight", "TotalCharge", "UpdatedBy", "UpdatedAt"); e != nil {
		o.Rollback()
		return nil, e
	}
	if r.DebitNote != nil {
		r.DebitNote.TotalPrice = r.TotalPriceDebitNote

		if _, e = o.Update(r.DebitNote, "TotalPrice"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	for _, v := range r.PurchaseOrderItems {
		var unitPrice float64

		if v.TaxAmount == 0 {
			unitPrice = v.UnitPrice
		} else {
			unitPrice = v.UnitPriceTax
		}
		if _, e = o.Update(v, "OrderQty", "UnitPrice", "TaxPercentage", "TaxAmount", "UnitPriceTax", "Subtotal", "Weight"); e != nil {
			o.Rollback()
			return nil, e
		}

		if val, ok := r.MapDebitNoteItem[v.Product.ID]; ok {
			val.UnitPrice = unitPrice
			val.Subtotal = unitPrice * val.ReturnQty

			if _, e = o.Update(val, "UnitPrice", "Subtotal"); e != nil {
				o.Rollback()
				return nil, e
			}
		}

		if v.PurchasePlanItem != nil {
			v.PurchasePlanItem.PurchaseQty += v.OrderQty
			r.PurchasePlan.TotalPurchaseQty += v.OrderQty
			v.PurchaseQty = v.OrderQty

			if _, e = o.Update(v.PurchasePlanItem, "PurchaseQty"); e != nil {
				o.Rollback()
				return nil, e
			}

			if _, e = o.Update(v, "PurchaseQty"); e != nil {
				o.Rollback()
				return nil, e
			}
		}
	}

	if r.PurchasePlan != nil {
		if _, e = o.Update(r.PurchasePlan, "TotalPurchaseQty"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if len(r.Images) > 0 {
		for _, v := range r.Images {
			if v.ID == "" {
				image := &model.PurchaseOrderImage{
					PurchaseOrder: r.PurchaseOrder,
					ImageURL:      v.Url,
					CreatedAt:     time.Now(),
					CreatedBy:     r.Session.Staff.ID,
				}

				if _, e = o.Insert(image); e != nil {
					o.Rollback()
					return
				}
			}
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.PurchaseOrder.ID, "purchase_order", "update product", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return r.PurchaseOrder, e
}

// PrintDeliverySlipXls : function to print purchase delivery slip
func PrintDeliverySlipXls(data *model.PurchaseOrder) (filePath string, err error) {
	var (
		file  *xlsx.File
		sheet *xlsx.Sheet
		row   *xlsx.Row
	)

	dir := env.GetString("EXPORT_DIRECTORY", "")

	filename := fmt.Sprintf("PurchaseDelivery_%s_%s.xlsx", data.Code, util.GenerateRandomDoc(5))
	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "Supplier Name : " + data.Supplier.Name

		row = sheet.AddRow()
		row.AddCell().Value = "Purchase Order Code : " + data.Code

		row = sheet.AddRow()
		row.AddCell().Value = "ETA Date : " + data.EtaDate.Format("02/01/2006")

		row = sheet.AddRow()
		row.AddCell().Value = "ETA Time : " + data.EtaTime

		sheet.AddRow()

		row = sheet.AddRow()
		row.Sheet.SetColWidth(0, 0, 5)
		row.Sheet.SetColWidth(1, 1, 45)
		row.Sheet.SetColWidth(3, 6, 10)
		row.Sheet.SetColWidth(7, 7, 35)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Product Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Order Qty"
		row.AddCell().Value = "Purchase Qty"
		row.AddCell().Value = "Received Qty"
		row.AddCell().Value = "Reject Qty"
		row.AddCell().Value = "Note"

		for index, i := range data.PurchaseOrderItems {
			row = sheet.AddRow()
			row.SetHeight(30)
			row.AddCell().SetInt(index + 1)                         // No
			row.AddCell().Value = i.Product.Name                    // Product Name
			row.AddCell().Value = i.Product.Uom.Name                // UOM
			row.AddCell().SetFloatWithFormat(i.OrderQty, "0.00")    // Order Qty
			row.AddCell().SetFloatWithFormat(i.PurchaseQty, "0.00") // Purchase Qty
		}

		boldStyle := xlsx.NewStyle()
		boldFont := xlsx.NewFont(10, "Liberation Sans")
		boldFont.Bold = true
		boldStyle.Font = *boldFont
		boldStyle.ApplyFont = true

		// looping to get column range 0-7. making BOLD font header
		for col := 0; col < 28; col++ {
			sheet.Cell(5, col).SetStyle(boldStyle)
		}

		err = file.Save(fileDir)
		filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
		// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
		os.Remove(fileDir)

	}

	return
}

// AddMarketPurchase : function to insert market purchase data
func AddMarketPurchase(r marketPurchaseRequest) (po *model.PurchaseOrder, e error) {
	o := orm.NewOrm()
	o.Begin()

	if _, e = o.Update(r.PurchaseOrder, "TotalPrice", "TotalCharge"); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, v := range r.PurchaseOrderItemArr {
		if _, e = o.Update(v, "PurchaseQty", "UnitPrice", "Subtotal", "MarketPurchaseStr"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.PurchaseOrder.ID, "purchase_order", "add market", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return r.PurchaseOrder, e
}

// Assign : function to assign purchase order to field purchaser
func Assign(r assignRequest) (m *model.PurchaseOrder, err error) {
	o := orm.NewOrm()
	o.Begin()

	m = &model.PurchaseOrder{
		ID:         r.ID,
		AssignedTo: r.FieldPurchaser,
		AssignedBy: r.Session.Staff,
		AssignedAt: time.Now(),
	}

	if _, err = o.Update(m, "AssignedTo", "AssignedBy", "AssignedAt"); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, m.ID, "purchase_order", "assign", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	mn := &util.ModelPurchaserNotification{
		SendTo:    r.FieldPurchaser.User.PurchaserNotifToken,
		Title:     r.MessageNotif.Title,
		Message:   r.MessageNotif.Message,
		Type:      "6",
		RefID:     r.PurchaseOrder.ID,
		StaffID:   r.FieldPurchaser.ID,
		ServerKey: util.PurchaserServerKeyFireBase,
	}
	util.PostPurchaserModelNotification(mn)

	o.Commit()

	return m, err
}

// Lock : function to change purchase order locked into 1
func Lock(r lockRequest) (gt *model.PurchaseOrder, e error) {
	o := orm.NewOrm()
	o.Begin()

	if r.CancelReq != 1 {
		r.PurchaseOrder.Locked = 1
		r.PurchaseOrder.LockedBy = r.Session.Staff.ID
		if _, e = o.Update(r.PurchaseOrder, "Locked", "LockedBy"); e != nil {
			o.Rollback()
			return nil, e
		}
	} else {
		r.PurchaseOrder.Locked = 2
		if _, e = o.Update(r.PurchaseOrder, "Locked"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	o.Commit()
	return r.PurchaseOrder, nil
}

// CountPrint : function to count copy of print
func CountPrint(r countPrintRequest) (po *model.PurchaseOrder, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.PurchaseOrder.DeltaPrint += 1
	if _, err := o.Update(r.PurchaseOrder, "DeltaPrint"); err != nil {
		o.Rollback()
		return nil, err
	}

	if err := log.AuditLogByUser(r.Session.Staff, r.PurchaseOrder.ID, "purchase order", "print", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return r.PurchaseOrder, e
}

// Sign : function to sign purchase deliver
func Sign(r signRequest) (po *model.PurchaseOrder, err error) {
	o := orm.NewOrm()
	o.Begin()

	pos := &model.PurchaseOrderSignature{
		PurchaseOrder: r.PurchaseOrder,
		JobFunction:   r.JobFunction,
		Name:          r.Name,
		SignatureURL:  r.SignatureURL,
		CreatedAt:     time.Now(),
		CreatedBy:     r.Session.Staff.ID,
	}

	if _, err := o.Insert(pos); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, r.PurchaseOrder.ID, "purchase order", "signed", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return r.PurchaseOrder, err
}
