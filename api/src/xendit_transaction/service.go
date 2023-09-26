// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package xendit_transaction

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
)

func Save(c fixedVaRequest) (txnXendit *model.TxnXendit, e error) {
	o := orm.NewOrm()

	o.Begin()

	merchant := &model.Merchant{ID: c.MerchantAccNum.Merchant.ID}
	txnXendit = &model.TxnXendit{
		Merchant:        merchant,
		PaymentChannel:  c.MerchantAccNum.PaymentChannel,
		Type:            1,
		AccountNumber:   c.AccountNumber,
		Amount:          c.PaidAmount,
		TransactionDate: c.TransactionDate,
		TransactionTime: c.TransactionTime,
		CreatedAt:       time.Now(),
	}

	e = txnXendit.Save()
	if e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return
}

func SaveInvoice(r invoicePaidRequest) (si *model.SalesInvoice, e error) {
	o := orm.NewOrm()

	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	o.Begin()

	var status int
	if e = orSelect.Raw("SELECT status FROM sales_order WHERE id = ?", r.SalesOrder.ID).QueryRow(&status); e != nil {
		o.Rollback()
		return nil, e
	}

	if status == 1 {
		docCode, e := util.GenerateDocCode("SI", r.SalesOrder.Branch.Code, "sales_invoice")
		if e != nil {
			o.Rollback()
			return nil, e
		}

		docPaymentCode, e := util.GenerateDocCode("SP", r.SalesOrder.Branch.Code, "sales_payment")
		if e != nil {
			o.Rollback()
			return nil, e
		}

		si = &model.SalesInvoice{
			SalesOrder:        r.SalesOrder,
			PaymentGroup:      r.SalesOrder.PaymentGroup,
			SalesTerm:         r.SalesOrder.SalesTerm,
			InvoiceTerm:       r.SalesOrder.InvoiceTerm,
			Code:              docCode,
			RecognitionDate:   time.Now(),
			Status:            2,
			BillingAddress:    r.SalesOrder.BillingAddress,
			AdjNote:           "-",
			TotalPrice:        r.SalesOrder.TotalPrice,
			TotalCharge:       r.SalesOrder.TotalCharge,
			DeliveryFee:       r.SalesOrder.DeliveryFee,
			Adjustment:        1,
			AdjAmount:         0,
			DeltaPrint:        0,
			DueDate:           time.Now(),
			CodeExt:           docCode,
			CreatedAt:         time.Now(),
			CreatedBy:         222,
			Note:              "Auto-generated invoice after Xendit Payment'",
			PointRedeemAmount: r.SalesOrder.PointRedeemAmount,
		}

		if r.SalesOrder.Voucher != nil {
			si.VoucherID = r.SalesOrder.Voucher.ID
			si.VouRedeemCode = r.SalesOrder.VouRedeemCode
			si.VouDiscAmount = r.SalesOrder.VouDiscAmount
		}

		if e = si.Save(); e == nil {
			for _, row := range r.SalesOrder.SalesOrderItems {
				item := &model.SalesInvoiceItem{
					SalesInvoice:   &model.SalesInvoice{ID: si.ID},
					Product:        row.Product,
					InvoiceQty:     row.OrderQty,
					UnitPrice:      row.UnitPrice,
					Subtotal:       row.Subtotal,
					SalesOrderItem: row,
					Note:           "Auto-generated invoice after Xendit Payment'",
				}
				e = item.Save()
			}

		}

		if e = log.AuditLogByUser(&model.Staff{ID: 222}, si.ID, "sales_invoice", "create", ""); e != nil {
			o.Rollback()
			return nil, e
		}

		// ========================================Payment====================================
		sp := &model.SalesPayment{
			Code:            docPaymentCode,
			RecognitionDate: time.Now(),
			Amount:          r.Amount,
			PaymentChannel:  r.PaymentChannel,
			PaidOff:         1,
			PaymentMethod:   &model.PaymentMethod{ID: 2},
			Note:            "Auto-generated payment after Xendit Payment",
			SalesInvoice:    si,
			CreatedAt:       time.Now(),
			CreatedBy:       222,
			Status:          2,
		}

		if e = sp.Save(); e != nil {
			o.Rollback()
			return nil, e
		}

		r.SalesOrder.Status = 12
		if e = r.SalesOrder.Save("Status"); e != nil {
			o.Rollback()
			return nil, e
		}

		tx := &model.TxnXendit{
			SalesOrder:      r.SalesOrder,
			Merchant:        r.SalesOrder.Branch.Merchant,
			PaymentChannel:  r.PaymentChannel,
			Type:            2,
			AccountNumber:   r.VaNumber,
			Amount:          r.Amount,
			TransactionDate: r.TransactionDate,
			TransactionTime: r.TransactionTime,
			CreatedAt:       time.Now(),
		}

		if e = tx.Save(); e != nil {
			o.Rollback()
			return nil, e
		}

		if e = log.AuditLogByUser(&model.Staff{ID: 222}, sp.ID, "sales_payment", "create", ""); e != nil {
			o.Rollback()
			return nil, e
		}

		if r.IsCreateCreditLimitLog {
			if e = log.CreditLimitLogByStaff(r.SalesOrder.Branch.Merchant, r.SalesOrder.ID, "sales_order", r.CreditLimitBefore, r.CreditLimitAfter, 0, "auto paid sales order"); e != nil {
				o.Rollback()
				return nil, e
			}
			r.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount = r.CreditLimitAfter
			if _, e = o.Update(r.SalesOrder.Branch.Merchant, "credit_limit_remaining"); e != nil {
				o.Rollback()
				return nil, e
			}
		}

		o.Commit()

		messageNotif := &util.MessageNotification{}

		orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0002'").QueryRow(&messageNotif)
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
	} else {
		dt := &model.DocumentTemp{
			SalesOrderID:   r.SalesOrder.ID,
			SalesOrderCode: r.SalesOrder.Code,
			CreatedAt:      time.Now(),
			FromCronJob:    r.FromCronJob,
		}

		if e = dt.Save(); e != nil {
			o.Rollback()
			return nil, e
		}

		o.Commit()
	}

	return
}

