// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetSalesInvoice find a single data price set using field and value condition.
func GetSalesInvoice(field string, values ...interface{}) (si *model.SalesInvoice, err error) {
	m := new(model.SalesInvoice)
	o := orm.NewOrm()
	o.Using("read_only")

	if err = o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "SalesInvoiceItems", 2)
	for _, v := range m.SalesInvoiceItems {
		o.Raw("select doi.* from delivery_order_item doi join delivery_order do on doi.delivery_order_id = do.id where do.status in (1,2,5,6,7) and doi.sales_order_item_id = ?", v.SalesOrderItem.ID).QueryRow(&v.SalesOrderItem.DeliveryOrderItem)
		if v.SalesOrderItem.SkuDiscountItem != nil {
			o.LoadRelated(v.SalesOrderItem, "SkuDiscountItem")
			o.LoadRelated(v.SalesOrderItem.SkuDiscountItem, "SkuDiscountItemTiers")
		}
	}
	o.LoadRelated(m, "SalesPayment", 1)

	m.RemainingAmount, err = CheckRemainingSalesInvoiceAmount(m.ID)

	o.Raw("SELECT type from voucher where id = ?", m.VoucherID).QueryRow(&m.VoucherType)
	o.Raw("SELECT * from merchant_acc_num where merchant_id = ?", m.SalesOrder.Branch.Merchant.ID).QueryRows(&m.MerchantAccNum)
	if m.MerchantAccNum != nil {
		for _, man := range m.MerchantAccNum {
			man.PaymentChannel.Read("ID")
			if man.PaymentChannel.ID == 6 {
				m.XenditBCA = man.AccountNumber
			} else if man.PaymentChannel.ID == 7 {
				m.XenditPermata = man.AccountNumber
			}
		}
	}

	return m, nil
}

// GetSalesInvoices : function to get data from database based on parameters
func GetSalesInvoices(rq *orm.RequestQuery, date []string, merchant ...int64) (m []*model.SalesInvoice, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SalesInvoice))

	// if QueryParam merchant is exist
	condMerchant := q.GetCond()
	if len(merchant) > 0 {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("salesorder__branch__merchant__id__in", merchant)

		condMerchant = condMerchant.AndCond(cond1)
	}
	q = q.SetCond(condMerchant)

	// if QueryParam delivery date is exist
	condDeliveryDate := q.GetCond()
	if len(date) > 0 {
		cond2 := orm.NewCondition()
		cond2 = cond2.And("salesorder__delivery_date__between", date)

		condDeliveryDate = condDeliveryDate.AndCond(cond2)
	}
	q = q.SetCond(condDeliveryDate)

	if total, err = q.Exclude("status", 4).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SalesInvoice
	if _, err = q.Exclude("status", 4).All(&mx, rq.Fields...); err == nil {

		for _, v := range mx {
			v.TotalPaid, err = CheckTotalPaidPaymentAmount(v.ID)
			v.RemainingAmount, err = CheckRemainingSalesInvoiceAmount(v.ID)
		}
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterSalesInvoices : function to get data from database based on parameters with filtered permission
func GetFilterSalesInvoices(rq *orm.RequestQuery) (m []*model.SalesInvoice, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SalesInvoice))

	if total, err = q.Exclude("status", 4).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SalesInvoice
	if _, err = q.Exclude("status", 4).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidSalesInvoice : function to check if id is valid in database
func ValidSalesInvoice(id int64) (salesInvoice *model.SalesInvoice, e error) {
	salesInvoice = &model.SalesInvoice{ID: id}
	e = salesInvoice.Read("ID")

	return
}

// CheckRemainingSalesInvoiceAmount : function to check remaining invoice amount
func CheckRemainingSalesInvoiceAmount(invoiceID int64) (remainingAmount float64, e error) {
	o := orm.NewOrm()
	o.Using("read_only")

	if e = o.Raw("select case when sum(amount) is null then si.total_charge else si.total_charge - sum(amount) end remainingAmount from sales_invoice si "+
		"left join "+
		"sales_payment sp on si.id = sp.sales_invoice_id and sp.status IN (2,5) "+
		"where si.id = ? ", invoiceID).QueryRow(&remainingAmount); e != nil {
		return 0, e
	}

	return
}

// CheckRemainingSalesInvoiceAmountByStatus: function to check remaining invoice amount by status
func CheckRemainingSalesInvoiceAmountByStatus(invoiceID int64, status []int64) (remainingAmount float64, e error) {
	o := orm.NewOrm()
	o.Using("read_only")

	query := "select case when sum(amount) is null then si.total_charge else si.total_charge - sum(amount) end remainingAmount from sales_invoice si " +
		"left join sales_payment sp on si.id = sp.sales_invoice_id and sp.status IN ("

	for _, v := range status {
		query = query + strconv.FormatInt(v, 10) + ","
	}

	query = strings.TrimSuffix(query, ",") + ") where si.id = ? "

	if e = o.Raw(query, invoiceID).QueryRow(&remainingAmount); e != nil {
		return 0, e
	}

	return
}

