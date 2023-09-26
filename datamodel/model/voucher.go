// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/util"
)

func init() {
	orm.RegisterModel(new(Voucher))
}

// Voucher model for voucher table.
type Voucher struct {
	ID                     int64      `orm:"column(id);auto" json:"-"`
	Area                   *Area      `orm:"column(area_id);null;rel(fk)" json:"area,omitempty"`
	Archetype              *Archetype `orm:"column(archetype_id);null;rel(fk)" json:"archetype,omitempty"`
	TagCustomer            string     `orm:"column(tag_customer);size(100);null" json:"tag_customer"`
	Code                   string     `orm:"column(code);size(50);null" json:"code"`
	RedeemCode             string     `orm:"column(redeem_code);size(20);null" json:"redeem_code"`
	Type                   int8       `orm:"column(type);null" json:"type"`
	Name                   string     `orm:"column(name);size(100);null" json:"name"`
	StartTimestamp         time.Time  `orm:"column(start_timestamp);type(timestamp);null" json:"start_timestamp"`
	EndTimestamp           time.Time  `orm:"column(end_timestamp);type(timestamp);null" json:"end_timestamp"`
	OverallQuota           int64      `orm:"column(overall_quota);" json:"overall_quota"`
	UserQuota              int64      `orm:"column(user_quota);" json:"user_quota"`
	RemOverallQuota        int64      `orm:"column(rem_overall_quota);" json:"rem_overall_quota"`
	MinOrder               float64    `orm:"column(min_order);null;digits(20);decimals(2)" json:"min_order"`
	DiscAmount             float64    `orm:"column(disc_amount);null;digits(20);decimals(2)" json:"disc_amount"`
	VoidReason             int8       `orm:"column(void_reason);null" json:"void_reason"`
	Note                   string     `orm:"column(note)" json:"note"`
	Status                 int8       `orm:"column(status);null" json:"status"`
	TagCustomerName        string     `orm:"-" json:"tag_customer_name,omitempty"`
	ChannelVoucher         string     `orm:"column(channel_voucher);null" json:"channel_voucher"`
	VoucherItem            int8       `orm:"column(voucher_item);null" json:"voucher_item"`
	MerchantID             int64      `orm:"column(merchant_id);null" json:"-"`
	MembershipLevelID      int64      `orm:"column(membership_level_id)" json:"-"`
	MembershipCheckpointID int64      `orm:"column(membership_checkpoint_id)" json:"-"`

	Merchant             *Merchant             `orm:"-" json:"merchant,omitempty"`
	VoucherContent       *VoucherContent       `orm:"reverse(one)" json:"voucher_content,omitempty"`
	VoucherItems         []*VoucherItem        `orm:"reverse(many)" json:"voucher_items,omitempty"`
	MembershipLevel      *MembershipLevel      `orm:"-" json:"membership_level,omitempty"`
	MembershipCheckpoint *MembershipCheckpoint `orm:"-" json:"membership_checkpoint,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *Voucher) MarshalJSON() ([]byte, error) {
	type Alias Voucher

	alias := &struct {
		ID            string `json:"id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		StatusConvert: util.ConvertStatusMaster(m.Status),
		Alias:         (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating User struct into user table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *Voucher) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting user data
// this also will truncated all data from all table
// that have relation with this user.
func (m *Voucher) Delete() (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		var i int64
		if i, err = o.Delete(m); i == 0 && err == nil {
			err = orm.ErrNoAffected
		}
		return
	}
	return orm.ErrNoRows
}

// Read execute select based on data struct that already
// assigned.
func (m *Voucher) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
