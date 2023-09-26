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

// GetConsolidatedShipment find a single data consolidated_shipment set using field and value condition.
func GetConsolidatedShipment(field string, values ...interface{}) (*model.ConsolidatedShipment, error) {
	m := new(model.ConsolidatedShipment)
	o := orm.NewOrm()
	o.Using("read_only")
	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	o.LoadRelated(m, "PurchaseOrders", 2)
	for _, v := range m.PurchaseOrders {
		o.LoadRelated(v, "PurchaseOrderItems", 2)
	}

	if len(m.PurchaseOrders) > 0 {
		m.WarehouseName = m.PurchaseOrders[0].Warehouse.Name
	}

	o.LoadRelated(m, "ConsolidatedShipmentSignatures", 2)
	if len(m.ConsolidatedShipmentSignatures) > 0 {
		//// Initialize minio client object.
		minioClient, _ := minio.New(util.S3endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(util.S3accessKeyID, util.S3secretAccessKey, ""),
			Secure: true,
		})
		reqParams := make(url.Values)
		//// Retrieve URL valid for 60 second
		for _, item := range m.ConsolidatedShipmentSignatures {
			tempSignatureImage := strings.Split(item.SignatureURL, "/")
			preSignedURLSignatureImage, _ := minioClient.PresignedGetObject(context.Background(), util.S3bucketNameImage, tempSignatureImage[4]+"/"+tempSignatureImage[5], time.Second*60, reqParams)
			item.SignatureURL = preSignedURLSignatureImage.String()
		}
	}

	o.Raw("SELECT p.name AS product_name, u.name AS uom_name, SUM(poi.purchase_qty) AS total_qty "+
		"FROM purchase_order po "+
		"JOIN purchase_order_item poi ON poi.purchase_order_id = po.id "+
		"JOIN product p ON poi.product_id = p.id "+
		"JOIN uom u ON p.uom_id = u.id "+
		"WHERE po.consolidated_shipment_id = ? "+
		"GROUP BY p.id", m.ID).QueryRows(&m.SkuSummaries)

	for _, sku := range m.SkuSummaries {
		o.Raw("SELECT po.code AS purchase_order_code, poi.purchase_qty AS qty "+
			"FROM purchase_order po "+
			"JOIN purchase_order_item poi ON poi.purchase_order_id = po.id "+
			"JOIN product p ON poi.product_id = p.id "+
			"WHERE po.consolidated_shipment_id = ? AND p.name = ? "+
			"GROUP BY poi.id", m.ID, sku.ProductName).QueryRows(&sku.PurchaseOrders)
	}

	return m, nil
}

// GetConsolidatedShipments : function to get data from database based on parameters
func GetConsolidatedShipments(rq *orm.RequestQuery, warehouseID int64) (mx []*model.ConsolidatedShipment, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")
	warehouse := new(model.Warehouse)
	staffWarehouse := new(model.Warehouse)
	q, _ := rq.QueryReadOnly(new(model.ConsolidatedShipment))

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

	for _, v := range mx {
		o.Raw("SELECT w.name FROM warehouse w JOIN purchase_order po ON w.id = po.warehouse_id WHERE po.consolidated_shipment_id = ? AND po.status IN (?,?) LIMIT 1", v.ID, 1, 2).QueryRow(&v.WarehouseName)
	}

	return mx, total, nil
}

// ValidConsolidatedShipment : function to check if id is valid in database
func ValidConsolidatedShipment(id int64) (consolidatedShipment *model.ConsolidatedShipment, e error) {
	consolidatedShipment = &model.ConsolidatedShipment{ID: id}
	e = consolidatedShipment.Read("ID")

	return
}
