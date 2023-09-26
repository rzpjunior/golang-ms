// Copyright 2020 PT Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"
	"net/url"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// GetProspectiveCustomer find a single data item using field and value condition.
func GetProspectiveCustomer(field string, values ...interface{}) (*model.ProspectCustomer, error) {
	m := new(model.ProspectCustomer)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel(2).Limit(1).One(m); err != nil {
		return nil, err
	}

	m.SubDistrict.Read("ID")
	m.SubDistrict.Area.Read("ID")
	m.SubDistrict.District.Read("ID")
	m.SubDistrict.District.City.Read("ID")
	m.SubDistrict.District.City.Province.Read("ID")

	//// Initialize minio client object.
	minioClient, _ := minio.New(util.S3endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(util.S3accessKeyID, util.S3secretAccessKey, ""),
		Secure: true,
	})
	reqParams := make(url.Values)
	//// Retrieve URL valid for 60 second

	if m.IDCardImage != "" {
		tempIdCardImage := strings.Split(m.IDCardImage, "/")
		preSignedURLIdCardImage, _ := minioClient.PresignedGetObject(context.Background(), util.S3bucketNameImage, tempIdCardImage[4]+"/"+tempIdCardImage[5], time.Second*60, reqParams)
		m.IDCardImage = preSignedURLIdCardImage.String()
	}
	if m.SelfieImage != "" {
		tempSelfieImage := strings.Split(m.SelfieImage, "/")
		preSignedURLSelfieImage, _ := minioClient.PresignedGetObject(context.Background(), util.S3bucketNameImage, tempSelfieImage[4]+"/"+tempSelfieImage[5], time.Second*60, reqParams)
		m.SelfieImage = preSignedURLSelfieImage.String()

	}
	if m.TaxpayerImage != "" {
		tempTaxImage := strings.Split(m.TaxpayerImage, "/")
		preSignedURLTaxPayerImage, _ := minioClient.PresignedGetObject(context.Background(), util.S3bucketNameImage, tempTaxImage[4]+"/"+tempTaxImage[5], time.Second*60, reqParams)
		m.TaxpayerImage = preSignedURLTaxPayerImage.String()
	}

	// Returning images of outlet
	m.OutletPhotoArr = []string{}

	if m.OutletPhoto != "" {
		m.OutletPhotoArr = strings.Split(m.OutletPhoto, ",")
	}

	for _, op := range m.OutletPhotoArr {
		tempOutletImage := strings.Split(op, "/")
		preSignedURLOutletImage, _ := minioClient.PresignedGetObject(context.Background(), util.S3bucketNameImage, tempOutletImage[4]+"/"+tempOutletImage[5], time.Second*60, reqParams)
		m.OutletPhotoList = append(m.OutletPhotoList, preSignedURLOutletImage.String())
	}

	o.Raw("SELECT s.name FROM staff s where s.id = ?", m.SalespersonID).QueryRow(&m.Salesperson)
	o.Raw("SELECT g.value_name FROM glossary g WHERE g.table = 'prospect_customer' and g.attribute = 'decline_type' and g.value_int = ?", m.DeclineTypeID).QueryRow(&m.DeclineType)

	glossaryRegChannel, err := GetGlossaryMultipleValue("table", "prospect_customer", "attribute", "reg_channel", "value_int", m.RegChannel)
	if err != nil {
		m.RegChannelName = ""
	} else {
		m.RegChannelName = glossaryRegChannel.Note
	}

	return m, nil
}

// GetProspectiveCustomers get all data item that matched with query request parameters.
// returning slices of Item, total data without limit and error.
func GetProspectiveCustomers(rq *orm.RequestQuery) (m []*model.ProspectCustomer, total int64, err error) {
	// make new orm query
	o := orm.NewOrm()
	o.Using("read_only")

	q, _ := rq.QueryReadOnly(new(model.ProspectCustomer))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.ProspectCustomer
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			o.Raw("SELECT s.name from staff s WHERE s.id = ?", v.SalespersonID).QueryRow(&v.Salesperson)
		}
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// GetFilterProspectiveCustomer get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetFilterProspectiveCustomer(rq *orm.RequestQuery) (m []*model.ProspectCustomer, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.ProspectCustomer))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.ProspectCustomer
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func ValidProspectiveCustomer(id int64) (pc *model.ProspectCustomer, e error) {
	pc = &model.ProspectCustomer{ID: id}
	e = pc.Read("ID")

	return
}
