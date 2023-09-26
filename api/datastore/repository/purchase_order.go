// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"
	"net/url"
	"strings"
	"time"

	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"git.edenfarm.id/cuxs/orm"
)

// GetPurchaseOrder find a single data price set using field and value condition.
func GetPurchaseOrder(field string, values ...interface{}) (*model.PurchaseOrder, error) {
	m := new(model.PurchaseOrder)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel(2).Limit(1).One(m); err != nil {
		return nil, err
	}
	o.LoadRelated(m, "PurchaseOrderItems", 2)
	o.LoadRelated(m, "GoodsReceipt", 1)
	o.LoadRelated(m, "PurchaseInvoice", 1)
	o.LoadRelated(m, "PurchaseOrderImage", 2)
	if len(m.PurchaseOrderImage) > 0 {
		//// Initialize minio client object.
		minioClient, _ := minio.New(util.S3endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(util.S3accessKeyID, util.S3secretAccessKey, ""),
			Secure: true,
		})
		reqParams := make(url.Values)
		//// Retrieve URL valid for 600 second
		for _, item := range m.PurchaseOrderImage {
			tempSignatureImage := strings.Split(item.ImageURL, "/")
			preSignedURLSignatureImage, _ := minioClient.PresignedGetObject(context.Background(), util.S3bucketNameImage, tempSignatureImage[4]+"/"+tempSignatureImage[5], time.Second*600, reqParams)
			item.ImageURL = preSignedURLSignatureImage.String()
		}
	}

	o.Raw("select note from audit_log where `type` = 'purchase_order' and `function` = 'cancel' and ref_id = ?", m.ID).QueryRow(&m.CancellationNote)

	tax := float64(m.TaxPct)
	m.Tax = m.TotalPrice * tax

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

// GetPurchaseOrders : function to get data from database based on parameters
func GetPurchaseOrders(rq *orm.RequestQuery, isInbound string) (m []*model.PurchaseOrder, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PurchaseOrder))
	o := orm.NewOrm()
	o.Using("read_only")

	if total, err = q.Exclude("status", 4).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PurchaseOrder
	if _, err = q.Exclude("status", 4).All(&mx, rq.Fields...); err != nil {
		return nil, total, err
	}

	var arrPo []*model.PurchaseOrder
	var isExist bool
	for _, v := range mx {
		o.LoadRelated(v, "GoodsReceipt", 1)
		totalSku, _ := o.QueryTable(new(model.PurchaseOrderItem)).Filter("purchase_order_id", v.ID).Count()
		v.TotalSku = int8(totalSku)
		if isInbound == "1" {
			if isExist = o.QueryTable(new(model.GoodsReceipt)).Filter("purchase_order_id", v.ID).Filter("status", 2).Exist(); !isExist {
				arrPo = append(arrPo, v)
			}
		}
	}
	if isInbound == "1" {
		return arrPo, total, nil
	}

	return mx, total, nil

}

// GetFilterPurchaseOrders : function to get data from database based on parameters with filtered permission
func GetFilterPurchaseOrders(rq *orm.RequestQuery) (m []*model.PurchaseOrder, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PurchaseOrder))
	o1 := orm.NewOrm()
	o1.Using("read_only")

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PurchaseOrder
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err != nil {
		return nil, total, err
	}

	var arrPo []*model.PurchaseOrder
	var isExist bool
	for _, v := range mx {
		if isExist = o1.QueryTable(new(model.GoodsReceipt)).Filter("purchase_order_id", v.ID).Filter("status__in", 1, 2).Exist(); !isExist {
			arrPo = append(arrPo, v)
		}
	}

	return arrPo, total, nil
}

// ValidPurchaseOrder : function to check if id is valid in database
func ValidPurchaseOrder(id int64) (purchaseOrder *model.PurchaseOrder, e error) {
	purchaseOrder = &model.PurchaseOrder{ID: id}
	e = purchaseOrder.Read("ID")

	return
}

// CheckPurchaseOrderProductStatus : function to check if product is valid in table
func CheckPurchaseOrderProductStatus(productID int64, status int8, warehouseArr ...string) (*model.PurchaseOrderItem, int64, error) {
	var err error
	o := orm.NewOrm()
	o.Using("read_only")

	poi := new(model.PurchaseOrderItem)
	q := o.QueryTable(poi).RelatedSel("PurchaseOrder").Filter("PurchaseOrder__Status", status).Filter("product_id", productID)

	if len(warehouseArr) > 0 {
		q = q.Exclude("PurchaseOrder__Warehouse__id__in", warehouseArr)
	}

	if total, err := q.All(poi); err == nil {
		return poi, total, nil
	}

	return nil, 0, err
}

