// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"strings"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetPurchaseInvoice find a single data price set using field and value condition.
func GetPurchaseInvoice(field string, values ...interface{}) (pi *model.PurchaseInvoice, err error) {
	m := new(model.PurchaseInvoice)
	o := orm.NewOrm()
	var qMark string
	o.Using("read_only")

	var paidAmount float64

	if err = o.QueryTable(m).Filter(field, values...).RelatedSel("PurchaseOrder").RelatedSel("PurchaseOrder__SupplierBadge").RelatedSel("PurchaseOrder__TermPaymentPur").RelatedSel("PurchaseOrder__Supplier").RelatedSel("PurchaseTerm").Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "PurchaseInvoiceItems", 1)
	for _, v := range m.PurchaseInvoiceItems {
		o.Raw("select * from goods_receipt_item gri where gri.purchase_order_item_id =?", v.PurchaseOrderItem.ID).QueryRow(&v.GoodsReceiptItem)
	}

	dnIds := strings.Split(m.DebitNoteIDs, ",")
	for _, _ = range dnIds {
		qMark = qMark + "?,"
	}

	qMark = strings.TrimSuffix(qMark, ",")

	o.Raw("SELECT dn.id, dn.code, dn.status, dn.total_price FROM debit_note dn where dn.id in ("+qMark+")", dnIds).QueryRows(&m.DebitNote)

	m.IsPaid, paidAmount, err = CheckPurchasePaymentAmount(m.ID)
	m.RemainingAmount = m.TotalCharge - paidAmount

	return m, nil
}

// GetPurchaseInvoices : function to get data from database based on parameters
func GetPurchaseInvoices(rq *orm.RequestQuery, ataDate ...string) (pis []*model.PurchaseInvoice, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PurchaseInvoice))

	if len(ataDate) > 0 && ataDate[0] != "" {
		q = q.RelatedSel(1).Filter("purchaseorder__goodsreceipt__ata_date__between", ataDate).Exclude("purchaseorder__goodsreceipt__status", 3)
	}

	if total, err = q.Filter("status__in", 1, 2, 3, 6).All(&pis, rq.Fields...); err == nil {
		return pis, total, nil
	}

	return nil, total, err
}

// GetFilterPurchaseInvoices : function to get data from database based on parameters with filtered permission
func GetFilterPurchaseInvoices(rq *orm.RequestQuery) (pis []*model.PurchaseInvoice, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.GoodsReceipt))

	if total, err = q.Filter("status__in", 1, 2, 3, 6).All(&pis, rq.Fields...); err == nil {
		return pis, total, nil
	}

	return nil, total, err
}

// ValidPurchaseInvoice : function to check if id is valid in database
func ValidPurchaseInvoice(id int64) (purchaseInvoice *model.PurchaseInvoice, e error) {
	purchaseInvoice = &model.PurchaseInvoice{ID: id}
	e = purchaseInvoice.Read("ID")

	return
}

// CheckPurchaseInvoiceProductStatus : function to check if product is valid in table
func CheckPurchaseInvoiceProductStatus(productID int64, status int8, warehouseArr ...string) (*model.PurchaseInvoiceItem, error) {
	o := orm.NewOrm()
	o.Using("read_only")

	gri := new(model.PurchaseInvoiceItem)
	var err error

	q := o.QueryTable(gri).RelatedSel("GoodsReceipt").Filter("GoodsReceipt__Status", status).Filter("product_id", productID)

	if len(warehouseArr) > 0 {
		q.Filter("GoodsReceipt__Warehouse__id__in", warehouseArr)
	}

	if err = q.One(gri); err == nil {
		return gri, nil
	}

	return nil, err
}

// CheckPurchasePaymentAmount : function to check payment that had been paid
func CheckPurchasePaymentAmount(invoiceID int64) (isPaid int8, paidAmount float64, e error) {
	o := orm.NewOrm()
	o.Using("read_only")

	if e = o.Raw("SELECT CASE WHEN pp.id IS NOT NULL THEN 1 ELSE 2 END isPaid, SUM(CASE WHEN pp.status = 2 THEN pp.amount ELSE 0 END) paid_amount "+
		"FROM purchase_invoice pi "+
		"LEFT JOIN purchase_payment pp ON pi.id = pp.purchase_invoice_id "+
		"WHERE pi.id = ? "+
		"GROUP BY pi.id, CASE WHEN pp.id IS NOT NULL THEN 1 ELSE 2 END "+
		"ORDER BY pi.id , pp.status", invoiceID).QueryRow(&isPaid, &paidAmount); e == nil {

		return isPaid, paidAmount, nil
	}

	return 2, 0, e
}

// GetDataPurchaseInvoice : function to get data based on filter and exclude parameters
func GetDataPurchaseInvoice(filter map[string]interface{}, exclude map[string]interface{}) (pi []*model.PurchaseInvoice, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	q := o.QueryTable(new(model.PurchaseInvoice))

	for k, v := range filter {
		q = q.Filter(k, v)
	}

	for k, v := range exclude {
		q = q.Exclude(k, v)
	}

	if total, err := q.All(&pi); err == nil {
		return pi, total, nil
	}

	return nil, 0, err
}
