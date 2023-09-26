// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// GetMerchant find a single data division using field and value condition.
func GetMerchant(field string, values ...interface{}) (*model.Merchant, error) {
	var err error
	m := new(model.Merchant)
	o := orm.NewOrm()
	o.Using("read_only")

	if err = o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	o.Raw("select group_concat(distinct tc.name order by tc.id separator ',') "+
		"from merchant m "+
		"join tag_customer tc on concat(',', m.tag_customer, ',') like concat('%,', tc.id, ',%') "+
		"where m.id = ? "+
		"group by m.id", values[0]).QueryRow(&m.TagCustomerName)

	m.TagCustomer = util.DecryptIdInStr(m.TagCustomer)
	o.LoadRelated(m, "MerchantAccNum", 2)
	o.LoadRelated(m, "MerchantPriceSet", 1)

	if m.CreditLimit, err = CheckSingleCreditLimitData(m.BusinessType.ID, m.PaymentTerm.ID, m.BusinessTypeCreditLimit); err != nil {
		return nil, err
	}

	//// Initialize minio client object.
	minioClient, _ := minio.New(util.S3endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(util.S3accessKeyID, util.S3secretAccessKey, ""),
		Secure: true,
	})
	reqParams := make(url.Values)

	// extract ktp photos
	if m.KTPPhotosUrl != "" {
		ktpArr := strings.Split(m.KTPPhotosUrl, ",")
		for _, tp := range ktpArr {
			tempImage := strings.Split(tp, "/")
			preSignedURLImage, _ := minioClient.PresignedGetObject(context.Background(), util.S3bucketNameImage, tempImage[4]+"/"+tempImage[5], time.Second*60, reqParams)
			m.KTPPhotosUrlArr = append(m.KTPPhotosUrlArr, preSignedURLImage.String())
		}
	}

	// extract merchant photos
	if m.MerchantPhotosUrl != "" {
		ktpArr := strings.Split(m.MerchantPhotosUrl, ",")
		for _, tp := range ktpArr {
			m.MerchantPhotosUrlArr = append(m.MerchantPhotosUrlArr, tp)
		}
	}

	return m, nil
}

// GetMerchants : function to get data from database based on parameters
func GetMerchants(rq *orm.RequestQuery) (m []*model.Merchant, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Merchant))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Merchant
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			//// Initialize minio client object.
			minioClient, _ := minio.New(util.S3endpoint, &minio.Options{
				Creds:  credentials.NewStaticV4(util.S3accessKeyID, util.S3secretAccessKey, ""),
				Secure: true,
			})
			reqParams := make(url.Values)

			// extract ktp photos
			if v.KTPPhotosUrl != "" {
				ktpArr := strings.Split(v.KTPPhotosUrl, ",")
				for _, tp := range ktpArr {
					tempImage := strings.Split(tp, "/")
					preSignedURLImage, _ := minioClient.PresignedGetObject(context.Background(), util.S3bucketNameImage, tempImage[4]+"/"+tempImage[5], time.Second*60, reqParams)
					v.KTPPhotosUrlArr = append(v.KTPPhotosUrlArr, preSignedURLImage.String())
				}
			}

			// extract merchant photos
			if v.MerchantPhotosUrl != "" {
				ktpArr := strings.Split(v.MerchantPhotosUrl, ",")
				for _, tp := range ktpArr {
					v.MerchantPhotosUrlArr = append(v.MerchantPhotosUrlArr, tp)
				}
			}
		}
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterMerchants : function to get data from database based on parameters with filtered permission
func GetFilterMerchants(rq *orm.RequestQuery) (m []*model.Merchant, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Merchant))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Merchant
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			//// Initialize minio client object.
			minioClient, _ := minio.New(util.S3endpoint, &minio.Options{
				Creds:  credentials.NewStaticV4(util.S3accessKeyID, util.S3secretAccessKey, ""),
				Secure: true,
			})
			reqParams := make(url.Values)

			// extract ktp photos
			if v.KTPPhotosUrl != "" {
				ktpArr := strings.Split(v.KTPPhotosUrl, ",")
				for _, tp := range ktpArr {
					tempImage := strings.Split(tp, "/")
					preSignedURLImage, _ := minioClient.PresignedGetObject(context.Background(), util.S3bucketNameImage, tempImage[4]+"/"+tempImage[5], time.Second*60, reqParams)
					v.KTPPhotosUrlArr = append(v.KTPPhotosUrlArr, preSignedURLImage.String())
				}
			}

			// extract merchant photos
			if v.MerchantPhotosUrl != "" {
				ktpArr := strings.Split(v.MerchantPhotosUrl, ",")
				for _, tp := range ktpArr {
					v.MerchantPhotosUrlArr = append(v.MerchantPhotosUrlArr, tp)
				}
			}
		}
		return mx, total, nil
	}

	return nil, total, err
}