func ExpiredInvoice(r invoiceExpiredRequest) (u *model.TxnXendit, e error) {
	if r.SalesInvoiceExternal.SalesOrder.Status == 1 {
		orSelect := orm.NewOrm()
		orSelect.Using("read_only")

		o := orm.NewOrm()
		o.Begin()

		so := &model.SalesOrder{
			ID:         r.SalesInvoiceExternal.SalesOrder.ID,
			Status:     3,
			CancelType: 2,
		}

		if e = so.Save("Status", "CancelType"); e != nil {
			o.Rollback()
			return nil, e
		}

		if r.SalesInvoiceExternal.SalesOrder.Voucher != nil {
			if _, e = orm.NewOrm().Raw("UPDATE voucher_log SET status = 3 WHERE sales_order_id = ? AND voucher_id = ?;", r.SalesInvoiceExternal.SalesOrder.ID, r.SalesInvoiceExternal.SalesOrder.Voucher.ID).Exec(); e != nil {
				o.Rollback()
				return
			}

			if _, e = orm.NewOrm().Raw("UPDATE voucher set rem_overall_quota = rem_overall_quota + 1 where id = ?", r.SalesInvoiceExternal.SalesOrder.Voucher.ID).Exec(); e != nil {
				o.Rollback()
				return
			}
		}

		if r.SalesInvoiceExternal.SalesOrder.PointRedeemID != 0 && r.SalesInvoiceExternal.SalesOrder.PointRedeemAmount != 0 {
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

				// check the point expired or not
				isPointExpired := time.Now().After(v.ExpiredDate)

				// This condition for back up data existing that didn't record current point used and next point used
				isDataExisting := v.CurrentPointUsed == 0 && v.NextPointUsed == 0
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
							Merchant:        r.SalesInvoiceExternal.SalesOrder.Branch.Merchant,
							SalesOrder:      r.SalesInvoiceExternal.SalesOrder,
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
				recentPoint := v.PointValue + float64(r.RecentPoint)

				// Insert to merchant point log if point value not 0
				if v.PointValue != 0 {
					mpl := &model.MerchantPointLog{
						PointValue:      v.PointValue,
						RecentPoint:     recentPoint,
						Status:          1,
						Note:            "Point Issued From Cancellation Redeem",
						Merchant:        r.SalesInvoiceExternal.SalesOrder.Branch.Merchant,
						SalesOrder:      r.SalesInvoiceExternal.SalesOrder,
						CreatedDate:     currentDate,
						TransactionType: 7,
					}
					//start trying to insert new row for adding point back because of redeem
					if _, e = o.Insert(mpl); e != nil {
						o.Rollback()
						return
					}
					r.SalesInvoiceExternal.SalesOrder.Branch.Merchant.TotalPoint = recentPoint
					//update total point in table merchant
					if _, e = o.Update(r.SalesInvoiceExternal.SalesOrder.Branch.Merchant, "TotalPoint"); e != nil {
						o.Rollback()
						return
					}

					mps[r.SalesInvoiceExternal.SalesOrder.Branch.Merchant.ID] += v.PointValue

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

		ivx := &model.SalesInvoiceExternal{
			ID:          r.SalesInvoiceExternal.ID,
			CancelledAt: time.Now(),
		}

		if e = ivx.Save("CancelledAt"); e != nil {
			o.Rollback()
			return
		}

		if r.IsCreateCreditLimitLog {
			if e = log.CreditLimitLogByStaff(r.SalesInvoiceExternal.SalesOrder.Branch.Merchant, r.SalesInvoiceExternal.SalesOrder.ID, "sales_order", r.CreditLimitBefore, r.CreditLimitAfter, 0, "auto cancel sales order"); e != nil {
				o.Rollback()
				return nil, e
			}
			r.SalesInvoiceExternal.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount = r.CreditLimitAfter
			if _, e = o.Update(r.SalesInvoiceExternal.SalesOrder.Branch.Merchant, "credit_limit_remaining"); e != nil {
				o.Rollback()
				return nil, e
			}
		}

		if e = log.AuditLogByUser(&model.Staff{ID: 222}, so.ID, "sales_order", "cancel", "Auto-cancelled due to expired xenInvoice"); e != nil {
			o.Rollback()
			return nil, e
		}

		o.Commit()

		messageNotif := &util.MessageNotification{}
		orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0005'").QueryRow(&messageNotif)
		messageNotif.Message = util.ReplaceNotificationSalesOrder(messageNotif.Message, "#sales_order_code#", r.SalesInvoiceExternal.SalesOrder.Code)

		mn := &util.ModelNotification{
			SendTo:     r.SalesInvoiceExternal.SalesOrder.Branch.Merchant.UserMerchant.FirebaseToken,
			Title:      messageNotif.Title,
			Message:    messageNotif.Message,
			Type:       "1",
			RefID:      r.SalesInvoiceExternal.SalesOrder.ID,
			MerchantID: r.SalesInvoiceExternal.SalesOrder.Branch.Merchant.ID,
			ServerKey:  util.ServerKeyFireBase,
		}
		util.PostModelNotification(mn)
	}

	return
}
