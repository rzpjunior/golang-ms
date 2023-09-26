package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(TransferSkuItem))
}

// TransferSkuItem  : struct to hold TransferSkuItem  model data for database
type TransferSkuItem struct {
	ID               int64        `orm:"column(id);auto" json:"-"`
	TransferSku      *TransferSku `orm:"column(transfer_sku_id);null;rel(fk)" json:"transfer_sku"`
	Product          *Product     `orm:"column(product_id);null;rel(fk)" json:"product"`
	TransferProduct  *Product     `orm:"column(transfer_product_id);null;rel(fk)" json:"transfer_product"`
	TransferQty      float64      `orm:"column(transfer_qty)" json:"transfer_qty"`
	WasteQty         float64      `orm:"column(waste_qty)" json:"waste_qty"`
	Discrepancy      float64      `orm:"column(discrepancy)" json:"discrepancy"`
	PurchaseOrderQty float64      `orm:"-" json:"po_qty"`
	GoodsReceiptQty  float64      `orm:"-" json:"gr_qty"`
	WasteReason      int8         `orm:"column(waste_reason);null" json:"waste_reason,omitempty"`
	WasteReasonValue string       `orm:"-" json:"waste_reason_value"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *TransferSkuItem) MarshalJSON() ([]byte, error) {
	type Alias TransferSkuItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *TransferSkuItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}
