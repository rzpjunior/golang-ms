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
	orm.RegisterModel(new(GoodsTransfer))
}

// GoodsTransfer: struct to hold model data for database
type GoodsTransfer struct {
	ID                 int64     `orm:"column(id);auto" json:"-"`
	Code               string    `orm:"column(code);size(50);null" json:"code"`
	RequestDate        time.Time `orm:"column(request_date)" json:"request_date"`
	RecognitionDate    time.Time `orm:"column(recognition_date)" json:"recognition_date"`
	EtaDate            time.Time `orm:"column(eta_date)" json:"eta_date"`
	EtaTime            string    `orm:"column(eta_time)" json:"eta_time"`
	AtaDate            time.Time `orm:"column(ata_date)" json:"ata_date"`
	AtaTime            string    `orm:"column(ata_time)" json:"ata_time"`
	AdditionalCost     float64   `orm:"column(addtl_cost)" json:"additional_cost"`
	AdditionalCostNote string    `orm:"column(addtl_cost_note)" json:"additional_cost_note"`
	StockType          int8      `orm:"column(stock_type)" json:"stock_type"`
	TotalCost          float64   `orm:"column(total_cost)" json:"total_cost"`
	TotalCharge        float64   `orm:"column(total_charge)" json:"total_charge"`
	TotalWeight        float64   `orm:"column(total_weight)" json:"total_weight"`
	Note               string    `orm:"column(note)" json:"note"`
	Status             int8      `orm:"column(status);null" json:"status"`
	Locked             int8      `orm:"column(locked)" json:"locked"`
	LockedBy           int64     `orm:"column(locked_by)" json:"-"`
	LockedByObj        *Staff    `orm:"-" json:"locked_by"`
	TotalSku           int64     `orm:"-" json:"total_sku"`
	UpdatedAt          time.Time `orm:"column(updated_at)" json:"updated_at"`
	UpdatedBy          int64     `orm:"column(updated_by)" json:"-"`
	UpdatedByObj       *Staff    `orm:"-" json:"updated_by"`

	Origin            *Warehouse           `orm:"column(origin_id);null;rel(fk)" json:"origin"`
	Destination       *Warehouse           `orm:"column(destination_id);null;rel(fk)" json:"destination"`
	GoodsTransferItem []*GoodsTransferItem `orm:"reverse(many)" json:"goods_transfer_item,omitempty"`
	GoodsReceipt      []*GoodsReceipt      `orm:"reverse(many)" json:"goods_receipt"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *GoodsTransfer) MarshalJSON() ([]byte, error) {
	type Alias GoodsTransfer

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
func (m *GoodsTransfer) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *GoodsTransfer) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
