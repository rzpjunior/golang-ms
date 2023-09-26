package dto

import "git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/model"

type CustomerPointLogResponse struct {
	ID               int64   `json:"id"`
	CustomerID       int64   `json:"merchant_id"`
	SalesOrderID     int64   `json:"sales_oder_id"`
	EPCampaignID     int64   `json:"eden_point_campaign_id"`
	PointValue       float64 `json:"point_value"`
	RecentPoint      float64 `json:"recent_point"`
	Status           int8    `json:"status"`
	CreatedDate      string  `json:"created_date"`
	ExpiredDate      string  `json:"expired_date"`
	Note             string  `json:"note"`
	CurrentPointUsed float64 `json:"current_point_used"`
	NextPointUsed    float64 `json:"next_point_used"`
	TransactionType  int8    `json:"transaction_type"`
}

type CustomerPointLogRequestGet struct {
	Limit           int64  `json:"limit"`
	Offset          int64  `json:"offset"`
	CustomerID      int64  `json:"merchant_id"`
	SalesOrderID    int64  `json:"sales_oder_id"`
	Status          int8   `json:"status"`
	OrderBy         string `json:"order_by"`
	TransactionType int8   `json:"transaction_type"`
	CreatedDate     string `json:"created_date"`
}

type CustomerPointLogRequestGetDetail struct {
	ID              int64  `json:"id"`
	CustomerID      int64  `json:"merchant_id"`
	SalesOrderID    int64  `json:"sales_oder_id"`
	Status          int8   `json:"status"`
	CreatedDate     string `json:"created_date"`
	TransactionType int8   `json:"transaction_type"`
}

type CustomerPointLogRequestCreate struct {
	CustomerID       int64   `json:"merchant_id"`
	SalesOrderID     int64   `json:"sales_oder_id"`
	EPCampaignID     int64   `json:"eden_point_campaign_id"`
	PointValue       float64 `json:"point_value"`
	RecentPoint      float64 `json:"recent_point"`
	Status           int8    `json:"status"`
	CreatedDate      string  `json:"created_date"`
	ExpiredDate      string  `json:"expired_date"`
	Note             string  `json:"note"`
	CurrentPointUsed float64 `json:"current_point_used"`
	NextPointUsed    float64 `json:"next_point_used"`
	TransactionType  int8    `json:"transaction_type"`
}

type PointHistoryList struct {
	CreatedDate string  ` json:"created_date"`
	PointValue  float64 ` json:"point_value"`
	StatusType  string  ` json:"status_type"`
	Status      int8    ` json:"status"`
}

// ReferralHistoryReturn : struct to hold data that will return referral point history
type ReferralHistoryReturn struct {
	Summary struct {
		TotalPoint    float64 `json:"total_point"`
		TotalReferral int64   `json:"total_referral"`
	} `json:"summary"`
	Detail struct {
		ReferralList      []*model.ReferralList      `json:"referral_list"`
		ReferralPointList []*model.ReferralPointList `json:"referral_point_list"`
	} `json:"detail"`
}

type CustomerPointLogRequestCancel struct {
	CustomerID       int64   `json:"merchant_id"`
	SalesOrderID     int64   `json:"sales_oder_id"`
	PointValue       float64 `json:"point_value"`
	RecentPoint      float64 `json:"recent_point"`
	Status           int8    `json:"status"`
	CreatedDate      string  `json:"created_date"`
	ExpiredDate      string  `json:"expired_date"`
	Note             string  `json:"note"`
	CurrentPointUsed float64 `json:"current_point_used"`
	NextPointUsed    float64 `json:"next_point_used"`
	TransactionType  int8    `json:"transaction_type"`
}
