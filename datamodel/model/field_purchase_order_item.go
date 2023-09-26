package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(FieldPurchaseOrderItem))
}

// FieldPurchaseOrderItem  : struct to hold field_puchase_order_item model data for database
type FieldPurchaseOrderItem struct {
	ID                 int64               `orm:"column(id);auto" json:"-"`
	FieldPurchaseOrder *FieldPurchaseOrder `orm:"column(field_purchase_order_id);null;rel(fk)" json:"field_purchase_order,omitempty"`
	PurchaseOrderItem  *PurchaseOrderItem  `orm:"column(purchase_order_item_id);null;rel(fk)" json:"purchase_order_item,omitempty"`
	Product            *Product            `orm:"column(product_id);null;rel(fk)" json:"product,omitempty"`
	PurchaseQty        float64             `orm:"column(purchase_qty)" json:"purchase_qty"`
	UnitPrice          float64             `orm:"column(unit_price)" json:"unit_price"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *FieldPurchaseOrderItem) MarshalJSON() ([]byte, error) {
	type Alias FieldPurchaseOrderItem

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *FieldPurchaseOrderItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *FieldPurchaseOrderItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
