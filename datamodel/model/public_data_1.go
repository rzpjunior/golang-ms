// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

type PublicData1 struct {
	Data struct {
		Sellings struct {
			Count       int64 `json:"count"`
			TotalPages  int64 `json:"totalPages"`
			CurrentPage int64 `json:"currentPage"`
			Params      struct {
				From      int64       `json:"from"`
				Size      int64       `json:"size"`
				SortField string      `json:"sortField"`
				SortOrder string      `json:"sortOrder"`
				RegionID  int64       `json:"regionId"`
				SearchKey interface{} `json:"searchKey"`
			} `json:"params"`
			Items []*PublicData1Items `json:"items"`
		} `json:"sellings"`
	} `json:"data"`
}

type PublicData1Items struct {
	ID      int64 `json:"id"`
	Product struct {
		ID                   int64  `json:"id"`
		Name                 string `json:"name"`
		Slug                 string `json:"slug"`
		CommercialSkuContent string `json:"commercialSkuContent"`
		Brand                struct {
			Name string `json:"name"`
		} `json:"brand"`
		ProductImages []struct {
			IsDefault bool   `json:"isDefault"`
			ImageURL  string `json:"imageURL"`
		} `json:"productImages"`
		Unit struct {
			Description string `json:"description"`
		} `json:"unit"`
		ProductPackaging struct {
			ID         int64  `json:"id"`
			Name       string `json:"name"`
			Multiplier int64  `json:"multiplier"`
		} `json:"productPackaging"`
		Groups []struct {
			ID       int64  `json:"id"`
			Name     string `json:"name"`
			ImageURL string `json:"imageUrl"`
		} `json:"groups"`
		Grade struct {
			ID int64 `json:"id"`
		} `json:"grade"`
	} `json:"product"`
	ProductPrices []struct {
		ID           int64       `json:"id"`
		Discount     int64       `json:"discount"`
		DiscountType string      `json:"discountType"`
		MinQty       int64       `json:"minQty"`
		MaxQty       interface{} `json:"maxQty"`
		Price        int64       `json:"price"`
	} `json:"productPrices"`
	Discount                 float64     `json:"discount"`
	DiscountInRupiah         float64     `json:"discountInRupiah"`
	DiscountType             string      `json:"discountType"`
	MinOrder                 int64       `json:"minOrder"`
	MaxOrder                 int64       `json:"maxOrder"`
	IsActive                 bool        `json:"isActive"`
	ShowedPrice              float64     `json:"showedPrice"`
	ShowedPriceAfterDiscount float64     `json:"showedPriceAfterDiscount"`
	StockLowerLimit          interface{} `json:"stockLowerLimit"`
	StockQty                 float64     `json:"stockQty"`
}
