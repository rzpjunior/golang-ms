package dto

import "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"

type MembershipLevelResponse struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	Level    int8   `json:"level"`
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
	Status   int8   `json:"status"`
}

type MembershipCheckpointResponse struct {
	ID                int64   `json:"id"`
	Checkpoint        int8    `json:"checkpoint"`
	TargetAmount      float64 `json:"target_amount"`
	Status            int8    `json:"status"`
	MembershipLevelID int64   `json:"membership_level_id"`
}

// requestGetList : struct to hold request data
type RequestGetMembershipList struct {
	Platform string               `json:"platform" valid:"required"`
	Data     dataGetMembership    `json:"data" valid:"required"`
	Session  *SessionDataCustomer `json:"-"`
}

type dataGetMembership struct {
	Level      string `json:"level"`
	Checkpoint string `json:"checkpoint"`

	DataResponse []*model.MembershipLevel
}

type ResponseMembershipLevelList struct {
	ID       int64  `json:"-"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Level    string `json:"level"`
	ImageUrl string `json:"image_url"`
	Status   string `json:"status"`

	MembershipAdvantages  []*MembershipAdvantage  `orm:"-" json:"membership_advantage,omitempty"`
	MembershipCheckpoints []*MembershipCheckpoint `orm:"-" json:"membership_checkpoint,omitempty"`
}

// requestGetRewardList : struct to hold request data
type RequestGetRewardList struct {
	Platform string               `json:"platform" valid:"required"`
	Session  *SessionDataCustomer `json:"-"`
}

// MembershipReward : struct to hold membership checkpoint data
type MembershipRewardList struct {
	ID                 string `json:"-"`
	OpenedImageUrl     string `json:"opened_image_url"`
	ClosedImageUrl     string `json:"closed_image_url"`
	BackgroundImageUrl string `json:"background_image_url"`
	RewardLevel        string `json:"reward_level"`
	MaxAmount          string `json:"max_amount"`
	Status             string `json:"status"`
	Description        string `json:"description"`
	IsPassed           string `json:"is_passed"`
	CurrentPercentage  string `json:"current_percentage,omitempty"`
	RemainingAmount    string `json:"remaining_amount,omitempty"`
}

// MembershipAdvantage : struct to hold membership advantage data
type MembershipAdvantage struct {
	ID          int64  `orm:"column(id);auto" json:"-"`
	Name        string `orm:"column(name)" json:"name"`
	Description string `orm:"column(description)" json:"description"`
	ImageUrl    string `orm:"column(image_url)" json:"image_url"`
	LinkUrl     string `orm:"column(link_url)" json:"link_url"`
	Status      string `orm:"column(status)" json:"status"`
}

// MembershipCheckpoint : struct to hold membership checkpoint data
type MembershipCheckpoint struct {
	ID                int64  `orm:"column(id);auto" json:"-"`
	Checkpoint        string `orm:"column(checkpoint)" json:"checkpoint"`
	TargetAmount      string `orm:"column(target_amount)" json:"target_amount"`
	Status            string `orm:"column(status)" json:"status"`
	MembershipLevelID string `orm:"column(membership_level_id)" json:"-"`
}

// CustomerMembership : struct to hold membership data of customer
type CustomerMembership struct {
	MembershipLevel      string `json:"membership_level"`
	MembershipLevelName  string `json:"membership_level_name"`
	MembershipCheckpoint string `json:"membership_checkpoint"`
	CheckpointPercentage string `json:"percentage_checkpoint"`
	CurrentAmount        string `json:"current_amount"`
	TargetAmount         string `json:"target_amount"`
}
