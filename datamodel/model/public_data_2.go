// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

type PublicData2 struct {
	Data struct {
		CatalogVariantList struct {
			Limit       int  `json:"limit"`
			Page        int  `json:"page"`
			Size        int  `json:"size"`
			HasNextPage bool `json:"hasNextPage"`
			Category    struct {
				DisplayName string `json:"displayName"`
			} `json:"category"`
			List []*PublicData2Items `json:"list"`
		} `json:"catalogVariantList"`
	} `json:"data"`
}

type PublicData2Items struct {
	Key          string   `json:"key"`
	Availability bool     `json:"availability"`
	Categories   []string `json:"categories"`
	Farmers      []struct {
		Image string `json:"image"`
		Name  string `json:"name"`
	} `json:"farmers"`
	Image struct {
		Md string `json:"md"`
		Sm string `json:"sm"`
		Lg string `json:"lg"`
	} `json:"image"`
	IsDiscount           bool        `json:"isDiscount"`
	Discount             int64       `json:"discount"`
	LabelDesc            string      `json:"labelDesc"`
	LabelName            string      `json:"labelName"`
	MaxQty               int64       `json:"maxQty"`
	Name                 string      `json:"name"`
	DisplayName          string      `json:"displayName"`
	NextAvailableDates   []string    `json:"nextAvailableDates"`
	PackDesc             string      `json:"packDesc"`
	PackNote             string      `json:"packNote"`
	Price                string      `json:"price"`
	PriceFormatted       string      `json:"priceFormatted"`
	ActualPrice          string      `json:"actualPrice"`
	ActualPriceFormatted string      `json:"actualPriceFormatted"`
	ShortDesc            string      `json:"shortDesc"`
	StockAvailable       int         `json:"stockAvailable"`
	Type                 string      `json:"type"`
	EmptyMessageHTML     interface{} `json:"emptyMessageHtml"`
	PromoMessageHTML     string      `json:"promoMessageHtml"`
}

type TokenData struct {
	Kind         string `json:"kind"`
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
}
