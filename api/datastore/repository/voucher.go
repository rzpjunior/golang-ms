// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"strings"

	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/util"
)

// GetVoucher find a single data voucher using field and value condition.
func GetVoucher(field string, values ...interface{}) (*model.Voucher, error) {
	var (
		tagCustomerArr []string
		err            error
		qMark          string
	)
	m := new(model.Voucher)
	o := orm.NewOrm()
	o.Using("read_only")

	if err = o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	if m.TagCustomer != "" {
		tagCustomerArr = strings.Split(m.TagCustomer, ",")
		for range tagCustomerArr {
			qMark += "?,"
		}
		qMark = strings.TrimSuffix(qMark, ",")

		if err = o.Raw("select group_concat(name) from tag_customer where id in ("+qMark+") order by id asc", tagCustomerArr).QueryRow(&m.TagCustomerName); err != nil {
			return nil, err
		}
		m.TagCustomerName = strings.ReplaceAll(m.TagCustomerName, ",", ", ")
	}

	o.Raw("SELECT * FROM merchant WHERE id = ?", m.MerchantID).QueryRow(&m.Merchant)
	m.TagCustomer = util.DecryptIdInStr(m.TagCustomer)

	if m.MembershipLevelID != 0 {
		m.MembershipLevel, _ = ValidMembershipLevel(m.MembershipLevelID)

		if m.MembershipCheckpointID != 0 {
			m.MembershipCheckpoint, _ = ValidMembershipCheckpoint(m.MembershipCheckpointID, m.MembershipLevelID)
		}
	}

	// get order channel name from glossary
	orderChannelArr := strings.Split(m.ChannelVoucher, ",")
	qMark = ""
	for range orderChannelArr {
		qMark += "?,"
	}
	qMark = strings.TrimSuffix(qMark, ",")
	o.Raw("SELECT GROUP_CONCAT(note) FROM glossary WHERE attribute = 'order_channel' AND value_int in ("+qMark+")", orderChannelArr).QueryRow(&m.ChannelVoucher)

	o.LoadRelated(m, "VoucherItems", 2)
	o.LoadRelated(m, "VoucherContent", 0)

	return m, nil
}

// GetAreas get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetVouchers(rq *orm.RequestQuery) (m []*model.Voucher, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Voucher))

	// get data requested
	if total, err = q.Exclude("status", 3).All(&m, rq.Fields...); err != nil {
		return nil, 0, err
	}

	for _, v := range m {
		if v.MembershipLevelID != 0 {
			v.MembershipLevel = &model.MembershipLevel{ID: v.MembershipLevelID}
			v.MembershipLevel.Read("ID")
		}

		if v.MembershipCheckpointID != 0 {
			v.MembershipCheckpoint = &model.MembershipCheckpoint{ID: v.MembershipCheckpointID}
			v.MembershipCheckpoint.Read("ID")
		}

		if v.MerchantID != 0 {
			v.Merchant = &model.Merchant{ID: v.MerchantID}
			v.Merchant.Read("ID")
		}
	}

	return m, total, nil
}

func ValidVoucher(id int64) (v *model.Voucher, e error) {
	v = &model.Voucher{ID: id}
	e = v.Read("ID")

	return
}

// GetAreas get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetFilterVoucher(rq *orm.RequestQuery) (m []*model.Voucher, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Voucher))

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Voucher
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// CheckVoucherData : function to check data based on filter and exclude parameters
func CheckVoucherData(filter, exclude map[string]interface{}) (voucher []*model.Voucher, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.Voucher))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if countResult, err := o.All(&voucher); err == nil {
		return voucher, countResult, nil
	}

	return nil, 0, err
}
