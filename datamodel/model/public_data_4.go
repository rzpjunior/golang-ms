// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

type PublicData4 struct {
	Status  string      `json:"status"`
	Message interface{} `json:"message"`
	Data    struct {
		PriceDate   string `json:"price_date"`
		CompareDate string `json:"compare_date"`
		Location    string `json:"location"`
		Prices      []struct {
			CommodityID        string `json:"commodity_id"`
			Name               string `json:"name"`
			ImagePath          string `json:"image_path"`
			ImageURL           string `json:"image_url"`
			Unit               string `json:"unit"`
			Price              int    `json:"price"`
			PriceCompare       int    `json:"price_compare"`
			Changed            int    `json:"changed"`
			Status             string `json:"status"`
			NotifPrice         int    `json:"notif_price"`
			NotifPercentage    int    `json:"notif_percentage"`
			NotifDuration      int    `json:"notif_duration"`
			NotifLimit         int    `json:"notif_limit"`
			NotifDecreases     int    `json:"notif_decreases"`
			RisePercentage     string `json:"rise_percentage"`
			RiseDuration       int    `json:"rise_duration"`
			DecreasePercentage string `json:"decrease_percentage"`
		} `json:"prices"`
	} `json:"data"`
}
