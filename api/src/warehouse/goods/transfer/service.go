// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package transfer

import (
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

func Save(r createRequest) (gt *model.GoodsTransfer, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.Code, e = util.GenerateDocCode("GT", r.WarehouseOrigin.Code, "goods_transfer")
	gt = &model.GoodsTransfer{
		Code:               r.Code,
		Origin:             r.WarehouseOrigin,
		Destination:        r.WarehouseDestination,
		RequestDate:        r.RequestDate,
		AdditionalCost:     r.AdditionalCost,
		AdditionalCostNote: r.AdditionalCostNote,
		StockType:          r.StockTypeID,
		Note:               r.Note,
		Status:             int8(5),
	}

	if _, e = o.Insert(gt); e != nil {
		o.Rollback()
		return nil, e
	}
	for _, row := range r.GoodsTransferItems {
		gti := &model.GoodsTransferItem{
			GoodsTransfer: gt,
			RequestQty:    row.RequestQty,
			Product:       row.Product,
			Note:          row.Note,
		}

		if _, e = o.Insert(gti); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, gt.ID, "goods_transfer", "request", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return
}

func Commit(r commitRequest) (gt *model.GoodsTransfer, e error) {
	o := orm.NewOrm()
	o1 := orm.NewOrm()
	o1.Using("read_only")

	o.Begin()

	var keepItemsId []int64
	var isProductExist = make(map[int64]bool)
	var isItemCreated bool
	var addedProductsGT []*model.GoodsTransferItem

	r.GoodsTransfer.RecognitionDate = r.RecognitionDate
	r.GoodsTransfer.EtaDate = r.EtaDate
	r.GoodsTransfer.EtaTime = r.EtaTimeStr
	r.GoodsTransfer.TotalWeight = r.TotalWeight
	r.GoodsTransfer.Note = r.Note
	r.GoodsTransfer.TotalCost = r.TotalCost
	r.GoodsTransfer.TotalCharge = r.TotalCharge
	r.GoodsTransfer.AdditionalCost = r.AdditionalCost
	r.GoodsTransfer.AdditionalCostNote = r.AdditionalCostNote
	r.GoodsTransfer.Status = 1
	if r.AdditionalCost == 0 {
		r.GoodsTransfer.AdditionalCostNote = ""
	}

	if _, e = o.Update(r.GoodsTransfer, "RecognitionDate", "EtaDate", "EtaTime", "TotalWeight", "Note", "TotalCost", "TotalCharge", "AdditionalCost", "AdditionalCostNote", "Status"); e != nil {
		util.PostErrorToSentry(e,"updateGoodsTransfer", "update goodsTransfer at function commit Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
		o.Rollback()
		return nil, e
	}
	for _, row := range r.GoodsTransferItems {
		gti := &model.GoodsTransferItem{
			GoodsTransfer: r.GoodsTransfer,
			Product:       row.Product,
			DeliverQty:    row.TransferQty,
			UnitCost:      row.UnitCost,
			Subtotal:      row.TransferQty * row.UnitCost,
			Weight:        row.TransferQty * row.Product.UnitWeight,
			Note:          row.Note,
		}

		if isItemCreated, gti.ID, e = o.ReadOrCreate(gti, "GoodsTransfer", "Product"); e == nil {
			if !isItemCreated {
				gti.DeliverQty = row.TransferQty
				gti.UnitCost = row.UnitCost
				gti.Subtotal = row.TransferQty * row.UnitCost
				gti.Weight = row.TransferQty * row.Product.UnitWeight
				gti.Note = row.Note
				if _, e = o.Update(gti, "DeliverQty", "UnitCost", "Subtotal", "Note", "Weight"); e != nil {
					util.PostErrorToSentry(e,"updateGoodsTransferItem", "update goodsTransferItem at function commit Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
					o.Rollback()
					return nil, e
				}
			} else {
				addedProductsGT = append(addedProductsGT, gti)
				isProductExist[row.Product.ID] = true
			}
		} else {
			util.PostErrorToSentry(e,"ReadOrCreate", "ReadOrCreate goodsTransfer at function commit Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
			o.Rollback()
			return nil, e
		}

		keepItemsId = append(keepItemsId, gti.ID)

		// check if the product is added while update Goods Transfer
		if _, exist := isProductExist[row.Product.ID]; exist {
			continue
		}

		if r.StockType.ValueName == "good stock" {
			initialStock := row.Stock.AvailableStock
			row.Stock.AvailableStock = initialStock - row.TransferQty
			if _, e = o.Update(row.Stock, "AvailableStock"); e != nil {
				util.PostErrorToSentry(e,"updateAvailableStock", "update AvailableStock goodsTransfer at function commit Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
				o.Rollback()
				return nil, e
			}

			sl := &model.StockLog{
				Warehouse:    r.GoodsTransfer.Origin,
				Product:      row.Product,
				Ref:          r.GoodsTransfer.ID,
				RefType:      4,
				Type:         2,
				InitialStock: initialStock,
				Quantity:     row.TransferQty,
				FinalStock:   row.Stock.AvailableStock,
				UnitCost:     row.UnitCost,
				Status:       1,
				DocNote:      r.GoodsTransfer.Note,
				ItemNote:     gti.Note,
				CreatedAt:    time.Now(),
			}

			if _, e = o.Insert(sl); e != nil {
				util.PostErrorToSentry(e,"insertStockLog", "insert StockLog goodsTransfer at function commit Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
				o.Rollback()
				return nil, e
			}
		} else if r.StockType.ValueName == "waste stock" {
			initialStock := row.Stock.WasteStock
			row.Stock.WasteStock = initialStock - row.TransferQty
			if _, e = o.Update(row.Stock, "WasteStock"); e != nil {
				util.PostErrorToSentry(e,"updateWasteStock", "update WasteStock goodsTransfer at function commit Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
				o.Rollback()
				return nil, e
			}

			wl := &model.WasteLog{
				Warehouse:    r.GoodsTransfer.Origin,
				Product:      row.Product,
				Ref:          r.GoodsTransfer.ID,
				RefType:      5,
				Type:         2,
				InitialStock: initialStock,
				Quantity:     row.TransferQty,
				FinalStock:   row.Stock.WasteStock,
				UnitCost:     row.UnitCost,
				Status:       1,
				DocNote:      r.GoodsTransfer.Note,
				ItemNote:     gti.Note,
			}

			if _, e = o.Insert(wl); e != nil {
				util.PostErrorToSentry(e,"insertWasteLog", "insert WasteLog goodsTransfer at function commit Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
				o.Rollback()
				return nil, e
			}
		} else {
			util.PostErrorToSentry(e,"notGoodStockOrWasteStock", "goes to condition Else from waste stock and goods stock at func commit Goods Transfer GoodsTransferID: "+strconv.FormatInt( r.GoodsTransfer.ID,10))
			o.Rollback()
			return nil, e
		}
	}

	// get all items deleted
	if _, e = o.QueryTable(new(model.GoodsTransferItem)).Filter("goods_transfer_id", r.GoodsTransfer.ID).Exclude("ID__in", keepItemsId).Delete(); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.GoodsTransfer.ID, "goods_transfer", "commit", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return r.GoodsTransfer, e
}

func Update(r updateRequest) (gt *model.GoodsTransfer, e error) {
	o := orm.NewOrm()
	o1 := orm.NewOrm()
	o1.Using("read_only")

	o.Begin()

	var keepItemsId []int64
	var isProductExist = make(map[int64]bool)
	var isItemCreated bool
	var addedProductsGT []*model.GoodsTransferItem
	var deletedProductsGT []*model.GoodsTransferItem

	r.GoodsTransfer.RecognitionDate = r.RecognitionDate
	r.GoodsTransfer.RequestDate = r.RequestDate
	r.GoodsTransfer.EtaDate = r.EtaDate
	r.GoodsTransfer.EtaTime = r.EtaTimeStr
	r.GoodsTransfer.TotalWeight = r.TotalWeight
	r.GoodsTransfer.Note = r.Note
	r.GoodsTransfer.TotalCost = r.TotalCost
	r.GoodsTransfer.TotalCharge = r.TotalCharge
	r.GoodsTransfer.AdditionalCost = r.AdditionalCost
	r.GoodsTransfer.AdditionalCostNote = r.AdditionalCostNote
	r.GoodsTransfer.UpdatedAt = time.Now()
	r.GoodsTransfer.UpdatedBy = r.Session.Staff.ID

	if r.AdditionalCost == 0 {
		r.GoodsTransfer.AdditionalCostNote = ""
	}

	if _, e = o.Update(r.GoodsTransfer, "RecognitionDate", "RequestDate", "EtaDate", "EtaTime", "TotalWeight", "Note", "TotalCost", "TotalCharge", "AdditionalCost", "AdditionalCostNote", "UpdatedAt", "UpdatedBy"); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, row := range r.GoodsTransferItems {
		gti := &model.GoodsTransferItem{
			GoodsTransfer: r.GoodsTransfer,
			Product:       row.Product,
			DeliverQty:    row.TransferQty,
			UnitCost:      row.UnitCost,
			Subtotal:      row.TransferQty * row.UnitCost,
			Weight:        row.TransferQty * row.Product.UnitWeight,
			Note:          row.Note,
			RequestQty:    row.RequestQty,
		}

		if isItemCreated, gti.ID, e = o.ReadOrCreate(gti, "GoodsTransfer", "Product"); e != nil {
			o.Rollback()
			return nil, e
		}

		if !isItemCreated {
			gti.DeliverQty = row.TransferQty
			gti.UnitCost = row.UnitCost
			gti.Subtotal = row.TransferQty * row.UnitCost
			gti.Weight = row.TransferQty * row.Product.UnitWeight
			gti.Note = row.Note
			gti.RequestQty = row.RequestQty
			if _, e = o.Update(gti, "DeliverQty", "UnitCost", "Subtotal", "Note", "Weight", "RequestQty"); e != nil {
				o.Rollback()
				return nil, e
			}
		} else {
			addedProductsGT = append(addedProductsGT, gti)
			isProductExist[row.Product.ID] = true
		}

		keepItemsId = append(keepItemsId, gti.ID)

		// check if the product is added while update Goods Transfer
		if _, exist := isProductExist[row.Product.ID]; exist {
			continue
		}

		// if qty is not updated, so add record to stock log is unnecessary
		var gtiTemp = new(model.GoodsTransferItem)
		if e = o1.QueryTable(gtiTemp).Filter("goods_transfer_id", r.GoodsTransfer.ID).Filter("product_id", row.Product.ID).One(gtiTemp); e != nil {
			o.Rollback()
			return nil, e
		}

		if gtiTemp.DeliverQty == row.TransferQty {
			continue
		}

		// region if status Goods Transfer is active so added to stock log
		if r.GoodsTransfer.Status == 1 {
			if r.StockType.ValueName == "good stock" {
				var sl *model.StockLog

				if e = o1.Raw("SELECT sl.quantity , sl.created_at FROM stock_log sl WHERE ref_id = ? AND warehouse_id = ? AND ref_type = 4 AND product_id = ? and `type` = 2 and status = 1 ORDER BY id DESC LIMIT 1", r.GoodsTransfer.ID, r.GoodsTransfer.Origin.ID, row.Product.ID).QueryRow(&sl); e != nil {
					util.PostErrorToSentry(e,"readStockLog", "read stockLog at function Update Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
					o.Rollback()
					return nil, e
				}

				s := &model.Stock{
					Product:   row.Product,
					Warehouse: r.GoodsTransfer.Origin,
				}
				s.Read("Product", "Warehouse")

				slIn := &model.StockLog{
					Warehouse:    r.GoodsTransfer.Origin,
					Product:      row.Product,
					Ref:          r.GoodsTransfer.ID,
					RefType:      4,
					Type:         1,
					InitialStock: s.AvailableStock,
					Quantity:     sl.Quantity,
					FinalStock:   s.AvailableStock + sl.Quantity,
					UnitCost:     row.UnitCost,
					Status:       1,
					DocNote:      "",
					ItemNote:     "",
					CreatedAt:    sl.CreatedAt.Add(time.Second * 1),
				}

				if _, e = o.Insert(slIn); e != nil {
					util.PostErrorToSentry(e,"insertStockLog", "insert stockLog `in` at function Update Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
					o.Rollback()
					return nil, e
				}

				slOut := &model.StockLog{
					Warehouse:    r.GoodsTransfer.Origin,
					Product:      row.Product,
					Ref:          r.GoodsTransfer.ID,
					RefType:      4,
					Type:         2,
					InitialStock: slIn.FinalStock,
					Quantity:     row.TransferQty,
					FinalStock:   slIn.FinalStock - row.TransferQty,
					UnitCost:     row.UnitCost,
					Status:       1,
					DocNote:      r.GoodsTransfer.Note,
					ItemNote:     row.Note,
					CreatedAt:    time.Now(),
				}

				if _, e = o.Insert(slOut); e != nil {
					util.PostErrorToSentry(e,"insertStockLog", "insert stockLog `out` at function Update Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
					o.Rollback()
					return nil, e
				}

				s.AvailableStock = slOut.FinalStock

				if _, e = o.Update(s, "AvailableStock"); e != nil {
					util.PostErrorToSentry(e,"updateAvailableStock", "update AvailableStock at function Update Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
					o.Rollback()
					return nil, e
				}
			} else if r.StockType.ValueName == "waste stock" {
				var wl *model.WasteLog

				if e = o1.Raw("SELECT wl.quantity , wl.created_at FROM waste_log wl WHERE ref_id = ? AND warehouse_id = ? AND ref_type = 5 AND product_id = ? and `type` = 2 and status = 1 ORDER BY id DESC LIMIT 1", r.GoodsTransfer.ID, r.GoodsTransfer.Origin.ID, row.Product.ID).QueryRow(&wl); e != nil {
					util.PostErrorToSentry(e,"readwaste_log", "read waste log at function Update Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
					o.Rollback()
					return nil, e
				}

				s := &model.Stock{
					Product:   row.Product,
					Warehouse: r.GoodsTransfer.Origin,
				}
				if e =s.Read("Product", "Warehouse"); e!= nil{
					util.PostErrorToSentry(e,"readStock", "read stock from r.GoodsTransfer.Origin and row.Product at function Update Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
				}

				wlIn := &model.WasteLog{
					Warehouse:    r.GoodsTransfer.Origin,
					Product:      row.Product,
					Ref:          r.GoodsTransfer.ID,
					RefType:      5,
					Type:         1,
					InitialStock: s.WasteStock,
					Quantity:     wl.Quantity,
					FinalStock:   s.WasteStock + wl.Quantity,
					UnitCost:     row.UnitCost,
					Status:       1,
					DocNote:      "",
					ItemNote:     "",
				}

				if _, e = o.Insert(wlIn); e != nil {
					util.PostErrorToSentry(e,"insertWasteStock", "insert WasteStock `in` GoodsTransfer at function Update Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
					o.Rollback()
					return nil, e
				}

				wlOut := &model.WasteLog{
					Warehouse:    r.GoodsTransfer.Origin,
					Product:      row.Product,
					Ref:          r.GoodsTransfer.ID,
					RefType:      5,
					Type:         2,
					InitialStock: wlIn.FinalStock,
					Quantity:     row.TransferQty,
					FinalStock:   wlIn.FinalStock - row.TransferQty,
					UnitCost:     row.UnitCost,
					Status:       1,
					DocNote:      r.GoodsTransfer.Note,
					ItemNote:     row.Note,
				}

				if _, e = o.Insert(wlOut); e != nil {
					util.PostErrorToSentry(e,"insertWasteStock", "insert WasteStock `out` GoodsTransfer at function Update Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
					o.Rollback()
					return nil, e
				}

				s.WasteStock = wlOut.FinalStock

				if _, e = o.Update(s, "WasteStock"); e != nil {
					util.PostErrorToSentry(e,"updateWasteStock", "update WasteStock GoodsTransfer at function Update Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
					o.Rollback()
					return nil, e
				}
			} else {
				util.PostErrorToSentry(e,"notGoodStockOrWasteStock", "goes to condition Else from waste stock and goods stock at func update Goods Transfer GoodsTransferID: "+strconv.FormatInt( r.GoodsTransfer.ID,10))
				o.Rollback()
				return nil, e
			}
		}
	}

	if r.GoodsTransfer.Status == 1 {
		// get all items deleted
		if _, e = o1.QueryTable(new(model.GoodsTransferItem)).Filter("goods_transfer_id", r.GoodsTransfer.ID).Exclude("ID__in", keepItemsId).All(&deletedProductsGT); e != nil {
			o.Rollback()
			return nil, e
		}

		if e = InsertStockLogForModifiedProductInGTI(2, r.GoodsTransfer.Origin, addedProductsGT); e != nil {
			o.Rollback()
			return nil, e
		}

		if e = InsertStockLogForModifiedProductInGTI(1, r.GoodsTransfer.Origin, deletedProductsGT); e != nil {
			o.Rollback()
			return nil, e
		}

		if e = log.AuditLogByUser(r.Session.Staff, r.GoodsTransfer.ID, "goods_transfer", "update", ""); e != nil {
			o.Rollback()
			return nil, e
		}
	} else {
		if e = log.AuditLogByUser(r.Session.Staff, r.GoodsTransfer.ID, "goods_transfer", "update draft", ""); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if _, e = o.QueryTable(new(model.GoodsTransferItem)).Filter("goods_transfer_id", r.GoodsTransfer.ID).Exclude("ID__in", keepItemsId).Delete(); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return r.GoodsTransfer, e
}

func Cancel(r cancelRequest) (gt *model.GoodsTransfer, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.GoodsTransfer.Status = 3

	if e = r.GoodsTransfer.Save("Status"); e != nil {
		util.PostErrorToSentry(e,"cancelGoodsTransfer", "update status cancel GoodsTransfer at function Cancel Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
		o.Rollback()
		return nil, e
	}

	if r.StockType.ValueName == "good stock" {
		for i, v := range r.Stock {
			r.StockLog[i].Status = 2
			if _, e = o.Update(r.StockLog[i], "Status"); e != nil {
				util.PostErrorToSentry(e,"updateStockLog", "update StockLog goods stock status GoodsTransfer at function Cancel Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
				o.Rollback()
				return nil, e
			}

			slIn := &model.StockLog{
				Warehouse:    r.GoodsTransfer.Origin,
				Ref:          r.GoodsTransfer.ID,
				Product:      v.Product,
				Quantity:     r.GoodsTransferItem[i].DeliverQty,
				RefType:      4,
				Type:         1,
				InitialStock: v.AvailableStock,
				FinalStock:   v.AvailableStock + r.GoodsTransferItem[i].DeliverQty,
				UnitCost:     r.GoodsTransferItem[i].UnitCost,
				DocNote:      r.GoodsTransfer.Note,
				Status:       1,
				CreatedAt:    time.Now(),
			}
			if _, e = o.Insert(slIn); e != nil {
				util.PostErrorToSentry(e,"insertStockLog", "insert StockLog goods stock at function Cancel Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
				o.Rollback()
				return nil, e
			}

			v.AvailableStock = slIn.FinalStock
			if _, e = o.Update(v, "AvailableStock"); e != nil {
				util.PostErrorToSentry(e,"updateAvailableStock", "update AvailableStock goods stock at function Cancel Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
				o.Rollback()
				return nil, e
			}
		}
	} else if r.StockType.ValueName == "waste stock" {
		for i, v := range r.Stock {
			r.WasteLog[i].Status = 2
			if _, e = o.Update(r.WasteLog[i], "Status"); e != nil {
				util.PostErrorToSentry(e,"updateWasteLog", "update WasteLog waste stock status GoodsTransfer at function Cancel Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
				o.Rollback()
				return nil, e
			}

			wlIn := &model.WasteLog{
				Warehouse:    r.GoodsTransfer.Origin,
				Ref:          r.GoodsTransfer.ID,
				Product:      v.Product,
				Quantity:     r.GoodsTransferItem[i].DeliverQty,
				RefType:      5,
				Type:         1,
				InitialStock: v.WasteStock,
				FinalStock:   v.WasteStock + r.GoodsTransferItem[i].DeliverQty,
				UnitCost:     r.GoodsTransferItem[i].UnitCost,
				DocNote:      r.GoodsTransfer.Note,
				Status:       1,
			}
			if _, e = o.Insert(wlIn); e != nil {
				util.PostErrorToSentry(e,"InsertWasteLog", "insert WasteLog waste stock status GoodsTransfer at function Cancel Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
				o.Rollback()
				return nil, e
			}

			v.WasteStock = wlIn.FinalStock
			if _, e = o.Update(v, "WasteStock"); e != nil {
				util.PostErrorToSentry(e,"updateStock", "update stock waste stock status GoodsTransfer at function Cancel Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
				o.Rollback()
				return nil, e
			}
		}
	} else {
		util.PostErrorToSentry(e,"notGoodStockOrWasteStock", "goes to condition Else from waste stock and goods stock at func cancel Goods Transfer GoodsTransferID: "+strconv.FormatInt( r.GoodsTransfer.ID,10))
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.GoodsTransfer.ID, "goods_transfer", "cancel", r.Note); e != nil {
		util.PostErrorToSentry(e,"insertAuditLogByUser", "auditLogByUser goods stock at func cancel Goods Transfer GoodsTransferID: "+strconv.FormatInt( r.GoodsTransfer.ID,10))
		o.Rollback()
		return nil, e
	}
	if r.GoodsReceipt != nil {
		r.GoodsReceipt.Status = 3
		if _, e = o.Update(r.GoodsReceipt, "Status"); e != nil {
			util.PostErrorToSentry(e,"updateGoodsReceipt", "update status GoodsReceipt at func cancel Goods Transfer GoodsTransferID: "+strconv.FormatInt( r.GoodsTransfer.ID,10))
			o.Rollback()
			return nil, e
		}
		if e = log.AuditLogByUser(r.Session.Staff, r.GoodsReceipt.ID, "goods_receipt", "cancel", r.Note); e != nil {
			util.PostErrorToSentry(e,"insertAuditLogByUser", "auditLogByUser GoodsReceipt at func cancel Goods Transfer GoodsTransferID: "+strconv.FormatInt( r.GoodsTransfer.ID,10))
			o.Rollback()
			return nil, e
		}

	}

	if e =o.Commit(); e != nil{
		util.PostErrorToSentry(e,"cancelGoodsTransfer", "cancel GoodsTransfer at function cancel Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
		o.Rollback()
	}
	return r.GoodsTransfer, e
}

func Confirm(r confirmRequest) (gt *model.GoodsTransfer, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.GoodsTransfer.Status = 2
	r.GoodsTransfer.AtaDate = r.AtaDate
	r.GoodsTransfer.AtaTime = r.AtaTime
	if _, e = o.Update(r.GoodsTransfer, "Status", "AtaDate", "AtaTime"); e == nil {
		for i, v := range r.GoodsTransferItems {
			v.GoodsTransferItem.ReceiveQty = v.ReceiveQty
			v.GoodsTransferItem.ReceiveNote = v.ReceiveNote
			if _, e = o.Update(v.GoodsTransferItem, "ReceiveQty", "ReceiveNote"); e != nil {
				util.PostErrorToSentry(e,"updateGoodsTransferItem", "update GoodsTransferItem at function Confirm Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
				o.Rollback()
				return nil, e
			}

			if v.ReceiveQty > 0 {
				if r.StockType.ValueName == "good stock" {
					initialStock := r.Stocks[i].AvailableStock
					r.Stocks[i].AvailableStock = initialStock + v.ReceiveQty
					if _, e = o.Update(r.Stocks[i], "AvailableStock"); e != nil {
						util.PostErrorToSentry(e,"updateAvailableStock", "update AvailableStock at function Confirm Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
						o.Rollback()
						return nil, e
					}

					sl := &model.StockLog{
						Warehouse:    r.GoodsTransfer.Destination,
						Product:      v.Product,
						Ref:          r.GoodsTransfer.ID,
						RefType:      4,
						Type:         1,
						InitialStock: initialStock,
						Quantity:     v.ReceiveQty,
						FinalStock:   r.Stocks[i].AvailableStock,
						UnitCost:     v.UnitCost,
						Status:       1,
						DocNote:      r.GoodsTransfer.Note,
						ItemNote:     v.Note,
						CreatedAt:    time.Now(),
					}
					if _, e = o.Insert(sl); e != nil {
						util.PostErrorToSentry(e,"insertStockLog", "insert stock Log at function Confirm Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
						o.Rollback()
						return nil, e
					}
				} else if r.StockType.ValueName == "waste stock" {
					initialStock := r.Stocks[i].WasteStock
					r.Stocks[i].WasteStock = initialStock + v.ReceiveQty
					if _, e = o.Update(r.Stocks[i], "WasteStock"); e != nil {
						util.PostErrorToSentry(e,"updateWasteStock", "update waste stock at function Confirm Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
						o.Rollback()
						return nil, e
					}

					wl := &model.WasteLog{
						Warehouse:    r.GoodsTransfer.Destination,
						Product:      v.Product,
						Ref:          r.GoodsTransfer.ID,
						RefType:      5,
						Type:         1,
						InitialStock: initialStock,
						Quantity:     v.ReceiveQty,
						FinalStock:   r.Stocks[i].WasteStock,
						UnitCost:     v.UnitCost,
						Status:       1,
						DocNote:      r.GoodsTransfer.Note,
						ItemNote:     v.Note,
					}
					if _, e = o.Insert(wl); e != nil {
						util.PostErrorToSentry(e,"insertWasteStock", "insert waste stock at function Confirm Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
						o.Rollback()
						return nil, e
					}
				} else {
					util.PostErrorToSentry(e,"notGoodStockOrWasteStock", "goes to condition Else from waste stock and goods stock at func Confirm Goods Transfer GoodsTransferID: "+strconv.FormatInt( r.GoodsTransfer.ID,10))
					o.Rollback()
					return nil, e
				}
			}
		}

		if e = log.AuditLogByUser(r.Session.Staff, r.GoodsTransfer.ID, "goods_transfer", "confirm", ""); e != nil {
			util.PostErrorToSentry(e,"insertAuditLogByUser", "insert AuditLogByUser GoodsTransfer at function Confirm Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
			o.Rollback()
			return nil, e
		}

	} else {
		util.PostErrorToSentry(e,"updateGoodsTransfer", "update GoodsTransfer at function Confirm Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
		o.Rollback()
		return nil, e
	}

	if e =o.Commit(); e != nil{
		util.PostErrorToSentry(e,"commitGoodsTransfer", "commit GoodsTransfer at function Confirm Goods Transfer GoodsTransferID: "+strconv.FormatInt(r.GoodsTransfer.ID,10))
		o.Rollback()
	}
	return r.GoodsTransfer, e
}

func InsertStockLogForModifiedProductInGTI(typeSl int8, w *model.Warehouse, items []*model.GoodsTransferItem) (e error) {

	o := orm.NewOrm()
	if len(items) == 0 {
		return e
	}

	for _, row := range items {
		stockType, e := repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_int", row.GoodsTransfer.StockType)
		if e != nil {
			return e
		}

		if stockType.ValueName == "good stock" {
			s := &model.Stock{
				Product:   row.Product,
				Warehouse: w,
			}
			s.Read("Product", "Warehouse")

			var finalStock float64
			if typeSl == 1 {
				finalStock = s.AvailableStock + row.DeliverQty
			} else {
				finalStock = s.AvailableStock - row.DeliverQty
			}

			sl := &model.StockLog{
				Warehouse:    w,
				Product:      row.Product,
				Ref:          row.GoodsTransfer.ID,
				RefType:      4,
				Type:         typeSl,
				InitialStock: s.AvailableStock,
				Quantity:     row.DeliverQty,
				FinalStock:   finalStock,
				UnitCost:     row.UnitCost,
				Status:       1,
				DocNote:      row.GoodsTransfer.Note,
				ItemNote:     row.Note,
				CreatedAt:    time.Now(),
			}

			if _, e = o.Insert(sl); e != nil {
				util.PostErrorToSentry(e,"insertStockLog", "insert StockLog GoodsTransfer at function InsertStockLogForModifiedProductInGTI GoodsTransferID: "+strconv.FormatInt(row.GoodsTransfer.ID,10))
				return e
			}

			s.AvailableStock = sl.FinalStock

			if _, e = o.Update(s, "AvailableStock"); e != nil {
				util.PostErrorToSentry(e,"updateAvailableStock", "update AvailableStock GoodsTransfer at function InsertStockLogForModifiedProductInGTI GoodsTransferID: "+strconv.FormatInt(row.GoodsTransfer.ID,10))
				return e
			}
		} else if stockType.ValueName == "waste stock" {
			s := &model.Stock{
				Product:   row.Product,
				Warehouse: w,
			}
			s.Read("Product", "Warehouse")

			var finalStock float64
			if typeSl == 1 {
				finalStock = s.WasteStock + row.DeliverQty
			} else {
				finalStock = s.WasteStock - row.DeliverQty
			}

			wl := &model.WasteLog{
				Warehouse:    w,
				Product:      row.Product,
				Ref:          row.GoodsTransfer.ID,
				RefType:      5,
				Type:         typeSl,
				InitialStock: s.WasteStock,
				Quantity:     row.DeliverQty,
				FinalStock:   finalStock,
				UnitCost:     row.UnitCost,
				Status:       1,
				DocNote:      row.GoodsTransfer.Note,
				ItemNote:     row.Note,
			}

			if _, e = o.Insert(wl); e != nil {
				util.PostErrorToSentry(e,"insertWasteLog", "insert WasteLog GoodsTransfer at function InsertStockLogForModifiedProductInGTI GoodsTransferID: "+strconv.FormatInt(row.GoodsTransfer.ID,10))
				return e
			}

			s.WasteStock = wl.FinalStock

			if _, e = o.Update(s, "WasteStock"); e != nil {
				util.PostErrorToSentry(e,"updateWasteLog", "update WasteStock GoodsTransfer at function InsertStockLogForModifiedProductInGTI GoodsTransferID: "+strconv.FormatInt(row.GoodsTransfer.ID,10))
				return e
			}
		} else {
			return e
		}
	}
	return e

}

// Lock : function to change goods transfer locked into 1
func Lock(r lockRequest) (gt *model.GoodsTransfer, e error) {
	o := orm.NewOrm()
	o.Begin()

	if r.CancelReq != 1 {
		r.GoodsTransfer.Locked = 1
		r.GoodsTransfer.LockedBy = r.Session.Staff.ID
		if _, e = o.Update(r.GoodsTransfer, "Locked", "LockedBy"); e != nil {
			o.Rollback()
			return nil, e
		}
	} else {
		r.GoodsTransfer.Locked = 2
		if _, e = o.Update(r.GoodsTransfer, "Locked"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	o.Commit()
	return r.GoodsTransfer, nil
}