// ValidMerchant : function to check if id is valid in database
func ValidMerchant(id int64) (merchant *model.Merchant, e error) {
	merchant = &model.Merchant{ID: id}
	e = merchant.Read("ID")

	return
}

// ValidUserMerchant : function to check if id is valid in database
func ValidUserMerchant(id int64) (userMerchant *model.UserMerchant, e error) {
	userMerchant = &model.UserMerchant{ID: id}
	e = userMerchant.Read("ID")

	return
}

// CountNonDeletedMerchantBySalesTermId : function to check whether sales term id is still used by any active or archive merchant
func CountNonDeletedMerchantBySalesTermId(id int64) (countDeletedMerchant int64, e error) {
	m := new(model.Merchant)
	o := orm.NewOrm()
	o.Using("read_only")
	o1 := o.QueryTable(m)

	countDeletedMerchant, err := o1.Filter("term_payment_sls_id", id).Exclude("status", 3).Count()
	if err != nil {
		return 0, err
	}

	return countDeletedMerchant, nil
}

// CheckMerchantData : function to check data based on filter and exclude parameters
func CheckMerchantData(filter, exclude map[string]interface{}) (m []*model.Merchant, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	o1 := o.QueryTable(new(model.Merchant))

	for k, v := range filter {
		o1 = o1.Filter(k, v)
	}

	for k, v := range exclude {
		o1 = o1.Exclude(k, v)
	}

	if total, err = o1.All(&m); err == nil {
		return m, total, err
	}

	return nil, 0, err
}

// CountMerchantTagCustomer : function to count wether there are merchant that still using a tag customer
func CountMerchantTagCustomer(tagCustomerID int64) (countResult int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	q := "select count(id) from merchant where status != 3 and find_in_set(?, tag_customer)"
	if err = o.Raw(q, tagCustomerID).QueryRow(&countResult); err != nil {
		return 0, nil
	}

	return
}

// Get Remaining Bill of Merchant with calculate SO, SI and SP existing
func GetTotalRemainingBySOAndSI(merchantId int64) (remainingBill float64, e error) {
	var totalChargeSO float64
	var totalChargeSI float64
	var totalAmount float64

	if totalChargeSO, e = GetGrandTotalChargeSO(merchantId); e != nil {
		return remainingBill, e
	}
	if totalChargeSI, e = GetGrandTotalChargeSI(merchantId); e != nil {
		return remainingBill, e
	}
	if totalAmount, e = GetGrandTotalAmountSP(merchantId); e != nil {
		return remainingBill, e
	}

	remainingBill = totalChargeSO + totalChargeSI - totalAmount

	return remainingBill, nil
}

// GetCreditLimitRemainingMerchant: function to get credit limit remaining of Merchant
func GetCreditLimitRemainingMerchant(merchantID int64) (creditLimitRemaining float64, e error) {
	o := orm.NewOrm()

	// This query 'select' only use in this unique case, in other case you must use 'orm read only' to get data;
	q := "SELECT credit_limit_remaining FROM merchant WHERE id = ?"
	if e = o.Raw(q, merchantID).QueryRow(&creditLimitRemaining); e != nil {
		return 0, nil
	}

	return
}

// GetMerchantsDistributionNetwork : function to get data from database based on parameters
func GetMerchantsDistributionNetwork(rq *orm.RequestQuery) (m []*model.Merchant, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Merchant))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Merchant
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err != nil {
		return nil, total, err
	}

	for _, v := range mx {
		v.RemainingOutstanding = v.CreditLimitAmount - v.RemainingCreditLimitAmount
	}
	return mx, total, nil
}

