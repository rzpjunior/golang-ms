// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetGoodsReceipt find a single data price set using field and value condition.
func GetGoodsReceipt(field string, values ...interface{}) (*model.GoodsReceipt, error) {
	m := new(model.GoodsReceipt)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel("PurchaseOrder").RelatedSel("PurchaseOrder__Supplier").RelatedSel("Warehouse").RelatedSel("Warehouse__Area").RelatedSel("GoodsTransfer").RelatedSel("GoodsTransfer__Origin").RelatedSel("GoodsTransfer__Destination").Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "GoodsReceiptItems", 1)

	return m, nil
}

// GetGoodsReceipts : function to get data from database based on parameters
func GetGoodsReceipts(rq *orm.RequestQuery) (grs []*model.GoodsReceipt, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.GoodsReceipt))

	if total, err = q.Filter("status__in", 1, 2, 3).All(&grs, rq.Fields...); err == nil {
		return grs, total, nil
	}

	return nil, total, err
}

// GetFilterGoodsReceipts : function to get data from database based on parameters with filtered permission
func GetFilterGoodsReceipts(rq *orm.RequestQuery) (grs []*model.GoodsReceipt, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")
	q, _ := rq.QueryReadOnly(new(model.GoodsReceipt))

	if total, err = q.Filter("status", 2).All(&grs, rq.Fields...); err != nil {
		return nil, total, err

	}

	return grs, total, nil

}

// ValidGoodsReceipt : function to check if id is valid in database
func ValidGoodsReceipt(id int64) (goodsReceipt *model.GoodsReceipt, e error) {
	goodsReceipt = &model.GoodsReceipt{ID: id}
	e = goodsReceipt.Read("ID")

	return
}

// CheckGoodsReceiptProductStatus : function to check if product is valid in table
func CheckGoodsReceiptProductStatus(productID int64, status int8, warehouseArr ...string) (*model.GoodsReceiptItem, int64, error) {
	var err error
	o := orm.NewOrm()
	o.Using("read_only")

	gri := new(model.GoodsReceiptItem)
	q := o.QueryTable(gri).RelatedSel("GoodsReceipt").Filter("GoodsReceipt__Status", status).Filter("product_id", productID)

	if len(warehouseArr) > 0 {
		q = q.Exclude("GoodsReceipt__Warehouse__id__in", warehouseArr)
	}

	if total, err := q.All(gri); err == nil {
		return gri, total, nil
	}

	return nil, 0, err
}

// CheckGoodsReceiptData : function to check data based on filter and exclude parameters
func CheckGoodsReceiptData(filter, exclude map[string]interface{}) (gr []*model.GoodsReceipt, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.GoodsReceipt))

	for i, v := range filter {
		o = o.Filter(i, v)
	}

	for i, v := range exclude {
		o = o.Exclude(i, v)
	}

	if total, err := o.All(&gr); err == nil {
		return gr, total, nil
	}

	return nil, 0, err
}

// GetGoodsReceipts : function to get data from database based on parameters
func GetGoodsReceiptItems(rq *orm.RequestQuery) (gris []*model.GoodsReceiptItem, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.GoodsReceiptItem))

	if total, err = q.All(&gris, rq.Fields...); err == nil {
		return gris, total, nil
	}

	return nil, total, err
}

// GetGoodsReceiptItem : function to get single goods receipt item from database based on goods receipt id & product id
func GetGoodsReceiptItem(filter map[string]interface{}, exclude map[string]interface{}) (gri *model.GoodsReceiptItem, err error) {
	m := new(model.GoodsReceiptItem)
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(m)

	if len(filter) > 0 {
		for i, v := range filter {
			o = o.Filter(i, v)
		}
	}

	if len(exclude) > 0 {
		for i, v := range exclude {
			o = o.Filter(i, v)
		}
	}

	if _, err := o.All(m); err == nil {
		return m, nil
	}

	return nil, err
}

// GetGoodsReceipts : function to get data from database based on parameters with filtered permission
func GetFilterGoodsReceiptItems(rq *orm.RequestQuery) (gris []*model.GoodsReceiptItem, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.GoodsReceiptItem))

	if total, err = q.All(&gris, rq.Fields...); err == nil {
		return gris, total, nil
	}

	return nil, total, err
}

func GetGoodsReceiptWithProductGroup(field string, values ...interface{}) (*model.GoodsReceipt, error) {
	m := new(model.GoodsReceipt)
	o := orm.NewOrm()

	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel("PurchaseOrder").RelatedSel("PurchaseOrder__Supplier").RelatedSel("PurchaseOrder__Supplier__PaymentTerm").RelatedSel("PurchaseOrder__TermPaymentPur").RelatedSel("Warehouse").RelatedSel("Warehouse__Area").RelatedSel("GoodsTransfer").RelatedSel("GoodsTransfer__Origin").RelatedSel("GoodsTransfer__Destination").Limit(1).One(m); err != nil {
		return nil, err
	}

	// region get supplier return data
	if _, err := o.QueryTable(new(model.SupplierReturn)).Filter("goods_receipt_id", m.ID).All(&m.SupplierReturn); err != nil {
		m.SupplierReturn = nil
	}
	// endregion

	// region get tranfer SKU return data
	if _, err := o.QueryTable(new(model.TransferSku)).Filter("goods_receipt_id", m.ID).All(&m.TransferSKU); err != nil {
		m.TransferSKU = nil
	}
	// endregion

	// region get supplier return data
	for _, v := range m.SupplierReturn {
		dn := &model.DebitNote{
			SupplierReturn: v,
		}
		dn.Read("SupplierReturn")
		dn.Read("ID")

		m.DebitNote = append(m.DebitNote, dn)
	}
	// endregion

	o.LoadRelated(m, "GoodsReceiptItems", 1)

	for _, v := range m.GoodsReceiptItems {
		if err := v.Product.Uom.Read("ID"); err != nil {
			return nil, err
		}
		var isExist int8
		o.Raw("select exists(select ts.id from transfer_sku ts join transfer_sku_item tsi on ts.id = tsi.transfer_sku_id "+
			"where ts.goods_receipt_id = ? and ts.status in(1,2) and tsi.product_id = ?)", m.ID, v.Product.ID).QueryRow(&isExist)
		if isExist > 0 {
			v.IsDisabled = 1
		}

		productGroupItem, _ := GetProductGroupItem("product_id", v.Product.ID)
		if productGroupItem != nil {
			v.ProductGroup = productGroupItem.ProductGroup
		}
	}

	if m.PurchaseOrder != nil {
		o.LoadRelated(m.PurchaseOrder, "PurchaseOrderItems", 1)
	}

	if m.GoodsTransfer != nil {
		o.LoadRelated(m.GoodsTransfer, "GoodsTransferItem", 1)
		for i, v := range m.GoodsReceiptItems {
			v.GoodsTransferItem = m.GoodsTransfer.GoodsTransferItem[i]
		}
	}

	var err error
	if m.LockedBy != 0 {
		if m.LockedByObj, err = ValidStaff(m.LockedBy); err != nil {
			return nil, err
		}
	}

	if m.UpdatedBy != 0 {
		if m.UpdatedByObj, err = ValidStaff(m.UpdatedBy); err != nil {
			return nil, err
		}
	}

	return m, nil
}
