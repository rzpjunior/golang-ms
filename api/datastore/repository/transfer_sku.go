// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetTransferSku : find a single transfer sku using field and value condition.
func GetTransferSku(field string, values ...interface{}) (*model.TransferSku, error) {
	var err error
	m := new(model.TransferSku)
	o := orm.NewOrm()
	o.Using("read_only")

	if err = o.QueryTable(m).Filter(field, values...).RelatedSel("PurchaseOrder").RelatedSel("GoodsTransfer").RelatedSel("Warehouse").RelatedSel("Warehouse__Area").RelatedSel("GoodsReceipt").RelatedSel("GoodsReceipt__PurchaseOrder__Supplier").Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "TransferSkuItems", 2)

	if m.GoodsReceipt != nil {

		for _, v := range m.TransferSkuItems {

			filter := map[string]interface{}{"product_id": v.Product.ID, "goods_receipt_id": m.GoodsReceipt.ID}
			exclude := map[string]interface{}{}

			goodReceiptItem, err := GetGoodsReceiptItem(filter, exclude)

			if err != nil {
				return nil, err
			}

			v.PurchaseOrderQty = goodReceiptItem.DeliverQty
			v.GoodsReceiptQty = goodReceiptItem.ReceiveQty

			o.Raw("select value_name from glossary where `table` = ? and `attribute` = ? and `value_int` = ?", "all", "waste_reason", v.WasteReason).QueryRow(&v.WasteReasonValue)
		}
	} else {
		for _, v := range m.TransferSkuItems {
			o.Raw("select value_name from glossary where `table` = ? and `attribute` = ? and `value_int` = ?", "all", "waste_reason", v.WasteReason).QueryRow(&v.WasteReasonValue)
		}
	}

	return m, nil
}

// GetListTransferSku : function to get list transfer sku from database based on request query
func GetListTransferSku(rq *orm.RequestQuery) (ts []*model.TransferSku, total int64, err error) {
	m := new(model.TransferSku)
	q, _ := rq.QueryReadOnly(m)
	o := orm.NewOrm()
	o.Using("read_only")

	if total, err = q.Filter("status__in", 1, 2, 3).All(&ts); err != nil {
		return nil, total, err
	}

	for _, row := range m.TransferSkuItems {
		o.Raw("select value_name from glossary where `table` = ? and `attribute` = ? and `value_int` = ?", "all", "waste_reason", row.WasteReason).QueryRow(&row.WasteReasonValue)
	}

	for _, row2 := range ts {
		o.Raw("select SUM(tsi.discrepancy) discrepancy from transfer_sku_item tsi where tsi.transfer_sku_id = ?", row2.ID).QueryRow(&row2.TotalDiscrepancy)
	}

	return ts, total, nil

}

// GetTransferSkuItem : function to get single transfer sku item from database based on transfer sku data parameters
func GetTransferSkuItem(field string, values ...interface{}) (*model.TransferSku, error) {
	m := new(model.TransferSku)
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(m)

	if err := o.Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}
