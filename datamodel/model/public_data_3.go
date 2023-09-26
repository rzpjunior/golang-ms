// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

type PublicData3 struct {
	Products []struct {
		ID          string   `json:"_id"`
		Activated   bool     `json:"activated"`
		GlobalSku   string   `json:"global_sku"`
		Sku         string   `json:"sku"`
		ProductID   string   `json:"product_id"`
		ProductName string   `json:"product_name"`
		Description string   `json:"description"`
		Images      []string `json:"images"`
		Thumbnail   struct {
			Original string `json:"original"`
			Small    string `json:"small"`
		} `json:"thumbnail"`
		Category struct {
			CategoryID string      `json:"category_id"`
			Title      string      `json:"title"`
			TitleIDN   interface{} `json:"title_IDN"`
		} `json:"category"`
		IsFresh              bool     `json:"is_fresh"`
		IsBulky              bool     `json:"is_bulky"`
		IsGift               bool     `json:"is_gift"`
		IsTaxed              bool     `json:"is_taxed"`
		IsKino               bool     `json:"is_kino"`
		Tags                 []string `json:"tags"`
		Weight               int      `json:"weight"`
		Length               int      `json:"length"`
		Width                int      `json:"width"`
		Height               int      `json:"height"`
		Location             string   `json:"location"`
		ProductCount         int      `json:"product_count"`
		B2CDailyLimit        int      `json:"b2c_daily_limit"`
		B2BDailyLimit        int      `json:"b2b_daily_limit"`
		B2CLabel             string   `json:"b2c_label"`
		B2BLabel             string   `json:"b2b_label"`
		OriginalPrice        int      `json:"original_price"`
		ResalePrice          int      `json:"resale_price"`
		FlashSalePrice       int      `json:"flash_sale_price"`
		B2BOriginalPrice     int      `json:"b2b_original_price"`
		B2BResalePrice       int      `json:"b2b_resale_price"`
		DailySalePrice       int      `json:"daily_sale_price"`
		CommissionPercent    float64  `json:"commission_percent"`
		UseCommissionPercent int      `json:"use_commission_percent"`
		DiscountDate         struct {
			StartDate interface{} `json:"startDate"`
			EndDate   interface{} `json:"endDate"`
		} `json:"discountDate"`
		FlashSaleDate struct {
			StartDate interface{} `json:"startDate"`
			EndDate   interface{} `json:"endDate"`
		} `json:"flashSaleDate"`
		B2BDiscountDate struct {
			StartDate interface{} `json:"startDate"`
			EndDate   interface{} `json:"endDate"`
		} `json:"b2b_discountDate"`
		DailySaleDate struct {
			StartDate interface{} `json:"startDate"`
			EndDate   interface{} `json:"endDate"`
		} `json:"dailySaleDate"`
		PricingSchemeB2C []interface{} `json:"pricing_scheme_b2c"`
		PricingSchemeB2B []interface{} `json:"pricing_scheme_b2b"`
		SortOrder        int           `json:"sortOrder"`
	} `json:"products"`
	Pages int `json:"pages"`
	Total int `json:"total"`
}
