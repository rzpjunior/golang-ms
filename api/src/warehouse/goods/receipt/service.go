// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receipt

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

func Save(r createRequest) (gr *model.GoodsReceipt, e error) {
	o := orm.NewOrm()
	o.Begin()

	switch r.InboundType {
	case "purchase_order":
		r.Code, e = util.GenerateDocCode("GR", r.PurchaseOrder.Supplier.Code, "goods_receipt")
		gr = &model.GoodsReceipt{
			Code:                r.Code,
			Warehouse:           r.Warehouse,
			PurchaseOrder:       r.PurchaseOrder,
			AtaDate:             r.AtaDate,
			AtaTime:             r.AtaTimeStr,
			TotalWeight:         r.TotalWeight,
			Note:                r.Note,
			Status:              int8(1),
			ValidSupplierReturn: int8(2),
			CreatedAt:           time.Now(),
			CreatedBy:           r.Session.Staff.ID,
			InboundType:         1,
			StockType:           1,
		}
	case "goods_transfer":
		r.Code, e = util.GenerateDocCode("GR", r.GoodsTransfer.Origin.Code, "goods_receipt")
		gr = &model.GoodsReceipt{
			Code:                r.Code,
			Warehouse:           r.Warehouse,
			GoodsTransfer:       r.GoodsTransfer,
			AtaDate:             r.AtaDate,
			AtaTime:             r.AtaTimeStr,
			TotalWeight:         r.TotalWeight,
			Note:                r.Note,
			Status:              int8(1),
			ValidSupplierReturn: int8(2),
			CreatedAt:           time.Now(),
			CreatedBy:           r.Session.Staff.ID,
			InboundType:         2,
			StockType:           r.GoodsTransfer.StockType,
		}
	}
	if _, e = o.Insert(gr); e != nil {
		o.Rollback()
		return nil, e
	}
	var arrGri []*model.GoodsReceiptItem
	for _, row := range r.GoodsReceiptItems {
		gri := &model.GoodsReceiptItem{
			GoodsReceipt:      gr,
			PurchaseOrderItem: row.PurchaseOrderItem,
			Product:           row.Product,
			DeliverQty:        row.DeliveryQty,
			RejectQty:         row.RejectQty,
			ReceiveQty:        row.ReceiveQty,
			Weight:            row.ReceiveQty * row.Product.UnitWeight,
			Note:              row.Note,
			RejectReason:      row.RejectReason,
		}
		arrGri = append(arrGri, gri)

		if r.InboundType == "goods_transfer" {
			gtiReceiveQty := &model.GoodsTransferItem{
				ID:         row.GoodsTransferItem.ID,
				ReceiveQty: row.ReceiveQty,
			}
			if _, e = o.Update(gtiReceiveQty, "ReceiveQty"); e != nil {
				o.Rollback()
				return nil, e
			}
		}

	}
	if _, e := o.InsertMulti(100, &arrGri); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, gr.ID, "goods_receipt", "create", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return
}

func Update(r updateRequest) (gr *model.GoodsReceipt, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.GoodsReceipt.AtaDate = r.AtaDate
	r.GoodsReceipt.AtaTime = r.AtaTimeStr
	r.GoodsReceipt.TotalWeight = r.TotalWeight
	r.GoodsReceipt.Note = r.Note
	r.GoodsReceipt.UpdatedBy = r.Session.Staff.ID
	r.GoodsReceipt.UpdatedAt = time.Now()

	if _, e = o.Update(r.GoodsReceipt, "AtaDate", "AtaTime", "TotalWeight", "Note", "UpdatedBy", "UpdatedAt"); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, row := range r.GoodsReceiptItems {
		row.GoodsReceiptItem.DeliverQty = row.DeliveryQty
		row.GoodsReceiptItem.RejectQty = row.RejectQty
		row.GoodsReceiptItem.ReceiveQty = row.ReceiveQty
		row.GoodsReceiptItem.Weight = row.ReceiveQty * row.GoodsReceiptItem.Product.UnitWeight
		row.GoodsReceiptItem.Note = row.Note
		row.GoodsReceiptItem.RejectReason = row.RejectReason

		if _, e = o.Update(row.GoodsReceiptItem, "DeliverQty", "RejectQty", "ReceiveQty", "Weight", "Note", "RejectReason"); e != nil {
			o.Rollback()
			return nil, e
		}

		if r.InboundType == "goods_transfer" {
			gtiReceiveQty := &model.GoodsTransferItem{
				ID:         row.GoodsTransferItem.ID,
				ReceiveQty: row.ReceiveQty,
			}
			if _, e = o.Update(gtiReceiveQty, "ReceiveQty"); e != nil {
				o.Rollback()
				return nil, e
			}
		}
	}

	if r.IsAtaChanged {
		e = log.AuditLogByUser(r.Session.Staff, r.GoodsReceipt.ID, "goods_receipt", "update ATA", r.NoteChanged)
	}

	e = log.AuditLogByUser(r.Session.Staff, r.GoodsReceipt.ID, "goods_receipt", "update", "")

	o.Commit()
	return r.GoodsReceipt, e
}

