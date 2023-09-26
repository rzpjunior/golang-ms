package transfer_sku

// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

func Save(r createRequest) (resp *createRequest, err error) {
	o := orm.NewOrm()
	o.Begin()

	r.Code, err = util.GenerateDocCode(r.Code, r.Warehouse.Code, "transfer_sku")

	if err != nil {
		o.Rollback()
		return nil, err
	}

	ts := &model.TransferSku{
		Code:             r.Code,
		GoodsReceipt:     r.GoodsReceipt,
		Warehouse:        r.Warehouse,
		PurchaseOrder:    r.PurchaseOrder,
		GoodsTransfer:    r.GoodsTransfer,
		Status:           1,
		CreatedAt:        time.Now(),
		CreatedBy:        r.Session.Staff,
		TotalTransferQty: r.TotalTransferQty,
		TotalWasteQty:    r.TotalWasteQty,
		RecognitionDate:  r.RecognitionDateAt,
	}

	if _, err = o.Insert(ts); err != nil {
		o.Rollback()
		return nil, err
	}

	var transferSkuItems []*model.TransferSkuItem

	for _, product := range r.Products {

		for _, transferTo := range product.TransferTo {
			var discrepancy float64

			if product.Product.ID == transferTo.Product.ID {
				discrepancy = product.Discrepancy
			} else {
				discrepancy = 0
			}
			tsi := &model.TransferSkuItem{
				TransferSku:     ts,
				Product:         product.Product,
				TransferProduct: transferTo.Product,
				TransferQty:     transferTo.TransferQty,
				Discrepancy:     discrepancy,
				WasteQty:        transferTo.WasteQty,
				WasteReason:     transferTo.WasteReason,
			}

			transferSkuItems = append(transferSkuItems, tsi)

		}

	}

	if _, err = o.InsertMulti(100, &transferSkuItems); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, ts.ID, "transfer_sku", "create", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()
	return &r, nil
}
