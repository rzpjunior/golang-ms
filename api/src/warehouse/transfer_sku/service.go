package transfer_sku

// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

func Save(r createRequest) (resp *createRequest, err error) {
	o := orm.NewOrm()
	o.Begin()

	r.Code, err = util.GenerateDocCode(r.Code, r.Warehouse.Code, "transfer_sku")

	if err != nil {
		o.Rollback()
		return nil, err
	}

	ts := &model.TransferSku{
		Code:             r.Code,
		Warehouse:        r.Warehouse,
		Status:           1,
		CreatedAt:        time.Now(),
		CreatedBy:        r.Session.Staff,
		TotalTransferQty: r.TotalTransferQty,
		TotalWasteQty:    r.TotalWasteQty,
		RecognitionDate:  r.RecognitionDateAt,
		Note:             r.Note,
	}

	if _, err = o.Insert(ts); err != nil {
		o.Rollback()
		return nil, err
	}

	var transferSkuItems []*model.TransferSkuItem

	for _, product := range r.Products {

		for _, transferTo := range product.TransferTo {

			var discrepancy float64
			if product.Product.ID == transferTo.Product.ID {
				discrepancy = product.Discrepancy
			} else {
				discrepancy = 0
			}

			tsi := &model.TransferSkuItem{
				TransferSku:     ts,
				Product:         product.Product,
				TransferProduct: transferTo.Product,
				TransferQty:     transferTo.TransferQty,
				Discrepancy:     discrepancy,
				WasteQty:        transferTo.WasteQty,
				WasteReason:     transferTo.WasteReason,
			}

			transferSkuItems = append(transferSkuItems, tsi)
		}
	}

	if _, err = o.InsertMulti(100, &transferSkuItems); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, ts.ID, "transfer_sku", "create", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()
	return &r, nil
}

