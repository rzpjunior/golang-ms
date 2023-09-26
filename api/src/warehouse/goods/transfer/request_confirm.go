// Copyright 2020 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package transfer

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

type confirmRequest struct {
	ID         int64     `json:"-"`
	AtaDataStr string    `json:"ata_date" valid:"required"`
	AtaTime    string    `json:"ata_time" valid:"required"`
	AtaDate    time.Time `json:"-"`

	GoodsTransfer      *model.GoodsTransfer        `json:"-"`
	GoodsTransferItems []*goodsTransferItemRequest `json:"items" valid:"required"`
	Stocks             []*model.Stock              `json:"-"`
	StockType          *model.Glossary             `json:"-"`

	Session *auth.SessionData `json:"-"`
}

func (r *confirmRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var filter, exclude map[string]interface{}

	if r.GoodsTransfer, err = repository.ValidGoodsTransfer(r.ID); err == nil {
		if r.GoodsTransfer.Status != 1 {
			o.Failure("status.inactive", util.ErrorDocStatus("goods receipt", "active"))
		} else {
			r.GoodsTransfer.Origin.Read("ID")
		}
	} else {
		o.Failure("id.invalid", util.ErrorInvalidData("id"))
	}

	layout := "2006-01-02"
	if r.AtaDate, err = time.Parse(layout, r.AtaDataStr); err != nil {
		o.Failure("ata_date.invalid", util.ErrorInvalidData("actual arrival date"))
	} else {
		if r.AtaDate.Before(r.GoodsTransfer.RecognitionDate) {
			o.Failure("ata_date.invalid", util.ErrorEqualLater("actual arrival date", "departure date"))
		}
	}

	if _, err = time.Parse("15:04", r.AtaTime); err != nil {
		o.Failure("ata_time.invalid", util.ErrorInvalidData("actual arrival time"))
	}

	r.StockType, err = repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_int", r.GoodsTransfer.StockType)
	if err != nil {
		o.Failure("stock_type_id.invalid", util.ErrorInvalidData("stock type"))
		return o
	}

	filter = map[string]interface{}{"status": 1, "warehouse_id": r.GoodsTransfer.Origin.ID, "stock_type": r.StockType.ValueInt}
	if _, countStockOpname, e := repository.CheckStockOpnameData(filter, exclude); e == nil && countStockOpname > 0 {
		o.Failure("id.invalid", util.ErrorRelated("active ", "stock opname", r.GoodsTransfer.Origin.Code+"-"+r.GoodsTransfer.Origin.Name))
	}

	for i, v := range r.GoodsTransferItems {
		if goodsTransferItemID, err := common.Decrypt(v.GoodsTransferItemID); err == nil {
			v.GoodsTransferItem = &model.GoodsTransferItem{ID: goodsTransferItemID}
			if err = v.GoodsTransferItem.Read("ID"); err != nil {
				o.Failure("goods_transfer_item_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("goods transfer item"))
			}
		}

		if productID, err := common.Decrypt(v.ProductID); err == nil {
			if v.Product, err = repository.ValidProduct(productID); err == nil {
				stock := &model.Stock{
					Warehouse: r.GoodsTransfer.Destination,
					Product:   v.Product,
					Status:    1,
				}

				if err = stock.Read("Warehouse", "Product", "Status"); err != nil {
					o.Failure("stock"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("stock"))
				}

				r.Stocks = append(r.Stocks, stock)
			} else {
				o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
			}
		} else {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
		}
	}

	return o
}

func (r *confirmRequest) Messages() map[string]string {
	return map[string]string{}
}