// GetPurchaseOrdersFieldPurchaser : function to get list data of purchase_order & load purchase_order_items
func GetPurchaseOrdersFieldPurchaser(rq *orm.RequestQuery, warehouseID int64) (mx []*model.PurchaseOrder, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")
	warehouse := new(model.Warehouse)
	q, _ := rq.QueryReadOnly(new(model.PurchaseOrder))

	o.QueryTable(warehouse).Filter("name", "All Warehouse").Limit(1).One(warehouse)

	if warehouseID == warehouse.ID || warehouseID == 0 {
		if total, err = q.Exclude("status", 4).Count(); err != nil || total == 0 {
			return nil, total, err
		}

		if _, err = q.RelatedSel(2).Exclude("status", 4).All(&mx, rq.Fields...); err != nil {
			return nil, total, err
		}
	} else {
		if total, err = q.Filter("warehouse_id", warehouseID).Exclude("status", 4).Count(); err != nil || total == 0 {
			return nil, total, err
		}

		if _, err = q.RelatedSel(2).Filter("warehouse_id", warehouseID).Exclude("status", 4).All(&mx, rq.Fields...); err != nil {
			return nil, total, err
		}
	}

	for _, item := range mx {

		err = o.Raw("SELECT count(id) FROM purchase_order_item WHERE purchase_order_id = ?", item.ID).QueryRow(&item.TotalSku)

		if err != nil {
			return nil, total, err
		}

		_, err = o.Raw("SELECT u.name as uom_name, sum(poi.order_qty) as total_weight FROM purchase_order_item poi JOIN product p ON p.id =  poi.product_id JOIN uom u ON u.id = p.uom_id WHERE purchase_order_id = ? GROUP BY p.uom_id", item.ID).QueryRows(&item.TonasePurchaseOrder)

		if err != nil {
			return nil, total, err
		}
	}

	return mx, total, nil
}

// GetPurchaseOrderFieldPurchaser find a single data price set using field and value condition for field purchaser mobile app.
func GetPurchaseOrderFieldPurchaser(field string, values ...interface{}) (*model.PurchaseOrder, error) {
	m := new(model.PurchaseOrder)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel(2).Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "PurchaseOrderItems", 2)
	for _, items := range m.PurchaseOrderItems {
		o.LoadRelated(items, "FieldPurchaseOrderItems", 2)
		for _, v := range items.FieldPurchaseOrderItems {
			o.LoadRelated(v, "FieldPurchaseOrder", 2)
		}
	}

	tax := float64(m.TaxPct)
	m.Tax = m.TotalPrice * tax

	o.LoadRelated(m, "PurchaseOrderSignature", 2)
	o.LoadRelated(m, "PurchaseOrderImage", 2)
	if len(m.PurchaseOrderSignature) > 0 || len(m.PurchaseOrderImage) > 0 {
		//// Initialize minio client object.
		minioClient, _ := minio.New(util.S3endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(util.S3accessKeyID, util.S3secretAccessKey, ""),
			Secure: true,
		})
		reqParams := make(url.Values)
		//// Retrieve URL valid for 60 second
		if len(m.PurchaseOrderSignature) > 0 {
			for _, item := range m.PurchaseOrderSignature {
				tempSignatureImage := strings.Split(item.SignatureURL, "/")
				preSignedURLSignatureImage, _ := minioClient.PresignedGetObject(context.Background(), util.S3bucketNameImage, tempSignatureImage[4]+"/"+tempSignatureImage[5], time.Second*60, reqParams)
				item.SignatureURL = preSignedURLSignatureImage.String()
			}
		}

		//// Retrieve URL valid for 600 second
		if len(m.PurchaseOrderImage) > 0 {
			for _, item := range m.PurchaseOrderImage {
				tempSignatureImage := strings.Split(item.ImageURL, "/")
				preSignedURLSignatureImage, _ := minioClient.PresignedGetObject(context.Background(), util.S3bucketNameImage, tempSignatureImage[4]+"/"+tempSignatureImage[5], time.Second*600, reqParams)
				item.ImageURL = preSignedURLSignatureImage.String()
			}
		}
	}

	return m, nil
}

// CheckPurchaseOrderData : function to get all purchase order data based on filter and exclude parameters
func CheckPurchaseOrderData(filter, exclude map[string]interface{}) (m []*model.PurchaseOrder, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.PurchaseOrder))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err = o.All(&m); err != nil {
		return nil, 0, err
	}

	return m, total, nil
}

// GetFilterPurchaseOrdersForConsolidatedShipment : function to get data from database based on parameters with filtered permission
func GetFilterPurchaseOrdersForConsolidatedShipment(rq *orm.RequestQuery) (mx []*model.PurchaseOrder, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PurchaseOrder))
	o := orm.NewOrm()
	o.Using("read_only")

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err != nil {
		return nil, total, err
	}

	for _, item := range mx {

		o.LoadRelated(item, "PurchaseOrderItems", 2)

		err = o.Raw("SELECT count(id) FROM purchase_order_item WHERE purchase_order_id = ?", item.ID).QueryRow(&item.TotalSku)

		if err != nil {
			return nil, total, err
		}

		_, err = o.Raw("SELECT u.name as uom_name, sum(poi.order_qty) as total_weight FROM purchase_order_item poi JOIN product p ON p.id =  poi.product_id JOIN uom u ON u.id = p.uom_id WHERE purchase_order_id = ? GROUP BY p.uom_id", item.ID).QueryRows(&item.TonasePurchaseOrder)

		if err != nil {
			return nil, total, err
		}

	}

	return mx, total, err
}