// Confirm : To Generate Stock Log & Waste Log for Transfer SKU that created From List Transfer SKU
func Confirm(r confirmRequest) (ts *model.TransferSku, err error) {
	o := orm.NewOrm()
	o.Begin()

	r.TransferSku.Status = 2
	r.TransferSku.ConfirmedAt = time.Now()
	r.TransferSku.ConfirmedBy = r.Session.Staff

	if _, err = o.Update(r.TransferSku, "Status", "ConfirmedAt", "ConfirmedBy"); err != nil {
		o.Rollback()
		return nil, err
	}

	for _, v := range r.TransferSku.TransferSkuItems {
		var (
			stockParent  *model.Stock
			stockChild   *model.Stock
			initialStock float64
			finalStock   float64
			wasteStock   float64
			finalWaste   float64
		)

		// region read stock
		stockParent = &model.Stock{
			Warehouse: r.TransferSku.Warehouse,
			Product:   v.Product,
			Status:    1,
		}

		if err = stockParent.Read("Warehouse", "Product", "Status"); err != nil {
			o.Rollback()
			return nil, err
		}

		if val, ok := r.MapAvailableStock[v.Product.ID]; ok {
			stockParent.AvailableStock = val
		}

		stockChild = &model.Stock{
			Warehouse: r.TransferSku.Warehouse,
			Product:   v.TransferProduct,
			Status:    1,
		}

		if err = stockChild.Read("Warehouse", "Product", "Status"); err != nil {
			o.Rollback()
			return nil, err
		}

		if val2, ok2 := r.MapAvailableStock[v.TransferProduct.ID]; ok2 {
			stockChild.AvailableStock = val2
		}
		// endregion

		if v.WasteQty > 0 {
			wasteStock = stockParent.WasteStock
			finalWaste = stockParent.WasteStock + v.WasteQty
			stockParent.WasteStock = finalWaste
			if _, err = o.Update(stockParent, "WasteStock"); err != nil {
				o.Rollback()
				return nil, err
			}
			wasteLogIn := &model.WasteLog{
				Warehouse:    r.TransferSku.Warehouse,
				Product:      v.Product,
				Ref:          r.TransferSku.ID,
				RefType:      4,
				Type:         1,
				InitialStock: wasteStock,
				Quantity:     v.WasteQty,
				FinalStock:   finalWaste,
				Status:       1,
				WasteReason:  v.WasteReason,
			}

			if _, err = o.Insert(wasteLogIn); err != nil {
				o.Rollback()
				return nil, err
			}
		}

		if r.GrTransferSKU == false {
			if v.Product.ID == v.TransferProduct.ID {
				initialStock = stockParent.AvailableStock
				finalStock = v.TransferQty
				stockParent.AvailableStock = finalStock
				if _, err = o.Update(stockParent, "AvailableStock"); err != nil {
					o.Rollback()
					return nil, err
				}

				stockLogOut := &model.StockLog{
					Warehouse:    r.TransferSku.Warehouse,
					Product:      v.Product,
					Ref:          r.TransferSku.ID,
					RefType:      7,
					Type:         2,
					InitialStock: initialStock,
					Quantity:     initialStock,
					FinalStock:   0,
					UnitCost:     0,
					Status:       1,
					CreatedAt:    time.Now(),
				}

				if _, err = o.Insert(stockLogOut); err != nil {
					o.Rollback()
					return nil, err
				}

				r.MapAvailableStock[v.Product.ID] = finalStock
				initialStock = 0
			}

			if v.Product.ID != v.TransferProduct.ID {
				initialStock = stockChild.AvailableStock
				finalStock = initialStock + v.TransferQty
			}

			stockChild.AvailableStock = finalStock
			if _, err = o.Update(stockChild, "AvailableStock"); err != nil {
				o.Rollback()
				return nil, err
			}

			stockLogIn := &model.StockLog{
				Warehouse:    r.TransferSku.Warehouse,
				Product:      v.TransferProduct,
				Ref:          r.TransferSku.ID,
				RefType:      7,
				Type:         1,
				InitialStock: initialStock,
				Quantity:     v.TransferQty,
				FinalStock:   finalStock,
				UnitCost:     0,
				Status:       1,
				CreatedAt:    time.Now(),
			}

			if _, err = o.Insert(stockLogIn); err != nil {
				o.Rollback()
				return nil, err
			}
			r.MapAvailableStock[v.TransferProduct.ID] = finalStock
		} else {
			if v.Product.ID == v.TransferProduct.ID {
				gr := &model.GoodsReceiptItem{
					GoodsReceipt: r.TransferSku.GoodsReceipt,
					Product:      v.Product,
				}
				if err = o.Read(gr, "GoodsReceipt", "Product"); err != nil {
					o.Rollback()
					return nil, err
				}

				initialStock = stockParent.AvailableStock
				finalStock = initialStock - gr.ReceiveQty
				stockParent.AvailableStock = finalStock + v.TransferQty
				if _, err = o.Update(stockParent, "AvailableStock"); err != nil {
					o.Rollback()
					return nil, err
				}

				stockLogOut := &model.StockLog{
					Warehouse:    r.TransferSku.Warehouse,
					Product:      v.Product,
					Ref:          r.TransferSku.ID,
					RefType:      7,
					Type:         2,
					InitialStock: initialStock,
					Quantity:     gr.ReceiveQty,
					FinalStock:   finalStock,
					UnitCost:     0,
					Status:       1,
					CreatedAt:    time.Now(),
				}

				if _, err = o.Insert(stockLogOut); err != nil {
					o.Rollback()
					return nil, err
				}

				r.MapAvailableStock[v.Product.ID] = finalStock
			}

			if v.Product.ID == v.TransferProduct.ID {
				initialStock = finalStock
				finalStock += v.TransferQty
			} else {
				initialStock = stockChild.AvailableStock
				finalStock = initialStock + v.TransferQty

			}

			stockChild.AvailableStock = finalStock
			if _, err = o.Update(stockChild, "AvailableStock"); err != nil {
				o.Rollback()
				return nil, err
			}

			stockLogIn := &model.StockLog{
				Warehouse:    r.TransferSku.Warehouse,
				Product:      v.TransferProduct,
				Ref:          r.TransferSku.ID,
				RefType:      7,
				Type:         1,
				InitialStock: initialStock,
				Quantity:     v.TransferQty,
				FinalStock:   finalStock,
				UnitCost:     0,
				Status:       1,
				CreatedAt:    time.Now(),
			}

			if _, err = o.Insert(stockLogIn); err != nil {
				o.Rollback()
				return nil, err
			}
			r.MapAvailableStock[v.TransferProduct.ID] = finalStock
		}

		// Cogs
		if v.TransferSku.PurchaseOrder != nil {
			var unitPrice, totalSubtotalGR, actualTotalWeight, unitPricePerKG float64

			var isCreated bool

			purchaseOrderItem := &model.PurchaseOrderItem{}
			if purchaseOrderItem, err = repository.GetPurchaseOrderItemByProduct(v.TransferSku.PurchaseOrder.ID, v.Product.ID); err != nil {
				o.Rollback()
				return nil, err
			}

			unitPrice = purchaseOrderItem.UnitPrice
			// Use unit price tax, if there is price tax in the sku
			if purchaseOrderItem.UnitPriceTax > 0 {
				unitPrice = purchaseOrderItem.UnitPriceTax
			}

			// Recalculate unit price actual if there is transfer to other SKU
			if r.IsTransferToAnotherSku[v.Product.ID] {
				totalSubtotalGR = v.GoodsReceiptQty * unitPrice
				totalGoodStock := v.GoodsReceiptQty - r.TotalQtyWasteProduct[v.Product.ID]
				unitPrice = totalSubtotalGR / totalGoodStock
			}

			// Unit Price Per KG needed to calculate total subtotal per unit weight product
			unitPricePerKG = unitPrice / v.Product.UnitWeight
			actualTotalWeight = v.TransferQty * v.TransferProduct.UnitWeight

			cogs := &model.Cogs{
				Product:       v.TransferProduct,
				Warehouse:     v.TransferSku.Warehouse,
				EtaDate:       r.TransferSku.PurchaseOrder.EtaDate,
				TotalQty:      v.TransferQty,
				TotalSubtotal: actualTotalWeight * unitPricePerKG,
				TotalAvg:      actualTotalWeight * unitPricePerKG / v.TransferQty,
			}

			if isCreated, cogs.ID, err = o.ReadOrCreate(cogs, "Product", "Warehouse", "EtaDate"); err != nil {
				o.Rollback()
				return nil, err
			}
			if !isCreated {
				totalQty := cogs.TotalQty - v.WasteQty

				// Calculate qty and total sub total for child product
				if v.TransferProduct.ID != v.Product.ID {
					totalQty = cogs.TotalQty + v.TransferQty
					subTotalChildSKU := actualTotalWeight * unitPricePerKG
					cogs.TotalSubtotal = cogs.TotalSubtotal + subTotalChildSKU
				}

				// Recalculate totalsubtotal cogs parent product
				if r.IsTransferToAnotherSku[v.TransferProduct.ID] {
					selisihQtyTransfer := r.TotalQtyTransferProduct[v.TransferProduct.ID]
					selisihSubtotalTransfer := selisihQtyTransfer * unitPricePerKG
					cogs.TotalSubtotal = cogs.TotalSubtotal - selisihSubtotalTransfer
					actualTotalWeight = totalQty*v.TransferProduct.UnitWeight - selisihQtyTransfer
					totalQty = actualTotalWeight / v.TransferProduct.UnitWeight
				}

				cogs.TotalQty = totalQty
				cogs.TotalAvg = cogs.TotalSubtotal / totalQty

				if _, err = o.Update(cogs, "TotalQty", "TotalSubtotal", "TotalAvg"); err != nil {
					o.Rollback()
					return nil, err
				}
			}
		}
	}

	if err = log.AuditLogByUser(r.Session.Staff, r.TransferSku.ID, "transfer_sku", "confirm", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()
	return r.TransferSku, nil
}

// Cancel: Function to cancel transfer sku
func Cancel(r cancelRequest) (ts *model.TransferSku, err error) {
	o := orm.NewOrm()
	o.Begin()

	r.TransferSku.Status = 3

	if _, err = o.Update(r.TransferSku, "Status"); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, r.TransferSku.ID, "transfer_sku", "cancel", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()
	return r.TransferSku, nil
}
