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

// GetConsolidatedPurchaseDeliver find a single data purchase_deliver set using field and value condition.
func GetConsolidatedPurchaseDeliver(field string, values ...interface{}) (*model.ConsolidatedPurchaseDeliver, error) {
	m := new(model.ConsolidatedPurchaseDeliver)
	o := orm.NewOrm()
	o.Using("read_only")
	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	o.LoadRelated(m, "PurchaseDelivers", 2)

	o.Raw("SELECT p.name as product_name, u.name as uom_name "+
		"FROM purchase_deliver pd "+
		"JOIN field_purchase_order fpo ON pd.field_purchase_order_id = fpo.id "+
		"JOIN field_purchase_order_item fpoi ON fpoi.field_purchase_order_id = fpo.id "+
		"LEFT JOIN purchase_order po ON pd.purchase_order_id = po.id "+
		"JOIN product p ON fpoi.product_id = p.id "+
		"JOIN uom u ON p.uom_id = u.id "+
		"WHERE pd.consolidated_purchase_deliver_id = ? "+
		"GROUP BY p.id", m.ID).QueryRows(&m.Products)

	for _, product := range m.Products {
		o.Raw("SELECT po.code AS purchase_order_code FROM purchase_order po "+
			"JOIN field_purchase_order fpo ON po.id = fpo.purchase_order_id "+
			"JOIN field_purchase_order_item fpoi ON fpoi.field_purchase_order_id = fpo.id "+
			"JOIN product p ON fpoi.product_id = p.id "+
			"JOIN purchase_deliver pd ON pd.purchase_order_id = po.id "+
			"WHERE pd.consolidated_purchase_deliver_id = ? AND p.name = ? "+
			"GROUP BY po.id", m.ID, product.ProductName).QueryRows(&product.PurchaseOrders)

		for _, purchaseOrder := range product.PurchaseOrders {
			o.Raw("SELECT fpoi.purchase_qty AS purchase_qty, pd.code AS purchase_deliver_code "+
				"FROM field_purchase_order_item fpoi "+
				"JOIN field_purchase_order fpo ON fpoi.field_purchase_order_id = fpo.id "+
				"JOIN purchase_order po ON fpo.purchase_order_id = po.id "+
				"JOIN purchase_deliver pd ON pd.field_purchase_order_id = fpo.id "+
				"JOIN product p ON fpoi.product_id = p.id "+
				"WHERE pd.consolidated_purchase_deliver_id = ? AND p.name = ? AND po.code = ? "+
				"", m.ID, product.ProductName, purchaseOrder.PurchaseOrderCode).QueryRows(&purchaseOrder.Items)
		}
	}

	m.TotalProduct = len(m.Products)
	m.TotalPurchaseDeliver = len(m.PurchaseDelivers)
	if len(m.PurchaseDelivers) > 0 {
		m.SupplierName = m.PurchaseDelivers[0].PurchaseOrder.Supplier.Name
	}
	o.LoadRelated(m, "ConsolidatedPurchaseDeliverSignature", 2)
	if len(m.ConsolidatedPurchaseDeliverSignature) > 0 {
		//// Initialize minio client object.
		minioClient, _ := minio.New(util.S3endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(util.S3accessKeyID, util.S3secretAccessKey, ""),
			Secure: true,
		})
		reqParams := make(url.Values)
		//// Retrieve URL valid for 60 second
		for _, item := range m.ConsolidatedPurchaseDeliverSignature {
			tempSignatureImage := strings.Split(item.Signature, "/")
			preSignedURLSignatureImage, _ := minioClient.PresignedGetObject(context.Background(), util.S3bucketNameImage, tempSignatureImage[4]+"/"+tempSignatureImage[5], time.Second*60, reqParams)
			item.Signature = preSignedURLSignatureImage.String()
		}
	}

	return m, nil
}

// GetConsolidatedPurchaseDelivers : function to get data from database based on parameters
func GetConsolidatedPurchaseDelivers(rq *orm.RequestQuery, warehouseID int64) (mx []*model.ConsolidatedPurchaseDeliver, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")
	warehouse := new(model.Warehouse)
	staffWarehouse := new(model.Warehouse)
	q, _ := rq.QueryReadOnly(new(model.ConsolidatedPurchaseDeliver))

	o.QueryTable(warehouse).Filter("name", "All Warehouse").Limit(1).One(warehouse)
	o.QueryTable(staffWarehouse).RelatedSel().Filter("id", warehouseID).Limit(1).One(staffWarehouse)

	if warehouseID == 0 || warehouseID == warehouse.ID {
		if total, err = q.Count(); err != nil || total == 0 {
			return nil, total, err
		}

		if _, err = q.RelatedSel().All(&mx, rq.Fields...); err != nil {
			return nil, total, err
		}
	} else {
		if total, err = q.Filter("Code__contains", staffWarehouse.Code).Count(); err != nil || total == 0 {
			return nil, total, err
		}

		if _, err = q.RelatedSel().Filter("Code__contains", staffWarehouse.Code).All(&mx, rq.Fields...); err != nil {
			return nil, total, err
		}
	}

	return mx, total, nil
}

// ValidConsolidatedPurchaseDeliver : function to check if id is valid in database
func ValidConsolidatedPurchaseDeliver(id int64) (consolidatedPurchaseDeliver *model.ConsolidatedPurchaseDeliver, e error) {
	consolidatedPurchaseDeliver = &model.ConsolidatedPurchaseDeliver{ID: id}
	e = consolidatedPurchaseDeliver.Read("ID")

	return
}
