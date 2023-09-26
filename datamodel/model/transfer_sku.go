package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/util"
)

func init() {
	orm.RegisterModel(new(TransferSku))
}

// TransferSku : struct to hold TransferSku model data for database
type TransferSku struct {
	ID               int64     `orm:"column(id);auto" json:"-"`
	Code             string    `orm:"column(code);size(50);null" json:"code"`
	TotalTransferQty float64   `orm:"column(total_transfer_qty)" json:"total_transfer_qty"`
	TotalWasteQty    float64   `orm:"column(total_waste_qty)" json:"total_waste_qty"`
	Status           int8      `orm:"column(status)" json:"status"`
	CreatedAt        time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy        *Staff    `orm:"column(created_by);null;rel(fk)" json:"created_by"`
	ConfirmedAt      time.Time `orm:"column(confirmed_at);type(timestamp);null" json:"confirmed_at"`
	ConfirmedBy      *Staff    `orm:"column(confirmed_by);null;rel(fk)" json:"confirmed_by"`
	RecognitionDate  time.Time `orm:"column(recognition_date)" json:"recognition_date"`
	Note             string    `orm:"column(note);size(250);null" json:"note"`
	TotalDiscrepancy float64   `orm:"-" json:"total_discrepancy"`

	GoodsReceipt     *GoodsReceipt      `orm:"column(goods_receipt_id);null;rel(fk)" json:"goods_receipt"`
	Warehouse        *Warehouse         `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse"`
	PurchaseOrder    *PurchaseOrder     `orm:"column(purchase_order_id);null;rel(fk)" json:"purchase_order"`
	GoodsTransfer    *GoodsTransfer     `orm:"column(goods_transfer_id);null;rel(fk)" json:"goods_transfer"`
	TransferSkuItems []*TransferSkuItem `orm:"reverse(many)" json:"transfer_sku_items,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *TransferSku) MarshalJSON() ([]byte, error) {
	type Alias TransferSku

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
func (m *TransferSku) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *TransferSku) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
