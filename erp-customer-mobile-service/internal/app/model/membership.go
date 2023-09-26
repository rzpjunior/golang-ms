// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

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

// MembershipAdvantage : struct to hold membership advantage data
type MembershipAdvantage struct {
	ID          int64  `orm:"column(id);auto" json:"-"`
	Name        string `orm:"column(name)" json:"name"`
	Description string `orm:"column(description)" json:"description"`
	ImageUrl    string `orm:"column(image_url)" json:"image_url"`
	LinkUrl     string `orm:"column(link_url)" json:"link_url"`
	Status      int8   `orm:"column(status)" json:"status"`
}

// MembershipCheckpoint : struct to hold membership checkpoint data
type MembershipCheckpoint struct {
	ID                int64   `orm:"column(id);auto" json:"-"`
	Checkpoint        int8    `orm:"column(checkpoint)" json:"checkpoint"`
	TargetAmount      float64 `orm:"column(target_amount)" json:"target_amount"`
	Status            int8    `orm:"column(status)" json:"status"`
	MembershipLevelID int64   `orm:"column(membership_level_id)" json:"-"`
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
	IsPassed           int8    `orm:"-" json:"is_passed"`
	CurrentPercentage  float64 `orm:"-" json:"current_percentage,omitempty"`
	RemainingAmount    float64 `orm:"-" json:"remaining_amount,omitempty"`
}
