// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// GetSubmissionCA : function to get data from database based on parameters
func GetSubmissionCA(rq *orm.RequestQuery) (m []*model.CustomerAcquisition, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	q, _ := rq.QueryReadOnly(new(model.CustomerAcquisition))
	q = q.Exclude("status", 4)

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.CustomerAcquisition
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			o.Raw("SELECT g.value_name FROM glossary g WHERE g.table = 'sales_assignment_item' AND g.attribute = 'task' AND g.value_int = ?", v.Task).QueryRow(&v.TaskStr)

			if v.Salesperson != nil && v.Salesperson.SalesGroupID != 0 {
				o.Raw("SELECT * FROM sales_group sg WHERE sg.id = ?", v.Salesperson.SalesGroupID).QueryRow(&v.Salesgroup)
			}
		}
		return mx, total, nil
	}

	return nil, total, err
}

// GetSubmissionCustomerAcquisitionDetail : function to get submitted customer acquisition by id
func GetSubmissionCustomerAcquisitionDetail(field string, values ...interface{}) (*model.CustomerAcquisition, error) {
	m := new(model.CustomerAcquisition)
	o := orm.NewOrm()
	o.Using("read_only")
	if err := o.QueryTable(m).Filter(field, values...).RelatedSel("Salesperson").Limit(1).One(m); err != nil {
		return nil, err
	}

	//// Initialize minio client object.
	minioClient, _ := minio.New(util.S3endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(util.S3accessKeyID, util.S3secretAccessKey, ""),
		Secure: true,
	})
	reqParams := make(url.Values)
	//// Retrieve URL valid for 60 second

	m.TaskPhotoArr = strings.Split(m.TaskPhoto, ",")
	for _, tp := range m.TaskPhotoArr {
		tempImage := strings.Split(tp, "/")
		preSignedURLImage, _ := minioClient.PresignedGetObject(context.Background(), util.S3bucketNameImage, tempImage[4]+"/"+tempImage[5], time.Second*60, reqParams)
		m.TaskPhotoList = append(m.TaskPhotoList, preSignedURLImage.String())
	}

	o.Raw("SELECT g.value_name FROM glossary g WHERE g.table = 'sales_assignment_item' AND g.attribute = 'task' AND g.value_int = ?", m.Task).QueryRow(&m.TaskStr)

	return m, nil
}

// GetCustomerAcquisition find a single data Customer Acquisition using field and value condition.
func GetCustomerAcquisition(field string, values ...interface{}) (*model.CustomerAcquisition, error) {
	m := new(model.CustomerAcquisition)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	// get top product on cai
	if m != nil {
		_, err := o.Raw("SELECT * from customer_acquisition_item cai "+
			"where cai.customer_acquisition_id = ?", m.ID).QueryRows(&m.CustomerAcquisitionItem)
		if err != nil && !errors.Is(err, orm.ErrNoRows) {
			return nil, err
		}
		if m.CustomerAcquisitionItem != nil {
			for _, cai := range m.CustomerAcquisitionItem {
				err = o.Raw("SELECT id, code, name from product p "+
					"where p.id = ?", cai.Product.ID).QueryRow(&cai.Product)
				if err != nil && !errors.Is(err, orm.ErrNoRows) {
					return nil, err
				}
			}
		}
	}
	return m, nil
}

// GetCustomerAcquisitions get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetCustomerAcquisitions(rq *orm.RequestQuery) (m []*model.CustomerAcquisition, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.CustomerAcquisition))
	o := orm.NewOrm()
	o.Using("read_only")

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get top product on cai
	var mx []*model.CustomerAcquisition
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			_, err = o.Raw("SELECT * from customer_acquisition_item cai "+
				"where cai.customer_acquisition_id = ?", v.ID).QueryRows(&v.CustomerAcquisitionItem)
			if err != nil && !errors.Is(err, orm.ErrNoRows) {
				return nil, total, err
			}
			if v.CustomerAcquisitionItem != nil {
				for _, cai := range v.CustomerAcquisitionItem {
					err = o.Raw("SELECT id, code, name from product p "+
						"where p.id = ?", cai.Product.ID).QueryRow(&cai.Product)
					if err != nil && !errors.Is(err, orm.ErrNoRows) {
						return nil, total, err
					}
				}
			}
		}
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// GetCustomerAcquisitions get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetFilterCustomerAcquisition(rq *orm.RequestQuery) (m []*model.CustomerAcquisition, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.CustomerAcquisition))
	o := orm.NewOrm()
	o.Using("read_only")

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.CustomerAcquisition
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			_, err = o.Raw("SELECT * from customer_acquisition_item cai "+
				"where cai.customer_acquisition_id = ?", v.ID).QueryRows(&v.CustomerAcquisitionItem)
			if err != nil && !errors.Is(err, orm.ErrNoRows) {
				return nil, total, err
			}
			if v.CustomerAcquisitionItem != nil {
				for _, cai := range v.CustomerAcquisitionItem {
					err = o.Raw("SELECT id, code, name from product p "+
						"where p.id = ?", cai.Product.ID).QueryRow(&cai.Product)
					if err != nil && !errors.Is(err, orm.ErrNoRows) {
						return nil, total, err
					}
				}
			}
		}
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func ValidCA(id int64) (div *model.CustomerAcquisition, e error) {
	div = &model.CustomerAcquisition{ID: id}
	e = div.Read("ID")

	return
}
