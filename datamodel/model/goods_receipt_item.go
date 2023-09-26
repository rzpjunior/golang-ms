// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/util"
)

func init() {
	orm.RegisterModel(new(GoodsReceiptItem))
}

// GoodsReceiptItem: struct to hold model data for database
type GoodsReceiptItem struct {
	ID           int64   `orm:"column(id);auto" json:"-"`
	DeliverQty   float64 `orm:"column(deliver_qty)" json:"delivery_qty"`
	RejectQty    float64 `orm:"column(reject_qty)" json:"reject_qty"`
	ReceiveQty   float64 `orm:"column(receive_qty)" json:"receive_qty"`
	Weight       float64 `orm:"column(weight)" json:"weight"`
	Note         string  `orm:"column(note)" json:"note"`
	RejectReason int8    `orm:"column(reject_reason)" json:"reject_reason"`
	IsDisabled   int8    `orm:"-" json:"is_disabled"`

	GoodsReceipt      *GoodsReceipt      `orm:"column(goods_receipt_id);null;rel(fk)" json:"goods_receipt"`
	PurchaseOrderItem *PurchaseOrderItem `orm:"column(purchase_order_item_id);null;rel(fk)" json:"purchase_order_item"`
	GoodsTransferItem *GoodsTransferItem `orm:"-" json:"goods_transfer_item"`
	Product           *Product           `orm:"column(product_id);null;rel(fk)" json:"product"`
	ProductGroup      *ProductGroup      `orm:"-" json:"product_group"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *GoodsReceiptItem) MarshalJSON() ([]byte, error) {
	type Alias GoodsReceiptItem

	return json.Marshal(&struct {
		ID                  string `json:"id"`
		RejectReasonConvert string `json:"reject_reason_convert"`
		*Alias
	}{
		ID:                  common.Encrypt(m.ID),
		RejectReasonConvert: util.ConvertRejectReasonDoc(m.RejectReason),
		Alias:               (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *GoodsReceiptItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *GoodsReceiptItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
