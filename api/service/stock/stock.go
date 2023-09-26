package stock

import (
	"git.edenfarm.id/project-version2/api/util"
	"math"
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// StockLogMaker : save datas when order delivery created
type StockLogMaker struct {
	//Doc            interface{}
	DeliveryOrder  *model.DeliveryOrder
	SalesOrderItem *model.SalesOrderItem
	Product        *model.Product
	Warehouse      *model.Warehouse
	Quantity       float64
	Note           string
}

// DELIVERY ORDER CREATE
func makeItemLogDeliveryOutInsert(doi *model.DeliveryOrderItem) (e error) {
	var doiS []*model.DeliveryOrderItem
	var slm []*StockLogMaker

	doiS = append(doiS, doi)
	if slm, e = ItemStockLog(doiS); e == nil {
		for _, lm := range slm {
			o := orm.NewOrm()
			o.Using("read_only")

			if e =o.Raw("SELECT * from `product` where id = ?", lm.Product.ID).QueryRow(&lm.Product); e != nil{
				util.PostErrorToSentry(e,"getProduct", "get product at function makeItemLogDeliveryOutInsert DeliveryID: "+strconv.FormatInt(lm.DeliveryOrder.ID,10))
			}
			if _,_, e =createItemLog(lm);e !=nil{
				util.PostErrorToSentry(e,"createItemLog", "createItemLog at function makeItemLogDeliveryOutInsert DeliveryID: "+strconv.FormatInt(lm.DeliveryOrder.ID,10))
			}
		}
	}

	return
}

func ItemStockLog(items []*model.DeliveryOrderItem) (lms []*StockLogMaker, e error) {
	// looping setiap item var doi

	for _, row := range items {
		row.DeliveryOrder.Read("ID")
		row.SalesOrderItem.SalesOrder.Read("ID")
		row.SalesOrderItem.SalesOrder.Warehouse.Read("ID")

		l := &StockLogMaker{
			DeliveryOrder:  row.DeliveryOrder,
			Product:        row.Product,
			SalesOrderItem: row.SalesOrderItem,
			Quantity:       row.DeliverQty,
			Warehouse:      row.SalesOrderItem.SalesOrder.Warehouse,
			Note:           row.Note,
		}
		lms = append(lms, l)
	}

	return
}

func createItemLog(lm *StockLogMaker) (sl *model.StockLog, wl *model.WasteLog, e error) {
	if e = lm.Product.Read("ID"); e != nil {
		util.PostErrorToSentry(e,"getProduct", "lm.Product.Read DeliveryID: "+strconv.FormatInt(lm.DeliveryOrder.ID,10))
		return nil, nil, e
	}
	if e = lm.SalesOrderItem.Read("ID"); e != nil {
		util.PostErrorToSentry(e,"getSalesOrderItem", "lm.SalesOrderItem.Read DeliveryID: "+strconv.FormatInt(lm.DeliveryOrder.ID,10))
		return nil, nil, e
	}
	if e = lm.SalesOrderItem.SalesOrder.Read("ID"); e != nil {
		util.PostErrorToSentry(e,"getSalesOrder", "lm.SalesOrderItem.SalesOrder.Read DeliveryID: "+strconv.FormatInt(lm.DeliveryOrder.ID,10))
		return nil, nil, e
	}
	if e = lm.SalesOrderItem.SalesOrder.OrderType.Read("ID"); e != nil {
		util.PostErrorToSentry(e,"getOrderType", "lm.SalesOrderItem.SalesOrder.OrderType.Read DeliveryID: "+strconv.FormatInt(lm.DeliveryOrder.ID,10))
		return nil, nil, e
	}

	var stock *model.Stock
	o := orm.NewOrm()
	o.Using("read_only")

	if e =o.Raw("SELECT * FROM stock where warehouse_id = ? AND product_id = ?", lm.Warehouse.ID, lm.Product.ID).QueryRow(&stock); e!= nil{
		util.PostErrorToSentry(e,"getStock", "getStock at Function createItemLog DeliveryID: "+strconv.FormatInt(lm.DeliveryOrder.ID,10))
	}

	if lm.SalesOrderItem.SalesOrder.OrderType.Name == "Zero Waste" {
		wl = &model.WasteLog{
			Warehouse:    lm.Warehouse,
			Ref:          lm.DeliveryOrder.ID,
			Product:      lm.Product,
			Quantity:     lm.Quantity,
			RefType:      7,
			Type:         2,
			InitialStock: stock.WasteStock,
			FinalStock:   stock.WasteStock - lm.Quantity,
			DocNote:      lm.DeliveryOrder.Note,
			Status:       1,
		}
		if e = wl.Save(); e == nil {
			st := &model.Stock{
				ID:         stock.ID,
				WasteStock: stock.WasteStock - float64(lm.Quantity),
			}
			st.Save("WasteStock")
		}else{
			util.PostErrorToSentry(e,"insertWasteLog", "create WasteLog with DeliveryID: "+strconv.FormatInt(lm.DeliveryOrder.ID,10))
		}
	} else {
		sl = &model.StockLog{
			Warehouse:    lm.Warehouse,
			Ref:          lm.DeliveryOrder.ID,
			Product:      lm.Product,
			Quantity:     lm.Quantity,
			RefType:      1,
			Type:         2,
			InitialStock: stock.AvailableStock,
			FinalStock:   stock.AvailableStock - lm.Quantity,
			UnitCost:     0,
			DocNote:      lm.DeliveryOrder.Note,
			Status:       1,
			CreatedAt:    time.Now(),
		}
		if e = sl.Save(); e == nil {
			st := &model.Stock{
				ID:             stock.ID,
				AvailableStock: stock.AvailableStock - float64(lm.Quantity),
			}
			st.Save("AvailableStock")
		}else{
			util.PostErrorToSentry(e,"insertStockLog", "create StockLog with DeliveryID: "+strconv.FormatInt(lm.DeliveryOrder.ID,10))
		}
	}

	return
}

// STOCK OPNAME COMMITTED
func makeLogStockOpnameCommitted(soi *model.StockOpnameItem) (e error) {
	o := orm.NewOrm()
	o.Begin()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	var stck *model.Stock
	var sltype int8

	if e =	orSelect.Raw("SELECT * FROM stock WHERE warehouse_id = ? AND product_id = ?", soi.StockOpname.Warehouse.ID, soi.Product.ID).QueryRow(&stck); e!= nil{
		util.PostErrorToSentry(e,"getStock", "getStock at Function makeLogStockOpnameCommitted StockOpnameID: "+strconv.FormatInt( soi.StockOpname.ID,10))
	}

	if e = soi.StockOpname.Read("ID"); e != nil {
		util.PostErrorToSentry(e,"getStockOpname", "soi.StockOpname.Read at Function makeLogStockOpnameCommitted")
		return e
	}
	stockType, e := repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_int", soi.StockOpname.StockType)
	if e != nil {
		util.PostErrorToSentry(e,"getstockType", "GetGlossaryMultipleValue at Function makeLogStockOpnameCommitted StockOpnameID: "+strconv.FormatInt( soi.StockOpname.ID,10))
		return e
	}

	if stockType.ValueName == "good stock" {
		stock := &model.Stock{
			ID:             stck.ID,
			AvailableStock: soi.FinalStock,
		}
		if _, e = o.Update(stock, "AvailableStock"); e == nil {

			if soi.AdjustQty > 0 {
				sltype = 1
			} else {
				sltype = 2
			}

			sl := &model.StockLog{
				Warehouse:    soi.StockOpname.Warehouse,
				Ref:          soi.StockOpname.ID,
				Product:      soi.Product,
				Quantity:     math.Abs(soi.AdjustQty),
				RefType:      5,
				Type:         sltype,
				InitialStock: soi.InitialStock,
				FinalStock:   soi.FinalStock,
				UnitCost:     0,
				DocNote:      soi.StockOpname.Note,
				ItemNote:     soi.Note,
				Status:       1,
				CreatedAt:    time.Now(),
			}
			if _, e = o.Insert(sl); e != nil {
				util.PostErrorToSentry(e,"insertStockLog", "update AvailableStock (good stock) at func makeLogStockOpnameCommitted StockOpnameID: "+strconv.FormatInt( soi.StockOpname.ID,10))
				o.Rollback()
			}

		} else {
			util.PostErrorToSentry(e,"updateAvailableStock", "update AvailableStock (good stock) at func makeLogStockOpnameCommitted StockOpnameID: "+strconv.FormatInt( soi.StockOpname.ID,10))
			o.Rollback()
		}
	} else if stockType.ValueName == "waste stock" {
		stock := &model.Stock{
			ID:         stck.ID,
			WasteStock: soi.FinalStock,
		}
		if _, e = o.Update(stock, "WasteStock"); e == nil {

			if soi.AdjustQty > 0 {
				sltype = 1
			} else {
				sltype = 2
			}

			sl := &model.WasteLog{
				Warehouse:    soi.StockOpname.Warehouse,
				Ref:          soi.StockOpname.ID,
				Product:      soi.Product,
				Quantity:     math.Abs(soi.AdjustQty),
				RefType:      8,
				Type:         sltype,
				InitialStock: soi.InitialStock,
				FinalStock:   soi.FinalStock,
				DocNote:      soi.StockOpname.Note,
				ItemNote:     soi.Note,
				Status:       1,
				WasteReason:  soi.OpnameReason,
			}
			if _, e = o.Insert(sl); e != nil {
				util.PostErrorToSentry(e,"insertWasteStock", "insert WasteLog (waste stock) at func makeLogStockOpnameCommitted StockOpnameID: "+strconv.FormatInt( soi.StockOpname.ID,10))
				o.Rollback()
			}

		} else {
			util.PostErrorToSentry(e,"updateWasteStock", "update WasteStock (waste stock) at func makeLogStockOpnameCommitted StockOpnameID: "+strconv.FormatInt( soi.StockOpname.ID,10))
			o.Rollback()
		}
	} else {
		util.PostErrorToSentry(e,"notGoodStockOrWasteStock", "goes to condition Else from waste stock and goods stock at func makeLogStockOpnameCommitted StockOpnameID: "+strconv.FormatInt( soi.StockOpname.ID,10))
		o.Rollback()
	}

	 if e = o.Commit(); e!= nil{
		 util.PostErrorToSentry(e,"FailCommit", "fail commit at func makeLogStockOpnameCommitted StockOpnameID: "+strconv.FormatInt( soi.StockOpname.ID,10))
	 }

	return
}

// WASTE ENTRY COMMITTED
func makeLogWasteEntryCommitted(wei *model.WasteEntryItem) (e error) {
	o := orm.NewOrm()
	o.Begin()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	var stck *model.Stock
	var sl *model.StockLog
	var wl *model.WasteLog

	orSelect.Raw("SELECT * FROM stock WHERE warehouse_id = ? AND product_id = ?", wei.WasteEntry.Warehouse.ID, wei.Product.ID).QueryRow(&stck)

	sl = &model.StockLog{
		Warehouse:    wei.WasteEntry.Warehouse,
		Ref:          wei.WasteEntry.ID,
		Product:      wei.Product,
		Quantity:     wei.WasteStock,
		RefType:      6,
		Type:         2,
		InitialStock: stck.AvailableStock,
		FinalStock:   stck.AvailableStock - wei.WasteStock,
		UnitCost:     0,
		DocNote:      wei.Note,
		Status:       1,
		CreatedAt:    time.Now(),
	}
	if _, e = o.Insert(sl); e == nil {
		wl = &model.WasteLog{
			Warehouse:    wei.WasteEntry.Warehouse,
			Ref:          wei.WasteEntry.ID,
			Product:      wei.Product,
			RefType:      2,
			Type:         2,
			InitialStock: stck.WasteStock,
			Quantity:     wei.WasteStock,
			FinalStock:   stck.WasteStock - wei.WasteStock,
			Status:       1,
			WasteReason:  wei.WasteReason,
		}
		if _, e = o.Insert(wl); e != nil {
			o.Rollback()
		}
	} else {
		o.Rollback()
	}

	o.Commit()

	return
}
