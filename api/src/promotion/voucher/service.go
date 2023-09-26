// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package voucher

import (
	"strings"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (v *model.Voucher, e error) {
	r.Code, e = util.GenerateCode(r.Code, "voucher")
	if e == nil {
		v = &model.Voucher{
			Area:            r.Area,
			Archetype:       r.ArcheType,
			TagCustomer:     r.CustomerTagStr,
			Code:            r.Code,
			RedeemCode:      r.RedeemCode,
			Type:            r.Type,
			Name:            r.Name,
			StartTimestamp:  r.StartTimestamp,
			EndTimestamp:    r.EndTimestamp,
			OverallQuota:    r.OverallQuotaInt,
			RemOverallQuota: r.OverallQuotaInt,
			UserQuota:       r.UserQuotaInt,
			MinOrder:        r.MinOrderFloat,
			DiscAmount:      r.DiscAmountFloat,
			Note:            r.Note,
			Status:          int8(1),
			ChannelVoucher:  strings.Join(r.ChannelVoucher, ","),
			VoucherItem:     r.HasVoucherItem,
		}
		if r.Merchant != nil {
			v.MerchantID = r.Merchant.ID
		}
		if r.MembershipLevel != nil {
			v.MembershipLevelID = r.MembershipLevel.ID
			v.MembershipCheckpointID = r.MembershipCheckpoint.ID
		}
		if e = v.Save(); e == nil {
			if r.IsMobile {
				vc := &model.VoucherContent{
					Voucher:        v,
					ImageUrl:       r.ImageUrl,
					TermConditions: r.TermConditions,
				}
				vc.Save()
			}
			if r.HasVoucherItem == 1 {
				for _, vi := range r.VoucherItem {
					mvi := &model.VoucherItem{
						Voucher:      v,
						Product:      vi.Product,
						MinOrderDisc: vi.MinQtyDisc,
					}
					mvi.Save()
				}
			}

			e = log.AuditLogByUser(r.Session.Staff, v.ID, "voucher", "create", "")
		}
	}

	return v, e
}

// Archive : function to update status data into archive
func Archive(r archiveRequest) (v *model.Voucher, e error) {
	v = &model.Voucher{
		ID:         r.ID,
		Status:     int8(2),
		VoidReason: r.VoidReason,
	}

	if e = v.Save("Status", "VoidReason"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, v.ID, "voucher", "archive", "")
	}

	return v, e
}

func CreateBulky(c bulkyRequest) (e error) {
	var arrVoucher []*model.Voucher

	for _, v := range c.Sheet {
		v.Code, e = util.GenerateCode(v.Code, "voucher")
		vou := &model.Voucher{
			Area:            v.Area,
			Archetype:       v.Archetype,
			TagCustomer:     v.CustomerTagID,
			Code:            v.Code,
			RedeemCode:      v.RedeemCode,
			Type:            v.VoucherType,
			Name:            v.VoucherName,
			StartTimestamp:  v.StartTimestamp,
			EndTimestamp:    v.EndTimestamp,
			OverallQuota:    v.OverallQuota,
			RemOverallQuota: v.OverallQuota,
			UserQuota:       v.UserQuota,
			MinOrder:        v.MinOrder,
			DiscAmount:      v.DiscountAmount,
			Note:            v.Note,
			Status:          int8(1),
		}
		if v.Merchant != nil {
			vou.MerchantID = v.Merchant.ID
		}
		if v.MembershipLevelModel != nil {
			vou.MembershipLevelID = v.MembershipLevelModel.ID
			vou.MembershipCheckpointID = v.MembershipCheckpointModel.ID
		}
		arrVoucher = append(arrVoucher, vou)
		e = log.AuditLogByUser(c.Session.Staff, 0, "voucher", "create_bulky", v.Code)
	}
	o := orm.NewOrm()
	o.Begin()
	if _, e = o.InsertMulti(100, &arrVoucher); e != nil {
		o.Rollback()
	}
	o.Commit()
	return
}
