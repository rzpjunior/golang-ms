// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sales_group

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"strconv"
)

// Save : function to insert data requested into database
func Save(r createRequest) (sg *model.SalesGroup, e error) {
	r.Code, e = util.GenerateCode(r.Code, "sales_group")
	o := orm.NewOrm()
	o.Begin()
	var arrSgi []*model.SalesGroupItem

	if e == nil {
		sg = &model.SalesGroup{
			Code:         r.Code,
			Name:         r.Name,
			BusinessType: r.BusinessType,
			SalesManager: r.SlsMan,
			Area:         r.Area,
			City:         r.CityStr,
			Status:       int8(1),
			CreatedAt:    time.Now(),
			CreatedBy:    r.Session.Staff.ID,
		}

		if _, e = o.Insert(sg); e != nil {
			o.Rollback()
			return nil, e
		}

		for _, v := range r.SubDistrict {
			idSubDist, _ := strconv.ParseInt(v, 10, 64)
			sgi := &model.SalesGroupItem{
				SalesGroup:  sg,
				SubDistrict: &model.SubDistrict{ID: idSubDist},
			}
			arrSgi = append(arrSgi, sgi)
		}

		if _, e = o.InsertMulti(100, &arrSgi); e != nil {
			o.Rollback()
			return
		}

		if e = log.AuditLogByUser(r.Session.Staff, sg.ID, "sales_group", "create", ""); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	o.Commit()

	return sg, e
}

// Archive : function to update status data into archive
func Archive(r archiveRequest) (sg *model.SalesGroup, e error) {
	o := orm.NewOrm()
	o.Begin()

	sg = &model.SalesGroup{
		ID:     r.ID,
		Status: int8(2),
	}

	if _, e = o.Update(sg, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, t := range r.Staff {
		stf := &model.Staff{
			ID:           t.ID,
			SalesGroupID: 0,
		}
		if _, e = o.Update(stf, "SalesGroupID"); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, sg.ID, "sales_group", "archive", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return sg, e
}

// Update : function to insert data requested into database
func Update(r updateRequest) (sg *model.SalesGroup, e error) {
	o := orm.NewOrm()
	o.Begin()

	var keepItemsId []int64
	var isCreated bool

	sg = &model.SalesGroup{
		ID:            r.ID,
		Name:          r.Name,
		BusinessType:  r.BusinessType,
		SalesManager:  r.SlsMan,
		City:          r.CityStr,
		LastUpdatedAt: time.Now(),
		LastUpdatedBy: r.Session.Staff.ID,
	}

	if _, e = o.Update(sg, "Name", "BusinessType", "SalesManager", "City", "LastUpdatedAt", "LastUpdatedBy"); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, v := range r.SubDistrict {
		idSubDist, _ := strconv.ParseInt(v, 10, 64)
		sgi := &model.SalesGroupItem{
			SalesGroup:  sg,
			SubDistrict: &model.SubDistrict{ID: idSubDist},
		}
		if isCreated, sgi.ID, e = o.ReadOrCreate(sgi, "SalesGroup", "SubDistrict"); e == nil {
			if !isCreated {
				sgi.SalesGroup = sg
				sgi.SubDistrict = &model.SubDistrict{ID: idSubDist}
				if _, e = o.Update(sgi, "SalesGroup", "SubDistrict"); e != nil {
					o.Rollback()
					return nil, e
				}
			}
		} else {
			o.Rollback()
			return nil, e
		}
		keepItemsId = append(keepItemsId, sgi.ID)
	}

	if _, e := o.QueryTable(new(model.SalesGroupItem)).Filter("sales_group_id", r.SalesGroup.ID).Exclude("ID__in", keepItemsId).Delete(); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, sg.ID, "sales_group", "update", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return sg, e
}
