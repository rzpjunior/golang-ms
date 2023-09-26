// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receipt

import (
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createRequest : struct to hold goods receipt request data
type createRequest struct {
	Code        string    `json:"-"`
	AreaID      string    `json:"area_id" valid:"required"`
	WarehouseID string    `json:"warehouse_id" valid:"required"`
	InboundID   string    `json:"inbound_id" valid:"required"`
	AtaDateStr  string    `json:"ata_date"`
	AtaTimeStr  string    `json:"ata_time"`
	Note        string    `json:"note"`
	InboundType string    `json:"inbound_type" valid:"required"`
	AtaDate     time.Time `json:"-"`
	AtaTime     time.Time `json:"-"`
	TotalWeight float64   `json:"-"`

	Area              *model.Area                `json:"-"`
	Warehouse         *model.Warehouse           `json:"-"`
	PurchaseOrder     *model.PurchaseOrder       `json:"-"`
	GoodsTransfer     *model.GoodsTransfer       `json:"-"`
	GoodsReceiptItems []*goodsReceiptItemRequest `json:"items" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

type goodsReceiptItemRequest struct {
	GoodsReceiptItemID string  `json:"goods_receipt_item_id"`
	ProductID          string  `json:"product_id"`
	InboundItemID      string  `json:"inbound_item_id"`
	DeliveryQty        float64 `json:"delivery_qty"`
	RejectQty          float64 `json:"reject_qty"`
	RejectReason       int8    `json:"reject_reason"`
	ReceiveQty         float64 `json:"-"`
	Note               string  `json:"note"`

	Product           *model.Product           `json:"-"`
	PurchaseOrderItem *model.PurchaseOrderItem `json:"-"`
	GoodsReceiptItem  *model.GoodsReceiptItem  `json:"-"`
	GoodsTransferItem *model.GoodsTransferItem `json:"-"`
}

// Validate : function to validate goods receipt request data
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var isProductExist = make(map[string]bool)
	var filter, exclude map[string]interface{}
	var weight float64
	var inboundID int64
	var countGR int64
	warehouseID, _ := common.Decrypt(r.WarehouseID)
	r.Warehouse = &model.Warehouse{ID: warehouseID}
	r.Warehouse.Read("ID")

	if inboundID, err = common.Decrypt(r.InboundID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("inbound"))
		return o
	}
	switch r.InboundType {
	case "purchase_order":
		r.PurchaseOrder = &model.PurchaseOrder{ID: inboundID}
		if err = r.PurchaseOrder.Read("ID"); err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("purchase order"))

		}
		if r.PurchaseOrder.Status != 1 {
			o.Failure("id.invalid", util.ErrorDocStatus("purchase order", "active"))
		}

		if err = r.PurchaseOrder.Supplier.Read("ID"); err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("supplier"))
		}

		r.AtaDate = time.Now()
		r.AtaTimeStr = time.Now().Format("15:04")

		filter = map[string]interface{}{"purchase_order_id": inboundID, "status__in": []int{1, 2}}
		exclude = map[string]interface{}{}
		if _, countGR, err = repository.CheckGoodsReceiptData(filter, exclude); err == nil {
			if countGR > 0 {
				o.Failure("id.invalid", util.ErrorCreateDoc("goods receipt", "purchase order"))
			}
		}

	case "goods_transfer":
		r.GoodsTransfer = &model.GoodsTransfer{ID: inboundID}

		if err = r.GoodsTransfer.Read("ID"); err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("goods transfer"))

		}
		if r.GoodsTransfer.Status != 1 {
			o.Failure("id.invalid", util.ErrorDocStatus("goods transfer", "active"))
		}

		if err = r.GoodsTransfer.Origin.Read("ID"); err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("warehouse origin"))
		}

		r.AtaDate = time.Now()
		r.AtaTimeStr = time.Now().Format("15:04")

		filter = map[string]interface{}{"goods_transfer_id": inboundID, "status__in": []int{1, 2}}
		exclude = map[string]interface{}{}
		if _, countGR, err = repository.CheckGoodsReceiptData(filter, exclude); err == nil {
			if countGR > 0 {
				o.Failure("id.invalid", util.ErrorCreateDoc("goods receipt", "goods transfer"))
			}
		}
	default:
		layout := "2006-01-02"
		if r.AtaDate, err = time.Parse(layout, r.AtaDateStr); err != nil {
			o.Failure("ata_date.invalid", util.ErrorInvalidData("ata date"))
		}

		if r.AtaTime, err = time.Parse("15:04", r.AtaTimeStr); err != nil {
			o.Failure("ata_time.invalid", util.ErrorInvalidData("ata time"))
		}

	}

	for i, v := range r.GoodsReceiptItems {
		if _, exist := isProductExist[v.ProductID]; exist {
			o.Failure("product_id"+strconv.Itoa(i)+".duplicate", util.ErrorDuplicate("product"))
		} else {
			var productID int64
			if productID, err = common.Decrypt(v.ProductID); err != nil {
				o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
				return o
			}
			if v.Product, err = repository.ValidProduct(productID); err != nil {
				o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
				return o
			}

			if err = v.Product.Uom.Read("ID"); err != nil {
				o.Failure("uom_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("uom"))
				return o
			}

			isProductExist[v.ProductID] = true

			if v.Product.Uom.DecimalEnabled == 2 {
				if v.DeliveryQty != float64((int64(v.DeliveryQty))) {
					o.Failure("deliver_qty"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("deliver quantity"))
				}

				if v.RejectQty != float64((int64(v.RejectQty))) {
					o.Failure("reject_qty"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("reject quantity"))
				}
			}

			if v.DeliveryQty < 0 {
				o.Failure("delivery_qty"+strconv.Itoa(i)+".gte", util.ErrorEqualGreater("delivery quantity", "0"))
			}

			if len(v.Note) > 100 {
				o.Failure("note"+strconv.Itoa(i), util.ErrorCharLength("note", 100))
			}

			switch r.InboundType {
			case "purchase_order":
				purchaseOrderItemID, _ := common.Decrypt(v.InboundItemID)
				if v.PurchaseOrderItem, err = repository.ValidPurchaseOrderItem(purchaseOrderItemID); err != nil {
					o.Failure("purchase_order_item"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("purchase order item"))
				}

				if v.RejectQty < 0 {
					o.Failure("reject_qty"+strconv.Itoa(i)+".gte", util.ErrorEqualGreater("reject quantity", "0"))
				}

				v.ReceiveQty = v.DeliveryQty - v.RejectQty
				weight = v.ReceiveQty * v.Product.UnitWeight

				r.TotalWeight = r.TotalWeight + weight

				if v.ReceiveQty < 0 {
					o.Failure("receive_qty"+strconv.Itoa(i)+".gte", util.ErrorEqualGreater("receive quantity", "0"))
				}
				if v.ReceiveQty > v.PurchaseOrderItem.OrderQty {
					o.Failure("receive_qty"+strconv.Itoa(i)+".lt", util.ErrorGreater("order quantity", "receive quantity"))
				}

			case "goods_transfer":
				goodsTransferItemID, _ := common.Decrypt(v.InboundItemID)
				if v.GoodsTransferItem, err = repository.ValidGoodsTransferItem(goodsTransferItemID); err != nil {
					o.Failure("goods_transfer_item"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("goods transfer item"))
				}
				v.ReceiveQty = v.DeliveryQty

				if v.DeliveryQty > v.GoodsTransferItem.DeliverQty {
					o.Failure("deliver_qty"+strconv.Itoa(i)+".lt", util.ErrorGreater("transfer quantity", "deliver quantity"))
				}

				if v.DeliveryQty != 0 && v.DeliveryQty < v.GoodsTransferItem.DeliverQty {
					if v.RejectReason == 0 {
						o.Failure("reject_reason"+strconv.Itoa(i)+".invalid", util.ErrorInputRequired("reject reason"))
					}
				}
				weight = v.DeliveryQty * v.Product.UnitWeight

				r.TotalWeight = r.TotalWeight + weight
			default:

				purchaseOrderItemID, _ := common.Decrypt(v.InboundItemID)
				if v.PurchaseOrderItem, err = repository.ValidPurchaseOrderItem(purchaseOrderItemID); err != nil {
					o.Failure("purchase_order_item"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("purchase order item"))
				}

				if v.RejectQty < 0 {
					o.Failure("reject_qty"+strconv.Itoa(i)+".gte", util.ErrorEqualGreater("reject quantity", "0"))
				}

				if v.ReceiveQty < 0 {
					o.Failure("receive_qty"+strconv.Itoa(i)+".gte", util.ErrorEqualGreater("receive quantity", "0"))
				}
				if v.ReceiveQty > v.PurchaseOrderItem.OrderQty {
					o.Failure("receive_qty"+strconv.Itoa(i)+".lt", util.ErrorGreater("order quantity", "receive quantity"))
				}

				v.ReceiveQty = v.DeliveryQty - v.RejectQty
				weight = v.ReceiveQty * v.Product.UnitWeight

				r.TotalWeight = r.TotalWeight + weight

			}

		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"area_id.required":           util.ErrorInputRequired("area"),
		"warehouse_id.required":      util.ErrorInputRequired("warehouse"),
		"supplier_id.required":       util.ErrorInputRequired("supplier"),
		"purchase_order_id.required": util.ErrorInputRequired("purchase order"),
	}

	return messages
}
