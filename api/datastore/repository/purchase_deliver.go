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

// GetPurchaseDeliver find a single data purchase_deliver set using field and value condition.
func GetPurchaseDeliver(field string, values ...interface{}) (*model.PurchaseDeliver, error) {
	m := new(model.PurchaseDeliver)
	o := orm.NewOrm()
	o.Using("read_only")
	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	o.LoadRelated(m, "PurchaseDeliverSignature", 2)
	if len(m.PurchaseDeliverSignature) > 0 {
		//// Initialize minio client object.
		minioClient, _ := minio.New(util.S3endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(util.S3accessKeyID, util.S3secretAccessKey, ""),
			Secure: true,
		})
		reqParams := make(url.Values)
		//// Retrieve URL valid for 60 second
		for _, item := range m.PurchaseDeliverSignature {
			tempSignatureImage := strings.Split(item.Signature, "/")
			preSignedURLSignatureImage, _ := minioClient.PresignedGetObject(context.Background(), util.S3bucketNameImage, tempSignatureImage[4]+"/"+tempSignatureImage[5], time.Second*60, reqParams)
			item.Signature = preSignedURLSignatureImage.String()
		}
	}
	o.LoadRelated(m.FieldPurchaseOrder, "FieldPurchaseOrderItems", 2)

	return m, nil
}

// GetPurchaseDelivers : function to get data from database based on parameters
func GetPurchaseDelivers(rq *orm.RequestQuery, staffID int64) (mx []*model.PurchaseDeliver, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")
	warehouse := new(model.Warehouse)
	staff := new(model.Staff)
	q, _ := rq.QueryReadOnly(new(model.PurchaseDeliver))

	o.QueryTable(warehouse).Filter("name", "All Warehouse").Limit(1).One(warehouse)
	o.QueryTable(staff).Filter("id", staffID).Limit(1).One(staff)

	if staffID == 0 {
		if total, err = q.Count(); err != nil || total == 0 {
			return nil, total, err
		}

		if _, err = q.RelatedSel().All(&mx, rq.Fields...); err != nil {
			return nil, total, err
		}
	} else {
		if staff.Warehouse.ID == warehouse.ID {
			if total, err = q.Count(); err != nil || total == 0 {
				return nil, total, err
			}

			if _, err = q.RelatedSel().All(&mx, rq.Fields...); err != nil {
				return nil, total, err
			}
		} else {
			if total, err = q.Filter("PurchaseOrder__Warehouse__ID", staff.Warehouse.ID).Count(); err != nil || total == 0 {
				return nil, total, err
			}

			if _, err = q.RelatedSel().Filter("PurchaseOrder__Warehouse__ID", staff.Warehouse.ID).All(&mx, rq.Fields...); err != nil {
				return nil, total, err
			}
		}
	}

	return mx, total, nil
}

// GetFilterPurchaseDelivers : function to get data from database based on parameters
func GetFilterPurchaseDelivers(rq *orm.RequestQuery, codes []string) (mx []*model.PurchaseDeliver, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PurchaseDeliver))
	o := orm.NewOrm()
	o.Using("read_only")

	if total, err = q.Filter("code__in", codes).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	if _, err = q.RelatedSel().Filter("code__in", codes).All(&mx, rq.Fields...); err != nil {
		return nil, total, err
	}

	for _, v := range mx {
		o.LoadRelated(v.FieldPurchaseOrder, "FieldPurchaseOrderItems", 2)
	}
	return mx, total, nil
}

// ValidPurchaseDeliver : function to check if id is valid in database
func ValidPurchaseDeliver(id int64) (purchaseDeliver *model.PurchaseDeliver, e error) {
	purchaseDeliver = &model.PurchaseDeliver{ID: id}
	e = purchaseDeliver.Read("ID")

	return
}
