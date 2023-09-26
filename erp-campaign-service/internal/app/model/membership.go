package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type MembershipLevel struct {
	ID       int64  `orm:"column(id);auto" json:"-"`
	Code     string `orm:"column(code)" json:"code"`
	Level    int8   `orm:"column(level)" json:"level"`
	Name     string `orm:"column(name)" json:"name"`
	ImageUrl string `orm:"column(image_url)" json:"image_url"`
	Status   int8   `orm:"column(status)" json:"status"`
}

func init() {
	orm.RegisterModel(new(MembershipLevel))
}

func (m *MembershipLevel) TableName() string {
	return "membership_level"
}

type MembershipCheckpoint struct {
	ID                int64   `orm:"column(id);auto" json:"id"`
	Checkpoint        int8    `orm:"column(checkpoint)" json:"checkpoint"`
	TargetAmount      float64 `orm:"column(target_amount)" json:"target_amount"`
	Status            int8    `orm:"column(status)" json:"status"`
	MembershipLevelID int64   `orm:"column(membership_level_id)" json:"membership_level_id"`
}

func init() {
	orm.RegisterModel(new(MembershipCheckpoint))
}

func (m *MembershipCheckpoint) TableName() string {
	return "membership_checkpoint"
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

func init() {
	orm.RegisterModel(new(MembershipAdvantage))
}

func (m *MembershipAdvantage) TableName() string {
	return "membership_advantage"
}

// MembershipAdvantage : struct to hold membership advantage data
type MembershipLevelAdvantage struct {
	ID                    int64 `orm:"column(id);auto" json:"-"`
	MembershipLevelID     int64 `orm:"column(membership_level_id)" json:"membership_level_id"`
	MembershipAdvantageID int64 `orm:"column(membership_advantage_id)" json:"membership_advantage_id"`
}

func init() {
	orm.RegisterModel(new(MembershipLevelAdvantage))
}

func (m *MembershipLevelAdvantage) TableName() string {
	return "membership_level_advantage"
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

func init() {
	orm.RegisterModel(new(MembershipReward))
}

func (m *MembershipReward) TableName() string {
	return "membership_reward"
}

type CustomerProfileData struct {
	Profile Profile `json:"profile"`
}

type Profile struct {
	ID                int         `json:"id"`
	CreatedDate       time.Time   `json:"created"`
	IntegrationID     string      `json:"integrationId"`
	Attributes        interface{} `json:"attributes"`
	AccountID         int         `json:"accountId"`
	ClosedSessions    int         `json:"closedSessions"`
	TotalSales        int         `json:"totalSales"`
	LoyaltyMembership string      `json:"loyaltyMemberships"`
	LastActivity      time.Time   `json:"lastActivity"`
}
