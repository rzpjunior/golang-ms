package _return

import (
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (dr *model.DeliveryReturn, e error) {
	//generate codes for document
	r.Code, _ = util.GenerateDocCode("DR", r.Branch.Code, "delivery_return")
	o := orm.NewOrm()
	o.Begin()
	dr = &model.DeliveryReturn{
		Warehouse:       r.Warehouse,
		DeliveryOrder:   r.DeliveryOrder,
		Code:            r.Code,
		Status:          int8(1),
		RecognitionDate: r.RecognitionDateAt,
		Note:            r.Note,
		CreatedBy:       r.Session.Staff.ID,
		CreatedAt:       time.Now(),
	}

	if _, e = o.Insert(dr); e == nil {
		var arrDri []*model.DeliveryReturnItem
		for _, row := range r.DeliveryReturnItems {
			item := &model.DeliveryReturnItem{
				DeliveryReturn:    dr,
				Product:           row.Product,
				DeliveryOrderItem: row.DeliveryOrderItem,
				ReturnWasteQty:    common.Rounder(row.ReturnWasteStockQty, 0.5, 2),
				Note:              row.Note,
			}
			if r.DeliveryOrder.SalesOrder.OrderType.Name != "Zero Waste" {
				item.ReturnGoodQty = common.Rounder(row.ReturnGoodStockQty, 0.5, 2)
				if row.WasteReason != 0 {
					item.WasteReason = row.WasteReason
				}
			}

			arrDri = append(arrDri, item)
		}
		if _, e = o.InsertMulti(100, &arrDri); e == nil {
			e = log.AuditLogByUser(r.Session.Staff, dr.ID, "delivery_return", "create", "")
		} else {
			o.Rollback()
		}
	} else {
		o.Rollback()
	}
	o.Commit()
	return dr, e
}

func Update(r updateRequest) (dr *model.DeliveryReturn, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.DeliveryReturn.RecognitionDate = r.RecognitionDateAt
	r.DeliveryReturn.Note = r.Note

	if e = r.DeliveryReturn.DeliveryOrder.Read("ID"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = r.DeliveryReturn.DeliveryOrder.SalesOrder.Read("ID"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = r.DeliveryReturn.DeliveryOrder.SalesOrder.OrderType.Read("ID"); e != nil {
		o.Rollback()
		return nil, e
	}

	if _, e = o.Update(r.DeliveryReturn, "RecognitionDate", "Note"); e == nil {
		for _, row := range r.DeliveryReturnItems {
			row.DeliveryReturnItem.Product = row.Product
			row.DeliveryReturnItem.DeliveryOrderItem = row.DeliveryOrderItem
			row.DeliveryReturnItem.ReturnWasteQty = row.ReturnWasteStockQty
			row.DeliveryReturnItem.Note = row.Note

			if r.DeliveryReturn.DeliveryOrder.SalesOrder.OrderType.Name != "Zero Waste" {
				row.DeliveryReturnItem.ReturnGoodQty = row.ReturnGoodStockQty
				row.DeliveryReturnItem.WasteReason = row.WasteReason
			}

			if _, e = o.Update(row.DeliveryReturnItem, "Product", "DeliveryOrderItem", "ReturnWasteQty", "ReturnGoodQty", "WasteReason", "Note"); e != nil {
				o.Rollback()
				return nil, e
			}
		}
		//
		e = log.AuditLogByUser(r.Session.Staff, r.DeliveryReturn.ID, "delivery_return", "update", "")
	} else {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return r.DeliveryReturn, e
}

// Cancel: function to change data status into 3
func Cancel(r cancelRequest) (dr *model.DeliveryReturn, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.DeliveryReturn.Status = int8(3)
	if _, e = o.Update(r.DeliveryReturn, "status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, r.DeliveryReturn.ID, "delivery_return", "cancel", r.CancellationNote)
	} else {
		o.Rollback()
	}
	o.Commit()
	return r.DeliveryReturn, nil
}

// Confirm: function to change data status into 2
func Confirm(r confirmRequest) (dr *model.DeliveryReturn, e error) {
	o := orm.NewOrm()
	o.Begin()

	if len(r.StocksAva) > 0 {
		for k, v := range r.StocksAva {
			prevAvailableStock := v.AvailableStock
			returnGoodQty := r.DeliveryReturn.DeliveryReturnItems[k].ReturnGoodQty

			v.AvailableStock = prevAvailableStock + returnGoodQty

			if _, e = o.Update(v, "AvailableStock"); e == nil {
				ls := &model.StockLog{
					Warehouse:    r.DeliveryReturn.Warehouse,
					Product:      v.Product,
					Ref:          r.DeliveryReturn.ID,
					RefType:      2,
					Type:         1,
					InitialStock: prevAvailableStock,
					Quantity:     returnGoodQty,
					FinalStock:   v.AvailableStock,
					UnitCost:     0,
					DocNote:      r.DeliveryReturn.Note,
					ItemNote:     r.DeliveryReturn.DeliveryReturnItems[k].Note,
					Status:       1,
					CreatedAt:    time.Now(),
				}

				if _, e = o.Insert(ls); e != nil {
					o.Rollback()
					return nil, e
				}
			} else {
				o.Rollback()
				return nil, e
			}
		}
	}

	if len(r.StocksWaste) > 0 {
		for k, v := range r.StocksWaste {
			prevWasteStock := v.WasteStock
			returnWasteQty := r.DeliveryReturn.DeliveryReturnItems[k].ReturnWasteQty

			v.WasteStock = prevWasteStock + returnWasteQty

			if _, e = o.Update(v, "WasteStock"); e == nil {

				wl := &model.WasteLog{
					Warehouse:    r.DeliveryReturn.Warehouse,
					Product:      v.Product,
					Ref:          r.DeliveryReturn.ID,
					RefType:      1,
					Type:         1,
					InitialStock: prevWasteStock,
					Quantity:     returnWasteQty,
					FinalStock:   v.WasteStock,
					WasteReason:  r.DeliveryReturn.DeliveryReturnItems[k].WasteReason,
					DocNote:      r.DeliveryReturn.Note,
					ItemNote:     r.DeliveryReturn.DeliveryReturnItems[k].Note,
					Status:       1,
				}

				if _, e = o.Insert(wl); e != nil {
					o.Rollback()
					return nil, e
				}
			} else {
				o.Rollback()
				return nil, e
			}
		}
	}

	r.DeliveryReturn.Status = 2
	r.DeliveryReturn.ConfirmedBy = r.Session.Staff.ID
	r.DeliveryReturn.ConfirmedAt = time.Now()
	if _, e = o.Update(r.DeliveryReturn, "Status", "ConfirmedBy", "ConfirmedAt"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.DeliveryReturn.ID, "delivery_return", "confirm", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return r.DeliveryReturn, nil
}
