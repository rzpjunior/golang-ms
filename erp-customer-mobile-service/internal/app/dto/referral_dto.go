package dto

// ReferralHistoryReturn : struct to hold data that will return referral point history
type ReferralHistoryReturn struct {
	Summary struct {
		TotalPoint    float64 `json:"total_point"`
		TotalReferral int64   `json:"total_referral"`
	} `json:"summary"`
	Detail struct {
		ReferralList      []*ReferralList      `json:"referral_list"`
		ReferralPointList []*ReferralPointList `json:"referral_point_list"`
	} `json:"detail"`
}

// ReferralList : struct to hold referral list data
type ReferralList struct {
	Name      string `orm:"column(name)" json:"name"`
	CreatedAt string `orm:"column(created_at)" json:"created_at"`
}

// ReferralPointList : struct to hold referral point list data
type ReferralPointList struct {
	Name       string `orm:"column(name)" json:"name"`
	CreatedAt  string `orm:"column(created_date)" json:"created_at"`
	PointValue string `orm:"column(point_value)" json:"point"`
}