// GetMerchantTopProduct find a single data top product using field and value condition.
func GetMerchantTopProduct(fromDate, toDate string, field string, values ...interface{}) (*model.Product, error) {
	var (
		err        error
		topProduct *model.Product
	)
	m := new(model.Merchant)
	o := orm.NewOrm()
	o.Using("read_only")

	if err = o.QueryTable(m).Filter(field, values...).Limit(1).One(m); err != nil {
		return nil, err
	}

	query := ""
	if fromDate != "" && toDate != "" {
		query = fmt.Sprintf("AND so.recognition_date BETWEEN '%s' AND '%s' ", fromDate, toDate)
	}

	err = o.Raw(
		"SELECT p.*, count(soi.product_id) total_product "+
			"FROM sales_order_item soi "+
			"JOIN sales_order so ON so.id = soi.sales_order_id "+
			"JOIN branch b ON b.id = so.branch_id "+
			"JOIN product p ON p.id = soi.product_id "+
			"WHERE b.merchant_id = ? AND so.status = 2 AND so.order_type_sls_id = 13 "+query+
			"GROUP BY soi.product_id "+
			"ORDER BY total_product DESC", m.ID).QueryRow(&topProduct)
	if err != nil && !errors.Is(err, orm.ErrNoRows) {
		return nil, err
	}

	return topProduct, nil
}

// GetMerchantOrderPerformance find summary order performance using field and value condition.
func GetMerchantOrderPerformance(field string, productId int, fromDate, toDate time.Time, values ...interface{}) (*model.MerchantOrderPerformance, error) {
	var (
		err                        error
		so                         []*model.SalesOrder
		product                    *model.Product
		totalWeights, totalCharges float64
		totalOrder                 int64
	)
	m := new(model.Merchant)
	o := orm.NewOrm()
	o.Using("read_only")

	if err = o.QueryTable(m).Filter(field, values...).Limit(1).One(m); err != nil {
		return nil, err
	}

	// validation product
	product, err = ValidProduct(int64(productId))
	if err != nil {
		return nil, err
	}

	_, err = o.Raw(
		"SELECT so.total_charge, so.total_weight FROM sales_order_item soi "+
			"JOIN sales_order so ON so.id = soi.sales_order_id "+
			"JOIN branch b ON b.id = so.branch_id "+
			"WHERE so.status = 2 AND so.order_type_sls_id = 13 AND soi.product_id = ? AND "+
			"b.merchant_id = ? AND so.recognition_date BETWEEN ? AND ?",
		productId, m.ID, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).QueryRows(&so)
	if err != nil && errors.Is(err, orm.ErrNoRows) {
		return nil, err
	}

	for _, order := range so {
		totalCharges += order.TotalCharge
		totalWeights += order.TotalWeight
		totalOrder += 1
	}
	avgSales := float64(0)
	if totalOrder > 0 {
		avgSales = totalCharges / float64(totalOrder)
	}

	return &model.MerchantOrderPerformance{
		ProductId:    common.Encrypt(product.ID),
		ProductName:  product.Name,
		QtySell:      totalWeights,
		AverageSales: avgSales,
		OrderTotal:   totalOrder,
	}, nil
}

