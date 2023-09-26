// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock

import (
	"net/http"
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/project-version2/api/log"
)

// SaveStockDO : function to create delivery order data requested into database
func SaveStockDO(r createRequest) (doi []model.DeliveryOrderItem, e error) {
	//generate codes for document
	o := orm.NewOrm()
	o.Begin()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	var sl *model.StockLog
	var doID int64
	for _, row := range r.DeliveryOrderItems {
		e = row.Product.Read("ID")
		if e != nil {
			o.Rollback()
			return nil, e
		}

		doiID, _ := strconv.Atoi(row.ID)
		row.DeliveryOrderItem = &model.DeliveryOrderItem{ID: int64(doiID)}
		e = row.DeliveryOrderItem.Read("ID")
		if e != nil {
			o.Rollback()
			return nil, e
		}

		e = row.DeliveryOrderItem.DeliveryOrder.Read("ID")
		if e != nil {
			o.Rollback()
			return nil, e
		}

		e = row.DeliveryOrderItem.DeliveryOrder.SalesOrder.Read("ID")
		if e != nil {
			o.Rollback()
			return nil, e
		}

		e = row.DeliveryOrderItem.DeliveryOrder.Warehouse.Read("ID")
		if e != nil {
			o.Rollback()
			return nil, e
		}

		if e = row.DeliveryOrderItem.DeliveryOrder.SalesOrder.OrderType.Read("ID"); e != nil {
			o.Rollback()
			return nil, e
		}

		doi = append(doi, *row.DeliveryOrderItem)

		var stock *model.Stock

		e = o.Raw("SELECT * FROM stock where warehouse_id = ? AND product_id = ?", row.DeliveryOrderItem.DeliveryOrder.Warehouse.ID, row.Product.ID).QueryRow(&stock)
		if e != nil {
			o.Rollback()
			return nil, e
		}

		if row.DeliveryOrderItem.DeliveryOrder.SalesOrder.OrderType.Name == "Zero Waste" {
			var wl *model.WasteLog
			wl = &model.WasteLog{
				Warehouse:    row.DeliveryOrderItem.DeliveryOrder.Warehouse,
				Ref:          row.DeliveryOrderItem.DeliveryOrder.ID,
				Product:      row.Product,
				Quantity:     row.DeliverQty,
				RefType:      1,
				Type:         2,
				InitialStock: stock.WasteStock,
				FinalStock:   stock.WasteStock - row.DeliverQty,
				DocNote:      row.DeliveryOrderItem.DeliveryOrder.Note,
				Status:       1,
			}
			doID = wl.Ref
			if _, e = o.Insert(wl); e != nil {
				o.Rollback()
				return nil, e
			}
			st := &model.Stock{
				ID:         stock.ID,
				WasteStock: stock.WasteStock - float64(row.DeliverQty),
			}
			if _, e = o.Update(st, "WasteStock"); e != nil {
				o.Rollback()
				return nil, e
			}
		} else {
			sl = &model.StockLog{
				Warehouse:    row.DeliveryOrderItem.DeliveryOrder.Warehouse,
				Ref:          row.DeliveryOrderItem.DeliveryOrder.ID,
				Product:      row.Product,
				Quantity:     row.DeliverQty,
				RefType:      1,
				Type:         2,
				InitialStock: stock.AvailableStock,
				FinalStock:   stock.AvailableStock - row.DeliverQty,
				UnitCost:     0,
				DocNote:      row.DeliveryOrderItem.DeliveryOrder.Note,
				Status:       1,
				CreatedAt:    time.Now(),
			}
			doID = sl.Ref
			if _, e = o.Insert(sl); e != nil {
				o.Rollback()
				return nil, e
			}
			st := &model.Stock{
				ID:             stock.ID,
				AvailableStock: stock.AvailableStock - float64(row.DeliverQty),
			}
			if _, e = o.Update(st, "AvailableStock"); e != nil {
				o.Rollback()
				return nil, e
			}
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, doID, "delivery_order", "create", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return doi, e
}

// UpdateStockDO : function to update delivery order data requested into database
func UpdateStockDO(r updateRequest) (u *model.DeliveryOrder, e error) {
	o := orm.NewOrm()
	o.Begin()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	u = r.DeliveryOrder
	//UPDATE DO

	if e = r.DeliveryOrder.SalesOrder.Read("ID"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = r.DeliveryOrder.SalesOrder.OrderType.Read("ID"); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, row := range r.DeliveryOrderItems {

		if e = row.Product.Read("ID"); e != nil {
			o.Rollback()
			return nil, e
		}
		doiID, _ := strconv.Atoi(row.ID)
		DOEncryptedID, e := common.Decrypt(doiID)
		if e != nil {
			o.Rollback()
			return nil, e
		}

		row.DeliveryOrderItem = &model.DeliveryOrderItem{ID: int64(DOEncryptedID)}
		if e = row.DeliveryOrderItem.Read("ID"); e != nil {
			o.Rollback()
			return nil, e
		}

		if e = row.DeliveryOrderItem.DeliveryOrder.Read("ID"); e != nil {
			o.Rollback()
			return nil, e
		}

		if e = row.DeliveryOrderItem.DeliveryOrder.Warehouse.Read("ID"); e != nil {
			o.Rollback()
			return nil, e
		}

		stock := &model.Stock{
			Product:   row.Product,
			Warehouse: row.DeliveryOrderItem.DeliveryOrder.Warehouse,
		}
		if e = o.Read(stock, "Product", "Warehouse"); e != nil {
			o.Rollback()
			return nil, e
		}

		if r.DeliveryOrder.SalesOrder.OrderType.Name == "Zero Waste" {
			var wl *model.WasteLog
			if e = o.Raw("SELECT sl.quantity , sl.created_at FROM waste_log sl WHERE ref_id = ? AND warehouse_id = ? "+
				"AND ref_type = 1 AND product_id = ? and `type` = 2 and status = 1 ORDER BY id DESC LIMIT 1", row.DeliveryOrderItem.DeliveryOrder.ID, row.DeliveryOrderItem.DeliveryOrder.Warehouse.ID, row.Product.ID).QueryRow(&wl); e != nil {
				o.Rollback()
				return nil, e
			}

			if wl == nil || stock == nil {
				o.Rollback()
				e = echo.NewHTTPError(http.StatusBadRequest, "Stock not found")
				return nil, e
			}

			wlIn := &model.WasteLog{
				Warehouse:    row.DeliveryOrderItem.DeliveryOrder.Warehouse,
				Product:      row.Product,
				Ref:          row.DeliveryOrderItem.DeliveryOrder.ID,
				RefType:      1,
				Type:         1,
				InitialStock: stock.WasteStock,
				Quantity:     wl.Quantity,
				FinalStock:   stock.WasteStock + wl.Quantity,
				Status:       1,
				DocNote:      "",
				ItemNote:     "",
			}
			if _, e = o.Insert(wlIn); e != nil {
				o.Rollback()
				return nil, e
			}

			wlOut := &model.WasteLog{
				Warehouse:    row.DeliveryOrderItem.DeliveryOrder.Warehouse,
				Product:      row.Product,
				Ref:          row.DeliveryOrderItem.DeliveryOrder.ID,
				RefType:      1,
				Type:         2,
				InitialStock: wlIn.FinalStock,
				Quantity:     row.DeliverQty,
				FinalStock:   wlIn.FinalStock - row.DeliverQty,
				Status:       1,
				DocNote:      row.DeliveryOrderItem.DeliveryOrder.Note,
				ItemNote:     row.Note,
			}
			if _, e = o.Insert(wlOut); e != nil {
				o.Rollback()
				return nil, e
			}

			stock.WasteStock = wlOut.FinalStock

			if _, e = o.Update(stock, "WasteStock"); e != nil {
				o.Rollback()
				return nil, e
			}
		} else {
			var sl *model.StockLog
			if e = o.Raw("SELECT sl.quantity , sl.created_at FROM stock_log sl WHERE ref_id = ? AND warehouse_id = ? "+
				"AND ref_type = 1 AND product_id = ? and `type` = 2 and status = 1 ORDER BY id DESC LIMIT 1", row.DeliveryOrderItem.DeliveryOrder.ID, row.DeliveryOrderItem.DeliveryOrder.Warehouse.ID, row.Product.ID).QueryRow(&sl); e != nil {
				o.Rollback()
				return nil, e
			}

			if sl == nil || stock == nil {
				o.Rollback()
				e = echo.NewHTTPError(http.StatusBadRequest, "Stock not found")
				return nil, e
			}

			slIn := &model.StockLog{
				Warehouse:    row.DeliveryOrderItem.DeliveryOrder.Warehouse,
				Product:      row.Product,
				Ref:          row.DeliveryOrderItem.DeliveryOrder.ID,
				RefType:      1,
				Type:         1,
				InitialStock: stock.AvailableStock,
				Quantity:     sl.Quantity,
				FinalStock:   stock.AvailableStock + sl.Quantity,
				UnitCost:     0,
				Status:       1,
				DocNote:      "",
				ItemNote:     "",
				CreatedAt:    sl.CreatedAt.Add(time.Second * 1),
			}
			if _, e = o.Insert(slIn); e != nil {
				o.Rollback()
				return nil, e
			}

			slOut := &model.StockLog{
				Warehouse:    row.DeliveryOrderItem.DeliveryOrder.Warehouse,
				Product:      row.Product,
				Ref:          row.DeliveryOrderItem.DeliveryOrder.ID,
				RefType:      1,
				Type:         2,
				InitialStock: slIn.FinalStock,
				Quantity:     row.DeliverQty,
				FinalStock:   slIn.FinalStock - row.DeliverQty,
				UnitCost:     row.UnitPrice,
				Status:       1,
				DocNote:      row.DeliveryOrderItem.DeliveryOrder.Note,
				ItemNote:     row.Note,
				CreatedAt:    time.Now(),
			}
			if _, e = o.Insert(slOut); e != nil {
				o.Rollback()
				return nil, e
			}

			stock.AvailableStock = slOut.FinalStock

			if _, e = o.Update(stock, "AvailableStock"); e != nil {
				o.Rollback()
				return nil, e
			}
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.DeliveryOrder.ID, "delivery_order", "update", ""); e != nil {
		o.Rollback()
		return nil, e
	}
	o.Commit()
	return u, e
}

// CancelDOStock : function to cancel delivery order data requested into database
func CancelDOStock(r cancelRequest) (u *model.DeliveryOrder, e error) {
	o := orm.NewOrm()
	o.Begin()
	u = r.DeliveryOrder

	o.LoadRelated(r.DeliveryOrder, "DeliveryOrderItems")

	if e = r.DeliveryOrder.SalesOrder.Read("ID"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = r.DeliveryOrder.SalesOrder.OrderType.Read("ID"); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, v := range r.DeliveryOrder.DeliveryOrderItems {
		r.DeliveryOrderItem = append(r.DeliveryOrderItem, v)

		if r.DeliveryOrder.SalesOrder.OrderType.Name == "Zero Waste" {
			wasteLog := &model.WasteLog{
				Warehouse: r.DeliveryOrder.Warehouse,
				Product:   v.Product,
				Ref:       r.DeliveryOrder.ID,
				Status:    1,
			}

			if e = o.Read(wasteLog, "Warehouse", "Product", "Ref", "Status"); e == nil {
				r.WasteLog = append(r.WasteLog, wasteLog)
			} else {
				o.Rollback()
				return nil, e
			}
		} else {
			stockLog := &model.StockLog{
				Warehouse: r.DeliveryOrder.Warehouse,
				Product:   v.Product,
				Ref:       r.DeliveryOrder.ID,
				Status:    1,
			}

			if e = o.Read(stockLog, "Warehouse", "Product", "Ref", "Status"); e == nil {
				r.StockLog = append(r.StockLog, stockLog)
			} else {
				o.Rollback()
				return nil, e
			}
		}

		stock := &model.Stock{
			Warehouse: r.DeliveryOrder.Warehouse,
			Product:   v.Product,
		}

		if e = o.Read(stock, "Warehouse", "Product"); e == nil {
			r.Stock = append(r.Stock, stock)
		} else {
			o.Rollback()
			return nil, e
		}
	}

	if r.DeliveryOrder.SalesOrder.OrderType.Name == "Zero Waste" {
		for i, v := range r.Stock {
			r.WasteLog[i].Status = 2
			if _, e = o.Update(r.WasteLog[i], "Status"); e != nil {
				o.Rollback()
				return nil, e
			}

			wlIn := &model.WasteLog{
				Warehouse:    r.DeliveryOrder.Warehouse,
				Ref:          r.DeliveryOrder.ID,
				Product:      v.Product,
				Quantity:     r.DeliveryOrderItem[i].DeliverQty,
				RefType:      1,
				Type:         1,
				InitialStock: v.WasteStock,
				FinalStock:   v.WasteStock + r.DeliveryOrderItem[i].DeliverQty,
				DocNote:      r.DeliveryOrder.Note,
				Status:       1,
			}
			if _, e = o.Insert(wlIn); e != nil {
				o.Rollback()
				return nil, e
			}

			v.WasteStock = wlIn.FinalStock
			if _, e = o.Update(v, "WasteStock"); e != nil {
				o.Rollback()
				return nil, e
			}
		}
	} else {
		for i, v := range r.Stock {
			r.StockLog[i].Status = 2
			if _, e = o.Update(r.StockLog[i], "Status"); e != nil {
				o.Rollback()
				return nil, e
			}

			slIn := &model.StockLog{
				Warehouse:    r.DeliveryOrder.Warehouse,
				Ref:          r.DeliveryOrder.ID,
				Product:      v.Product,
				Quantity:     r.DeliveryOrderItem[i].DeliverQty,
				RefType:      1,
				Type:         1,
				InitialStock: v.AvailableStock,
				FinalStock:   v.AvailableStock + r.DeliveryOrderItem[i].DeliverQty,
				UnitCost:     0,
				DocNote:      r.DeliveryOrder.Note,
				Status:       1,
				CreatedAt:    time.Now(),
			}
			if _, e = o.Insert(slIn); e != nil {
				o.Rollback()
				return nil, e
			}

			v.AvailableStock = slIn.FinalStock
			if _, e = o.Update(v, "AvailableStock"); e != nil {
				o.Rollback()
				return nil, e
			}
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, int64(r.DeliveryOrder.ID), "delivery_order", "cancel", r.DeliveryOrder.CancellationNote); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return u, e
}
