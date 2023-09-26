// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/util"
)

func init() {
	orm.RegisterModel(new(GoodsReceipt))
}

// GoodsReceipt: struct to hold model data for database
type GoodsReceipt struct {
	ID                  int64     `orm:"column(id);auto" json:"-"`
	Code                string    `orm:"column(code);size(50);null" json:"code"`
	AtaDate             time.Time `orm:"column(ata_date)" json:"ata_date"`
	AtaTime             string    `orm:"column(ata_time)" json:"ata_time"`
	TotalWeight         float64   `orm:"column(total_weight)" json:"total_weight"`
	Note                string    `orm:"column(note)" json:"note"`
	Status              int8      `orm:"column(status);null" json:"status"`
	InboundType         int8      `orm:"column(inbound_type);null" json:"inbound_type"`
	ValidSupplierReturn int8      `orm:"column(valid_supplier_return);null" json:"valid_supplier_return"`
	CreatedAt           time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy           int64     `orm:"column(created_by)" json:"created_by"`
	ConfirmedAt         time.Time `orm:"column(confirmed_at);type(timestamp);null" json:"confirmed_at"`
	ConfirmedBy         int64     `orm:"column(confirmed_by)" json:"confirmed_by"`
	Locked              int8      `orm:"column(locked)" json:"locked"`
	StockType           int8      `orm:"column(stock_type)" json:"stock_type"`
	LockedBy            int64     `orm:"column(locked_by)" json:"-"`
	LockedByObj         *Staff    `orm:"-" json:"locked_by"`
	UpdatedAt           time.Time `orm:"column(updated_at)" json:"updated_at"`
	UpdatedBy           int64     `orm:"column(updated_by)" json:"-"`
	UpdatedByObj        *Staff    `orm:"-" json:"updated_by"`

	Warehouse         *Warehouse          `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse"`
	PurchaseOrder     *PurchaseOrder      `orm:"column(purchase_order_id);null;rel(fk)" json:"purchase_order"`
	GoodsTransfer     *GoodsTransfer      `orm:"column(goods_transfer_id);null;rel(fk)" json:"goods_transfer"`
	SupplierReturn    []*SupplierReturn   `orm:"-" json:"supplier_return"`
	DebitNote         []*DebitNote        `orm:"-" json:"debit_note"`
	GoodsReceiptItems []*GoodsReceiptItem `orm:"reverse(many)" json:"goods_receipt_items,omitempty"`
	TransferSKU       []*TransferSku      `orm:"-" json:"transfer_SKU"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *GoodsReceipt) MarshalJSON() ([]byte, error) {
	type Alias GoodsReceipt

	return json.Marshal(&struct {
		ID            string `json:"id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		StatusConvert: util.ConvertStatusDoc(m.Status),
		Alias:         (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *GoodsReceipt) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *GoodsReceipt) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
