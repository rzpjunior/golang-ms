package dto

import "time"

type PackingOrderResponse struct {
	ID                    int64                            `json:"id"`
	Code                  string                           `json:"code"`
	Note                  string                           `json:"note"`
	DeliveryDate          time.Time                        `json:"delivery_date"`
	Status                int8                             `json:"status"`
	StatusConvert         string                           `json:"status_convert"`
	Site                  *SiteResponse                    `json:"site,omitempty"`
	Region                *RegionResponse                  `json:"region,omitempty"`
	PackingRecommendation []*PackingRecommendationResponse `json:"packing_recommendation"`
}

type PackingOrderResponseExport struct {
	Url string `json:"url"`
}

type RegionResponse struct {
	ID            string `json:"id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Status        int8   `json:"status"`
	StatusConvert string `json:"status_convert"`
}

type UomResponse struct {
	ID             string `json:"id"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	DecimalEnabled int    `json:"decimal_enabled"`
	Note           string `json:"note"`
	Status         int    `json:"status"`
}

type PackingOrderRequestGenerate struct {
	DeliveryDate string `json:"delivery_date" valid:"required"`
	RegionID     string `json:"region_id" valid:"required"`
	SiteID       string `json:"site_id" valid:"required"`
	Note         string `json:"note" valid:"lte:350"`
}
