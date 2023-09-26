// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package edenpoint

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (epc *model.EdenPointCampaign, err error) {
	if r.Code, err = util.GenerateCode(r.Code, "eden_point_campaign"); err != nil {
		return nil, err
	}

	o := orm.NewOrm()
	o.Begin()

	epc = &model.EdenPointCampaign{
		Code:               r.Code,
		Name:               r.Name,
		CampaignFilterType: r.CampaignFilterType,
		StartDate:          r.StartTimestamp,
		EndDate:            r.EndTimestamp,
		ImageUrl:           r.ImageUrl,
		Multiple:           r.Multiplier,
		Status:             1,
		Note:               r.Note,
		CreatedAt:          time.Now(),
		CreatedBy:          r.Session.Staff.ID,
	}

	if r.CampaignFilterType == 1 {
		epc.Area = r.AreaStr
		epc.Archetype = r.ArchetypeStr
	} else if r.CampaignFilterType == 2 {
		epc.TagCustomer = r.CustomerTagStr
	}

	if _, err = o.Insert(epc); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, epc.ID, "eden_point_campaign", "create", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	epc.Area = util.EncIdInStr(epc.Area)
	epc.Archetype = util.EncIdInStr(epc.Archetype)
	epc.TagCustomer = util.EncIdInStr(epc.TagCustomer)

	return epc, nil
}

// Archive : function to update status data into archive
func Archive(r archiveRequest) (epc *model.EdenPointCampaign, err error) {
	o := orm.NewOrm()
	o.Begin()

	if err = r.EdenPointCampaign.Save("Status"); err != nil {
		o.Rollback()
		return nil, err
	}

	if err = log.AuditLogByUser(r.Session.Staff, r.EdenPointCampaign.ID, "eden_point_campaign", "archive", r.Note); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()
	return r.EdenPointCampaign, nil
}
