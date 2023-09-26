// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"time"

	"git.edenfarm.id/cuxs/dbredis"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetMerchantPointLog find a single data division using field and value condition.
func GetMerchantPointLog(field string, values ...interface{}) (*model.MerchantPointLog, error) {
	m := new(model.MerchantPointLog)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetMerchantPointLogs : function to get data from database based on parameters
func GetMerchantPointLogs(rq *orm.RequestQuery) (man []*model.MerchantPointLog, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.MerchantPointLog))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.MerchantPointLog
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterMerchantPointLogs : function to get data from database based on parameters with filtered permission
func GetFilterMerchantPointLogs(rq *orm.RequestQuery) (man []*model.MerchantPointLog, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.MerchantPointLog))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.MerchantPointLog
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetTotalPoint : function to get total point from loyalty point table
func GetTotalPoint(rq *orm.RequestQuery, m int64) (active float64, used float64, err error) {
	q, _ := rq.QueryReadOnly(new(model.MerchantPointLog))

	var mx []*model.MerchantPointLog
	var mz []*model.MerchantPointLog
	if _, err = q.Filter("status", 1).Filter("merchant_id", m).All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			active = active + v.PointValue
		}
	}

	if _, err = q.Filter("status", 2).Filter("merchant_id", m).All(&mz, rq.Fields...); err == nil {
		for _, v := range mz {
			used = used + v.PointValue
		}
	}

	return active, used, err
}

// ValidMerchantPointLog : function to check if id is valid in database
func ValidMerchantPointLog(id int64) (merchant *model.MerchantPointLog, e error) {
	merchant = &model.MerchantPointLog{ID: id}
	e = merchant.Read("ID")

	return
}

// CheckMerchantPointLogData : function to check data based on filter and exclude parameters
func CheckMerchantPointLogData(filter, exclude map[string]interface{}) (mpl []*model.MerchantPointLog, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.MerchantPointLog))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err = o.All(&mpl); err != nil {
		return nil, 0, err
	}

	return mpl, total, err
}

// GetMerchantPointLogsOfMerchant : function to get points data of each merchant
func GetMerchantPointLogsOfMerchant(areaID, businessTypeID, merchantID, perpage, page int64, period []string) (merchant []*model.Merchant, totalPoint float64, lastUpdated time.Time, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	type MerchantLogList struct {
		MerchantID       int64   `orm:"column(id)"`
		MerchantCode     string  `orm:"column(code)"`
		MerchantName     string  `orm:"column(name)"`
		BusinessTypeID   int64   `orm:"column(business_type_id)"`
		BusinessTypeName string  `orm:"column(business_type_name)"`
		AreaID           int64   `orm:"column(area_id)"`
		AreaName         string  `orm:"column(area_name)"`
		CurrentPoint     float64 `orm:"column(total_point)"`
		LogStatus        int8    `orm:"column(status)"`
		PointValue       float64 `orm:"column(point_value)"`
	}

	var (
		q, whereArea, whereBusinessType, whereMerchant, wherePeriod string
		values                                                      []interface{}
		merchantLogs                                                []*MerchantLogList
		merchantIDTemp, counter                                     int64
	)

	if areaID > 0 {
		whereArea = "and m.finance_area_id = ? "
		values = append(values, areaID)
	}

	if businessTypeID > 0 {
		whereBusinessType = "and m.business_type_id = ? "
		values = append(values, businessTypeID)
	}

	if merchantID > 0 {
		whereMerchant = "and m.id = ? "
		values = append(values, merchantID)
	}

	if len(period) > 0 {
		wherePeriod = "and mpl.created_date between ? and ? "
		values = append(values, period)
	}

	q = "select m.id, m.code, m.name, bt.id business_type_id, bt.name business_type_name, a.id area_id, a.name area_name, m.total_point, mpl.status, mpl.point_value " +
		"from merchant m " +
		"join merchant_point_log mpl on m.id = mpl.merchant_id " +
		"join business_type bt on m.business_type_id = bt.id " +
		"join area a on m.finance_area_id = a.id " +
		"where 1=1 " +
		whereArea +
		whereBusinessType +
		whereMerchant +
		wherePeriod +
		"order by m.id "
	if _, err = o.Raw(q, values).QueryRows(&merchantLogs); err != nil {
		return nil, 0, lastUpdated, err
	}

	merchantIDTemp = 0
	counter = -1
	for _, v := range merchantLogs {
		earnedPoint := float64(0)
		redeemedPoint := float64(0)
		if v.LogStatus == 1 {
			earnedPoint = v.PointValue
		} else if v.LogStatus == 2 {
			redeemedPoint = v.PointValue
		}

		if merchantIDTemp != v.MerchantID {
			merchantData := &model.Merchant{
				ID:            v.MerchantID,
				Code:          v.MerchantCode,
				Name:          v.MerchantName,
				TotalPoint:    v.CurrentPoint,
				EarnedPoint:   earnedPoint,
				RedeemedPoint: redeemedPoint,
			}

			merchantData.BusinessType = &model.BusinessType{ID: v.BusinessTypeID, Name: v.BusinessTypeName}
			merchantData.FinanceArea = &model.Area{ID: v.AreaID, Name: v.AreaName}

			merchant = append(merchant, merchantData)
			merchantIDTemp = v.MerchantID
			counter++
			continue
		}

		merchant[counter].EarnedPoint += earnedPoint
		merchant[counter].RedeemedPoint += redeemedPoint
	}

	dataLength := int64(len(merchant))
	startIdx := (page - 1) * perpage
	endIdx := startIdx + perpage

	if startIdx <= dataLength && endIdx <= dataLength {
		merchant = merchant[startIdx:endIdx]
	} else {
		if dataLength < startIdx {
			merchant = nil
		} else if dataLength < endIdx {
			merchant = merchant[startIdx:]
		}
	}

	totalPoint, lastUpdated, _ = GetTotalCurrentPoint()

	return merchant, totalPoint, lastUpdated, nil
}

func GetTotalCurrentPoint() (totalPoint float64, lastUpdated time.Time, err error) {
	var (
		pointMerchant []float64
		q             string
	)

	key := "total_current_eden_point"
	if !dbredis.Redis.CheckExistKey(key) {
		o := orm.NewOrm()
		o.Using("read_only")

		q = "select total_point from merchant"
		if _, err = o.Raw(q).QueryRows(&pointMerchant); err != nil {
			return 0, lastUpdated, err
		}

		for _, v := range pointMerchant {
			totalPoint += v
		}

		wib, _ := time.LoadLocation("Asia/Jakarta")
		currentTime, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), wib)

		dbredis.Redis.SetCache(key, totalPoint, 0)
		dbredis.Redis.SetCache(key+"_updated_date", currentTime, 0)

		return totalPoint, currentTime, nil
	}

	dbredis.Redis.GetCache(key, &totalPoint)
	dbredis.Redis.GetCache(key+"_updated_date", &lastUpdated)

	return totalPoint, lastUpdated, nil
}
