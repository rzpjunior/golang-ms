// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import "time"

// Profile : struct to hold profile data
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

// CustomerProfileReturn : struct to hold data that will be returned by talon.one's update customer profile endpoint
type CustomerProfileReturn struct {
	CustomerProfile Profile `json:"customerProfile"`
}

// CustomerSessionReturn : struct to hold data that will be returned by talon.one's update customer session endpoint
type CustomerSessionReturn struct {
	CustomerSession struct {
		ID              int       `json:"id"`
		CreatedDate     time.Time `json:"created"`
		IntegrationCode string    `json:"integrationId"`
		ApplicationID   int       `json:"applicationId"`
		ProfileCode     string    `json:"profileId"`
		Attributes      struct {
			PointEarned      float64 `json:"eden_point_earned"`
			CountGetCampaign int     `json:"count_get_campaign"`
		} `json:"attributes"`
		TotalCharge   float64 `json:"total"`
		Subtotal      float64 `json:"cartItemTotal"`
		AdditionalFee float64 `json:"additionalCostTotal"`
	} `json:"customerSession"`
	CustomerProfile *Profile  `json:"customerProfile"`
	Effects         []Effects `json:"effects"`
}

// Effects : struct to hold effects of campaign from talon
type Effects struct {
	CampaignID int    `json:"campaignId"`
	EffectType string `json:"effectType"`
	Props      struct {
		Name                   string      `json:"name"`
		Value                  interface{} `json:"value"`
		RecipientIntegrationID string      `json:"recipientIntegrationId"`
		SubLedgerID            string      `json:"subLedgerId"`
	}
}

// SessionItemData : struct to hold item data for customer session request
type SessionItemData struct {
	ProductName  string
	ProductCode  string
	CategoryName string
	UnitPrice    float64
	OrderQty     float64
	UnitWeight   float64
	Attributes   map[string]string
}

// CampaignList : struct to hold campaign list data of talon.one
type CampaignList struct {
	Data []struct {
		ID            int    `json:"id"`
		ApplicationID int    `json:"applicationId"`
		State         string `json:"state"`
	} `json:"data"`
}

// CampaignDetail : struct to hold campaign detail data of talon.one
type CampaignDetail struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Tags       []string `json:"tags"`
	Attributes struct {
		Multiplier      int `json:"eden_point_multiplier"`
		MaxEarnPerTrans int `json:"max_earn_per_trans"`
	} `json:"attributes"`
}

// SessionResponse : struct to hold response of create session in talon.one
type SessionResponse struct {
	UserID    int       `json:"userId"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created"`
}

// ErrorResponse : struct to hold error response when talon.one's api
type ErrorResponse struct {
	Message string `json:"message"`
	Errors  []struct {
		Title   string `json:"title"`
		Details string `json:"details"`
	} `json:"errors"`
	Code int `json:"StatusCode"`
}

// CustomerProfileData : struct to hold data that will be returned by talon.one's get customer profile endpoint
type CustomerProfileData struct {
	Profile Profile `json:"profile"`
}

// Attribute : struct to hold data of custom attribute in talon.one
type Attribute struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created"`
	Entity      string    `json:"entity"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
}
