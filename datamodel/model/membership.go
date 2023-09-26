// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(MembershipLevel))
	orm.RegisterModel(new(MembershipAdvantage))
	orm.RegisterModel(new(MembershipCheckpoint))
	orm.RegisterModel(new(MembershipReward))
}

// MembershipLevel : struct to hold membership level data
type MembershipLevel struct {
	ID       int64  `orm:"column(id);auto" json:"-"`
	Code     string `orm:"column(code)" json:"code"`
	Name     string `orm:"column(name)" json:"name"`
	Level    int8   `orm:"column(level)" json:"level"`
	ImageUrl string `orm:"column(image_url)" json:"image_url"`
	Status   int8   `orm:"column(status)" json:"status"`

	MembershipAdvantages  []*MembershipAdvantage  `orm:"-" json:"membership_advantage,omitempty"`
	MembershipCheckpoints []*MembershipCheckpoint `orm:"-" json:"membership_checkpoint,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *MembershipLevel) MarshalJSON() ([]byte, error) {
	type Alias MembershipLevel

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *MembershipLevel) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *MembershipLevel) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}

// MembershipAdvantage : struct to hold membership advantage data
type MembershipAdvantage struct {
	ID          int64  `orm:"column(id);auto" json:"-"`
	Name        string `orm:"column(name)" json:"name"`
	Description string `orm:"column(description)" json:"description"`
	ImageUrl    string `orm:"column(image_url)" json:"image_url"`
	LinkUrl     string `orm:"column(link_url)" json:"link_url"`
	Status      int8   `orm:"column(status)" json:"status"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *MembershipAdvantage) MarshalJSON() ([]byte, error) {
	type Alias MembershipAdvantage

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Read : function to get data from database
func (m *MembershipAdvantage) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}

// MembershipCheckpoint : struct to hold membership checkpoint data
type MembershipCheckpoint struct {
	ID                int64   `orm:"column(id);auto" json:"-"`
	Checkpoint        int8    `orm:"column(checkpoint)" json:"checkpoint"`
	TargetAmount      float64 `orm:"column(target_amount)" json:"target_amount"`
	Status            int8    `orm:"column(status)" json:"status"`
	MembershipLevelID int64   `orm:"column(membership_level_id)" json:"-"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *MembershipCheckpoint) MarshalJSON() ([]byte, error) {
	type Alias MembershipCheckpoint

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Read : function to get data from database
func (m *MembershipCheckpoint) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}

// MembershipReward : struct to hold membership checkpoint data
type MembershipReward struct {
	ID                 int64   `orm:"column(id);auto" json:"-"`
	OpenedImageUrl     string  `orm:"column(opened_image_url)" json:"opened_image_url"`
	ClosedImageUrl     string  `orm:"column(closed_image_url)" json:"closed_image_url"`
	BackgroundImageUrl string  `orm:"column(background_image_url)" json:"background_image_url"`
	RewardLevel        int8    `orm:"column(reward_level)" json:"reward_level"`
	MaxAmount          float64 `orm:"column(max_amount)" json:"max_amount"`
	Status             int8    `orm:"column(status)" json:"status"`
	Description        string  `orm:"column(description)" json:"description"`
	IsPassed           int8    `orm:"-" json:"is_passed,omitempty"`
	CurrentPercentage  float64 `orm:"-" json:"current_percentage,omitempty"`
	RemainingAmount    float64 `orm:"-" json:"remaining_amount,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *MembershipReward) MarshalJSON() ([]byte, error) {
	type Alias MembershipReward

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Read : function to get data from database
func (m *MembershipReward) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
