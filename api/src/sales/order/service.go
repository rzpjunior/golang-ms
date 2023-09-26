// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"errors"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"

	"git.edenfarm.id/project-version2/api/log"
	xendit2 "git.edenfarm.id/project-version2/api/service/xendit"
	"git.edenfarm.id/project-version2/api/util"
)

// Save : function to insert data requested into database
func Save(r createRequest) (so *model.SalesOrder, e error) {
	r.Code, e = util.GenerateDocCode("SO", r.Branch.Code, "sales_order")
	if e != nil {
		return nil, e
	}

	o := orm.NewOrm()
	o.Begin()

	var isCreated bool

	// insert sales order
	so = &model.SalesOrder{
		Code:               r.Code,
		Branch:             r.Branch,
		SalesTerm:          r.SalesTerm,
		InvoiceTerm:        r.InvoiceTerm,
		Salesperson:        r.Salesperson,
		SalesGroupID:       r.Salesperson.SalesGroupID,
		Warehouse:          r.Warehouse,
		Wrt:                r.Wrt,
		Area:               r.Branch.Area,
		Voucher:            r.Voucher,
		SubDistrict:        r.Branch.SubDistrict,
		PriceSet:           r.Branch.PriceSet,
		PaymentGroup:       r.PaymentGroup,
		Archetype:          r.Branch.Archetype,
		OrderType:          r.OrderType,
		DeliveryDate:       r.DeliveryDate,
		RecognitionDate:    r.RecognitionDate,
		BillingAddress:     r.BillingAddress,
		ShippingAddress:    r.ShippingAddress,
		DeliveryFee:        r.DeliveryFee,
		OrderChannel:       int8(1),
		HasExtInvoice:      int8(2),
		Note:               r.Note,
		Status:             int8(1),
		TotalPrice:         r.TotalPrice,
		TotalCharge:        r.TotalCharge,
		TotalWeight:        r.TotalWeight,
		CreatedAt:          time.Now(),
		CreatedBy:          r.Session.Staff.ID,
		TotalSkuDiscAmount: r.TotalSkuDiscAmount,
	}
	if r.Voucher != nil {
		so.Voucher = r.Voucher
		so.VouRedeemCode = r.Voucher.RedeemCode
		if r.Voucher.Type != 4 {
			so.VouDiscAmount = r.Voucher.DiscAmount
		}
	}
	if _, e = o.Insert(so); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, v := range r.Products {
		// insert sales order item
		soi := &model.SalesOrderItem{
			SalesOrder:    so,
			Product:       v.Product,
			OrderQty:      v.Quantity,
			UnitPrice:     float64(v.UnitPrice),
			ShadowPrice:   float64(v.Price.ShadowPrice),
			Subtotal:      v.Subtotal,
			Weight:        v.Weight,
			Note:          v.Note,
			ProductPush:   v.ProductPush,
			TaxableItem:   v.Product.Taxable,
			TaxPercentage: v.Product.TaxPercentage,
			DefaultPrice:  v.DefaultPrice,
		}
		if v.IsUseSkuDiscount == 1 {
			soi.SkuDiscountItem = v.SkuDiscountItem
			soi.DiscountQty = v.DiscQty
			soi.UnitPriceDiscount = v.UnitPriceDiscount
			soi.SkuDiscountAmount = v.DiscAmount
		}
		if _, e = o.Insert(soi); e != nil {
			o.Rollback()
			return nil, e
		}

		// insert packing order
		if r.PackingOrder != nil {
			poi := &model.PackingOrderItem{
				PackingOrder: r.PackingOrder,
				Product:      v.Product,
				TotalOrder:   v.Quantity,
				TotalWeight:  0,
				TotalPack:    0,
			}
			if isCreated, _, e = o.ReadOrCreate(poi, "PackingOrder", "Product"); e == nil {
				if !isCreated {
					poi.TotalOrder = poi.TotalOrder + v.Quantity
					if _, e = o.Update(poi, "TotalOrder"); e != nil {
						o.Rollback()
						return nil, e
					}
				}
			}
		}

		if v.IsUseSkuDiscount == 1 {
			// check newest data (remaining quota & budget) of sku discount item
			if e = v.SkuDiscountItem.Read("ID"); e != nil {
				return nil, e
			}

			if v.SkuDiscountItem.RemOverallQuota-v.DiscQty < 0 || (v.SkuDiscountItem.IsUseBudget == 1 && v.SkuDiscountItem.RemBudget <= 0) {
				e = errors.New("Failed to save quota")
				o.Rollback()
				return nil, e
			}

			v.SkuDiscountItem.RemOverallQuota = v.SkuDiscountItem.RemOverallQuota - v.DiscQty
			if v.SkuDiscountItem.IsUseBudget == 1 {
				v.SkuDiscountItem.RemBudget = v.SkuDiscountItem.RemBudget - v.DiscAmount
			}
			// update remaining quota and budget sku discount item
			if _, e = o.Update(v.SkuDiscountItem, "RemOverallQuota", "RemBudget"); e != nil {
				o.Rollback()
				return nil, e
			}

			// insert into sku discount log
			sdl := &model.SkuDiscountLog{
				SkuDiscount:     v.SkuDiscountItem.SkuDiscount,
				SkuDiscountItem: v.SkuDiscountItem,
				Merchant:        r.Branch.Merchant,
				Branch:          r.Branch,
				SalesOrderItem:  soi,
				Product:         v.Product,
				DiscountAmount:  v.DiscAmount,
				DiscountQty:     v.DiscQty,
				CreatedAt:       time.Now(),
				Status:          1,
			}
			if _, e = o.Insert(sdl); e != nil {
				o.Rollback()
				return nil, e
			}
		}
	}

	// insert voucher
	if r.Voucher != nil && r.Voucher.ID != 0 {
		r.Voucher.RemOverallQuota = r.Voucher.RemOverallQuota - 1
		if _, e = o.Update(r.Voucher, "rem_overall_quota"); e != nil {
			o.Rollback()
			return nil, e
		}

		vl := &model.VoucherLog{
			Voucher:           r.Voucher,
			Merchant:          r.Branch.Merchant,
			Branch:            r.Branch,
			SalesOrder:        so,
			TagCustomer:       r.SameTagCustomer,
			VoucherDiscAmount: r.Voucher.DiscAmount,
			Timestamp:         time.Now(),
			Status:            int8(1),
		}

		if _, e = o.Insert(vl); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if r.IsCreateCreditLimitLog == 1 {
		if e = log.CreditLimitLogByStaff(r.Branch.Merchant, so.ID, "sales_order", r.CreditLimitBefore, r.CreditLimitAfter, r.Session.Staff.ID, "create sales order"); e != nil {
			o.Rollback()
			return nil, e
		}
		r.Branch.Merchant.RemainingCreditLimitAmount = r.CreditLimitAfter
		if _, e = o.Update(r.Branch.Merchant, "credit_limit_remaining"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, so.ID, "sales_order", "create", r.NotePriceChange); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	if r.IsCreateMerchantVa["bca"] == 1 {
		xendit2.BCAXenditFixedVA(r.Branch.Merchant)
	}

	if r.IsCreateMerchantVa["permata"] == 1 {
		xendit2.PermataXenditFixedVA(r.Branch.Merchant)
	}

	return so, e
}

// Update : function to change data requested into database
func Update(r updateRequest) (so *model.SalesOrder, e error) {
	o := orm.NewOrm()
	o.Begin()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	var keepItemsId []int64
	var isCreated bool
	var soiPrevQty float64
	var deletedSoi []*model.SalesOrderItem

	if r.Voucher == nil {
		if r.SalesOrder.Voucher != nil {
			for _, v := range r.VoucherLog {
				v.Status = 3
				if _, e = o.Update(v, "Status"); e != nil {
					o.Rollback()
					return nil, e
				}
			}

			r.SalesOrder.Voucher.RemOverallQuota = r.SalesOrder.Voucher.RemOverallQuota + 1
			if _, e = o.Update(r.SalesOrder.Voucher, "RemOverallQuota"); e != nil {
				o.Rollback()
				return nil, e
			}

			r.SalesOrder.Voucher = nil
			r.SalesOrder.VouRedeemCode = ""
			r.SalesOrder.VouDiscAmount = 0
		} else {
			r.SalesOrder.Voucher = nil
			r.SalesOrder.VouRedeemCode = ""
			r.SalesOrder.VouDiscAmount = 0
		}
	} else {
		if r.SalesOrder.Voucher != nil {
			if r.Voucher.ID != r.SalesOrder.Voucher.ID {
				for _, v := range r.VoucherLog {
					v.Status = 3
					if _, e = o.Update(v, "Status"); e == nil {
						r.SalesOrder.Voucher.RemOverallQuota = r.SalesOrder.Voucher.RemOverallQuota + 1
						if _, e = o.Update(r.SalesOrder.Voucher, "RemOverallQuota"); e != nil {
							o.Rollback()
							return nil, e
						}
					} else {
						o.Rollback()
						return nil, e
					}
				}
				vl := &model.VoucherLog{
					Voucher:           r.Voucher,
					Merchant:          r.SalesOrder.Branch.Merchant,
					Branch:            r.SalesOrder.Branch,
					SalesOrder:        r.SalesOrder,
					TagCustomer:       r.SameTagCustomer,
					VoucherDiscAmount: r.Voucher.DiscAmount,
					Timestamp:         time.Now(),
					Status:            int8(1),
				}

				if _, e = o.Insert(vl); e != nil {
					o.Rollback()
					return nil, e
				}

				r.Voucher.RemOverallQuota = r.Voucher.RemOverallQuota - 1
				if _, e = o.Update(r.Voucher, "RemOverallQuota"); e != nil {
					o.Rollback()
					return nil, e
				}

				r.SalesOrder.Voucher = r.Voucher
				r.SalesOrder.VouRedeemCode = r.Voucher.RedeemCode
				if r.Voucher.Type != 4 {
					r.SalesOrder.VouDiscAmount = r.Voucher.DiscAmount
				}
			}
		} else {
			vl := &model.VoucherLog{
				Voucher:           r.Voucher,
				Merchant:          r.SalesOrder.Branch.Merchant,
				Branch:            r.SalesOrder.Branch,
				SalesOrder:        r.SalesOrder,
				TagCustomer:       r.SameTagCustomer,
				VoucherDiscAmount: r.Voucher.DiscAmount,
				Timestamp:         time.Now(),
				Status:            int8(1),
			}

			if _, e = o.Insert(vl); e != nil {
				o.Rollback()
				return nil, e
			}

			r.Voucher.RemOverallQuota = r.Voucher.RemOverallQuota - 1
			if _, e = o.Update(r.Voucher, "RemOverallQuota"); e != nil {
				o.Rollback()
				return nil, e
			}

			r.SalesOrder.Voucher = r.Voucher
			r.SalesOrder.VouRedeemCode = r.Voucher.RedeemCode
			if r.Voucher.Type != 4 {
				r.SalesOrder.VouDiscAmount = r.Voucher.DiscAmount
			}
		}
	}

	r.SalesOrder.Archetype = r.SalesOrder.Branch.Archetype
	r.SalesOrder.TotalWeight = r.TotalWeight
	r.SalesOrder.DeliveryDate = r.DeliveryDate
	r.SalesOrder.Wrt = r.Wrt
	r.SalesOrder.RecognitionDate = r.RecognitionDate
	r.SalesOrder.OrderType = r.OrderType
	r.SalesOrder.Salesperson = r.Salesperson
	r.SalesOrder.SalesTerm = r.SalesTerm
	r.SalesOrder.InvoiceTerm = r.InvoiceTerm
	r.SalesOrder.BillingAddress = r.SalesOrder.Branch.Merchant.BillingAddress
	r.SalesOrder.Note = r.Note
	r.SalesOrder.DeliveryFee = r.DeliveryFee
	r.SalesOrder.TotalCharge = r.TotalCharge
	r.SalesOrder.TotalPrice = r.TotalPrice
	r.SalesOrder.Warehouse = r.Warehouse
	r.SalesOrder.LastUpdatedAt = time.Now()
	r.SalesOrder.LastUpdatedBy = r.Session.Staff.ID
	r.SalesOrder.IsLocked = 2
	r.SalesOrder.SalesGroupID = r.Salesperson.SalesGroupID
	r.SalesOrder.TotalSkuDiscAmount = r.TotalSkuDiscAmount
	r.SalesOrder.PaymentGroup = r.PaymentGroup

	if _, e = o.Update(r.SalesOrder, "Archetype", "TotalWeight", "TotalCharge", "TotalPrice", "DeliveryDate", "Wrt", "RecognitionDate", "OrderType", "Salesperson", "SalesTerm",
		"InvoiceTerm", "BillingAddress", "Note", "Voucher", "VouRedeemCode", "VouDiscAmount", "DeliveryFee", "Warehouse", "LastUpdatedAt", "LastUpdatedBy", "IsLocked", "SalesGroupID", "TotalSkuDiscAmount", "ShippingAddress", "PaymentGroup"); e != nil {
		o.Rollback()
		return nil, e
	}

	if r.UpdateAll == 1 {
		for _, v := range r.Products {
			var (
				soiPrevSkuDiscountItem     *model.SkuDiscountItem
				skuDiscountRemOverallQuota float64
			)

			soi := &model.SalesOrderItem{
				SalesOrder:    r.SalesOrder,
				Product:       v.Product,
				OrderQty:      v.Quantity,
				UnitPrice:     float64(v.UnitPrice),
				Subtotal:      v.Subtotal,
				Note:          v.Note,
				ProductPush:   v.ProductPush,
				ShadowPrice:   v.Price.ShadowPrice,
				Weight:        v.Weight,
				TaxableItem:   v.TaxableItem,
				TaxPercentage: v.TaxPercentage,
				DefaultPrice:  v.DefaultPrice,
			}
			if v.IsUseSkuDiscount == 1 {
				soi.SkuDiscountItem = v.SkuDiscountItem
				soi.DiscountQty = v.DiscQty
				soi.UnitPriceDiscount = v.UnitPriceDiscount
				soi.SkuDiscountAmount = v.DiscAmount
			}
			if isCreated, soi.ID, e = o.ReadOrCreate(soi, "SalesOrder", "Product"); e != nil {
				o.Rollback()
				return nil, e
			}

			if !isCreated {
				soiPrevQty = soi.OrderQty

				soi.OrderQty = v.Quantity
				soi.UnitPrice = float64(v.UnitPrice)
				soi.Subtotal = v.Subtotal
				soi.Note = v.Note
				soi.ProductPush = v.ProductPush
				soi.ShadowPrice = v.Price.ShadowPrice
				soi.Weight = v.Weight
				soi.DefaultPrice = v.DefaultPrice

				if soi.SkuDiscountItem != nil {
					if e = soi.SkuDiscountItem.Read("ID"); e != nil {
						o.Rollback()
						return nil, e
					}
					soiPrevSkuDiscountItem = soi.SkuDiscountItem
				}

				if v.IsUseSkuDiscount == 1 {
					soi.SkuDiscountItem = v.SkuDiscountItem
					soi.DiscountQty = v.DiscQty
					soi.UnitPriceDiscount = v.UnitPriceDiscount
					soi.SkuDiscountAmount = v.DiscAmount
				} else {
					soi.SkuDiscountItem = nil
					soi.DiscountQty = 0
					soi.UnitPriceDiscount = 0
					soi.SkuDiscountAmount = 0
				}

				if _, e = o.Update(soi, "OrderQty", "UnitPrice", "Subtotal", "Note", "ShadowPrice", "Weight", "ProductPush", "SkuDiscountItem", "DiscountQty", "UnitPriceDiscount", "SkuDiscountAmount"); e != nil {
					o.Rollback()
					return nil, e
				}

				if soiPrevSkuDiscountItem != nil {
					sdl := &model.SkuDiscountLog{
						Branch:          r.SalesOrder.Branch,
						SalesOrderItem:  soi,
						SkuDiscountItem: soiPrevSkuDiscountItem,
						Status:          1,
					}

					if e = sdl.Read("Branch", "SalesOrderItem", "SkuDiscountItem", "Status"); e == nil {
						sdl.Status = 2
						if _, e = o.Update(sdl, "Status"); e != nil {
							o.Rollback()
							return nil, e
						}
					}

					skuDiscountRemOverallQuota = soiPrevSkuDiscountItem.RemOverallQuota
					skuDiscountRemOverallQuota += sdl.DiscountQty
					soiPrevSkuDiscountItem.RemOverallQuota = skuDiscountRemOverallQuota
					if soiPrevSkuDiscountItem.IsUseBudget == 1 {
						soiPrevSkuDiscountItem.RemBudget = soiPrevSkuDiscountItem.RemBudget + sdl.DiscountAmount
					}

					// update remaining quota and budget sku discount item
					if _, e = o.Update(soiPrevSkuDiscountItem, "RemOverallQuota", "RemBudget"); e != nil {
						o.Rollback()
						return nil, e
					}
				}
			}

			if v.IsUseSkuDiscount == 1 {
				// check newest data (remaining quota & budget) of sku discount item
				if e = v.SkuDiscountItem.Read("ID"); e != nil {
					return nil, e
				}

				if v.SkuDiscountItem.RemOverallQuota-v.DiscQty < 0 || (v.SkuDiscountItem.IsUseBudget == 1 && v.SkuDiscountItem.RemBudget <= 0) {
					e = errors.New("Failed to save quota")
					o.Rollback()
					return nil, e
				}

				if skuDiscountRemOverallQuota == 0 {
					skuDiscountRemOverallQuota = v.SkuDiscountItem.RemOverallQuota
				}
				skuDiscountRemOverallQuota -= v.DiscQty
				v.SkuDiscountItem.RemOverallQuota = skuDiscountRemOverallQuota
				if v.SkuDiscountItem.IsUseBudget == 1 {
					v.SkuDiscountItem.RemBudget = v.SkuDiscountItem.RemBudget - v.DiscAmount
				} else {
					v.SkuDiscountItem.RemBudget = 0
				}
				// update remaining quota and budget sku discount item
				if _, e = o.Update(v.SkuDiscountItem, "RemOverallQuota", "RemBudget"); e != nil {
					o.Rollback()
					return nil, e
				}

				sdl := &model.SkuDiscountLog{
					Merchant:        r.SalesOrder.Branch.Merchant,
					Branch:          r.SalesOrder.Branch,
					SalesOrderItem:  soi,
					SkuDiscountItem: v.SkuDiscountItem,
					SkuDiscount:     v.SkuDiscountItem.SkuDiscount,
					DiscountAmount:  v.DiscAmount,
					DiscountQty:     v.DiscQty,
					Product:         v.Product,
					Status:          1,
					CreatedAt:       time.Now(),
				}
				if _, e = o.Insert(sdl); e != nil {
					o.Rollback()
					return nil, e
				}
			}

			if r.PackingOrder != nil {
				poi := &model.PackingOrderItem{
					PackingOrder: r.PackingOrder,
					Product:      v.Product,
					TotalOrder:   v.Quantity,
					TotalWeight:  0,
					TotalPack:    0,
				}
				if isCreated, _, e = o.ReadOrCreate(poi, "PackingOrder", "Product"); e != nil {
					o.Rollback()
					return nil, e
				}

				if !isCreated {
					poi.TotalOrder = poi.TotalOrder - soiPrevQty + v.Quantity
					if _, e = o.Update(poi, "TotalOrder"); e != nil {
						o.Rollback()
						return nil, e
					}
				}
			}

			keepItemsId = append(keepItemsId, soi.ID)
		}

		if _, e := orSelect.QueryTable(new(model.SalesOrderItem)).Filter("sales_order_id", r.SalesOrder.ID).Exclude("ID__in", keepItemsId).All(&deletedSoi); e != nil {
			o.Rollback()
			return nil, e
		}

		for _, v := range deletedSoi {
			if r.PackingOrder != nil {
				poi := &model.PackingOrderItem{
					PackingOrder: r.PackingOrder,
					Product:      v.Product,
				}
				if e = poi.Read("PackingOrder", "Product"); e != nil {
					o.Rollback()
					return nil, e
				}

				poi.TotalOrder = poi.TotalOrder - v.OrderQty
				if _, e = o.Update(poi, "TotalOrder"); e != nil {
					o.Rollback()
					return nil, e
				}
			}

			if v.SkuDiscountItem != nil {
				// archive sku discount log
				sdl := &model.SkuDiscountLog{
					SkuDiscountItem: v.SkuDiscountItem,
					Branch:          r.SalesOrder.Branch,
					SalesOrderItem:  v,
					Status:          1,
				}
				sdl.Read("SkuDiscountItem", "Branch", "SalesOrderItem", "Status")
				sdl.Status = 2
				if _, e = o.Update(sdl, "Status"); e != nil {
					o.Rollback()
					return nil, e
				}

				// return quota into remaining quota
				v.SkuDiscountItem.Read("ID")
				v.SkuDiscountItem.RemOverallQuota = v.SkuDiscountItem.RemOverallQuota + sdl.DiscountQty
				if v.SkuDiscountItem.IsUseBudget == 1 {
					v.SkuDiscountItem.RemBudget = v.SkuDiscountItem.RemBudget + sdl.DiscountAmount
				}
				if _, e = o.Update(v.SkuDiscountItem, "RemOverallQuota", "RemBudget"); e != nil {
					o.Rollback()
					return nil, e
				}
			}

			if _, e = o.Delete(v); e != nil {
				o.Rollback()
				return nil, e
			}
		}
	}

	if r.IsCreateCreditLimitLog == 1 {
		if e = log.CreditLimitLogByStaff(r.SalesOrder.Branch.Merchant, r.SalesOrder.ID, "sales_order", r.CreditLimitBefore, r.CreditLimitAfter, r.Session.Staff.ID, "update sales order"); e != nil {
			o.Rollback()
			return nil, e
		}
		r.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount = r.CreditLimitAfter
		if _, e = o.Update(r.SalesOrder.Branch.Merchant, "credit_limit_remaining"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.SalesOrder.ID, "sales_order", "update", r.NotePriceChange); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	// notification FS Apps when updating SO
	messageNotif := &util.MessageNotification{}
	orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0013'").QueryRow(&messageNotif)
	messageNotif.Message = util.ReplaceNotificationSalesOrder(messageNotif.Message, "#customer#", r.SalesOrder.Branch.Name)
	messageNotif.Message = util.ReplaceNotificationSalesOrder(messageNotif.Message, "#sales_order_code#", r.SalesOrder.Code)

	r.SalesOrder.Salesperson.Read("ID")
	r.SalesOrder.Salesperson.User.Read("ID")
	mn := &util.ModelNotification{
		SendTo:    r.SalesOrder.Salesperson.User.SalesAppNotifToken,
		Title:     messageNotif.Title,
		Message:   messageNotif.Message,
		Type:      "1",
		RefID:     r.SalesOrder.ID,
		ServerKey: util.FieldSalesServerKeyFireBase,
		StaffID:   r.SalesOrder.Salesperson.ID,
	}
	util.PostModelNotificationFieldSales(mn)

	return r.SalesOrder, nil
}

// Cancel : function to change data status into 3
func Cancel(r cancelRequest) (so *model.SalesOrder, e error) {
	o := orm.NewOrm()
	o.Begin()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	var (
		isPointExpired, isDataExisting bool
	)

	if r.SalesOrder.Voucher != nil {
		for _, v := range r.VoucherLog {
			v.Status = 3
			if _, e = o.Update(v, "Status"); e != nil {
				o.Rollback()
				return nil, e
			}
		}

		r.SalesOrder.Voucher.RemOverallQuota = r.SalesOrder.Voucher.RemOverallQuota + 1
		if _, e = o.Update(r.SalesOrder.Voucher, "rem_overall_quota"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	r.SalesOrder.CancelType = r.CancelType
	r.SalesOrder.Status = 3
	//update status sales order first
	if _, e = o.Update(r.SalesOrder, "Status", "CancelType"); e != nil {
		o.Rollback()
		return nil, e
	}

	//if not error when update, check if this sales order using point redeemed or not
	if r.SalesOrder.PointRedeemID != 0 && r.SalesOrder.PointRedeemAmount != 0 {
		mps := map[int64]float64{}
		currentDate := time.Now()
		for _, v := range r.MerchantPointLog {
			//if using redemption,update the status to 4 in merchant point log and update notes
			v.Status = 4
			v.Note = "Cancellation due to cancel sales order"
			if _, e = o.Update(v, "Status", "Note"); e != nil {
				o.Rollback()
				return nil, e
			}

			// Calculate current and next period point
			r.MerchantPointExpiration.CurrentPeriodPoint += v.CurrentPointUsed
			r.MerchantPointExpiration.NextPeriodPoint += v.NextPointUsed

			// Change the time of date to midnight (23:59:59)
			v.ExpiredDate = time.Date(v.ExpiredDate.Year(), v.ExpiredDate.Month(), v.ExpiredDate.Day(), 23, 59, 59, 0, time.Local)
			// check the point expired or not
			isPointExpired = time.Now().After(v.ExpiredDate)

			// This condition for back up data existing that didn't record current point used and next point used
			isDataExisting = v.CurrentPointUsed == 0 && v.NextPointUsed == 0
			if isDataExisting {
				r.MerchantPointExpiration.CurrentPeriodPoint += v.PointValue
			}

			// Exclude data existing from validation expired
			if isPointExpired && !isDataExisting {
				// Reverse calculation current period & next period point
				r.MerchantPointExpiration.CurrentPeriodPoint -= v.CurrentPointUsed
				r.MerchantPointExpiration.NextPeriodPoint -= v.NextPointUsed
				// 'Next point used' move to current period point
				r.MerchantPointExpiration.CurrentPeriodPoint += v.NextPointUsed
				// reduce point value with current point that expired
				v.PointValue -= v.CurrentPointUsed

				// record to merchant point log, if there is current point used that expired
				if v.CurrentPointUsed != 0 {
					mpl := &model.MerchantPointLog{
						PointValue:      v.CurrentPointUsed,
						RecentPoint:     r.RecentPoint,
						Status:          6,
						Note:            "Point Issued From Cancellation Redeem And Current Point Used Have To Expired",
						Merchant:        r.SalesOrder.Branch.Merchant,
						SalesOrder:      r.SalesOrder,
						CreatedDate:     currentDate,
						TransactionType: 6,
					}
					//start trying to insert new row for adding point back because of redeem
					if _, e = o.Insert(mpl); e != nil {
						o.Rollback()
						return
					}
				}
			}

			if _, e = o.Update(r.MerchantPointExpiration, "CurrentPeriodPoint", "NextPeriodPoint"); e != nil {
				o.Rollback()
				return
			}

			//if not error when updating merchant point log,create new variable to insert it into merchant point log
			recentPoint := addRecentPointTotal(v.PointValue, r.RecentPoint)

			// Insert to merchant point log if point value not 0
			if v.PointValue != 0 {
				mpl := &model.MerchantPointLog{
					PointValue:      v.PointValue,
					RecentPoint:     recentPoint,
					Status:          1,
					Note:            "Point Issued From Cancellation Redeem",
					Merchant:        r.SalesOrder.Branch.Merchant,
					SalesOrder:      r.SalesOrder,
					CreatedDate:     currentDate,
					ExpiredDate:     r.MerchantPointExpiration.CurrentPeriodDate,
					TransactionType: 7,
				}
				//start trying to insert new row for adding point back because of redeem
				if _, e = o.Insert(mpl); e != nil {
					o.Rollback()
					return
				}
				r.SalesOrder.Branch.Merchant.TotalPoint = recentPoint
				//update total point in table merchant
				if _, e = o.Update(r.SalesOrder.Branch.Merchant, "TotalPoint"); e != nil {
					o.Rollback()
					return
				}

				mps[r.SalesOrder.Branch.Merchant.ID] += v.PointValue
			}
		}

		// start create or update merchant point summary
		for i, v := range mps {
			isExist := false
			if e = o.Raw("select exists(select id from merchant_point_summary mps where merchant_id = ? and summary_date = ?)", i, currentDate.Format("2006-01-02")).QueryRow(&isExist); e != nil || (e == nil && !isExist) {
				o.Raw("insert into merchant_point_summary (merchant_id, summary_date, earned_point, redeemed_point) values (?, ?, ?, 0)", i, currentDate.Format("2006-01-02"), v).Exec()
				continue
			}

			o.Raw("update merchant_point_summary set earned_point = earned_point + ? where merchant_id = ? and summary_date = ?", v, i, currentDate.Format("2006-01-02")).Exec()
		}
		// end create or update merchant point summary
	}

	if r.PackingOrder != nil {
		for _, v := range r.SalesOrderItem {
			poi := &model.PackingOrderItem{
				PackingOrder: r.PackingOrder,
				Product:      v.Product,
			}
			if e = poi.Read("PackingOrder", "Product"); e == nil {
				poi.TotalOrder = poi.TotalOrder - v.OrderQty
				if _, e = o.Update(poi, "TotalOrder"); e != nil {
					o.Rollback()
					return nil, e
				}
			}
		}
	}

	if r.PickingOrderAssign != nil {
		r.PickingOrderAssign.Status = 7
		if _, e = o.Update(r.PickingOrderAssign, "Status"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	// return disc qty to remaining quota & discount amount to remaining budget
	for i, v := range r.SkuDiscountItems {
		if _, e = o.Update(v, "RemOverallQuota", "RemBudget"); e != nil {
			o.Rollback()
			return nil, e
		}

		if _, e = o.Update(r.SkuDiscountLogs[i], "Status"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if r.IsCreateCreditLimitLog == 1 {
		if e = log.CreditLimitLogByStaff(r.SalesOrder.Branch.Merchant, r.SalesOrder.ID, "sales_order", r.CreditLimitBefore, r.CreditLimitAfter, r.Session.Staff.ID, "cancel sales order"); e != nil {
			o.Rollback()
			return nil, e
		}
		r.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount = r.CreditLimitAfter
		if _, e = o.Update(r.SalesOrder.Branch.Merchant, "credit_limit_remaining"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if r.Session.Staff.Role.Code == "ROL0008" {
		if e = log.AuditLogByUser(r.Session.Staff, r.SalesOrder.ID, "sales_order", "cancel SO draft", r.Note); e != nil {
			o.Rollback()
			return nil, e
		}
	} else {
		if e = log.AuditLogByUser(r.Session.Staff, r.SalesOrder.ID, "sales_order", "cancel", r.Note); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	o.Commit()
	if r.SalesOrder.HasExtInvoice == 2 {
		messageNotif := &util.MessageNotification{}
		orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0005'").QueryRow(&messageNotif)
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
	} else if r.SalesOrder.HasExtInvoice == 1 {
		sie := &model.SalesInvoiceExternal{SalesOrder: r.SalesOrder}
		sie.Read("SalesOrder")
		xendit.Opt.SecretKey = util.XenditKey

		data := invoice.ExpireParams{
			ID: sie.XenditInvoiceID,
		}
		invoice.Expire(&data)
	}

	// notification FS Apps when cancelling SO
	messageNotif := &util.MessageNotification{}
	orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0014'").QueryRow(&messageNotif)
	messageNotif.Message = util.ReplaceNotificationSalesOrder(messageNotif.Message, "#customer#", r.SalesOrder.Branch.Name)
	messageNotif.Message = util.ReplaceNotificationSalesOrder(messageNotif.Message, "#sales_order_code#", r.SalesOrder.Code)

	r.SalesOrder.Salesperson.Read("ID")
	r.SalesOrder.Salesperson.User.Read("ID")
	mn := &util.ModelNotification{
		SendTo:    r.SalesOrder.Salesperson.User.SalesAppNotifToken,
		Title:     messageNotif.Title,
		Message:   messageNotif.Message,
		Type:      "1",
		RefID:     r.SalesOrder.ID,
		ServerKey: util.FieldSalesServerKeyFireBase,
		StaffID:   r.SalesOrder.Salesperson.ID,
	}
	util.PostModelNotificationFieldSales(mn)

	if r.SalesOrder.OrderChannel == 5 && r.SalesOrder.Status == 3 {
		util.SendIDToOapi(common.Encrypt(r.SalesOrder.ID), r.Session.Token)
	}

	return r.SalesOrder, nil
}

// Lock : function to change sales order is_locked into 1
func Lock(r lockRequest) (so *model.SalesOrder, e error) {
	o := orm.NewOrm()
	o.Begin()

	if r.CancelReq != 1 {
		r.SalesOrder.IsLocked = 1
		r.SalesOrder.LockedBy = r.Session.Staff.ID
		if _, e = o.Update(r.SalesOrder, "IsLocked", "LockedBy"); e != nil {
			o.Rollback()
			return nil, e
		}
	} else {
		r.SalesOrder.IsLocked = 2
		if _, e = o.Update(r.SalesOrder, "IsLocked"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	o.Commit()
	return r.SalesOrder, nil
}
func addRecentPointTotal(pointValue float64, recentPoint float64) float64 {
	return pointValue + recentPoint
}

// updatePriceDeliveredSO : function to update price when SO status delivered
func updatePriceDeliveredSO(r updatePriceRequestDeliveredSO) (so *model.SalesOrder, e error) {
	o := orm.NewOrm()
	o.Begin()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	var isCreated bool

	r.SalesOrder.DeliveryFee = r.DeliveryFee
	r.SalesOrder.TotalCharge = r.TotalCharge
	r.SalesOrder.TotalPrice = r.TotalPrice
	r.SalesOrder.LastUpdatedAt = time.Now()
	r.SalesOrder.LastUpdatedBy = r.Session.Staff.ID
	r.SalesOrder.IsLocked = 2

	if _, e = o.Update(r.SalesOrder, "TotalCharge", "TotalPrice",
		"DeliveryFee", "LastUpdatedAt", "LastUpdatedBy", "IsLocked"); e != nil {
		o.Rollback()
		return nil, e
	}

	if r.UpdateAll == 1 {
		for _, v := range r.Products {

			soi := &model.SalesOrderItem{
				SalesOrder: r.SalesOrder,
				Product:    v.Product,
			}

			if isCreated, soi.ID, e = o.ReadOrCreate(soi, "SalesOrder", "Product"); e != nil {
				o.Rollback()
				return nil, e
			}

			if !isCreated {
				soi.UnitPrice = float64(v.UnitPrice)
				soi.Subtotal = v.Subtotal
				soi.DefaultPrice = v.DefaultPrice

				if _, e = o.Update(soi, "UnitPrice", "Subtotal", "DefaultPrice"); e != nil {
					o.Rollback()
					return nil, e
				}
			}
		}
	}

	if r.IsCreateCreditLimitLog == 1 {
		if e = log.CreditLimitLogByStaff(r.SalesOrder.Branch.Merchant, r.SalesOrder.ID, "sales_order", r.CreditLimitBefore, r.CreditLimitAfter, r.Session.Staff.ID, "update sales order"); e != nil {
			o.Rollback()
			return nil, e
		}
		r.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount = r.CreditLimitAfter
		if _, e = o.Update(r.SalesOrder.Branch.Merchant, "credit_limit_remaining"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.SalesOrder.ID, "sales_order", "Update Price", r.NotePriceChange); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	// notification FS Apps when updating SO
	messageNotif := &util.MessageNotification{}
	orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0013'").QueryRow(&messageNotif)
	messageNotif.Message = util.ReplaceNotificationSalesOrder(messageNotif.Message, "#customer#", r.SalesOrder.Branch.Name)
	messageNotif.Message = util.ReplaceNotificationSalesOrder(messageNotif.Message, "#sales_order_code#", r.SalesOrder.Code)

	r.SalesOrder.Salesperson.Read("ID")
	r.SalesOrder.Salesperson.User.Read("ID")
	mn := &util.ModelNotification{
		SendTo:    r.SalesOrder.Salesperson.User.SalesAppNotifToken,
		Title:     messageNotif.Title,
		Message:   messageNotif.Message,
		Type:      "1",
		RefID:     r.SalesOrder.ID,
		ServerKey: util.FieldSalesServerKeyFireBase,
		StaffID:   r.SalesOrder.Salesperson.ID,
	}
	util.PostModelNotificationFieldSales(mn)

	return r.SalesOrder, nil
}
