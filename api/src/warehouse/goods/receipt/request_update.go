// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receipt

import (
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/orm"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// updateRequest : struct to hold goods receipt request data
type updateRequest struct {
	ID           int64     `json:"-"`
	AtaDateStr   string    `json:"ata_date" valid:"required"`
	AtaTimeStr   string    `json:"ata_time" valid:"required"`
	Note         string    `json:"note"`
	InboundType  string    `json:"inbound_type" valid:"required"`
	AtaDate      time.Time `json:"-"`
	AtaTime      time.Time `json:"-"`
	TotalWeight  float64   `json:"-"`
	IsAtaChanged bool      `json:"-"`
	NoteChanged  string    `json:"-"`

	GoodsReceipt      *model.GoodsReceipt        `json:"-"`
	Warehouse         *model.Warehouse           `json:"-"`
	PurchaseOrder     *model.PurchaseOrder       `json:"-"`
	GoodsReceiptItems []*goodsReceiptItemRequest `json:"items" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate goods receipt request data
func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var weight float64
	o1 := orm.NewOrm()
	o1.Using("read_only")

	r.GoodsReceipt = &model.GoodsReceipt{ID: r.ID}
	if err = r.GoodsReceipt.Read("ID"); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("goods receipt"))
	}
	if r.GoodsReceipt.Status != 1 {
		o.Failure("id.inactive", util.ErrorDocStatus("goods receipt", "active"))
	}

	layout := "2006-01-02"
	if r.AtaDate, err = time.Parse(layout, r.AtaDateStr); err != nil {
		o.Failure("ata_date.invalid", util.ErrorInvalidData("ata date"))
	}

	if r.AtaTime, err = time.Parse("15:04", r.AtaTimeStr); err != nil {
		o.Failure("ata_time.invalid", util.ErrorInvalidData("ata time"))
	}

	if r.GoodsReceipt.AtaDate.Format("2006-01-02") != r.AtaDateStr ||
		r.GoodsReceipt.AtaTime != r.AtaTimeStr {
		r.IsAtaChanged = true
		r.NoteChanged = "ATA changed from: " + r.GoodsReceipt.AtaDate.Format("2006-01-02") + " " + r.GoodsReceipt.AtaTime + " to: " + r.AtaDateStr + " " + r.AtaTimeStr
	}

	switch r.InboundType {
	case "purchase_order":
		if err = r.GoodsReceipt.PurchaseOrder.Read("ID"); err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("purchase order"))
		}
	case "goods_transfer":
		if err = r.GoodsReceipt.GoodsTransfer.Read("ID"); err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("goods transfer"))
		}
	}

	for i, v := range r.GoodsReceiptItems {
		goodReceiptItemID, _ := common.Decrypt(v.GoodsReceiptItemID)

		v.GoodsReceiptItem = &model.GoodsReceiptItem{ID: goodReceiptItemID}
		if err = v.GoodsReceiptItem.Read("ID"); err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("goods receipt item"))
		}

		if err = v.GoodsReceiptItem.Product.Read("ID"); err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("product"))
		}

		if err = v.GoodsReceiptItem.Product.Uom.Read("ID"); err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("uom"))
		}

		if v.DeliveryQty < 0 {
			o.Failure("delivery_qty"+strconv.Itoa(i)+".gte", util.ErrorEqualGreater("delivery quantity", "0"))
		}

		if v.RejectQty < 0 {
			o.Failure("reject_qty"+strconv.Itoa(i)+".gte", util.ErrorEqualGreater("reject quantity", "0"))
		}

		if v.GoodsReceiptItem.Product.Uom.DecimalEnabled == 2 {
			if v.DeliveryQty != float64((int64(v.DeliveryQty))) {
				o.Failure("deliver_qty"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("deliver quantity"))
			}
			if v.RejectQty != float64((int64(v.RejectQty))) {
				o.Failure("reject_qty"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("reject quantity"))
			}
		}

		if len(v.Note) > 100 {
			o.Failure("note"+strconv.Itoa(i), util.ErrorCharLength("note", 100))
		}

		switch r.InboundType {
		case "purchase_order":
			if err = v.GoodsReceiptItem.PurchaseOrderItem.Read("ID"); err != nil {
				o.Failure("id.invalid", util.ErrorInvalidData("goods receipt item"))
			}

			v.ReceiveQty = v.DeliveryQty - v.RejectQty
			weight = v.ReceiveQty * v.GoodsReceiptItem.Product.UnitWeight

			r.TotalWeight = r.TotalWeight + weight

			if v.ReceiveQty < 0 {
				o.Failure("receive_qty"+strconv.Itoa(i)+".gte", util.ErrorEqualGreater("receive quantity", "0"))
			}

			if v.ReceiveQty > v.GoodsReceiptItem.PurchaseOrderItem.OrderQty {
				o.Failure("receive_qty"+strconv.Itoa(i)+".lt", util.ErrorGreater("order quantity", "receive quantity"))
			}
		case "goods_transfer":
			v.ReceiveQty = v.DeliveryQty

			var productId int64
			if productId, err = common.Decrypt(v.ProductID); err != nil {
				o.Failure("id.invalid", util.ErrorInvalidData("product"))
			}
			v.Product = &model.Product{ID: productId}
			if err = v.Product.Read("ID"); err != nil {
				o.Failure("id.invalid", util.ErrorInvalidData("product"))
			}

			if err = o1.Raw("select * from goods_transfer_item gti where gti.goods_transfer_id =? and gti.product_id =?", r.GoodsReceipt.GoodsTransfer.ID, v.Product.ID).QueryRow(&v.GoodsTransferItem); err != nil {
				o.Failure("id.invalid", util.ErrorInvalidData("goods transfer item"))
			}

			if v.DeliveryQty > v.GoodsTransferItem.DeliverQty {
				o.Failure("deliver_qty"+strconv.Itoa(i)+".lt", util.ErrorGreater("deliver quantity", "deliver quantity"))
			}

			if v.DeliveryQty != 0 && v.DeliveryQty < v.GoodsTransferItem.DeliverQty {
				if v.RejectReason == 0 {
					o.Failure("reject_reason"+strconv.Itoa(i)+".invalid", util.ErrorInputRequired("reject reason"))
				}
			}
			weight = v.DeliveryQty * v.Product.UnitWeight

			r.TotalWeight = r.TotalWeight + weight
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *updateRequest) Messages() map[string]string {
	messages := map[string]string{
		"warehouse_id.required":      util.ErrorInputRequired("warehouse"),
		"purchase_order_id.required": util.ErrorInputRequired("purchase order"),
		"ata_date.required":          util.ErrorInputRequired("ata date"),
		"ata_time.required":          util.ErrorInputRequired("ata time"),
	}

	return messages
}
