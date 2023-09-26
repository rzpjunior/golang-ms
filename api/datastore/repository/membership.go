// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetMembershipLevels : function to get membership level data from database
func GetMembershipLevels(rq *orm.RequestQuery) (ml []*model.MembershipLevel, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.MembershipLevel))

	if total, err = q.Exclude("status", 3).All(&ml, rq.Fields...); err != nil {
		return nil, 0, err
	}

	return ml, total, nil
}

// GetMembershipAdvantages : function to get membership advantage data from database
func GetMembershipAdvantages(rq *orm.RequestQuery, levelID int64) (ma []*model.MembershipAdvantage, total int64, err error) {
	if levelID == 0 {
		q, _ := rq.QueryReadOnly(new(model.MembershipAdvantage))

		if total, err = q.Exclude("status", 3).All(&ma, rq.Fields...); err != nil {
			return nil, 0, err
		}

		return ma, total, nil
	}

	var orderBy string

	if len(rq.OrderBy) > 0 {
		orderBy = rq.OrderBy[0]
		if orderBy[0:1] == "-" {
			orderBy = orderBy[1:] + " desc"
		}

		orderBy = "order by " + orderBy + " "
	}

	o := orm.NewOrm()
	o.Using("read_only")

	if total, err = o.Raw("select ma.* from membership_level_advantage mla join membership_advantage ma on mla.membership_advantage_id = ma.id where mla.membership_level_id = ? "+orderBy+"limit ?,?", levelID, rq.Offset, rq.Limit).QueryRows(&ma); err != nil {
		return nil, 0, err
	}

	return ma, total, nil
}

// GetMembershipCheckpoints : function to get membership checkpoint data from database
func GetMembershipCheckpoints(rq *orm.RequestQuery) (mc []*model.MembershipCheckpoint, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.MembershipCheckpoint))

	if total, err = q.Exclude("status", 3).All(&mc, rq.Fields...); err != nil {
		return nil, 0, err
	}

	return mc, total, nil
}

// GetMembershipLevel : function to get detail data of membership level
func GetMembershipLevel(field string, values ...interface{}) (*model.MembershipLevel, error) {
	var err error
	m := new(model.MembershipLevel)
	o := orm.NewOrm()
	o.Using("read_only")

	if err = o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetMembershipAdvantage : function to get detail data of membership advantage
func GetMembershipAdvantage(field string, values ...interface{}) (*model.MembershipAdvantage, error) {
	var err error
	m := new(model.MembershipAdvantage)
	o := orm.NewOrm()
	o.Using("read_only")

	if err = o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetMembershipCheckpoint : function to get detail data of membership checkpoint
func GetMembershipCheckpoint(field string, values ...interface{}) (*model.MembershipCheckpoint, error) {
	var err error
	m := new(model.MembershipCheckpoint)
	o := orm.NewOrm()
	o.Using("read_only")

	if err = o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// ValidMembershipLevel : function to check if id is valid in database
func ValidMembershipLevel(id int64) (m *model.MembershipLevel, e error) {
	m = &model.MembershipLevel{ID: id}
	e = m.Read("ID")

	return
}

// ValidMembershipAdvantage : function to check if id is valid in database
func ValidMembershipAdvantage(id int64) (m *model.MembershipAdvantage, e error) {
	m = &model.MembershipAdvantage{ID: id}
	e = m.Read("ID")

	return
}

// ValidMembershipCheckpoint : function to check if id is valid in database
func ValidMembershipCheckpoint(id, levelID int64) (m *model.MembershipCheckpoint, e error) {
	m = &model.MembershipCheckpoint{ID: id, MembershipLevelID: levelID}
	e = m.Read("ID", "MembershipLevelID")

	return
}