// CheckTotalPaidPaymentAmount : function to check total paid sales payment that related with sales invoice
func CheckTotalPaidPaymentAmount(invoiceID int64) (totalPaid float64, e error) {
	o := orm.NewOrm()

	var amount []float64
	// This query 'select' only use in this unique case, in other case you must use 'orm read only' to get data; example case bulk confirm payment
	if _, e = o.Raw("SELECT amount FROM sales_payment WHERE status IN (2,5) AND sales_invoice_id = ?", invoiceID).QueryRows(&amount); e != nil {
		return 0, nil
	}

	for _, v := range amount {
		totalPaid += v
	}

	return
}

// CheckSalesInvoicesData : function to get all sales invoice data based on filter and exclude parameters
func CheckSalesInvoicesData(filter, exclude map[string]interface{}) (m []*model.SalesInvoice, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	qtb := o.QueryTable(new(model.SalesInvoice))

	for k, v := range filter {
		qtb = qtb.Filter(k, v)
	}

	for k, v := range exclude {
		qtb = qtb.Exclude(k, v)
	}

	if total, err = qtb.OrderBy("created_at").All(&m); err != nil {
		return nil, 0, err
	}

	return m, total, nil
}

// GetSalesInvoiceForPrints : function to get data from database based on parameters
func GetSalesInvoiceForPrints(rq *orm.RequestQuery) (m []*model.SalesInvoice, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SalesInvoice))
	o := orm.NewOrm()
	o.Using("read_only")

	if total, err = q.Exclude("status__in", 4, 2, 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SalesInvoice
	if _, err = q.Exclude("status__in", 4, 2, 3).RelatedSel().All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			o.LoadRelated(v, "SalesInvoiceItems", 2)
			o.LoadRelated(v, "SalesPayment", 1)

			v.RemainingAmount, err = CheckRemainingSalesInvoiceAmount(v.ID)

			o.Raw("SELECT type from voucher where id = ?", v.VoucherID).QueryRow(&v.VoucherType)
			o.Raw("SELECT * from merchant_acc_num where merchant_id = ?", v.SalesOrder.Branch.Merchant.ID).QueryRows(&v.MerchantAccNum)
			if v.MerchantAccNum != nil {
				for _, man := range v.MerchantAccNum {
					man.PaymentChannel.Read("ID")
					if man.PaymentChannel.ID == 6 {
						v.XenditBCA = man.AccountNumber
					} else if man.PaymentChannel.ID == 7 {
						v.XenditPermata = man.AccountNumber
					}
				}
			}
		}
		return mx, total, nil
	}

	return nil, total, err
}

// Get Grand Total Charge Where SI Status 1, 6
func GetGrandTotalChargeSI(merchant_id int64) (totalCharge float64, e error) {
	o := orm.NewOrm()
	o.Using("read_only")
	q := "SELECT SUM(si.total_charge) total_charge FROM branch b " +
		"LEFT JOIN sales_order so ON so.branch_id = b.id " +
		"LEFT JOIN sales_invoice si ON si.sales_order_id = so.id " +
		"WHERE b.merchant_id = ? and si.status in (1,6);"

	if e = o.Raw(q, merchant_id).QueryRow(&totalCharge); e != nil {
		return totalCharge, e
	}

	return totalCharge, nil
}

// GetMerchantSalesInvoices : function to get data from database based on parameters
func GetMerchantSalesInvoices(rq *orm.RequestQuery, merchant int64) (m []*model.SalesInvoice, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SalesInvoice))

	var (
		loc, _      = time.LoadLocation("Asia/Jakarta")
		currentTime = time.Now().In(loc)
	)

	// set QueryParam merchant id
	condMerchant := q.GetCond()
	cond1 := orm.NewCondition()
	cond1 = cond1.And("salesorder__branch__merchant__id__in", merchant)

	condMerchant = condMerchant.AndCond(cond1)
	q = q.SetCond(condMerchant)

	if total, err = q.Exclude("status", 4).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SalesInvoice
	if _, err = q.Exclude("status", 4).All(&mx, rq.Fields...); err == nil {

		for _, v := range mx {
			// assign value to total paid
			v.TotalPaid, err = CheckTotalPaidPaymentAmount(v.ID)

			// assign value to payment percentage
			v.PaymentPercentage = v.TotalPaid / v.TotalCharge * 100

			// generate description status
			diff := int(v.DueDate.Sub(currentTime).Hours() / 24)
			if v.Status == 2 {
				v.StatusDescription = "paid"
			} else if diff > 0 {
				v.StatusDescription = fmt.Sprintf("%d days left", diff)
			} else if diff < 0 {
				v.StatusDescription = fmt.Sprintf("due %d days", diff)
			} else if diff == 0 {
				v.StatusDescription = "due today"
			}
		}
		return mx, total, nil
	}

	return nil, total, err
}
