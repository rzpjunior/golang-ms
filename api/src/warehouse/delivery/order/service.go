// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"net/http"
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs/event"
	"git.edenfarm.id/cuxs/orm"
	"github.com/labstack/echo/v4"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (u *model.DeliveryOrder, e error) {
	//generate codes for document
	r.SalesOrder.Branch.Read("ID")
	codeDO, _ := util.GenerateDocCode("DO", r.SalesOrder.Branch.Code, "delivery_order")
	checkStatus, e := repository.GetConfigApp("attribute", "kafka_mongo_delivery_order")
	if e != nil {
		return nil, e
	}
	o := orm.NewOrm()
	o.Begin()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	var so *model.SalesOrder

	if e == nil {
		u = &model.DeliveryOrder{
			Code:            codeDO,
			SalesOrder:      r.SalesOrder,
			Warehouse:       r.Warehouse,
			Wrt:             r.Wrt,
			Status:          1,
			RecognitionDate: r.OrderDate,
			ShippingAddress: r.ShippingAddress,
			ReceiptNote:     "",
			TotalWeight:     r.TotalWeight,
			DeltaPrint:      0,
			Note:            r.Note,
			CreatedAt:       time.Now(),
			CreatedBy:       r.Session.Staff.ID,
		}
		if e = r.SalesOrder.Branch.Read("ID"); e != nil {
			o.Rollback()
			return nil, e
		}
		if e = r.SalesOrder.Branch.Merchant.Read("ID"); e != nil {
			o.Rollback()
			return nil, e
		}
		if e = r.SalesOrder.Branch.Merchant.UserMerchant.Read("ID"); e != nil {
			o.Rollback()
			return nil, e
		}
		if r.SalesOrder.Branch.Merchant.InvoiceTerm.ID == 2 && r.SalesOrder.Status != 9 &&
			r.SalesOrder.Status != 10 && r.SalesOrder.Status != 11 &&
			r.SalesOrder.Status != 12 && r.SalesOrder.Status != 13 {
			u.TotalWeight = r.TotalWeightDirect
		}

		if _, e = o.Insert(u); e == nil {

			for i, row := range r.DeliveryOrderItems {
				// read soi
				soi := &model.SalesOrderItem{SalesOrder: u.SalesOrder, Product: row.Product}
				soi.Read("SalesOrder", "Product")

				item := &model.DeliveryOrderItem{
					DeliveryOrder:  &model.DeliveryOrder{ID: u.ID},
					SalesOrderItem: soi,
					Product:        row.Product,
					DeliverQty:     row.DeliverQty,
					Weight:         row.Weight,
					OrderItemNote:  row.Note,
				}

				if _, e = o.Insert(item); e == nil {
					u.DeliveryOrderItems = append(u.DeliveryOrderItems, item)
					r.DeliveryOrderItems[i].ID = strconv.Itoa(int(item.ID))
				} else {
					o.Rollback()
					return nil, e
				}
			}

			// create invoice directly when invoice_term == direct_invoice
			if r.SalesOrder.InvoiceTerm.ID == 2 && r.SalesOrder.Status != 2 &&
				r.SalesOrder.Status != 9 && r.SalesOrder.Status != 10 &&
				r.SalesOrder.Status != 11 && r.SalesOrder.Status != 12 &&
				r.SalesOrder.Status != 13 {

				if e = r.SalesOrder.SalesTerm.Read("ID"); e != nil {
					o.Rollback()
					return nil, e
				}

				dueDate, _ := time.Parse("2006-01-02", r.RecognitionDate)
				t2 := dueDate.AddDate(0, 0, int(r.SalesOrder.SalesTerm.DaysValue))

				//generate codes for document
				codeSI, _ := util.GenerateDocCode("SI", r.SalesOrder.Branch.Code, "sales_invoice")

				si := &model.SalesInvoice{
					SalesOrder:         r.SalesOrder,
					Code:               codeSI,
					SalesTerm:          r.SalesOrder.SalesTerm,
					InvoiceTerm:        r.SalesOrder.InvoiceTerm,
					PaymentGroup:       r.SalesOrder.PaymentGroup,
					CodeExt:            codeSI, //customer code
					RecognitionDate:    r.OrderDate,
					DueDate:            t2,
					BillingAddress:     r.SalesOrder.BillingAddress,
					DeliveryFee:        r.SalesOrder.DeliveryFee,
					TotalPrice:         r.TotalPrice,
					TotalCharge:        r.TotalCharge,
					Note:               r.SalesOrder.Note,
					AdjNote:            "-",
					Status:             1,
					CreatedAt:          time.Now(),
					CreatedBy:          r.Session.Staff.ID,
					PointRedeemAmount:  r.SalesOrder.PointRedeemAmount,
					TotalSkuDiscAmount: r.TotalSkuDiscAmount,
				}

				if r.SalesOrder.Voucher != nil {
					if e = r.SalesOrder.Voucher.Read("ID"); e != nil {
						o.Rollback()
						return nil, e
					}
					si.VouRedeemCode = r.SalesOrder.Voucher.RedeemCode
					si.VouDiscAmount = r.SalesOrder.VouDiscAmount
					si.VoucherID = r.SalesOrder.Voucher.ID
				}

				var salesInvoiceItemQty float64
				if _, e = o.Insert(si); e == nil {

					for i, rowDelivery := range u.DeliveryOrderItems {

						if rowDelivery.DeliverQty >= rowDelivery.SalesOrderItem.OrderQty {
							salesInvoiceItemQty = rowDelivery.SalesOrderItem.OrderQty
						} else {
							salesInvoiceItemQty = rowDelivery.DeliverQty
						}

						items := &model.SalesInvoiceItem{
							SalesInvoice:   &model.SalesInvoice{ID: si.ID},
							SalesOrderItem: rowDelivery.SalesOrderItem,
							Product:        rowDelivery.Product,
							TaxableItem:    rowDelivery.SalesOrderItem.TaxableItem,
							TaxPercentage:  rowDelivery.SalesOrderItem.TaxPercentage,
							InvoiceQty:     salesInvoiceItemQty,
							UnitPrice:      rowDelivery.SalesOrderItem.UnitPrice,
							Subtotal:       r.DeliveryOrderItems[i].Subtotal,
							SkuDiscAmount:  r.DeliveryOrderItems[i].SkuDiscAmount,
						}
						if _, e = o.Insert(items); e == nil {
							si.SalesInvoiceItems = append(si.SalesInvoiceItems, items)
						} else {
							o.Rollback()
							return nil, e
						}

					}
					so = &model.SalesOrder{
						ID:     r.SalesOrder.ID,
						Status: 10,
					}

					if _, e = o.Update(so, "Status"); e != nil {
						o.Rollback()
						return nil, e
					}

					if r.IsCreateCreditLimitLog == 1 {
						if e = log.CreditLimitLogByStaff(r.SalesOrder.Branch.Merchant, r.SalesOrder.ID, "sales_invoice", r.CreditLimitBefore, r.CreditLimitAfter, r.Session.Staff.ID, "create sales invoice direct"); e != nil {
							o.Rollback()
							return nil, e
						}
						r.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount = r.CreditLimitAfter
						if _, e = o.Update(r.SalesOrder.Branch.Merchant, "credit_limit_remaining"); e != nil {
							o.Rollback()
							return nil, e
						}
					}

					e = log.AuditLogByUser(r.Session.Staff, si.ID, "sales_invoice", "create", "")

				} else {
					o.Rollback()
					return nil, e
				}

			} else {

				so = &model.SalesOrder{
					ID: r.SalesOrder.ID,
				}
				if r.SalesOrder.Status == 1 {
					so.Status = 7
				} else if r.SalesOrder.Status == 9 {
					so.Status = 10
				} else if r.SalesOrder.Status == 12 {
					so.Status = 13
				}
				if _, e = o.Update(so, "Status"); e != nil {
					o.Rollback()
					return nil, e
				}
			}

			if checkStatus.Value != "1" {
				for _, row := range u.DeliveryOrderItems {

					row.Product.Read("ID")
					row.Product.Category.Read("ID")
					row.DeliveryOrder.Read("ID")
					row.SalesOrderItem.SalesOrder.Read("ID")
					row.SalesOrderItem.SalesOrder.Warehouse.Read("ID")

					go event.Call("delivery::delivery", row)

				}

				e = log.AuditLogByUser(r.Session.Staff, u.ID, "delivery_order", "create", "")
			}

		} else {
			o.Rollback()
			return nil, e

		}

	}

	o.Commit()

	if u.SalesOrder.OrderChannel == 5 && (so.Status == 7 || so.Status == 10 || so.Status == 13) {
		util.SendIDToOapi(r.SalesOrderID, r.Session.Token)
	}

	messageNotif := &util.MessageNotification{}

	// Get message notif for order type self pick up
	if r.SalesOrder.OrderType.ID == 6 {
		orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0034'").QueryRow(&messageNotif)
	} else {
		orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0003'").QueryRow(&messageNotif)
	}

	messageNotif.Message = util.ReplaceNotificationSalesOrder(messageNotif.Message, "#sales_order_code#", r.SalesOrder.Code)

	mn := &util.ModelNotification{
		SendTo:     r.SalesOrder.Branch.Merchant.UserMerchant.FirebaseToken,
		Title:      messageNotif.Title,
		Message:    messageNotif.Message,
		Type:       "1",
		RefID:      r.SalesOrder.ID,
		MerchantID: r.SalesOrder.Branch.Merchant.ID,
		ServerKey:  util.ServerKeyFireBase,
	}
	util.PostModelNotification(mn)

	return u, e
}

