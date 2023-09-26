// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/project-version2/api/util"
)

// GetBranch find a single data using field and value condition.
func GetBranch(field string, values ...interface{}) (*model.Branch, error) {
	m := new(model.Branch)
	o := orm.NewOrm()
	o.Using("read_only")
	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	m.SubDistrict.District.City.Province.Read()

	return m, nil
}

// GetBranchs : function to get data from database based on parameters
func GetBranchs(rq *orm.RequestQuery) (m []*model.Branch, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Branch))

	if total, err = q.Exclude("status", 3).Exclude("merchant__status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	if _, err = q.Exclude("status", 3).Exclude("merchant__status", 3).All(&m, rq.Fields...); err == nil {
		return m, total, nil
	}

	return m, total, err
}

// GetFilterBranchs : function to get data from database based on parameters with filtered permission
func GetFilterBranchs(rq *orm.RequestQuery) (m []*model.Branch, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Branch))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Branch
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidBranch : function to check if id is valid in database
func ValidBranch(id int64) (branch *model.Branch, e error) {
	branch = &model.Branch{ID: id}
	e = branch.Read("ID")

	return
}

// CheckBranchData : function to check data based on filter and exclude parameters
func CheckBranchData(filter map[string]interface{}, exclude map[string]interface{}) (branch []*model.Branch, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.Branch))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err := o.All(&branch); err == nil {
		return branch, total, nil
	}

	return nil, 0, err
}

// GetBranchsByMerchantId : get all branch data based on merchant id. Used for detail agent and create SO
func GetBranchsByMerchantId(id int64, isMain ...bool) (branch []*model.Branch, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	var status int
	o.Raw("select m.status from merchant m where m.id = ?", id).QueryRow(&status)

	q := o.QueryTable(new(model.Branch)).Filter("merchant_id", id)

	if status == 1 {
		q = q.Exclude("status__in", 2, 3)
	} else {
		q = q.Exclude("status", 3)
	}

	if len(isMain) > 0 && isMain[0] {
		q = q.Filter("main_branch", 1)
	}

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, 0, err
	}

	if _, err = q.RelatedSel().OrderBy("main_branch").All(&branch); err == nil {
		for k, v := range branch {
			err = o.Raw("select group_concat(distinct tc.name order by tc.id separator ',') "+
				"from merchant m "+
				"join tag_customer tc on concat(',', m.tag_customer, ',') like concat('%,', tc.id, ',%') "+
				"where m.id = ? "+
				"group by m.id", v.Merchant.ID).QueryRow(&branch[k].Merchant.TagCustomerName)

			branch[k].Merchant.TagCustomer = util.DecryptIdInStr(branch[k].Merchant.TagCustomer)

			o.Raw("select * from merchant_price_set mps join area a2 on mps.area_id = a2.id where mps.merchant_id = ? order by area_id asc", v.Merchant.ID).QueryRows(&v.Merchant.MerchantPriceSet)
			for _, v := range v.Merchant.MerchantPriceSet {
				o.Raw("select * from area a where a.id=?", v.Area.ID).QueryRow(v.Area)
				o.Raw("select * from price_set ps where ps.id =?", v.PriceSet.ID).QueryRow(v.PriceSet)
			}
			o.LoadRelated(v.Merchant, "MerchantAccNum", 1)

			if v.Merchant.CreditLimit, err = CheckSingleCreditLimitData(v.Merchant.BusinessType.ID, v.Merchant.PaymentTerm.ID, v.Merchant.BusinessTypeCreditLimit); err != nil {
				return nil, total, err
			}

		}
		return branch, total, nil
	}

	return nil, total, err
}

// GetBranchsByMerchantIdforUpdate : get all branch data based on merchant id. Used for update price set in module agent
func GetBranchsByMerchantIdforUpdate(id int64, isMain ...bool) (branch []*model.Branch, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	q := o.QueryTable(new(model.Branch)).Filter("merchant_id", id).Exclude("status", 3)

	if len(isMain) > 0 && isMain[0] {
		q = q.Filter("main_branch", 1)
	}

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, 0, err
	}

	if _, err = q.RelatedSel().OrderBy("main_branch").All(&branch); err == nil {
		for k, v := range branch {
			err = o.Raw("select group_concat(distinct tc.name order by tc.id separator ',') "+
				"from merchant m "+
				"join tag_customer tc on concat(',', m.tag_customer, ',') like concat('%,', tc.id, ',%') "+
				"where m.id = ? "+
				"group by m.id", v.Merchant.ID).QueryRow(&branch[k].Merchant.TagCustomerName)

			branch[k].Merchant.TagCustomer = util.DecryptIdInStr(branch[k].Merchant.TagCustomer)

			o.Raw("select * from merchant_price_set mps join area a2 on mps.area_id = a2.id where mps.merchant_id = ? order by area_id asc", v.Merchant.ID).QueryRows(&v.Merchant.MerchantPriceSet)
			for _, v := range v.Merchant.MerchantPriceSet {
				o.Raw("select * from area a where a.id=?", v.Area.ID).QueryRow(v.Area)
				o.Raw("select * from price_set ps where ps.id =?", v.PriceSet.ID).QueryRow(v.PriceSet)
			}
			o.LoadRelated(v.Merchant, "MerchantAccNum", 1)

		}

		return branch, total, nil
	}

	return nil, total, err
}

// CountActiveBranchByMerchantId : count active branch based on merchant id
func CountActiveBranchByMerchantId(id int64) (total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")
	if countActiveBranch, err := o.QueryTable(new(model.Branch)).Filter("merchant_id", id).Filter("status", 1).Count(); err == nil {
		return countActiveBranch, nil
	}

	return 0, err
}
