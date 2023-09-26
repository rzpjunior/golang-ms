// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetGoodsTransfer find a single data price set using field and value condition.
func GetGoodsTransfer(field string, values ...interface{}) (*model.GoodsTransfer, error) {
	m := new(model.GoodsTransfer)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "GoodsTransferItem", 2)
	o.LoadRelated(m, "GoodsReceipt", 1)

	var err error
	if m.LockedBy != 0 {
		if m.LockedByObj, err = ValidStaff(m.LockedBy); err != nil {
			return nil, err
		}
	}

	if m.UpdatedBy != 0 {
		if m.UpdatedByObj, err = ValidStaff(m.UpdatedBy); err != nil {
			return nil, err
		}
	}

	return m, nil
}

// GetGoodsTransfers : function to get data from database based on parameters
func GetGoodsTransfers(rq *orm.RequestQuery) (gt []*model.GoodsTransfer, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.GoodsTransfer))
	o := orm.NewOrm()
	o.Using("read_only")

	if total, err = q.Filter("status__in", 1, 2, 3, 5).All(&gt, rq.Fields...); err != nil {
		return nil, total, err
	}
	for _, v := range gt {
		o.LoadRelated(v, "GoodsReceipt", 1)
		v.TotalSku, _ = o.QueryTable(new(model.GoodsTransferItem)).Filter("goods_transfer_id", v.ID).Count()
	}

	return gt, total, nil
}

// GetFilterGoodsTransfers : function to get data from database based on parameters with filtered permission
func GetFilterGoodsTransfers(rq *orm.RequestQuery) (gt []*model.GoodsTransfer, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.GoodsTransfer))

	if total, err = q.Filter("status__in", 1, 2, 3, 5).All(&gt, rq.Fields...); err == nil {
		return gt, total, nil
	}

	return nil, total, err
}

// ValidGoodsTransfer : function to check if id is valid in database
func ValidGoodsTransfer(id int64) (gt *model.GoodsTransfer, e error) {
	gt = &model.GoodsTransfer{ID: id}
	e = gt.Read("ID")

	return
}

// CheckGoodsTransferProductStatus : function to check if product is valid in table
func CheckGoodsTransferProductStatus(productID int64, status int8, warehouseArr ...string) (*model.GoodsTransferItem, int64, error) {
	var err error
	o := orm.NewOrm()
	o.Using("read_only")

	gti := new(model.GoodsTransferItem)
	q := o.QueryTable(gti).RelatedSel("GoodsTransfer")

	cond := orm.NewCondition()
	cond = cond.And("GoodsTransfer__Status", status).And("product_id", productID)

	cond2 := orm.NewCondition()
	if len(warehouseArr) > 0 {
		cond2 = cond2.AndNot("GoodsTransfer__Origin__id__in", warehouseArr).AndNot("GoodsTransfer__Destination__id__in", warehouseArr)
		cond = cond.AndCond(cond2)
	}

	if total, err := q.SetCond(cond).All(gti); err == nil {
		return gti, total, nil
	}

	return nil, 0, err
}

// ValidGoodsTransferItem : function to check if id is valid in database
func ValidGoodsTransferItem(id int64) (gti *model.GoodsTransferItem, e error) {
	gti = &model.GoodsTransferItem{ID: id}
	e = gti.Read("ID")

	return
}
