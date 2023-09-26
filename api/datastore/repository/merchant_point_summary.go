// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetMerchantPointSummary find a single data division using field and value condition.
func GetMerchantPointSummary(field string, values ...interface{}) (*model.MerchantPointSummary, error) {
	m := new(model.MerchantPointSummary)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetMerchantPointSummarys : function to get data from database based on parameters
func GetMerchantPointSummarys(rq *orm.RequestQuery) (mps []*model.MerchantPointSummary, total int64, err error) {
	var mpsTemp []*model.MerchantPointSummary

	q, _ := rq.QueryReadOnly(new(model.MerchantPointSummary))
	o := orm.NewOrm()
	o.Using("read_only")

	cond := q.GetCond()

	cond1 := orm.NewCondition()
	cond1 = cond1.And("earnedpoint__gt", 0).Or("redeemedpoint__gt", 0)

	cond = cond.AndCond(cond1)

	q = q.SetCond(cond)

	if _, err = q.OrderBy("merchant").All(&mpsTemp, rq.Fields...); err != nil {
		return nil, total, err
	}

	merchantID := int64(0)
	counter := -1
	for _, v := range mpsTemp {
		if merchantID != v.Merchant.ID {
			o.LoadRelated(v, "Merchant", 1)
			mps = append(mps, v)
			counter++
			merchantID = v.Merchant.ID
			continue
		}

		mps[counter].EarnedPoint += v.EarnedPoint
		mps[counter].RedeemedPoint += v.RedeemedPoint
	}

	return mps, total, nil
}

// ValidMerchantPointSummary : function to check if id is valid in database
func ValidMerchantPointSummary(id int64) (mps *model.MerchantPointSummary, e error) {
	mps = &model.MerchantPointSummary{ID: id}
	e = mps.Read("ID")

	return
}

// GetMerchantPointSummaryData : function to check data based on filter and exclude parameters
func GetMerchantPointSummaryData(filter, exclude map[string]interface{}) (mps []*model.MerchantPointSummary, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.MerchantPointSummary))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err = o.All(&mps); err != nil {
		return nil, 0, err
	}

	return mps, total, err
}
