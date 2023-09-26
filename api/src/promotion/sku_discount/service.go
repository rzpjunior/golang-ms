// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sku_discount

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (sd *model.SkuDiscount, e error) {
	if r.Code, e = util.GenerateCode(r.Code, "sku_discount"); e != nil {
		return nil, e
	}

	o := orm.NewOrm()
	o.Begin()

	sd = &model.SkuDiscount{
		Code:           r.Code,
		Name:           r.Name,
		PriceSets:      r.PriceSet,
		Division:       r.Division,
		OrderChannels:  r.OrderChannel,
		StartTimestamp: r.StartTimestamp,
		EndTimestamp:   r.EndTimestamp,
		Note:           r.Note,
		Status:         1,
		CreatedAt:      time.Now(),
		CreatedBy:      r.Session.Staff.ID,
	}

	if _, e = o.Insert(sd); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, v := range r.Items {
		sdi := &model.SkuDiscountItem{
			SkuDiscount:         sd,
			Product:             v.Product,
			OverallQuota:        v.OverallQuota,
			OverallQuotaPerUser: v.OverallQuotaPerUser,
			DailyQuotaPerUser:   v.DailyQuotaPerUser,
			RemOverallQuota:     float64(v.OverallQuota),
			Budget:              v.Budget,
			RemBudget:           v.Budget,
			IsUseBudget:         v.IsUseBudget,
		}

		if _, e = o.Insert(sdi); e != nil {
			o.Rollback()
			return nil, e
		}

		for _, val := range v.Tiers {
			sdit := &model.SkuDiscountItemTier{
				SkuDiscountItem: sdi,
				TierLevel:       val.TierLevel,
				MinimumQty:      val.MinimumQty,
				DiscAmount:      val.Amount,
			}

			if _, e = o.Insert(sdit); e != nil {
				o.Rollback()
				return nil, e
			}
		}
	}

	o.Commit()

	log.AuditLogByUser(r.Session.Staff, sd.ID, "sku_discount", "create", "")

	return sd, e
}

func Archive(r archiveRequest) (*model.SkuDiscount, error) {
	o := orm.NewOrm()
	o.Begin()

	var err error

	r.SkuDiscount.ArchivedAt = time.Now()
	r.SkuDiscount.ArchivedBy = r.Session.Staff.ID
	if _, err = o.Update(r.SkuDiscount, "Status", "ArchivedAt", "ArchivedBy"); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	log.AuditLogByUser(r.Session.Staff, r.SkuDiscount.ID, "sku_discount", "archive", "")

	return r.SkuDiscount, nil
}
