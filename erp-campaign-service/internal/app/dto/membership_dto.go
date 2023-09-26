package dto

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

type MembershipLevelRequestGet struct {
	Limit   int64  `json:"limit"`
	Offset  int64  `json:"offset"`
	Status  int8   `json:"status"`
	Search  string `json:"search"`
	OrderBy string `json:"order_by"`
}

type MembershipCheckpointRequestGet struct {
	Limit             int64  `json:"limit"`
	Offset            int64  `json:"offset"`
	Status            int8   `json:"status"`
	OrderBy           string `json:"order_by"`
	MembershipLevelID int64  `json:"membership_level_id"`
	ID                int64  `json:"-"`
}

type MembershipAdvantage struct {
	ID          int64  `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
	LinkUrl     string `json:"link_url"`
	Status      int8   `json:"status"`
}

// MembershipAdvantage : struct to hold membership advantage data
type MembershipLevelAdvantage struct {
	ID                    int64 `json:"-"`
	MembershipLevelID     int64 `json:"membership_level_id"`
	MembershipAdvantageID int64 `json:"membership_advantage_id"`
}

// MembershipReward : struct to hold membership checkpoint data
type MembershipReward struct {
	ID                 int64   `json:"-"`
	OpenedImageUrl     string  `json:"opened_image_url"`
	ClosedImageUrl     string  `json:"closed_image_url"`
	BackgroundImageUrl string  `json:"background_image_url"`
	RewardLevel        int8    `json:"reward_level"`
	MaxAmount          float64 `json:"max_amount"`
	Status             int8    `json:"status"`
	Description        string  `json:"description"`
	IsPassed           int8    `json:"is_passed"`
	CurrentPercentage  float64 `json:"current_percentage,omitempty"`
	RemainingAmount    float64 `json:"remaining_amount,omitempty"`
}

// CustomerMembership : struct to hold membership data of customer
type CustomerMembership struct {
	MembershipLevel      int8    `json:"membership_level"`
	MembershipLevelName  string  `json:"membership_level_name"`
	MembershipCheckpoint int8    `json:"membership_checkpoint"`
	CheckpointPercentage float64 `json:"percentage_checkpoint"`
	CurrentAmount        float64 `json:"current_amount"`
	TargetAmount         float64 `json:"target_amount"`
}