func (r *confirmRequest) Confirm() (s *model.DeliveryOrder, e error) {
	o := orm.NewOrm()
	o.Begin()

	or := orm.NewOrm()
	or.Using("read_only")

	var so *model.SalesOrder

	s = &model.DeliveryOrder{
		ID:          r.DeliveryOrder.ID,
		Status:      2,
		ConfirmedAt: time.Now(),
		ConfirmedBy: r.Session.Staff.ID,
	}

	if _, e = o.Update(s, "Status", "ConfirmedAt", "ConfirmedBy"); e == nil {
		for _, row := range r.DeliveryOrderItems {

			product := &model.DeliveryOrderItem{
				ID:              row.DeliveryOrderItem.ID,
				ReceiveQty:      row.DeliverQty,
				ReceiptItemNote: row.Note,
			}

			if _, e = o.Update(product, "ReceiveQty", "ReceiptItemNote"); e != nil {
				o.Rollback()
				return nil, e
			}
		}

		if len(r.DeliveryReturnItem) > 0 {
			r.DeliveryReturn.Code, _ = util.GenerateDocCode("DR", r.DeliveryOrder.SalesOrder.Branch.Code, "delivery_order")
			r.DeliveryReturn.CreatedAt = time.Now()
			r.DeliveryReturn.CreatedBy = r.Session.Staff.ID
			if _, e = o.Insert(r.DeliveryReturn); e == nil {
				for _, row := range r.DeliveryReturnItem {
					row.DeliveryReturn = r.DeliveryReturn
					if _, e = o.Insert(row); e != nil {
						o.Rollback()
					}
				}
				e = log.AuditLogByUser(r.Session.Staff, r.DeliveryReturn.ID, "delivery_return", "create", "Confirm Delivery Order")
			} else {
				o.Rollback()
			}
		}
		so = &model.SalesOrder{
			ID: r.DeliveryOrder.SalesOrder.ID,
		}

		if r.DeliveryOrder.SalesOrder.Status == 7 {
			so.Status = 8
		} else if r.DeliveryOrder.SalesOrder.Status == 10 {
			so.Status = 11
		} else if r.DeliveryOrder.SalesOrder.Status == 13 {

			so.Status = 2
			so.FinishedAt = time.Now()

			messageNotif := &util.MessageNotification{}

			or.Raw("SELECT message, title FROM notification WHERE code= 'NOT0006'").QueryRow(&messageNotif)
			messageNotif.Message = util.ReplaceNotificationSalesOrder(messageNotif.Message, "#sales_order_code#", r.DeliveryOrder.SalesOrder.Code)

			mn := &util.ModelNotification{
				SendTo:     r.DeliveryOrder.SalesOrder.Branch.Merchant.UserMerchant.FirebaseToken,
				Title:      messageNotif.Title,
				Message:    messageNotif.Message,
				Type:       "1",
				RefID:      r.DeliveryOrder.SalesOrder.ID,
				MerchantID: r.DeliveryOrder.SalesOrder.Branch.Merchant.ID,
				ServerKey:  util.ServerKeyFireBase,
			}
			util.PostModelNotification(mn)
		}

		messageNotif := &util.MessageNotification{}

		// Get message notif for order type self pick up
		if r.DeliveryOrder.SalesOrder.OrderType.ID == 6 {
			or.Raw("SELECT message, title FROM notification WHERE code= 'NOT0035'").QueryRow(&messageNotif)
		} else {
			or.Raw("SELECT message, title FROM notification WHERE code= 'NOT0009'").QueryRow(&messageNotif)
		}

		mnd := &util.ModelNotification{
			SendTo:     r.DeliveryOrder.SalesOrder.Branch.Merchant.UserMerchant.FirebaseToken,
			Title:      messageNotif.Title,
			Message:    messageNotif.Message,
			Type:       "3",
			RefID:      r.DeliveryOrder.SalesOrder.ID,
			MerchantID: r.DeliveryOrder.SalesOrder.Branch.Merchant.ID,
			ServerKey:  util.ServerKeyFireBase,
		}
		util.PostModelNotification(mnd)

		if _, e = o.Update(so, "Status", "FinishedAt"); e == nil {
			e = log.AuditLogByUser(r.Session.Staff, s.ID, "delivery_order", "confirm", "")
		} else {
			o.Rollback()
			return nil, e
		}

	} else {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	//send webhook delivered and finished (for status 2)
	if r.DeliveryOrder.SalesOrder.OrderChannel == 5 && (so.Status == 8 || so.Status == 11 || so.Status == 2) {
		util.SendIDToOapi(common.Encrypt(r.DeliveryOrder.SalesOrder.ID), r.Session.Token)
	}

	return
}

func (r *updateRequest) Update() (s *model.DeliveryOrder, e error) {
	o := orm.NewOrm()
	o.Begin()
	checkStatus, e := repository.GetConfigApp("attribute", "kafka_mongo_delivery_order")
	if e != nil {
		return nil, e
	}
	s = &model.DeliveryOrder{
		ID:              r.DeliveryOrder.ID,
		Note:            r.Note,
		TotalWeight:     r.TotalWeight,
		RecognitionDate: r.OrderDate,
		UpdatedAt:       time.Now(),
		UpdatedBy:       r.Session.Staff.ID,
	}
	if _, e = o.Update(s, "Note", "TotalWeight", "RecognitionDate", "UpdatedAt", "UpdatedBy"); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, row := range r.DeliveryOrderItems {

		doiID := &model.DeliveryOrderItem{DeliveryOrder: r.DeliveryOrder, Product: row.Product}
		if e = doiID.Read("DeliveryOrder", "Product"); e != nil {
			o.Rollback()
			return nil, e
		}

		soi := &model.SalesOrderItem{SalesOrder: r.SalesOrder, Product: row.Product}
		if e = soi.Read("SalesOrder", "Product"); e != nil {
			o.Rollback()
			return nil, e
		}

		item := &model.DeliveryOrderItem{
			DeliveryOrder:  &model.DeliveryOrder{ID: r.DeliveryOrder.ID},
			SalesOrderItem: soi,
			Product:        row.Product,
			DeliverQty:     row.DeliverQty,
			Weight:         row.Weight,
		}

		if row.ID != "" {
			item.ID = doiID.ID
			if _, e = o.Update(item, "DeliveryOrder", "Product", "DeliverQty", "Weight"); e != nil {
				o.Rollback()
				return nil, e
			}

		} else {
			if _, e = o.Insert(item); e != nil {
				o.Rollback()
				return nil, e
			}
		}

		if checkStatus.Value != "1" {
			if item.Product.ID == 0 || item.DeliveryOrder.ID == 0 || item.SalesOrderItem.ID == 0 {
				o.Rollback()
				e = echo.NewHTTPError(http.StatusBadRequest, "Product or DO or SOI not found")
				return nil, e
			}
			if e = item.Product.Read("ID"); e != nil {
				o.Rollback()
				return nil, e
			}
			if e = item.DeliveryOrder.Read("ID"); e != nil {
				o.Rollback()
				return nil, e
			}
			if e = item.SalesOrderItem.Read("ID"); e != nil {
				o.Rollback()
				return nil, e
			}

			if item.SalesOrderItem.SalesOrder.ID == 0 {
				o.Rollback()
				e = echo.NewHTTPError(http.StatusBadRequest, "SO not found")
				return nil, e
			}
			if e = item.SalesOrderItem.SalesOrder.Read("ID"); e != nil {
				o.Rollback()
				return nil, e
			}

			if item.SalesOrderItem.SalesOrder.Warehouse.ID == 0 {
				o.Rollback()
				e = echo.NewHTTPError(http.StatusBadRequest, "Warehouse not found")
				return nil, e
			}
			if e = item.SalesOrderItem.SalesOrder.Warehouse.Read("ID"); e != nil {
				o.Rollback()
				return nil, e
			}

			if item.SalesOrderItem.SalesOrder.Warehouse.ID == 0 {
				o.Rollback()
				e = echo.NewHTTPError(http.StatusBadRequest, "Warehouse not found")
				return nil, e
			}

			if e = item.DeliveryOrder.SalesOrder.Read("ID"); e != nil {
				o.Rollback()
				return nil, e
			}

			if e = item.DeliveryOrder.SalesOrder.OrderType.Read("ID"); e != nil {
				o.Rollback()
				return nil, e
			}

			stock := &model.Stock{
				Product:   row.Product,
				Warehouse: r.DeliveryOrder.Warehouse,
			}
			if e = o.Read(stock, "Product", "Warehouse"); e != nil {
				o.Rollback()
				return nil, e
			}

			if item.DeliveryOrder.SalesOrder.OrderType.Name == "Zero Waste" {
				// stock in and out for waste stock
				var wl *model.WasteLog

				if e = o.Raw("SELECT wl.quantity , wl.created_at FROM waste_log wl WHERE ref_id = ? AND warehouse_id = ? "+
					"AND ref_type = 7 AND product_id = ? and `type` = 2 and status = 1 ORDER BY id DESC LIMIT 1", r.DeliveryOrder.ID, r.DeliveryOrder.Warehouse.ID, row.Product.ID).QueryRow(&wl); e != nil {
					o.Rollback()
					return nil, e
				}

				if wl == nil || stock == nil {
					o.Rollback()
					e = echo.NewHTTPError(http.StatusBadRequest, "Waste stock not found")
					return nil, e
				}

				wlIn := &model.WasteLog{
					Warehouse:    r.DeliveryOrder.Warehouse,
					Product:      row.Product,
					Ref:          r.DeliveryOrder.ID,
					RefType:      7,
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
					Warehouse:    r.DeliveryOrder.Warehouse,
					Product:      row.Product,
					Ref:          r.DeliveryOrder.ID,
					RefType:      7,
					Type:         2,
					InitialStock: wlIn.FinalStock,
					Quantity:     row.DeliverQty,
					FinalStock:   wlIn.FinalStock - row.DeliverQty,
					Status:       1,
					DocNote:      r.DeliveryOrder.Note,
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
				// stock in and out for goods stock
				var sl *model.StockLog

				if e = o.Raw("SELECT sl.quantity , sl.created_at FROM stock_log sl WHERE ref_id = ? AND warehouse_id = ? "+
					"AND ref_type = 1 AND product_id = ? and `type` = 2 and status = 1 ORDER BY id DESC LIMIT 1", r.DeliveryOrder.ID, r.DeliveryOrder.Warehouse.ID, row.Product.ID).QueryRow(&sl); e != nil {
					o.Rollback()
					return nil, e
				}

				if sl == nil || stock == nil {
					o.Rollback()
					e = echo.NewHTTPError(http.StatusBadRequest, "Stock not found")
					return nil, e
				}

				slIn := &model.StockLog{
					Warehouse:    r.DeliveryOrder.Warehouse,
					Product:      row.Product,
					Ref:          r.DeliveryOrder.ID,
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
					Warehouse:    r.DeliveryOrder.Warehouse,
					Product:      row.Product,
					Ref:          r.DeliveryOrder.ID,
					RefType:      1,
					Type:         2,
					InitialStock: slIn.FinalStock,
					Quantity:     row.DeliverQty,
					FinalStock:   slIn.FinalStock - row.DeliverQty,
					UnitCost:     row.UnitPrice,
					Status:       1,
					DocNote:      r.DeliveryOrder.Note,
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

			e = log.AuditLogByUser(r.Session.Staff, r.DeliveryOrder.ID, "delivery_order", "update", "")
		}

		// if term = direct_invoice && invoice status = active
		if r.CheckSalesInvoiceTermInvoice == 1 {
			if _, e = o.Update(row.SalesInvoiceItem, "InvoiceQty", "UnitPrice", "Subtotal", "SkuDiscAmount"); e != nil {
				o.Rollback()
				return nil, e
			}
		}

	}
	if checkStatus.Value != "1" {
		e = log.AuditLogByUser(r.Session.Staff, r.DeliveryOrder.ID, "delivery_order", "update", "")
	}

	// update sales invoice document if term = direct_invoice && invoice status = active
	if r.CheckSalesInvoiceTermInvoice == 1 {
		si := &model.SalesInvoice{
			ID:                 r.SalesInvoice[0].ID,
			TotalPrice:         r.TotalPrice,
			TotalCharge:        r.TotalCharge,
			LastUpdatedAt:      time.Now(),
			LastUpdatedBy:      r.Session.Staff.ID,
			TotalSkuDiscAmount: r.TotalSkuDiscAmount,
		}
		if _, e = o.Update(si, "TotalPrice", "TotalCharge", "LastUpdatedAt", "LastUpdatedBy", "TotalSkuDiscAmount"); e == nil {
			e = log.AuditLogByUser(r.Session.Staff, si.ID, "sales_invoice", "update", "Auto-update invoice by delivery order updated")
		}

	}

	o.Commit()
	s, _ = repository.GetDeliveryOrder("id", s.ID)

	return
}

func (r *cancelRequest) Cancel() (do *model.DeliveryOrder, e error) {
	var err error
	o := orm.NewOrm()
	o.Begin()
	checkStatus, e := repository.GetConfigApp("attribute", "kafka_mongo_delivery_order")
	if e != nil {
		return nil, e
	}

	do = &model.DeliveryOrder{
		ID: r.DeliveryOrder.ID,
	}

	do.Status = 3
	do.CancelledAt = time.Now()
	do.CancelledBy = r.Session.Staff.ID
	do.CancellationNote = r.Note

	if e = r.DeliveryOrder.Read("ID"); e != nil {
		o.Rollback()
		return nil, e
	}

	o.LoadRelated(r.DeliveryOrder, "DeliveryOrderItems")

	if e = r.DeliveryOrder.SalesOrder.Read("ID"); e != nil {
		o.Rollback()
		return nil, e
	}
	if e = r.DeliveryOrder.Warehouse.Read("ID"); e != nil {
		o.Rollback()
		return nil, e
	}

	if _, e = o.Update(do, "Status", "CancelledAt", "CancelledBy", "CancellationNote"); e != nil {
		o.Rollback()
		return nil, e
	}

	so := &model.SalesOrder{
		ID: r.DeliveryOrder.SalesOrder.ID,
	}

	if r.DeliveryOrder.SalesOrder.Status == 7 {
		so.Status = 1
	} else if r.DeliveryOrder.SalesOrder.Status == 10 {
		so.Status = 9
	} else if r.DeliveryOrder.SalesOrder.Status == 13 {
		so.Status = 12
	}

	if _, e = o.Update(so, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = r.DeliveryOrder.SalesOrder.OrderType.Read("ID"); e != nil {
		o.Rollback()
		return nil, e
	}

	if checkStatus.Value != "1" {
		for _, v := range r.DeliveryOrder.DeliveryOrderItems {
			r.DeliveryOrderItem = append(r.DeliveryOrderItem, v)

			if r.DeliveryOrder.SalesOrder.OrderType.Name == "Zero Waste" {
				wasteLog := &model.WasteLog{
					Warehouse: r.DeliveryOrder.Warehouse,
					Product:   v.Product,
					Ref:       r.DeliveryOrder.ID,
					Status:    1,
				}

				if err = wasteLog.Read("Warehouse", "Product", "Ref", "Status"); err == nil {
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
				if e = o.Read(stockLog, "Warehouse", "Product", "Ref", "Status"); e != nil {
					o.Rollback()
					return nil, e
				}
				r.StockLog = append(r.StockLog, stockLog)

			}

			stock := &model.Stock{
				Warehouse: r.DeliveryOrder.Warehouse,
				Product:   v.Product,
			}
			if err = o.Read(stock, "Warehouse", "Product"); err != nil {
				return nil, e
			}
			r.Stock = append(r.Stock, stock)

		}

		for i, v := range r.Stock {
			// check for sales order type for zero waste
			if r.DeliveryOrder.SalesOrder.OrderType.Name == "Zero Waste" {
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
					RefType:      7,
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

			} else {
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
		e = log.AuditLogByUser(r.Session.Staff, r.DeliveryOrder.ID, "delivery_order", "cancel", r.Note)

	}

	o.Commit()

	return r.DeliveryOrder, e
}