// GetMerchantPaymentPerformance find summary payment performance using field and value condition.
func GetMerchantPaymentPerformance(field string, values ...interface{}) (*model.MerchantPaymentPerformance, error) {
	var (
		err error
		invoiceAmount, paymentAmount, remainingOutstanding, orderDueAmount, overDuePercentage, creditLimitPercentage,
		totalPaymentPerc, avgPaymentAmount float64
		totalDaysDiffPayment, totalPaymentFinished int
		si                                         []*model.SalesInvoice
		daysDiffPayment                            = map[int64]int{}
		paymentFinished                            = map[int64]int{}

		siPaid      = map[int64]float64{}
		loc, _      = time.LoadLocation("Asia/Jakarta")
		currentTime = time.Now().In(loc)
	)

	m := new(model.Merchant)
	o := orm.NewOrm()
	o.Using("read_only")

	if err = o.QueryTable(m).Filter(field, values...).Limit(1).One(m); err != nil {
		return nil, err
	}

	if m.CreditLimitAmount > 0 {
		creditLimitPercentage = m.RemainingCreditLimitAmount / m.CreditLimitAmount * 100
	}

	// get invoices
	_, err = o.Raw(
		"SELECT si.id, si.total_charge, si.recognition_date, si.due_date, si.status "+
			"FROM sales_invoice si "+
			"JOIN sales_order so ON so.id = si.sales_order_id "+
			"JOIN branch b ON b.id = so.branch_id "+
			"JOIN merchant m ON m.id = b.merchant_id "+
			"WHERE m.id = ? AND si.status NOT IN (3,4)", m.ID).QueryRows(&si)
	if err != nil && errors.Is(err, orm.ErrNoRows) {
		return nil, err
	}

	type CustomSalesPayment struct {
		SalesPaymentceId        int64     `orm:"column(sales_payment_id)"`
		SalesInvoiceId          int64     `orm:"column(sales_invoice_id)"`
		PaymentRecoginitionDate time.Time `orm:"column(payment_recognition_date)"`
		InvoiceRecoginitionDate time.Time `orm:"column(invoice_recognition_date)"`
		PaymentAmount           float64   `orm:"column(payment_amount)"`
		PaymentStatus           int8      `orm:"column(payment_status)"`
		CreatedAt               time.Time `orm:"column(created_at)"`
	}

	var cSP []*CustomSalesPayment
	_, err = o.Raw(
		"SELECT sp.id sales_payment_id, sp.sales_invoice_id, sp.recognition_date payment_recognition_date, sp.amount payment_amount, "+
			"sp.status payment_status, si.recognition_date invoice_recognition_date, sp.created_at FROM sales_payment sp "+
			"JOIN sales_invoice si ON si.id = sp.sales_invoice_id "+
			"JOIN sales_order so ON so.id = si.sales_order_id "+
			"JOIN branch b ON b.id = so.branch_id "+
			"WHERE b.merchant_id = ? AND sp.status IN (2,5)",
		m.ID).QueryRows(&cSP)
	if err != nil && errors.Is(err, orm.ErrNoRows) {
		return nil, err
	}

	for _, payment := range cSP {
		siPaid[payment.SalesInvoiceId] += payment.PaymentAmount
		if payment.PaymentStatus == 5 {
			daysDiffPayment[payment.SalesInvoiceId] += int(payment.CreatedAt.Sub(payment.InvoiceRecoginitionDate).Hours() / 24)
		} else {
			daysDiffPayment[payment.SalesInvoiceId] += int(payment.PaymentRecoginitionDate.Sub(payment.InvoiceRecoginitionDate).Hours() / 24)
		}
		paymentFinished[payment.SalesInvoiceId] += 1
		paymentAmount += payment.PaymentAmount
	}

	for _, invoice := range si {
		totalPaymentPerc += siPaid[invoice.ID] / invoice.TotalCharge * 100
		totalDaysDiffPayment += daysDiffPayment[invoice.ID]
		totalPaymentFinished += paymentFinished[invoice.ID]
		invoiceAmount += invoice.TotalCharge

		if invoice.Status != 2 && invoice.DueDate.Before(currentTime) {
			orderDueAmount += (invoice.TotalCharge - siPaid[invoice.ID])
		}
	}

	remainingOutstanding = invoiceAmount - paymentAmount
	if totalPaymentFinished > 0 {
		avgPaymentAmount = paymentAmount / float64(totalPaymentFinished)
	}

	// calculation Overduedebt percentage
	if orderDueAmount > 0 {
		overDuePercentage = orderDueAmount / remainingOutstanding * 100
	}

	// calculation average payment
	if totalPaymentFinished > 0 {
		totalPaymentPerc = totalPaymentPerc / float64(totalPaymentFinished)
		totalDaysDiffPayment = totalDaysDiffPayment / totalPaymentFinished
	}

	return &model.MerchantPaymentPerformance{
		CreditLimitAmount:                   m.CreditLimitAmount,
		CreditLimitRemaining:                m.RemainingCreditLimitAmount,
		RemainingOutstanding:                remainingOutstanding,
		CreditLimitUsageRemainingPercentage: creditLimitPercentage,
		OverdueDebtAmount:                   orderDueAmount,
		OverdueDebtRemainingPercentage:      overDuePercentage,
		AveragePaymentAmount:                avgPaymentAmount,
		AveragePaymentPercentage:            totalPaymentPerc,
		AveragePaymentPeriod:                totalDaysDiffPayment,
	}, nil
}