func Cancel(r cancelRequest) (gr *model.GoodsReceipt, e error) {
	r.GoodsReceipt.Status = 3

	if e = r.GoodsReceipt.Save("Status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, r.GoodsReceipt.ID, "goods_receipt", "cancel", r.Note)
	}

	return
}

func Confirm(r confirmRequest) (gr *model.GoodsReceipt, e error) {
	o := orm.NewOrm()
	o.Begin()
	var isCreated bool

	if r.StockType.ValueName == "good stock" {
		for i, v := range r.Stocks {
			prevAvailableStock := v.AvailableStock
			receiveQty := r.GRQty[i].ReceiveQty

			v.AvailableStock = prevAvailableStock + receiveQty

			if _, e = o.Update(v, "AvailableStock"); e != nil {
				o.Rollback()
				return nil, e
			}
			switch r.InboundType {
			case "goods_transfer":
			case "purchase_order":
				var unitPrice float64

				unitPrice = r.GRQty[i].UnitPrice
				// Use unit price tax, if there is price tax in the sku
				if r.GRQty[i].UnitPriceTax > 0 {
					unitPrice = r.GRQty[i].UnitPriceTax
				}

				cogs := &model.Cogs{
					Product:       v.Product,
					Warehouse:     r.GoodsReceipt.Warehouse,
					EtaDate:       r.GoodsReceipt.PurchaseOrder.EtaDate,
					TotalQty:      receiveQty,
					TotalSubtotal: receiveQty * unitPrice,
					TotalAvg:      receiveQty * unitPrice / receiveQty,
				}

				if isCreated, cogs.ID, e = o.ReadOrCreate(cogs, "Product", "Warehouse", "EtaDate"); e != nil {
					o.Rollback()
					return nil, e
				}
				if !isCreated {
					totalQty := cogs.TotalQty + receiveQty
					totalSubtotal := cogs.TotalSubtotal + (receiveQty * unitPrice)

					cogs.TotalQty = totalQty
					cogs.TotalSubtotal = totalSubtotal
					cogs.TotalAvg = totalSubtotal / totalQty

					if _, e = o.Update(cogs, "TotalQty", "TotalSubtotal", "TotalAvg"); e != nil {
						o.Rollback()
						return nil, e
					}
				}
			}

			sl := &model.StockLog{
				Warehouse:    r.GoodsReceipt.Warehouse,
				Product:      v.Product,
				Ref:          r.GoodsReceipt.ID,
				RefType:      3,
				Type:         1,
				InitialStock: prevAvailableStock,
				Quantity:     receiveQty,
				FinalStock:   v.AvailableStock,
				UnitCost:     0,
				DocNote:      r.GoodsReceipt.Note,
				ItemNote:     r.GRQty[i].Note,
				Status:       1,
				CreatedAt:    time.Now(),
			}

			if _, e = o.Insert(sl); e != nil {
				o.Rollback()
				return nil, e
			}
		}

	} else if r.StockType.ValueName == "waste stock" {
		for i, v := range r.Stocks {
			prevWasteStock := v.WasteStock
			receiveQty := r.GRQty[i].ReceiveQty

			v.WasteStock = prevWasteStock + receiveQty

			if _, e = o.Update(v, "WasteStock"); e != nil {
				o.Rollback()
				return nil, e
			}
			switch r.InboundType {
			case "goods_transfer":
			case "purchase_order":
				cogs := &model.Cogs{
					Product:       v.Product,
					Warehouse:     r.GoodsReceipt.Warehouse,
					EtaDate:       r.GoodsReceipt.PurchaseOrder.EtaDate,
					TotalQty:      receiveQty,
					TotalSubtotal: receiveQty * r.GRQty[i].UnitPrice,
					TotalAvg:      receiveQty * r.GRQty[i].UnitPrice / receiveQty,
				}

				if isCreated, cogs.ID, e = o.ReadOrCreate(cogs, "Product", "Warehouse", "EtaDate"); e != nil {
					o.Rollback()
					return nil, e
				}
				if !isCreated {
					totalQty := cogs.TotalQty + receiveQty
					totalSubtotal := cogs.TotalSubtotal + (receiveQty * r.GRQty[i].UnitPrice)

					cogs.TotalQty = totalQty
					cogs.TotalSubtotal = totalSubtotal
					cogs.TotalAvg = totalSubtotal / totalQty

					if _, e = o.Update(cogs, "TotalQty", "TotalSubtotal", "TotalAvg"); e != nil {
						o.Rollback()
						return nil, e
					}
				}
			}

			wl := &model.WasteLog{
				Warehouse:    r.GoodsReceipt.Warehouse,
				Product:      v.Product,
				Ref:          r.GoodsReceipt.ID,
				RefType:      6,
				Type:         1,
				InitialStock: prevWasteStock,
				Quantity:     receiveQty,
				FinalStock:   v.WasteStock,
				UnitCost:     0,
				DocNote:      r.GoodsReceipt.Note,
				ItemNote:     r.GRQty[i].Note,
				Status:       1,
			}

			if _, e = o.Insert(wl); e != nil {
				o.Rollback()
				return nil, e
			}
		}

	} else {
		o.Rollback()
		return nil, e
	}

	r.GoodsReceipt.Status = 2
	r.GoodsReceipt.ConfirmedAt = time.Now()
	r.GoodsReceipt.ConfirmedBy = r.Session.Staff.ID
	if _, e = o.Update(r.GoodsReceipt, "Status", "ConfirmedBy", "ConfirmedAt"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.GoodsReceipt.ID, "goods_receipt", "confirm", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	switch r.InboundType {
	case "goods_transfer":
		r.GoodsReceipt.GoodsTransfer.Status = 2
		r.GoodsReceipt.GoodsTransfer.AtaDate = r.GoodsReceipt.AtaDate
		r.GoodsReceipt.GoodsTransfer.AtaTime = r.GoodsReceipt.AtaTime
		if _, e = o.Update(r.GoodsReceipt.GoodsTransfer, "Status", "AtaDate", "AtaTime"); e != nil {
			o.Rollback()
			return nil, e
		}
		if e = log.AuditLogByUser(r.Session.Staff, r.GoodsReceipt.GoodsTransfer.ID, "goods_transfer", "confirm", ""); e != nil {
			o.Rollback()
			return nil, e
		}
	case "purchase_order":
		r.GoodsReceipt.PurchaseOrder.HasFinishedGr = 1
		if r.PurchaseInvoice.ID != 0 && r.PurchaseInvoice.Status == 2 {
			r.GoodsReceipt.PurchaseOrder.Status = 2
		}

		if _, e = o.Update(r.GoodsReceipt.PurchaseOrder, "Status", "HasFinishedGr"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	o.Commit()
	return r.GoodsReceipt, nil
}

// Lock : function to change goods transfer locked into 1
func Lock(r lockRequest) (gt *model.GoodsReceipt, e error) {
	o := orm.NewOrm()
	o.Begin()

	if r.CancelReq != 1 {
		r.GoodsReceipt.Locked = 1
		r.GoodsReceipt.LockedBy = r.Session.Staff.ID
		if _, e = o.Update(r.GoodsReceipt, "Locked", "LockedBy"); e != nil {
			o.Rollback()
			return nil, e
		}
	} else {
		r.GoodsReceipt.Locked = 2
		if _, e = o.Update(r.GoodsReceipt, "Locked"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	o.Commit()
	return r.GoodsReceipt, nil
}
